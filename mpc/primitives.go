package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
)

func AddListOfMatrices(parties []*model.Party, id1, id2, id3 string) {
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

func MultiplyMatrices(a, b [][]byte) [][]byte {
	result := make([][]byte, 0)
	return result
}

func MultiplyMatrixWithConstantMatrix(a, b [][]byte) [][]byte {
	result := make([][]byte, 0)
	return result
}

func MatrixTranspose(a [][]byte) [][]byte {
	result := make([][]byte, 0)
	return result
}

func Coin(parties []*model.Party, lambda int) []byte {
	result := make([]byte, lambda+64)

	for i := 0; i < lambda+64; i++ {
		for _, _ = range parties {
			result[i] ^= rand.SampleFieldElement()
		}
	}

	return result
}

func RandMatrix(r, c int) [][]byte {
	result := make([][]byte, r)

	for i := 0; i < r; i++ {
		row := make([]byte, c)
		for j := 0; j < c; j++ {
			row[j] = rand.SampleFieldElement()
		}
		result[i] = row
	}

	return result
}
