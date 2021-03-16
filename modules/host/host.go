package host

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Host struct {
	auctionsBid      map[string]AuctionContract
	storageContracts map[string]StorageContract
}

type AuctionContract struct {
	Address    string `json : "address"`
	TaskID     string `json : "taskid"`
	Owner      string `json : "owner"`
	Duration   int    `json : "duration"`
	InitialBid int    `json : "initialbid"`
}

type StorageContract struct {
	Address  string
	TaskID   string
	Owner    string
	Duration int
}

func New() *Host {
	return &Host{
		auctionsBid:      make(map[string]AuctionContract),
		storageContracts: make(map[string]StorageContract),
	}
}

func (h *Host) FindContracts() {
	//fmt.Println("vim-go")
	//desiredMinimumBid := 1000
	//host must be a server listening to all new storage auctions

	resp, err := http.Get("http://localhost:8001/findAuction?maximumBid=1000")
	if err != nil {
		fmt.Println(err)
	}

	var auctionContract AuctionContract

	defer resp.Body.Close()
	text, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(text), &auctionContract)
	fmt.Println("diavasa to text mages mou  p prepei na einai OK")
	//fmt.Println("To text tou response einai : ", string(text))

	address := auctionContract.Address
	taskid := auctionContract.TaskID
	duration := auctionContract.Duration
	owner := auctionContract.Owner
	initialbid := auctionContract.InitialBid

	//h.auctionsBid = append(h.auctionsBid, auctionContract)
	h.auctionsBid[auctionContract.TaskID] = auctionContract

	//check if the auctionContract was encoded properly
	if address != "" && taskid != "" && duration != 0 && initialbid != 0 && owner != "" {
		fmt.Println("Auction found!!")
		fmt.Println("Address : ", auctionContract.Address)
		fmt.Println("TaskID : ", auctionContract.TaskID)
		fmt.Println("Owner(renter) : ", auctionContract.Owner)
		fmt.Println("InitialBid : ", auctionContract.InitialBid)
		fmt.Println("Duration(ms) : ", auctionContract.Duration)

		time.Sleep(15 * time.Second)
		fmt.Println("I slept for 15 seconds and i am gonna see who won the auction")

		resp2, _ := http.Get("http://localhost:8001/checkWhoWonAuction?auctionAddress=" + auctionContract.Address)
		text2, _ := ioutil.ReadAll(resp2.Body)
		fmt.Println("text2 is ", string(text2))
	}

}
