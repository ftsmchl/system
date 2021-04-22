package renterfile

import (
	"github.com/ftsmchl/system/modules"
	"os"
	"sync"
)

//Renterfile, contains all the necessary information to
//recover a file from its hosts.
type Renterfile struct {
	//metadata Metadata
	localPath   string
	fileSize    int64
	pieceSize   uint64
	erasureCode *modules.RSSubCode

	//hosts that this file's pieces are uploaded to
	publicKeyTable []HostPublicKey

	//chunks are the chunks the file was split into
	chunks []chunk

	//path of our file in our local
	//filePath string

	mu sync.Mutex
}

type chunk struct {

	//Pieces are the Pieces of the file the chunk consists of
	Pieces [][]piece
}

type piece struct {
	//offset of the host's key within the publicKeyTable
	HostOffset uint32

	//Hash is a Blake 256 bit digest
	MerkleRoot [32]byte
}

type HostPublicKey struct {
	PublicKey string
	Used      bool
}

func (rf *Renterfile) ErasureCode() *modules.RSSubCode {
	return rf.erasureCode
}

func (rf *Renterfile) ChunkSize() uint64 {
	return rf.pieceSize * uint64(rf.erasureCode.MinPieces())
}

func (rf *Renterfile) NumChunks() uint64 {
	rf.mu.Lock()
	numChunks := len(rf.chunks)
	rf.mu.Unlock()
	return uint64(numChunks)
}

//New creates a new Renterfile
func New(source string, erasureCode *modules.RSSubCode, fileSize uint64, fileMode os.FileMode) *Renterfile {
	file := &Renterfile{
		localPath:   source,
		fileSize:    int64(fileSize),
		pieceSize:   modules.SectorSize,
		erasureCode: erasureCode,
	}

	//Init chunks
	numChunks := fileSize / file.ChunkSize()
	if fileSize%file.ChunkSize() != 0 || numChunks == 0 {
		numChunks++
	}

	file.chunks = make([]chunk, erasureCode.NumPieces())

	for i := range file.chunks {
		file.chunks[i].Pieces = make([][]piece, erasureCode.NumPieces())
	}

	return file
}
