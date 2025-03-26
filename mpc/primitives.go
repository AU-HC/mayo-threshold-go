package mpc

import (
	"mayo-threshold-go/mock"
	"mayo-threshold-go/rand"
)

func AddListOfMatrices(parties []*mock.Party, id1, id2, id3 string) {
	for _, party := range parties {
		listOfMatrices1 := party.Shares[id1]
		listOfMatrices2 := party.Shares[id2]
		listOfMatrices3 := make([][][]byte, len(listOfMatrices1))

		for i := range listOfMatrices1 {
			listOfMatrices3[i] = make([][]byte, len(listOfMatrices1[i]))
			for j := range listOfMatrices1[i] {
				listOfMatrices3[i][j] = make([]byte, len(listOfMatrices1[i][j]))
				for k := range listOfMatrices1[i][j] {
					listOfMatrices3[i][j][k] = listOfMatrices1[i][j][k] ^ listOfMatrices2[i][j][k]
				}
			}
		}

		party.Shares[id3] = listOfMatrices3
	}
}

func Coin(parties []*mock.Party, lambda int) []byte {
	result := make([]byte, lambda+64)

	for i := 0; i < lambda+64; i++ {
		for _, _ = range parties {
			result[i] ^= rand.SampleFieldElement()
		}
	}

	return result
}
