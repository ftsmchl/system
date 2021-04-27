package main

import (
	"fmt"
	"github.com/ftsmchl/system/api/client"
	"github.com/ftsmchl/system/helper"
	"github.com/spf13/cobra"
	"os"
	//"strings"
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
	//resp, err := http.Get("http://localhost:8000/auctionCreate")
}

func hostregisterToMarketcmd(cmd *cobra.Command, _ []string) {
	fmt.Println("Registering to market...")
	err := httpClient.HostRegister()
	if err != nil {
		fmt.Println(err)
	}
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

func renteruploadfilecmd(cmd *cobra.Command, args []string) {
	fmt.Println("Preparing to upload ", args[0])
	//pathTrimmed := strings.TrimPrefix(args[0], "/")
	err := httpClient.RenterUploadFile(args[0])
	if err != nil {
		fmt.Println("file was not uploaded succesfully .. ", err)
	} else {
		fmt.Println("file was uploaded succesfully.. ")
	}

}

func accountaddcmd(cmd *cobra.Command, args []string) {
	fmt.Println("Connecting our node with an ethereum account..")

	//checks if the string given is a valid ethereum address
	if len(args[0]) == 42 && helper.IsHex(args[0]) {
		fmt.Println("args[0]", args[0])
	} else {
		fmt.Println("The address is not a valid ethereum address, try again giving a valid one!!")
	}

	err := httpClient.WalletAddAccount(args[0])
	if err != nil {
		//resp, err := http.Get("http://localhost:8000/auctionCreate")
		fmt.Println("Something went wrong while adding our ethereum account address :", err)
	} else {
		fmt.Println("The ethereum account address has been set succesfully")
	}
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

	accountAddCmd = &cobra.Command{
		Use:   "accountAdd",
		Short: "Add an account for this node",
		Long:  "Add an account for this node",
		Run:   accountaddcmd,
	}

	renterCreateContractsCmd = &cobra.Command{
		Use:   "createContracts",
		Short: "Create contracts",
		Long:  "Create Contracts",
		Run:   rentercreatecontractscmd,
	}

	renterUploadFileCmd = &cobra.Command{
		Use:   "uploadFile",
		Short: "upload a file",
		Long:  "upload a file using erasure coding",
		Run:   renteruploadfilecmd,
	}

	hostRegisterToMarketCmd = &cobra.Command{
		Use:   "register",
		Short: "Register to market",
		Long:  "Connects its ethereum address(public Key) with its personal ip",
		Run:   hostregisterToMarketcmd,
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
	root.AddCommand(accountAddCmd)

	renterCmd.AddCommand(renterCreateContractsCmd)
	renterCmd.AddCommand(renterUploadFileCmd)

	hostCmd.AddCommand(hostFindContractsCmd)
	hostCmd.AddCommand(hostRegisterToMarketCmd)

	//initialize client
	root.PersistentFlags().StringVarP(&httpClient.Address, "addr", "a", "localhost:8080", "host/port sysd is listening on ")

	//run
	if err := root.Execute(); err != nil {
		os.Exit(64)
	}
}
