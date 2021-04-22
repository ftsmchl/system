package renter

import (
	"github.com/ftsmchl/system/modules/renter/renterfile"
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

	mu              sync.Mutex
	pieceUsage      []bool
	piecesCompleted int
	released        bool
}
