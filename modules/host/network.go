package host

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/ftsmchl/system/my_merkleTree"
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

	SegmentSize := 2

	//get TaskID from renter
	taskid, err := r.ReadString('\n')
	taskID := strings.TrimRight(taskid, "\n")
	fmt.Println("TaskID : ", taskID)
	fmt.Println("err : ", err)

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

	//now we must sign the revision too
	privateKey := h.wallet.GetPrivateKey()
	fmt.Println("our privateKey is : ", privateKey)

	//calculating merkleRoot of our new added sector
	buf := bytes.NewBuffer(data)
	tree := my_merkleTree.New()
	var leaves int
	for buf.Len() > 0 {
		leaves++
		tree.Push(buf.Next(SegmentSize))
	}
	sectorRoot := tree.Root()
	fmt.Println("sectorRoot is : ", sectorRoot)

	h.mu.Lock()

	//update contractRoots for this taskID
	h.contractRoots[taskID].numMerkleRoots++
	fmt.Println("sectorRoots = ", h.contractRoots[taskID].numMerkleRoots)
	h.contractRoots[taskID].sectorRoots = append(h.contractRoots[taskID].sectorRoots, sectorRoot)

	//update contractRevision for this taskID
	h.fileContractRevisions[taskID].revisionNumber++
	h.fileContractRevisions[taskID].numLeaves += leaves
	h.fileContractRevisions[taskID].merkleRoot = sectorRoot

	fmt.Println("Before returning")

	h.mu.Unlock()

	merkleRootHex := hex.EncodeToString(sectorRoot)
	numLeavesNum := h.fileContractRevisions[taskID].numLeaves
	numLeaves := strconv.Itoa(numLeavesNum)
	fcRevisionNum := h.fileContractRevisions[taskID].revisionNumber
	fcRevision := strconv.Itoa(fcRevisionNum)

	fmt.Println("MerkleRootHex : ", merkleRootHex)
	fmt.Println("numLeaves : ", numLeaves)
	fmt.Println("fcRevision : ", fcRevision)

	resp1, err := http.Get("http://localhost:8001/signData?privateKey=" + privateKey + "&merkleRoot=" + merkleRootHex + "&numLeaves=" + numLeaves + "&fcRevision=" + fcRevision)
	if err != nil {
		fmt.Println("Err in http.Get : ", err)
	} else {
		text, _ := ioutil.ReadAll(resp1.Body)
		fmt.Println("text received : ", string(text))
	}

}
