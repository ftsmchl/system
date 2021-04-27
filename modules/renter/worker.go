package renter

import (
	"fmt"
	"sync"
)

//a worker is responsible for the communication with a specific
//host.
type worker struct {
	mu sync.Mutex

	contract StorageContract

	//Upload variables
	uploadChan chan struct{}

	unprocessedChunks []*unfinishedUploadChunk
}

func (r *Renter) updateWorkerPool() {
	r.storageContractsMu.Lock()
	fmt.Println("I am inside update workerPool")
	counter := 0

	for taskID, contract := range r.storageContracts {
		r.mu.Lock()
		_, exists := r.workers[taskID]
		if !exists {
			counter++
			fmt.Println("I am creating worker num. : ", counter)
			worker := &worker{
				contract:   contract,
				uploadChan: make(chan struct{}, 1),
			}
			r.workers[taskID] = worker

			go func() {
				worker.threadedWorkLoop()
			}()
		}
		r.mu.Unlock()
	}
	r.storageContractsMu.Unlock()
	fmt.Println("Exiting updateWorkerPool()")
	//we need to remove any worker that is connected with a contract that is no longer
	//live
}

func (w *worker) threadedWorkLoop() {
	for {
		select {
		case <-w.uploadChan:
			fmt.Println("I received a signal in my uploadChan")
		}
	}
}
