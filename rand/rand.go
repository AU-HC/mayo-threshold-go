package rand

import "C"
import (
	"crypto/rand"
	"crypto/sha3"
)

// SampleFieldElement TODO: sample random value
func SampleFieldElement() byte {
	buf := make([]byte, 1)

	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf[0] & 0xf
}

func Shake256(outputLength int, inputs ...[]byte) []byte {
	output := make([]byte, outputLength)

	h := sha3.NewSHAKE256()
	for _, input := range inputs {
		_, _ = h.Write(input[:])
	}
	_, _ = h.Read(output[:])

	return output
}
