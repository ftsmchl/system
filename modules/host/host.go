package host

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ftsmchl/system/modules/wallet"
)

type Host struct {
	mu sync.Mutex

	auctionsBid map[string]AuctionContract

	storageContracts map[string]StorageContract

	fileContractRevisions map[string]*contractRevision

	contractRoots map[string]*merkleRoots

	//mu sync.Mutex
	listener net.Listener

	//our wallet
	wallet *wallet.Wallet
}

type AuctionContract struct {
	Address    string `json : "address"`
	TaskID     string `json : "taskid"`
	Owner      string `json : "owner"`
	Duration   int    `json : "duration"`
	InitialBid int    `json : "initialbid"`
}

type Winner struct {
	WinningBidder string `json : "winningbidder"`
	Address       string `json : "address"`
}

type StorageContract struct {
	Address  string
	TaskID   string
	Owner    string
	Duration int
}

type contractRevision struct {
	revisionNumber  int
	numLeaves       int
	merkleRoot      []byte
	signatureRenter string
	signatureHost   string
}

type merkleRoots struct {
	//	merkleTree     *my_merkleTree.Tree
	sectorRoots    [][]byte
	numMerkleRoots int
}

func New(wal *wallet.Wallet) *Host {
	h := &Host{
		auctionsBid: make(map[string]AuctionContract),

		//[taskID]StorageContract maps a taskID with a StorageContract
		storageContracts: make(map[string]StorageContract),

		fileContractRevisions: make(map[string]*contractRevision),

		contractRoots: make(map[string]*merkleRoots),

		wallet: wal,
	}

	err := h.initNetworking("0.0.0.0:8087")
	if err != nil {
		fmt.Println("init networking could not be initiated properly")
	}

	return h
}

//connects its ethereum account address with its ip to the market
func (h *Host) Register(acc string) {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	var ourIP string
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ourIP = ipv4.String()
			fmt.Println("Our IPV4 is : ", ourIP)
		}
	}

	fmt.Println("I am inside Register func")

	resp, err := http.Get("http://localhost:8001/hostRegister?IP=" + ourIP + "&ethereumAddress=" + acc)
	if err != nil {
		fmt.Println("err from httpGet : ", err)
		return
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	if string(text) == "OK" {
		fmt.Println("Read OK !!")
	} else {
		fmt.Println("We did not read OK!!")
	}

}

func (h *Host) FindContracts(acc string) {
	//fmt.Println("vim-go")
	//desiredMinimumBid := 1000
	//host must be a server listening to all new storage auctions

	resp, err := http.Get("http://localhost:8001/findAuction?maximumBid=1000&ethereumAddress=" + acc)
	if err != nil {
		fmt.Println(err)
	}

	var auctionContract AuctionContract
	var storageContract StorageContract
	var winner Winner

	defer resp.Body.Close()
	text, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(text), &auctionContract)
	fmt.Println("diavasa to text mages mou  p prepei na einai OK")

	address := auctionContract.Address
	taskid := auctionContract.TaskID
	duration := auctionContract.Duration
	owner := auctionContract.Owner
	initialbid := auctionContract.InitialBid

	h.auctionsBid[auctionContract.TaskID] = auctionContract

	//check if the auctionContract was encoded properly
	if address != "" && taskid != "" && duration != 0 && initialbid != 0 && owner != "" {
		fmt.Println("----------------")
		fmt.Println("Auction found!!")
		fmt.Println("Address : ", auctionContract.Address)
		fmt.Println("TaskID : ", auctionContract.TaskID)
		fmt.Println("Owner(renter) : ", auctionContract.Owner)
		fmt.Println("InitialBid : ", auctionContract.InitialBid)
		fmt.Println("Duration(ms) : ", auctionContract.Duration)
		fmt.Println("----------------")

		//time.Sleep(15 * time.Second)
		//fmt.Println("I slept for 15 seconds and i am gonna see who won the auction")
		time.Sleep(65 * time.Second)
		fmt.Println("I slept for 65 seconds and i am gonna see who won the auction")

		//check if we actually won then auction
		resp2, _ := http.Get("http://localhost:8001/checkWhoWonAuction?auctionAddress=" + auctionContract.Address)
		text2, _ := ioutil.ReadAll(resp2.Body)
		json.Unmarshal([]byte(text2), &winner)

		if winner.WinningBidder == acc && winner.Address != "" {
			fmt.Println("we actually won the auction")
			//we need lock here
			storageContract.Address = winner.Address
			storageContract.TaskID = auctionContract.TaskID
			storageContract.Owner = auctionContract.Owner
			storageContract.Duration = auctionContract.Duration

			h.storageContracts[auctionContract.TaskID] = storageContract
			fmt.Println("Our storage Contract")
			fmt.Println("----------------")
			fmt.Println("Contract Address : ", h.storageContracts[auctionContract.TaskID].Address)
			fmt.Println("Duration : ", h.storageContracts[auctionContract.TaskID].Duration)
			fmt.Println("Owner : ", h.storageContracts[auctionContract.TaskID].Owner)
			fmt.Println("TaskID : ", h.storageContracts[auctionContract.TaskID].TaskID)
			fmt.Println("----------------")

			h.mu.Lock()

			//add our taskID to contractRoots map
			h.contractRoots[storageContract.TaskID] = &merkleRoots{}

			//add a file contractRevision for this taskID
			h.fileContractRevisions[storageContract.TaskID] = &contractRevision{}

			h.mu.Unlock()

		} else {
			fmt.Println("we did not win the auction!!!")
		}
	}

	//we are going to activate our StorageContract

	resp3, _ := http.Get("http://localhost:8001/activateContract?contractAddress=" + storageContract.Address + "&ethereumAddress=" + acc)
	text3, _ := ioutil.ReadAll(resp3.Body)

	if string(text3) == "OK" {
		fmt.Println("storage contract activated by us !!")
	} else {
		fmt.Println("Something went wrong while activating the storage contract")
	}

}
