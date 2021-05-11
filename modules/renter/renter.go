package renter

import (
	"encoding/json"
	"fmt"
	"github.com/ftsmchl/system/modules/renter/renterfile"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Renter struct {
	auctionContracts   map[string]AuctionContract
	auctionContractsMu sync.Mutex

	//[taskID]StorageContract
	storageContracts   map[string]StorageContract
	storageContractsMu sync.Mutex

	//editors map[taskID]Editor
	editors map[string]*Editor

	//[publicKey]IPV4
	hosts map[string]string

	//List of workers for uploading/downloading
	workers map[string]*worker //[taskID] *worker

	//
	renterFile *renterfile.Renterfile

	uploadHeap uploadHeap

	mu sync.Mutex
}

//constructor of renter module
func New() *Renter {
	renter := &Renter{
		auctionContracts: make(map[string]AuctionContract),
		storageContracts: make(map[string]StorageContract),
		hosts:            make(map[string]string),
		workers:          make(map[string]*worker),
		editors:          make(map[string]*Editor),

		uploadHeap: uploadHeap{
			newUploads: make(chan struct{}, 1),
		},
	}

	go renter.threadedUpload()
	return renter
}

type StorageContract struct {
	Address  string
	TaskID   string
	Duration int //(ms)
	Host     string
	IP       string
}

type FinalAuction struct {
	WinningBidder string `json : "winningBidder"`
	Address       string `json : "address"`
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
	counter := 1
	for _, value := range r.storageContracts {
		fmt.Println("---------------------")
		fmt.Println("STORAGE CONTRACT No : ", counter)
		fmt.Println("TaskID : ", value.TaskID)
		fmt.Println("Contract Address : ", value.Address)
		fmt.Println("Duration(ms) : ", value.Duration)
		fmt.Println("Host : ", value.Host)
		fmt.Println("Host's IP : ", value.IP)
		fmt.Println("---------------------")
		fmt.Println("")
		fmt.Println("")
		counter++
	}
	r.storageContractsMu.Unlock()
	r.updateWorkerPool()
}

//we will account address as an extra argument
func (r *Renter) AuctionCreate(wg *sync.WaitGroup, acc string) {
	defer wg.Done()

	var auctionContract AuctionContract
	var storageContract StorageContract

	var winner FinalAuction

	fmt.Println("I am inside r.AuctionCreate()")
	//sends an http request and receives an auctionContract without the host
	resp, err := http.Get("http://localhost:8000/auctionCreate?ethereumAddress=" + acc)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Phra to gamhmeno mhnuma")
	defer resp.Body.Close()

	text, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(text), &auctionContract)

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
		time.Sleep(60 * time.Second)
		fmt.Println("I slept for 60 seconds and i am about to finalize the auction with the host")

		//finalize the auction with the winningBidder if there is one
		resp2, err := http.Get("http://localhost:8000/auctionFinalize?auctionAddress=" + auctionContract.Address + "&ethereumAddress=" + acc)
		if err != nil {
			fmt.Println(err)
		}
		defer resp2.Body.Close()

		text2, err := ioutil.ReadAll(resp2.Body)
		json.Unmarshal([]byte(text2), &winner)

		//we need to check if text2 size is a valid ethereum account address
		if len(winner.WinningBidder) == 42 && winner.Address != "" {
			fmt.Println("O winningBidder tou auction to opoio molis egine finalized einai  : ", winner.WinningBidder)
			storageContract.Address = winner.Address
			storageContract.TaskID = auctionContract.TaskID
			storageContract.Duration = auctionContract.Duration
			storageContract.Host = string(winner.WinningBidder)

			/*
				r.storageContractsMu.Lock()
				r.storageContracts[storageContract.TaskID] = storageContract
				r.storageContractsMu.Unlock()
			*/
		} else {
			fmt.Println("There is not a winning bidder for our auction")
		}

	} else {
		fmt.Println("Something went wrong while creating the auctionContract")
		fmt.Println("text that was read from renterServer is ", string(text))
	}

	resp3, err := http.Get("http://localhost:8000/providerIP?publicKey=" + storageContract.Host)
	if err != nil {
		fmt.Println("Error in providerIP GET : ", err)
	}
	defer resp3.Body.Close()

	text3, err := ioutil.ReadAll(resp3.Body)
	if string(text3) == "!OK" {
		fmt.Println("We could not read the ip of our provider properly")
	} else {
		storageContract.IP = string(text3)
	}

	r.storageContractsMu.Lock()
	r.storageContracts[storageContract.TaskID] = storageContract
	r.storageContractsMu.Unlock()

}
