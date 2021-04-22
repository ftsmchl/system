package modules

import (
	"fmt"

	"github.com/klauspost/reedsolomon"
)

const (
	SectorSize = 90 //bytes
)

//RSCode is a ReedSolomon encoder/decoder
type RSCode struct {
	enc reedsolomon.Encoder

	numPieces  int
	dataPieces int
}

func (rs *RSCode) NumPieces() int { return rs.numPieces }

func (rs *RSCode) MinPieces() int { return rs.dataPieces }

//Encode splits data into equal-length pieces, some containing the
//original data and some containing parity data.
func (rs *RSCode) Encode(data []byte) ([][]byte, error) {
	pieces, err := rs.enc.Split(data)
	if err != nil {
		return nil, err
	}
	return rs.EncodeShards(pieces)
}

func (rs *RSCode) EncodeShards(pieces [][]byte) ([][]byte, error) {
	//Check that the caller provided the minimum amount of pieces
	if len(pieces) < rs.MinPieces() {
		return nil, fmt.Errorf("invalid number of pieces given %v < %v ", len(pieces), rs.MinPieces())
	}

	pieceSize := len(pieces[0])

	for len(pieces) < rs.MinPieces() {
		pieces = append(pieces, make([]byte, pieceSize))
	}
	err := rs.enc.Encode(pieces)
	if err != nil {
		return nil, err
	}
	return pieces, nil

}

func NewRSCode(nData, nParity int) (*RSCode, error) {
	return newRSCode(nData, nParity)
}

func newRSCode(nData, nParity int) (*RSCode, error) {
	enc, err := reedsolomon.New(nData, nParity)
	if err != nil {
		return nil, err
	}
	return &RSCode{
		enc:        enc,
		numPieces:  nData + nParity,
		dataPieces: nData,
	}, nil
}
