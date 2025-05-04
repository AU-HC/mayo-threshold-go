package mpc

import (
	"bytes"
	"crypto/sha3"
	"encoding/binary"
)

func Commit(m [][]uint64, r []uint64) []byte {
	output := make([]byte, 32)

	mByteSlice := make([]byte, 0, len(m)*len(m[0])*8)
	for _, row := range m {
		for _, val := range row {
			buf := make([]byte, 8)
			binary.BigEndian.PutUint64(buf, val)
			mByteSlice = append(mByteSlice, buf...)
		}
	}

	rByteSlice := make([]byte, 0, len(r)*8)
	for _, element := range r {

		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, element)
		rByteSlice = append(rByteSlice, buf...)
	}

	h := sha3.NewSHAKE256()

	_, _ = h.Write(mByteSlice)
	_, _ = h.Write(rByteSlice)

	_, _ = h.Read(output[:])

	return output
}

func VerifyCommitment(m [][]uint64, r []uint64, commitment []byte) bool {
	actualCommitment := Commit(m, r)
	return bytes.Equal(actualCommitment, commitment)
}
