package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func versionCmd(*cobra.Command, []string) {
	fmt.Println("Hello vriskesai sto pio prwtogono version")
}

func main() {
	root := &cobra.Command{
		Use:   os.Args[0],
		Short: "Start system",
		Long:  "Start system server and api",
		Run:   startDaemonCmd,
	}

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Long:  "Print version info about the system",
		Run:   versionCmd,
	})
	if err := root.Execute(); err != nil {
		fmt.Println("Kalhspera")
		os.Exit(64)
	}
}
