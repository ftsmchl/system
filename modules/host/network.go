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

	//file, _ := os.OpenFile("logs/data", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	/*
		message, _ := bufio.NewReader(conn).ReadString('\n')
		messageBytesN := []byte(message)
		messageBytes := messageBytesN[:len(messageBytesN)-1]
		fmt.Println("message received from renter : ", message, "message in bytes containing n : ", messageBytesN, "message in bytes without N : ", messageBytes)

	*/
	//receive actual data from renter

	/*
		data := make([]byte, 4)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("Wrong reading from conn : ", err)
		}
		fmt.Println("We read : ", n, " bytes")
	*/
	//dataBytesN := []byte(data)
	//dataBytes := dataBytesN[:len(dataBytesN)-1]
	//fmt.Println("data received : ", data)

	msg, err := bufio.NewReader(conn).ReadBytes('\n')
	fmt.Println("err : ", err)

	fmt.Println("msg read is ", msg)
	fmt.Println("msg.len is : ", len(msg))

	/*
		data2, _ := bufio.NewReader(conn).ReadString('\n')
		data2BytesN := []byte(data2)
		data2Bytes := data2BytesN[:len(data2BytesN)-1]
		fmt.Println("data received : ", data2Bytes)
	*/
}
