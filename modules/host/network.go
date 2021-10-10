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

type Challenge struct {
	address string `json : "address"`
}

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

		//I heard a challenge ....

		fmt.Println("challenge")
		msg2, _ := reader.ReadString('\n')
		address := strings.TrimRight(msg2, "\n")
		fmt.Println("address = ", address)

		//create a merkleproof

		fmt.Println("taskID = ", h.taskID)

		segmentSize := 2
		segmentIndex := 0
		sectorSize := 4

		sectorIndex := segmentIndex / (sectorSize / segmentSize)
		sectorSegment := segmentIndex % (sectorSize / segmentSize)

		taskID := h.taskID
		sectorBytes := h.sectors[taskID][sectorIndex]

		//Build the storage proof for jist the sector
		t := my_merkleTree.New()
		t.SetIndex(uint64(sectorSegment))

		buf := bytes.NewBuffer(sectorBytes)
		for buf.Len() > 0 {
			t.Push(buf.Next(segmentSize))
		}

		_, proof, _, _ := t.Prove()

		base := proof[0]
		hashSet := make([][32]byte, len(proof)-1)
		for i, p := range proof[1:] {
			copy(hashSet[i][:], p)
		}

		log2SectorSize := uint64(0)
		for 1<<log2SectorSize < (sectorSize / segmentSize) {
			log2SectorSize++
		}

		ct := my_merkleTree.NewCachedTree(log2SectorSize)
		ct.SetIndex(uint64(segmentIndex))

		for _, root := range h.contractRoots[taskID].sectorRoots {
			ct.Push(root)
		}

		cachedProofSet := make([][]byte, len(hashSet)+1)
		cachedProofSet[0] = base

		for i := range hashSet {
			cachedProofSet[i+1] = hashSet[i][:]
		}

		merkleRoot, proofSet, proofIndex, numLeaves := ct.Prove(cachedProofSet)
		nextApotelesma := my_merkleTree.VerifyProof(merkleRoot, proofSet, proofIndex, numLeaves)
		fmt.Println("Proof : ", nextApotelesma)
	}

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

	sector := data

	//read renter's signature
	renterSignature, err := r.ReadString('\n')
	fmt.Println("err : ", err)
	fmt.Println("Renter's Signature is : ", renterSignature)
	renterSignature = strings.TrimRight(renterSignature, "\n")

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
	h.sectors[taskID] = append(h.sectors[taskID], sector)

	//compute MerkleRoot until now
	log2SectorSize := uint64(0)
	for 1<<log2SectorSize < (4 / 2) {
		log2SectorSize++
	}
	ct := my_merkleTree.NewCachedTree(log2SectorSize)

	for _, root := range h.contractRoots[taskID].sectorRoots {
		ct.Push(root)
	}

	merkleRootUntilNow := ct.Root()

	//update contractRevision for this taskID
	h.fileContractRevisions[taskID].revisionNumber++
	h.fileContractRevisions[taskID].numLeaves += leaves
	h.fileContractRevisions[taskID].signatureRenter = renterSignature
	h.fileContractRevisions[taskID].merkleRoot = merkleRootUntilNow

	fmt.Println("Before returning")

	h.mu.Unlock()

	//merkleRootHex := hex.EncodeToString(sectorRoot)
	merkleRootHex := hex.EncodeToString(merkleRootUntilNow)
	numLeavesNum := h.fileContractRevisions[taskID].numLeaves
	numLeaves := strconv.Itoa(numLeavesNum)
	fcRevisionNum := h.fileContractRevisions[taskID].revisionNumber
	fcRevision := strconv.Itoa(fcRevisionNum)

	fmt.Println("MerkleRootHex : ", merkleRootHex)
	fmt.Println("numLeaves : ", numLeaves)
	fmt.Println("fcRevision : ", fcRevision)

	var hostSignature string

	resp1, err := http.Get("http://localhost:8001/signData?privateKey=" + privateKey + "&merkleRoot=" + merkleRootHex + "&numLeaves=" + numLeaves + "&fcRevision=" + fcRevision)
	if err != nil {
		fmt.Println("Err in http.Get : ", err)
	} else {
		hostSignatureBytes, _ := ioutil.ReadAll(resp1.Body)
		hostSignature = string(hostSignatureBytes)
		fmt.Println("Our Signature : ", string(hostSignature))
	}

	storageContractAddress := h.storageContracts[taskID].Address

	//check if renter's signature is OK
	resp2, err := http.Get("http://localhost:8001/checkSignatures?sigRenter=" + renterSignature + "&sigHost=" + hostSignature + "&merkleRoot=" + merkleRootHex + "&numLeaves=" + numLeaves + "&fcRevision=" + fcRevision + "&address=" + storageContractAddress)

	if err != nil {
		fmt.Println("Err from response in checkSignatures : ", err)
	} else {
		text2, _ := ioutil.ReadAll(resp2.Body)
		fmt.Println("Received : ", string(text2))
	}

	//send our Signature to renter
	fmt.Fprintf(c, hostSignature+"\n")

	h.mu.Lock()
	h.fileContractRevisions[taskID].signatureHost = hostSignature
	h.mu.Unlock()
}
