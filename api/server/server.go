package server

import (
	"context"
	"github.com/ftsmchl/system/api"
	"net"
	"net/http"
	"sync"
)

//serve is a wrapper of http.Serve()
func (srv *Server) serve() error {
	//Serve is a blocking function , will run until an error is encountered
	//or listener is closed via the close method.
	err := srv.apiServer.Serve(srv.listener)
	return err
}

//Close shutsdown the server
func (srv *Server) Close() error {
	srv.closeMu.Lock()
	defer srv.closeMu.Unlock()
	//Closing the server
	err := srv.apiServer.Shutdown(context.Background())
	//Waiting for serve to return
	<-srv.done
	//Check if server was closed properly
	if err != http.ErrServerClosed {
		return err
	}
	return nil

}

// Server is a collection of systemdaemon modules that can be
//communicated over an http api
type Server struct {
	api       *api.API
	apiServer *http.Server

	//channel that helps to understand when the server is closed
	done     chan struct{}
	listener net.Listener
	closeMu  sync.Mutex
}

func New() (*Server, error) {
	//Create the api which will create the renter and host modules as well
	api := api.New()
	//Create the routes of our handlers
	api.BuildRoutes()

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}

	srv := &Server{
		api: api,
		apiServer: &http.Server{
			Handler: api.Router,
		},
		listener: listener,
		done:     make(chan struct{}),
	}

	go func() {
		_ = srv.serve()
		//inform that the serve function has returned
		close(srv.done)
	}()
	return srv, nil
}
