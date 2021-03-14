package main

import (
	"fmt"
	"github.com/ftsmchl/system/api/server"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

//creates a signal handler for os.Interrupt, os.Kill, and syscall.SIGTERM, and returns
//a channel that is closed when one of them is caught.
func installKillSignalHandler() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	return sigChan
}

func startDaemonCmd(cmd *cobra.Command, _ []string) {
	fmt.Println("3ekinhse to system")
	sigChan := installKillSignalHandler()

	//Start and run the api server
	srv, err := server.New()
	if err != nil {
		fmt.Println("Something went wrong while creating the server : ", err)
		os.Exit(64)
	}

	//waiting for a signal to shutdown the server
	<-sigChan
	err = func() error {
		return srv.Close()
	}()
	if err != nil {
		fmt.Println("Server has closed with error : ", err)
		os.Exit(64)
	}

	//Server has closed without errors
	fmt.Println("Shutdown  complete")
}
