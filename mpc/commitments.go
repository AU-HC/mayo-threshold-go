package mpc

import (
	"bytes"
	"crypto/sha3"
)

func Commit(m, r byte) []byte {
	output := make([]byte, 32)

	h := sha3.NewSHAKE256()
	_, _ = h.Write([]byte{m})
	_, _ = h.Write([]byte{r})

	_, _ = h.Read(output[:])

	return output
}

func VerifyCommitment(m, r byte, commitment []byte) bool {
	actualCommitment := Commit(m, r)
	return bytes.Equal(actualCommitment, commitment)
}
