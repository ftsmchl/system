package host

import (
	"bufio"
	"fmt"
	"net"
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
	message, _ := bufio.NewReader(conn).ReadString('\n')
	messageBytesN := []byte(message)
	messageBytes := messageBytesN[:len(messageBytesN)-1]
	fmt.Println("message received from renter : ", message, "message in bytes containing n : ", messageBytesN, "message in bytes without N : ", messageBytes)

	//receive actual data from renter
	data, _ := bufio.NewReader(conn).ReadString('\n')
	dataBytesN := []byte(data)
	dataBytes := dataBytesN[:len(dataBytesN)-1]
	fmt.Println("data received : ", dataBytes)
	//data2, _ := bufio.NewReader(conn).ReadString('\n')
	//data2BytesN := []byte(data2)
	//data2Bytes := data2BytesN[:len(data2BytesN)-1]
	//fmt.Println("data received : ", data2Bytes)
}
