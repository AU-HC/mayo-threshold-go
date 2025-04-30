package mpc

import (
	"bytes"
	"crypto/sha3"
)

func Commit(m [][]byte, r []byte) []byte {
	output := make([]byte, 32)

	h := sha3.NewSHAKE256()

	for _, row := range m {
		_, _ = h.Write(row[:])
	}

	_, _ = h.Write(r[:])

	_, _ = h.Read(output[:])

	return output
}

func VerifyCommitment(m [][]byte, r []byte, commitment []byte) bool {
	actualCommitment := Commit(m, r)
	return bytes.Equal(actualCommitment, commitment)
}
