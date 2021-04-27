package renter

import (
	"container/heap"
	"fmt"
	"github.com/ftsmchl/system/modules/renter/renterfile"
	"sync"
)

// is a bunch of chunks that need to be uploaded.
type uploadChunkHeap []*unfinishedUploadChunk

func (uch uploadChunkHeap) Len() int { return len(uch) }
func (uch uploadChunkHeap) Less(i, j int) bool {
	return true
}
func (uch uploadChunkHeap) Swap(i, j int)       { uch[i], uch[j] = uch[j], uch[i] }
func (uch *uploadChunkHeap) Push(x interface{}) { *uch = append(*uch, x.(*unfinishedUploadChunk)) }
func (uch *uploadChunkHeap) Pop() interface{} {
	old := *uch
	n := len(old)
	x := old[n-1]
	*uch = old[0 : n-1]
	return x
}

type uploadHeap struct {
	heap uploadChunkHeap

	newUploads chan struct{}
	mu         sync.Mutex
}

func (uh *uploadHeap) managedPop() (uc *unfinishedUploadChunk) {
	uh.mu.Lock()
	if len(uh.heap) > 0 {
		uc = heap.Pop(&uh.heap).(*unfinishedUploadChunk)
	}
	uh.mu.Unlock()
	return uc
}

func (uh *uploadHeap) managedPush(uuc *unfinishedUploadChunk) {
	uh.mu.Lock()
	heap.Push(&uh.heap, uuc)
	uh.mu.Unlock()
}

func (uh *uploadHeap) managedLen() int {
	uh.mu.Lock()
	len := uh.heap.Len()
	uh.mu.Unlock()
	return len
}

func (r *Renter) buildAndPushChunks(file *renterfile.Renterfile) {
	r.mu.Lock()
	unfinishedChunks := r.buildUnfinishedChunks(file)
	r.mu.Unlock()

	fmt.Println("Number of unfinished chunks created : ", len(unfinishedChunks))

	//push each unfinished chunk to our upload heap
	for i := 0; i < len(unfinishedChunks); i++ {
		r.uploadHeap.managedPush(unfinishedChunks[i])
	}
}

func (r *Renter) buildUnfinishedChunks(file *renterfile.Renterfile) []*unfinishedUploadChunk {
	minPieces := file.ErasureCode().MinPieces()
	fmt.Println("Eimai sth buildUnfinishedChunks kai to minPieces einai : ", minPieces)

	if len(r.workers) < minPieces {
		return nil
	}

	var chunkIndexes []uint64

	for i := uint64(0); i < file.NumChunks(); i++ {
		chunkIndexes = append(chunkIndexes, i)

	}

	newUnfinishedChunks := make([]*unfinishedUploadChunk, len(chunkIndexes))

	for i, index := range chunkIndexes {
		newUnfinishedChunks[i] = r.buildUnfinishedChunk(file, uint64(index))
	}

	return newUnfinishedChunks

}

func (r *Renter) buildUnfinishedChunk(file *renterfile.Renterfile, chunkIndex uint64) *unfinishedUploadChunk {
	uuc := &unfinishedUploadChunk{
		file:              file,
		index:             chunkIndex,
		length:            file.ChunkSize(),
		offset:            int64(file.ChunkSize() * chunkIndex),
		minimumPieces:     file.ErasureCode().MinPieces(),
		piecesNeeded:      file.ErasureCode().NumPieces(),
		physicalChunkData: make([][]byte, file.ErasureCode().NumPieces()),
		pieceUsage:        make([]bool, file.ErasureCode().NumPieces()),
	}

	return uuc
}

//is a background process that maintains a queue of chunks to upload
func (r *Renter) threadedUpload() {
	for {
		//wait for an upload to be signaled
		select {
		case <-r.uploadHeap.newUploads:
			fmt.Println("Signal caught that we have chunks for uploading")
			//	err := r.repairLoop()
		}

		err := r.repairLoop()
		if err != nil {
			fmt.Println("Something went wrong in repairLoop : ", err)
		}

	}
}

//repairLoop works through the upload heap trying to process the
// unfinished chunks
func (r *Renter) repairLoop() error {
	counter := 0
	for r.uploadHeap.managedLen() > 0 {
		counter++
		nextChunk := r.uploadHeap.managedPop()
		fmt.Println("Chunk popped from upload heap num : ", counter)
		if nextChunk == nil {
			fmt.Println("uploadHeap is empty")
			return nil
		}
		r.mu.Lock()
		availableWorkers := len(r.workers)
		r.mu.Unlock()

		if availableWorkers < nextChunk.minimumPieces {
			fmt.Println("Not enough hosts to upload the chunk properly, reaching minimum redundancy")
			continue
		}

		err := r.prepareNextChunk(nextChunk)
		if err != nil {
			fmt.Println("Something went wrong in prepareNextChunk() : ", err)
		}

	}

	return nil

}

func (r *Renter) prepareNextChunk(uuc *unfinishedUploadChunk) error {
	go r.fetchChunk(uuc)
	return nil
}
