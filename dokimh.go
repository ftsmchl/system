package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			//fmt.Printf("type of IPV4\n", ipv4)
			fmt.Println("IPV4: ", ipv4)
			//fmt.Println("IPV4 : ", []byte(ipv4))
			//fmt.Println("IPV$ string : ", string(ipv4))
			fmt.Println("IPV4 string : ", ipv4.String())
		}
	}
}
