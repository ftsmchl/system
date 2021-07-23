package renter

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
	"time"

	"github.com/ftsmchl/system/my_merkleTree"
)

func (w *worker) queueUploadChunk(uc *unfinishedUploadChunk) {
	fmt.Println(" i am in queueuploadchunk ...")
	w.mu.Lock()
	w.unprocessedChunks = append(w.unprocessedChunks, uc)
	w.mu.Unlock()

	//Send a signal informing the work thread that there is a new chunk
	w.uploadChan <- struct{}{}
}

func (w *worker) nextUploadChunk() (nextChunk *unfinishedUploadChunk, pieceIndex uint64) {

	for {
		w.mu.Lock()
		if len(w.unprocessedChunks) == 0 {
			w.mu.Unlock()
			break
		}
		chunk := w.unprocessedChunks[0]
		w.unprocessedChunks = w.unprocessedChunks[1:]
		w.mu.Unlock()

		nextChunk, pieceIndex = w.processUploadChunk(chunk)
		if nextChunk != nil {
			return nextChunk, pieceIndex
		}
	}

	return nil, 0
}

func (w *worker) processUploadChunk(uc *unfinishedUploadChunk) (nextChunk *unfinishedUploadChunk, pieceIndex uint64) {
	uc.mu.Lock()
	index := -1
	for i := 0; i < len(uc.pieceUsage); i++ {
		if !uc.pieceUsage[i] {
			index = i
			uc.pieceUsage[i] = true
			break
		}
	}

	if index == -1 {
		fmt.Println("worker could not find an unused piece...")
		return nil, 0
	}

	uc.mu.Unlock()
	//fmt.Println("I chose to upload pieceIndex = ", index, " uc.index = ", uc.index)
	return uc, uint64(index)
}

//will perform some upload work
func (w *worker) upload(uc *unfinishedUploadChunk, pieceIndex uint64) {
	fmt.Println("I am currently a worker in upload")
	taskID := w.contract.TaskID
	/*
		_, exist := w.renter.host_in_use[hostPublicKey]

		if exist {
			w.renter.inUseMu.Lock()
			w.renter.host_in_use[hostPublicKey] = true
			w.renter.inUseMu.Unlock()
		} else {
			for {
				time.Sleep(3 * time.Second)
				if _, exists := w.renter.host_in_use[hostPublicKey]; !exists {
					break
				}
			}
		}
	*/

	/*
		//open an editing connection to the host
		e, err := w.renter.Editor(taskID)

		defer e.close()

		if err != nil {
			fmt.Println("Something went wrong from calling Editor() : ", err)
			w.mu.Unlock()
			return
		}

		fmt.Println("Right before e.Upload()")
		//upload pieceIndex to host
		data := uc.physicalChunkData[pieceIndex]
		e.upload(data)
		fmt.Println("Right after e.Upload()")

	*/

	conn, err := net.DialTimeout("tcp", w.renter.storageContracts[taskID].IP+":8087", 100*time.Second)
	//r := bufio.NewReader(conn)

	fmt.Println("Ip is ", w.renter.storageContracts[taskID].IP)
	if err != nil {
		fmt.Println("err sto conn : ", err)
		return
	}

	//Sending that we will use upload
	fmt.Fprintf(conn, "Upload\n")

	//Sending taskID
	fmt.Fprintf(conn, taskID+"\n")

	//Sending Data sector
	data := uc.physicalChunkData[pieceIndex]
	n, err := conn.Write(uc.physicalChunkData[pieceIndex])
	if err != nil {
		fmt.Println("Mistakes were made in conn.Write.. ", err)
	}
	fmt.Println("n of bytes sent is : ", n)
	fmt.Println("we sent : ", uc.physicalChunkData[pieceIndex])

	//Read Reply from host that he received Data
	r := bufio.NewReader(conn)
	reply, err2 := r.ReadString('\n')
	fmt.Println("We got reply : ", strings.TrimRight(reply, "\n"))
	fmt.Println("err : ", err2)

	defer conn.Close()

	buf := bytes.NewBuffer(data)
	tree := my_merkleTree.New()

	//calculating the merkle root of our adding sector
	w.renter.mu.Lock()
	var leaves int
	fmt.Println(" Mesa sto lock tou roots")
	for buf.Len() > 0 {
		leaves++
		//	w.renter.contractRoots[taskID].merkleTree.Push(buf.Next(SegmentSize))
		tree.Push(buf.Next(SegmentSize))
	}

	//Add sectorRoot to our sector Roots
	w.renter.contractRoots[taskID].numMerkleRoots++
	sectorRoot := tree.Root()
	w.renter.contractRoots[taskID].sectorRoots = append(w.renter.contractRoots[taskID].sectorRoots, sectorRoot)

	//compute MerkleRoot until now
	log2SectorSize := uint64(0)
	for 1<<log2SectorSize < (4 / 2) {
		log2SectorSize++
	}
	ct := my_merkleTree.NewCachedTree(log2SectorSize)

	for _, root := range w.renter.contractRoots[taskID].sectorRoots {
		ct.Push(root)
	}

	merkleRootUntilNow := ct.Root()

	//Update our fileContract Revision
	w.renter.fileContractRevisions[taskID].numLeaves += leaves
	w.renter.fileContractRevisions[taskID].merkleRoot = merkleRootUntilNow
	w.renter.fileContractRevisions[taskID].revisionNumber++
	privateKey := w.renter.wallet.GetPrivateKey() //Get our PrivateKey
	w.renter.mu.Unlock()
	fmt.Println("Outside the lock")
	fmt.Println("Sector's MerkleRoot : ", sectorRoot, "taskID : ", taskID)
	numLeavesNum := w.renter.fileContractRevisions[taskID].numLeaves
	fcRevisionNum := w.renter.fileContractRevisions[taskID].revisionNumber
	fmt.Println("taskID : ", taskID, " numLeaves : ", w.renter.fileContractRevisions[taskID].numLeaves)

	//Sign our Revision
	merkleRootHex := hex.EncodeToString(merkleRootUntilNow)
	numLeaves := strconv.Itoa(numLeavesNum)
	fcRevision := strconv.Itoa(fcRevisionNum)
	fmt.Println("merkeRootHex is ", merkleRootHex)
	fmt.Println("numLeaves is  ", numLeaves)
	fmt.Println("fcRevision is  ", fcRevision)
	fmt.Println("privateKey is ", privateKey)

	var ourSignature string

	resp1, err := http.Get("http://localhost:8000/signData?privateKey=" + privateKey + "&merkleRoot=" + merkleRootHex + "&numLeaves=" + numLeaves + "&fcRevision=" + fcRevision)
	if err != nil {
		fmt.Println("Error in sign.Get() : ", err)
	} else {
		text, _ := ioutil.ReadAll(resp1.Body)
		fmt.Println("Our Signature : ", string(text))
		ourSignature = string(text)
		//send our signature to host
		fmt.Fprintf(conn, ourSignature+"\n")

	}

	signatureHost, _ := r.ReadString('\n')
	signatureHost = strings.TrimRight(signatureHost, "\n")
	fmt.Println("Signature's Host : ", signatureHost)

	w.renter.mu.Lock()
	w.renter.fileContractRevisions[taskID].signatureHost = signatureHost
	w.renter.fileContractRevisions[taskID].signatureRenter = ourSignature
	w.renter.mu.Unlock()

	return

}
