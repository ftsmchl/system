package renter

import (
	"bytes"
	"fmt"
	"net"
	"time"
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
	//w.mu.Lock()
	taskID := w.contract.TaskID
	//hostPublicKey := w.renter.storageContracts[taskID].Host

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

		//w.renter.editors[taskID] = e
		w.mu.Unlock()
	*/

	conn, err := net.DialTimeout("tcp", w.renter.storageContracts[taskID].IP+":8087", 45*time.Second)
	fmt.Println("Ip is ", w.renter.storageContracts[taskID].IP)
	if err != nil {
		fmt.Println("err sto conn : ", err)
		return
	}

	//	fmt.Fprintf(conn, "Upload\n")
	//fmt.Fprintf(conn, string(uc.physicalChunkData[pieceIndex])+"\n")

	//	uc.mu.Lock()
	data := uc.physicalChunkData[pieceIndex]
	//	uc.mu.Unlock()

	n, err := conn.Write(uc.physicalChunkData[pieceIndex])
	if err != nil {
		fmt.Println("Mistakes were made in conn.Write.. ", err)
	}
	fmt.Println("n of bytes sent is : ", n)
	fmt.Println("we sent : ", uc.physicalChunkData[pieceIndex])

	//w.mu.Unlock()

	//	w.renter.inUseMu.Lock()
	_ = conn.Close()
	//	delete(w.renter.host_in_use, hostPublicKey)
	//	w.renter.inUseMu.Unlock()

	//data := uc.physicalChunkData[pieceIndex]
	buf := bytes.NewBuffer(data)
	w.renter.mu.Lock()
	fmt.Println(" Mesa sto lock tou roots")
	for buf.Len() > 0 {
		w.renter.roots.merkleTree.Push(buf.Next(SegmentSize))
	}
	w.renter.roots.numMerkleRoots++
	sectorRoot := w.renter.roots.merkleTree.Root()
	w.renter.roots.sectorRoots = append(w.renter.roots.sectorRoots, sectorRoot)
	w.renter.mu.Unlock()
	fmt.Println("Outside the lock")
	fmt.Println("To root tou sector p egine upload einai : ", sectorRoot)

	return

}
