package renter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Renter struct {
	auctionContracts   map[string]AuctionContract
	auctionContractsMu sync.Mutex

	storageContracts   map[string]StorageContract
	storageContractsMu sync.Mutex
}

//constructor of renter module
func New() *Renter {
	renter := &Renter{
		auctionContracts: make(map[string]AuctionContract),
		storageContracts: make(map[string]StorageContract),
	}
	return renter
}

type StorageContract struct {
	Address  string
	TaskID   string
	Duration int //(ms)
	Host     string
}

type AuctionContract struct {
	Address    string `json : "address"`
	TaskID     string `json : "taskid"`
	Owner      string `json : "owner"`
	Duration   int    `json : "duration"`
	InitialBid int    `json : "initialbid"`
	Host       string `json : "host"`
}

func (r *Renter) PrintContracts() {
	fmt.Println("Storage Contracts")
	r.storageContractsMu.Lock()
	for _, value := range r.storageContracts {
		fmt.Println("---------------------")
		fmt.Println("TaskID : ", value.TaskID)
		fmt.Println("Contract Address : ", value.Address)
		fmt.Println("Duration(ms) : ", value.Duration)
		fmt.Println("Host : ", value.Host)
		fmt.Println("---------------------")
		fmt.Println("")
	}
	r.storageContractsMu.Unlock()
}

//we will account address as an extra argument
func (r *Renter) AuctionCreate(wg *sync.WaitGroup, acc string) {
	defer wg.Done()

	var auctionContract AuctionContract
	var storageContract StorageContract

	//sends an http request and receives an auctionContract without the host
	//resp, err := http.Get("http://localhost:8000/auctionCreate")
	resp, err := http.Get("http://localhost:8000/auctionCreate?ethereumAddress=" + acc)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	text, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(text), &auctionContract)

	//r.auctionContracts = append(r.auctionContracts, auctionContract)
	//we need lock
	r.auctionContractsMu.Lock()
	r.auctionContracts[auctionContract.TaskID] = auctionContract
	r.auctionContractsMu.Unlock()

	address := auctionContract.Address
	taskid := auctionContract.TaskID
	duration := auctionContract.Duration
	initialbid := auctionContract.InitialBid
	owner := auctionContract.Owner

	//check if the auctionContract was encoded properly
	if address != "" && taskid != "" && duration != 0 && initialbid != 0 && owner != "" {
		fmt.Println("An auction contract has been created !!")
		fmt.Println("Address : ", auctionContract.Address)
		fmt.Println("TaskID : ", auctionContract.TaskID)
		fmt.Println("FileContract Duration(ms) : ", auctionContract.Duration)
		fmt.Println("Initial Bid : ", auctionContract.InitialBid)
		fmt.Println("Owner : ", auctionContract.Owner)
		fmt.Println("WinningBidder(host) : ", auctionContract.Host)

		//sleep until the auction is finished
		time.Sleep(10 * time.Second)
		fmt.Println("I slept for 10 seconds and i am about to finalize the auction with the host")

		//finalize the auction with the winningBidder if there is one
		resp2, err := http.Get("http://localhost:8000/auctionFinalize?auctionAddress=" + auctionContract.Address + "&ethereumAddress=" + acc)
		if err != nil {
			fmt.Println(err)
		}
		defer resp2.Body.Close()

		text2, err := ioutil.ReadAll(resp2.Body)
		//fmt.Println("O winningBidder tou auction to opoio molis egine finalized einai  : ", string(text2))

		//we need to check if text2 size is a valid ethereum account address
		fmt.Println("Size of ethereum address is : ", len(text2))

		//if len(text2) == 42 we have an auction winner
		if len(text2) == 42 {
			fmt.Println("O winningBidder tou auction to opoio molis egine finalized einai  : ", string(text2))
			storageContract.Address = auctionContract.Address
			storageContract.TaskID = auctionContract.TaskID
			storageContract.Duration = auctionContract.Duration
			storageContract.Host = string(text2)

			r.storageContractsMu.Lock()
			r.storageContracts[storageContract.TaskID] = storageContract
			r.storageContractsMu.Unlock()
		} else {
			fmt.Println("There is not a winning bidder for our auction")
		}

	} else {
		fmt.Println("Something went wrong while creating the auctionContract")
		fmt.Println("text that was read from renterServer is ", string(text))
	}

}
