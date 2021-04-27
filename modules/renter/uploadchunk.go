package renter

import (
	"fmt"
	"github.com/ftsmchl/system/modules/renter/renterfile"
	"io"
	"os"
	"sync"
)

type unfinishedUploadChunk struct {
	file *renterfile.Renterfile

	//information about the chunk, where it exists
	//within the file
	index         uint64
	length        uint64
	offset        int64
	piecesNeeded  int //number of pieces to achieve a 100% upload
	minimumPieces int //number of pieces required to recover the file

	logicalChunkData  [][]byte
	physicalChunkData [][]byte

	mu               sync.Mutex
	pieceUsage       []bool
	piecesCompleted  int
	released         bool
	workersRemaining int //number of inactive workers still able to upload a piece
}

func (r *Renter) fetchChunk(chunk *unfinishedUploadChunk) {
	//Fetch the logical data for the chunk
	err := r.fetchLogicalChunkData(chunk)
	if err != nil {
		fmt.Println("logical data could not be fetched exiting fetchChunk....")
		return
	}

	chunk.physicalChunkData, err = chunk.file.ErasureCode().EncodeShards(chunk.logicalChunkData)

	//
	if err != nil {
		fmt.Println("fetching physical data failed index of chunk is : ", chunk.index, "err : ", err)
		return
	}

	//Distribute the chunk to the workers
	r.distributeChunkToWorkers(chunk)

}

//will fetch the raw data for a chunk , pulling it from disk
func (r *Renter) fetchLogicalChunkData(chunk *unfinishedUploadChunk) error {
	fmt.Println("We are inside fetchLogicalChunkData")
	osFile, err := os.Open(chunk.file.LocalPath())
	if err != nil {
		return err
	}

	defer osFile.Close()
	sectionReader := io.NewSectionReader(osFile, chunk.offset, int64(chunk.length))
	buffer := NewBuffer(chunk.length, chunk.file.PieceSize())

	_, err = buffer.ReadFrom(sectionReader)
	if err != nil && err != io.EOF {
		fmt.Println("Err in ReadFrom : ", err)
		return err
	}
	chunk.logicalChunkData = buffer.buf
	fmt.Println("-----Chunk Logical Data-------")
	for _, data := range chunk.logicalChunkData {
		fmt.Println("logical Data : ", string(data))
	}
	fmt.Println("-----Chunk/Logical Data-------")

	//Data succesfully read from Disk
	return nil
}

func (r *Renter) distributeChunkToWorkers(uc *unfinishedUploadChunk) {
	fmt.Println("We are inside distributeChunkToWorkers")
	r.mu.Lock()
	uc.workersRemaining = len(r.workers)
	workers := make([]*worker, 0, len(r.workers))
	for _, worker := range r.workers {
		workers = append(workers, worker)
	}
	r.mu.Unlock()
	for _, worker := range workers {
		worker.queueUploadChunk(uc)
	}

}
