package api

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ftsmchl/system/modules/host"
	"github.com/ftsmchl/system/modules/renter"
	"github.com/ftsmchl/system/modules/wallet"
	"github.com/gorilla/mux"
)

func (api *API) scoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "To score tou api einai %d\n", api.score)
}

func (api *API) createAuctionHandler(w http.ResponseWriter, r *http.Request) {
	//io.WriteString(w, "3ekinaei o renter thn dhmiourgia twn contract agoraki mou\n")

	//get our account address
	account := api.wallet.GetPrimaryAccount()

	var wg sync.WaitGroup
	//we are creating 2 contracts
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go api.renter.AuctionCreate(&wg, account)
	}

	wg.Wait()
	io.WriteString(w, "Contract creation is finished\n")
	api.renter.PrintContracts()
}

func (api *API) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pathTrimmed := params["path"]
	err := api.renter.Upload(pathTrimmed)
	if err != nil {
		fmt.Fprintf(w, "There was an error uploading the file : %s", err)
	} else {
		io.WriteString(w, "File was uploaded succesfully...")
	}
}

func (api *API) challengeHostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pathTrimmed := params["hostPublicKey"]
	//fmt.Println("host to be challenged (renter) : ", pathTrimmed)
	//io.WriteString(w, "Host has been challenged succesfully")
	fmt.Fprintf(w, "Host %s has been challenged succesfully ...", pathTrimmed)
}

func (api *API) hostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	//io.WriteString(w, "Starting host registration\n")

	//get our ethereum address
	account := api.wallet.GetPrimaryAccount()

	api.host.Register(account)

	io.WriteString(w, "Host registration has finished\n")
}

func (api *API) findContractsHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "starting to find contract to bid\n")

	//get our account address
	account := api.wallet.GetPrimaryAccount()

	api.host.FindContracts(account)
}

func (api *API) setAccountAddressHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Connecting ethereum account address...\n")
	params := mux.Vars(r)
	pubKey := params["publicKey"]
	privKey := params["privateKey"]
	api.wallet.SetPrimaryAccount(pubKey, privKey)
}

func (api *API) BuildRoutes() {
	api.Router = mux.NewRouter()

	//renter commands
	api.Router.HandleFunc("/score", api.scoreHandler)
	api.Router.HandleFunc("/createAuction", api.createAuctionHandler)
	api.Router.HandleFunc("/uploadFile/{path}", api.uploadFileHandler)
	api.Router.HandleFunc("/challengeHost/{hostPublicKey}", api.challengeHostHandler)

	//host commands
	api.Router.HandleFunc("/hostRegister", api.hostRegisterHandler)
	api.Router.HandleFunc("/findContracts", api.findContractsHandler)

	//wallet commands
	api.Router.HandleFunc("/addAccount/{publicKey}/{privateKey}", api.setAccountAddressHandler)

}

type API struct {
	score  int
	renter *renter.Renter
	host   *host.Host
	wallet *wallet.Wallet
	Router *mux.Router
}

func New() *API {
	//creation of wallet
	w := wallet.New()
	//creation of the renter
	r := renter.New(w)
	//creation of host
	h := host.New(w)

	return &API{
		score:  42,
		renter: r,
		host:   h,
		wallet: w,
	}
}
