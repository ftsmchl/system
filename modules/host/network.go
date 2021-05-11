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
	fmt.Println("message received from renter : ", message)
}
