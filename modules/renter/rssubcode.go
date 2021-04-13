package renter

import (
	"bytes"
	"fmt"
)

type RSSubCode struct {
	RSCode
	staticSegmentSize uint64
	staticType        [4]byte
}

func (rs *RSSubCode) Encode(data []byte) ([][]byte, error) {
	pieces, err := rs.enc.Split(data)
	if err != nil {
		return nil, err
	}
	return rs.EncodeShards(pieces)
}

func (rs *RSSubCode) EncodeShards(pieces [][]byte) ([][]byte, error) {
	if len(pieces) < rs.MinPieces() {
		return nil, fmt.Errorf("not enough segments expectes %v but was %v",
			rs.MinPieces(), len(pieces))
	}
}
