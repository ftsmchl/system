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
}

func (r *Renter) updateWorkerPool() {
	r.storageContractsMu.Lock()
	fmt.Println("I am inside update workerPool")

	for taskID, contract := range r.storageContracts {
		_, exists := r.workers[taskID]
		if !exists {
			worker := &worker{
				contract:   contract,
				uploadChan: make(chan struct{}, 1),
			}
			r.workers[taskID] = worker
		}
	}
	r.storageContractsMu.Unlock()
}
