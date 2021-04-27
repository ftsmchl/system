package renter

import (
	"fmt"
)

func (w *worker) queueUploadChunk(uc *unfinishedUploadChunk) {
	fmt.Println(" i am in queue upload chunk ...")
	w.mu.Lock()
	w.unprocessedChunks = append(w.unprocessedChunks, uc)
	w.mu.Unlock()

	//Send a signal informing the work thread that there is a new chunk
	w.uploadChan <- struct{}{}
}
