package host

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Host struct {
	auctionsBid []AuctionContract
}

type AuctionContract struct {
	Address    string `json : "address"`
	TaskID     string `json : "taskid"`
	Owner      string `json : "owner"`
	Duration   int    `json : "duration"`
	InitialBid int    `json : "initialbid"`
}

func New() *Host {
	return &Host{}
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

	h.auctionsBid = append(h.auctionsBid, auctionContract)
	fmt.Println("To address tou auction einai : ", auctionContract.Address)
	fmt.Println("To taskID tou auction einai : ", auctionContract.TaskID)
	fmt.Println("O owner  tou auction einai : ", auctionContract.Owner)
	fmt.Println("To initialBid  tou auction einai : ", auctionContract.InitialBid)
	fmt.Println("To Duration  tou auction einai : ", auctionContract.Duration)

	time.Sleep(15 * time.Second)
	fmt.Println("I slept for 15 seconds and i am gonna see who won the auction")

	resp2, err := http.Get("http://localhost:8001/checkWhoWonAuction?auctionAddress=" + auctionContract.Address)
	text2, _ := ioutil.ReadAll(resp2.Body)
	fmt.Println("text2 is ", string(text2))

}
