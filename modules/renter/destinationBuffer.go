package renter

import (
	"io"
)

type destinationBuffer struct {
	buf       [][]byte
	pieceSize uint64
}

func NewBuffer(length, pieceSize uint64) destinationBuffer {
	if length%pieceSize != 0 {
		length += pieceSize - length%pieceSize
	}

	db := destinationBuffer{
		buf:       make([][]byte, length/pieceSize),
		pieceSize: pieceSize,
	}

	for length > 0 {
		db.buf = append(db.buf, make([]byte, pieceSize))
		length -= pieceSize
	}

	return db
}

func (d destinationBuffer) ReadFrom(r io.Reader) (int64, error) {
	var n int64
	for len(d.buf) > 0 {
		read, err := io.ReadFull(r, d.buf[0])
		if err != nil {
			return n, err
		}
		d.buf = d.buf[1:]
		n += int64(read)
	}

	return n, nil
}
