package renter

import (
	"fmt"
	"net"
	"sync"
)

type Editor struct {
	//endHeight int
	taskID string
	//ip:port
	netAddress string

	conn net.Conn

	invalid bool

	//how many clients are using the editor, safe to close when 0
	clients int

	mu sync.Mutex

	contractID string
}

func (r *Renter) Editor(tID string) (_ *Editor, err error) {
	r.mu.Lock()
	//Check if there is already an editor
	cachedEditor, haveEditor := r.editors[tID]
	if haveEditor {
		cachedEditor.mu.Lock()
		cachedEditor.clients++
		cachedEditor.mu.Unlock()
		r.mu.Unlock()
		return cachedEditor, nil
	}

	//Create the Editor
	//we initiate the contract revision with the host

	conn, err := initiateRevisionLoop(r.storageContracts[tID])
	if err != nil {
		r.mu.Unlock()
		return nil, err
	}

	editor := &Editor{
		clients:    1,
		conn:       conn,
		taskID:     tID,
		netAddress: r.storageContracts[tID].IP + ":8087",
	}

	r.editors[tID] = editor

	r.mu.Unlock()

	return editor, nil

}

func initiateRevisionLoop(contract StorageContract) (net.Conn, error) {
	c, err := net.Dial("tcp", contract.IP+":8087")
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(c, "Kalhspera\n")

	return c, err
}
