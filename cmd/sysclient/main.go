package main

import (
	"fmt"
	"github.com/ftsmchl/system/api/client"
	"github.com/spf13/cobra"
	"os"
)

var (
	httpClient client.Client
)

func statuscmd(cmd *cobra.Command, _ []string) {
	fmt.Println("This the first command of sysclient")
	fmt.Println("Http Client Address is : ", httpClient.Address)
}

func rentercmd(cmd *cobra.Command, _ []string) {
	fmt.Println("This is renter commmand answered")
}

func hostcmd(cmd *cobra.Command, _ []string) {
	fmt.Println("This is the host command answered")
}

func hostfindContractscmd(cmd *cobra.Command, _ []string) {
	fmt.Println("Listening for new contracts")
	err := httpClient.HostFindContracts()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func rentercreatecontractscmd(cmd *cobra.Command, _ []string) {
	fmt.Println("Creating Contracts, this may take a while..")
	err := httpClient.RenterCreateContracts()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Creating contracts operation has finished")
}

var (
	renterCmd = &cobra.Command{
		Use:   "renter",
		Short: "Perform renter actions",
		Long:  "Create Contracts, upload, download...",
		Run:   rentercmd,
	}

	hostCmd = &cobra.Command{
		Use:   "host",
		Short: "Perform host actions",
		Long:  "Find Contracts..",
		Run:   hostcmd,
	}

	renterCreateContractsCmd = &cobra.Command{
		Use:   "createContracts",
		Short: "Create contracts",
		Long:  "Create Contracts",
		Run:   rentercreatecontractscmd,
	}

	hostFindContractsCmd = &cobra.Command{
		Use:   "findContracts",
		Short: "Find Contracts",
		Long:  "Find Auction Contracts and bid them",
		Run:   hostfindContractscmd,
	}
)

func main() {
	root := &cobra.Command{
		Use:  os.Args[0],
		Long: "System Client v",
		Run:  statuscmd,
	}

	root.AddCommand(renterCmd)
	root.AddCommand(hostCmd)

	renterCmd.AddCommand(renterCreateContractsCmd)
	hostCmd.AddCommand(hostFindContractsCmd)
	//initialize client
	root.PersistentFlags().StringVarP(&httpClient.Address, "addr", "a", "localhost:8080", "host/port sysd is listening on ")

	//run
	if err := root.Execute(); err != nil {
		os.Exit(64)
	}
}
