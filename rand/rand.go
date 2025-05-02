package rand

import (
	"crypto/sha3"
	"math/rand"
)

func SampleFieldElement() byte {
	randomInt := rand.Int()
	return byte(randomInt) & 0xf
}

func SampleExtensionFieldElement() uint64 {
	randomInt := rand.Int()
	return uint64(randomInt) & 0xfffffff
}

func Shake256(outputLength int, inputs ...[]byte) []byte {
	output := make([]byte, outputLength)

	h := sha3.NewSHAKE256()
	for _, input := range inputs {
		_, _ = h.Write(input[:])
	}
	_, _ = h.Read(output[:])

	for index, elem := range output {
		output[index] = elem & 0xf
	}

	return output
}

func Coin(amountOfParties, lambda int) []byte {
	result := make([]byte, lambda+64)

	for i := 0; i < lambda+64; i++ {
		for j := 0; j < amountOfParties; j++ {
			result[i] ^= SampleFieldElement()
		}
	}

	return result
}

func CoinMatrix(amountOfParties, r, c int) [][]byte {
	matrix := make([][]byte, r)
	for i := range matrix {
		matrix[i] = make([]byte, c)
		for j := 0; j < c; j++ {
			for k := 0; k < amountOfParties; k++ {
				matrix[i][j] ^= SampleFieldElement()
			}
		}
	}
	return matrix
}

func Matrix(r, c int) [][]byte {
	result := make([][]byte, r)

	for i := 0; i < r; i++ {
		row := make([]byte, c)
		for j := 0; j < c; j++ {
			row[j] = SampleFieldElement()
		}
		result[i] = row
	}

	return result
}
