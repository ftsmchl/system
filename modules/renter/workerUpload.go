package renter

import (
	"fmt"
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

	//open an editing connection to the host
	e, err := w.renter.Editor(taskID)

	defer e.close()
	if err != nil {
		fmt.Println("Something went wrong from calling Editor() : ", err)
		w.mu.Unlock()
		return
	}

	//upload pieceIndex to host
	e.upload(uc.physicalChunkData[pieceIndex])

	w.renter.editors[taskID] = e
	//w.mu.Unlock()

	return

}
