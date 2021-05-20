package renter

import (
	"bytes"
	"fmt"
	//"github.com/ftsmchl/system/my_merkleTree"
	"net"
	"sync"
)

const (
	SegmentSize = 2 //bytes
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

	renter *Renter
}

func (e *Editor) close() {
	e.renter.mu.Lock()
	e.mu.Lock()
	e.clients--
	e.mu.Unlock()
	if e.clients == 0 {
		fmt.Println("I will delete the editor")
		delete(e.renter.editors, e.taskID)
	}
	e.renter.mu.Unlock()
}

func (r *Renter) Editor(tID string) (_ *Editor, err error) {
	r.mu.Lock()
	//Check if there is already an editor
	cachedEditor, haveEditor := r.editors[tID]
	if haveEditor {
		fmt.Println("The editor for this task exists")
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
		//r.mu.Unlock()
		return nil, err
	}

	editor := &Editor{
		clients:    1,
		conn:       conn,
		taskID:     tID,
		netAddress: r.storageContracts[tID].IP + ":8087",
		renter:     r,
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
	fmt.Fprintf(c, "Upload\n")

	return c, err
}

func (e *Editor) upload(data []byte) {
	//send the data to host
	fmt.Println("I am about to upload : ", data)
	fmt.Fprintf(e.conn, string(data)+"\n")

	//t := my_merkleTree.New()

	buf := bytes.NewBuffer(data)
	e.renter.mu.Lock()
	for buf.Len() > 0 {
		e.renter.roots.merkleTree.Push(buf.Next(SegmentSize))
	}
	e.renter.roots.numMerkleRoots++

	sectorRoot := e.renter.roots.merkleTree.Root()
	e.renter.roots.sectorRoots = append(e.renter.roots.sectorRoots, sectorRoot)
	e.renter.mu.Unlock()

	fmt.Println("To root tou sector p 8a ginei upload einai : ", sectorRoot)
}
