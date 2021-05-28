package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func installKillSignalHandler() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	return sigChan
}

func findContracts(w http.ResponseWriter, r *http.Request) {
	file, _ := os.OpenFile("logs/sysclient.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	args := []string{"host", "findContracts"}

	cmd := exec.Command("./sysclient", args[0], args[1])

	cmd.Stdout = file
	err := cmd.Run()
	fmt.Println("cmd : ", cmd)
	fmt.Println("err : ", err)
}

func createContracts(w http.ResponseWriter, r *http.Request) {

	file, _ := os.OpenFile("logs/sysclient.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	args := []string{"renter", "createContracts"}

	cmd := exec.Command("./sysclient", args[0], args[1])
	cmd.Stdout = file
	err := cmd.Run()
	fmt.Println("cmd : ", cmd)
	fmt.Println("err : ", err)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from init_server\n")
	fmt.Println("M hr8e hello")
}

func Router(r *mux.Router) {
	//router = mux.NewRouter()
	r.HandleFunc("/createContracts", createContracts)
	r.HandleFunc("/findAuctions", findContracts)
	r.HandleFunc("/hello", hello)

}

type Server struct {
	server   *http.Server
	listener net.Listener
}

func (srv *Server) serve(listener net.Listener, done chan struct{}) {
	err := srv.server.Serve(listener)
	if err != nil {
		fmt.Println("err in server init : ", err)
	}
	fmt.Println("Returning from server")
	close(done)
}

func main() {
	//var router *mux.Router
	//Router(router)

	router := mux.NewRouter()
	//router.HandleFunc("/hello", hello)
	Router(router)

	done := make(chan struct{})

	listener, err := net.Listen("tcp", "0.0.0.0:8089")

	sigChan := installKillSignalHandler()

	if err != nil {
		fmt.Println(err)
		return
	}

	srv := &http.Server{
		Handler: router,
	}

	myServer := &Server{
		server:   srv,
		listener: listener,
	}

	go myServer.serve(listener, done)

	fmt.Println("Server is running.....")
	fmt.Println("We are listening on :8089..")

	timh := <-sigChan

	myServer.server.Shutdown(context.Background())

	fmt.Println("Server is closing ..", timh)
	<-done

}
