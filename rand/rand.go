package rand

import "C"
import (
	"crypto/rand"
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
