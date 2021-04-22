package renter

import (
	"fmt"
	"github.com/ftsmchl/system/modules/renter/renterfile"
)

func (r *Renter) buildAndPushChunks(file *renterfile.Renterfile) {
	r.mu.Lock()
	unfinishedChunks := r.buildUnfinishedChunks(file)
	r.mu.Unlock()

	fmt.Println("Number of unfinished chunks created : ", len(unfinishedChunks))
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
