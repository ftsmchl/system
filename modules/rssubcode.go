package modules

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

	pieceSize := uint64(len(pieces[0]))

	//pieceSize must be divisible by SegmentSize
	if pieceSize%rs.staticSegmentSize != 0 {
		return nil, fmt.Errorf("piece size not divisible by segmentSize")
	}

	//Each piece should have pieceSize bytes
	for _, piece := range pieces {
		if uint64(len(piece)) != pieceSize {
			return nil, fmt.Errorf("pieces don't have right size expected %v, but was %v ", pieceSize, len(piece))
		}
	}

	//Flatten the pieces into a byte slice
	data := make([]byte, uint64(len(pieces))*pieceSize)
	for i, piece := range pieces {
		copy(data[uint64(i)*pieceSize:], piece)
		pieces[i] = pieces[i][:0]
	}

	//Add parity shards to pieces
	parityShards := make([][]byte, rs.NumPieces()-len(pieces))
	pieces = append(pieces, parityShards...)

	//Encode the pieces
	segmentOffset := uint64(0)
	for buf := bytes.NewBuffer(data); buf.Len() > 0; {
		s := buf.Next(int(rs.staticSegmentSize) * rs.MinPieces())
		segments := make([]byte, len(s))
		copy(segments, s)

		//Encode the segments
		encodedSegments, err := rs.RSCode.Encode(segments)
		if err != nil {
			return nil, err
		}

		//Write the encoded segments back to pieces
		for i, segment := range encodedSegments {
			pieces[i] = append(pieces[i], segment...)
		}

		segmentOffset += rs.staticSegmentSize
	}

	return pieces, nil
}

//NewRSSubCode creates a new Reed-Solomon encoder-decoder using the supplied
//parameters.

func NewRSSubCode(nData, nParity int) (*RSSubCode, error) {
	rs, err := newRSCode(nData, nParity)
	if err != nil {
		return nil, err
	}

	t := [4]byte{0, 0, 0, 2}

	return &RSSubCode{
		*rs,
		64,
		t,
	}, nil

}
