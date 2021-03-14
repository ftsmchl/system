package api

import (
	"fmt"
	"github.com/ftsmchl/system/modules/host"
	"github.com/ftsmchl/system/modules/renter"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (api *API) ScoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "To score tou api einai %d\n", api.score)
}

func (api *API) CreateAuctionHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "3ekinaei o renter thn dhmiourgia twn contract agoraki mou\n")
	api.renter.AuctionCreate()
}

func (api *API) findContractsHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "starting to find contract to bid\n")
	api.host.FindContracts()
}

func (api *API) BuildRoutes() {
	api.Router = mux.NewRouter()
	//renter commands
	api.Router.HandleFunc("/score", api.ScoreHandler)
	api.Router.HandleFunc("/createAuction", api.CreateAuctionHandler)
	api.Router.HandleFunc("/findContracts", api.findContractsHandler)

}

type API struct {
	score  int
	renter *renter.Renter
	host   *host.Host
	Router *mux.Router
}

func New() *API {
	//creation of the renter
	r := renter.New()
	//creation of host
	h := host.New()
	return &API{
		score:  42,
		renter: r,
		host:   h,
	}
}
