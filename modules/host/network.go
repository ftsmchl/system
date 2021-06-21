package host

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func (h *Host) initNetworking(address string) (err error) {
	h.listener, err = net.Listen("tcp", address)
	if err != nil {
		fmt.Println("listener could not be instantiated properly")
	}

	go h.threadedListen()

	return nil
}

func (h *Host) threadedListen() {
	for {
		conn, err := h.listener.Accept()
		if err != nil {
			return
		}

		go h.threadedHandleConn(conn)
	}
}

func (h *Host) threadedHandleConn(conn net.Conn) {

	reader := bufio.NewReader(conn)

	msg1, err := reader.ReadString('\n')
	fmt.Println("We read msg1 : ", strings.TrimRight(msg1, "\n"))
	fmt.Println("err : ", err)

	if strings.TrimRight(msg1, "\n") == "Upload" {
		fmt.Println("mesa sto if")
		h.uploadProtocol(conn, reader)
	} else {

		//	h.uploadProtocol(conn, reader)
	}

	/*

		msg, err := reader.ReadBytes('\n')

		fmt.Println("msg read is ", msg)
		fmt.Println("err : ", err)
		fmt.Println("msg.len is : ", len(msg))
	*/

}

func (h *Host) uploadProtocol(c net.Conn, r *bufio.Reader) {
	//reader := bufio.NewReader(c)

	data := make([]byte, 4)
	n, err := r.Read(data)
	fmt.Println("n ", n)
	fmt.Println("data received : ", data)
	fmt.Println("err : ", err)
	fmt.Fprintf(c, "Data\n")

	//read renter's signature
	renterSignature, err := r.ReadString('\n')
	fmt.Println("err : ", err)
	fmt.Println("Renter's Signature is : ", renterSignature)

}
