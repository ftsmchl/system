package api

import (
	"fmt"
	"github.com/ftsmchl/system/modules/host"
	"github.com/ftsmchl/system/modules/renter"
	"github.com/ftsmchl/system/modules/wallet"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"sync"
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
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go api.renter.AuctionCreate(&wg, account)
	}

	wg.Wait()
	api.renter.PrintContracts()
}

func (api *API) hostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Starting host registration\n")

	//get our ethereum address
	account := api.wallet.GetPrimaryAccount()

	api.host.Register(account)
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
	acc := params["address"]
	api.wallet.SetPrimaryAccount(acc)
}

func (api *API) BuildRoutes() {
	api.Router = mux.NewRouter()
	//renter commands
	api.Router.HandleFunc("/score", api.scoreHandler)
	api.Router.HandleFunc("/createAuction", api.createAuctionHandler)

	//host commands
	api.Router.HandleFunc("hostRegister", api.hostRegisterHandler)
	api.Router.HandleFunc("/findContracts", api.findContractsHandler)

	//wallet commands
	api.Router.HandleFunc("/addAccount/{address}", api.setAccountAddressHandler)

}

type API struct {
	score  int
	renter *renter.Renter
	host   *host.Host
	wallet *wallet.Wallet
	Router *mux.Router
}

func New() *API {
	//creation of the renter
	r := renter.New()
	//creation of host
	h := host.New()

	//creation of our wallet
	w := wallet.New()

	return &API{
		score:  42,
		renter: r,
		host:   h,
		wallet: w,
	}
}
