package mpc

import (
	"fmt"
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
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

func GenerateMultiplicationTriples(n int, r1, c1, r2, c2 int) ([][][]byte, [][][]byte, [][][]byte) {
	if c1 != r2 {
		panic(fmt.Errorf("dimensions not suitable for matrix multiplication"))
	}

	a := RandMatrix(r1, c1)
	b := RandMatrix(r2, c2)
	c := MultiplyMatrices(a, b)

	aShares := make([][][]byte, n)
	bShares := make([][][]byte, n)
	cShares := make([][][]byte, n)
	for i := 0; i < n-1; i++ {
		aShares[i] = RandMatrix(r1, c1)
		bShares[i] = RandMatrix(r2, c2)
		cShares[i] = RandMatrix(r1, c2) // TODO: check

		AddMatrices(a, aShares[i])
		AddMatrices(b, bShares[i])
		AddMatrices(c, cShares[i])
	}

	aShares[n-1] = a
	bShares[n-1] = b
	cShares[n-1] = c

	// Reconstruct a, b, c
	aReconstructed := RandMatrix(r1, c1)
	bReconstructed := RandMatrix(r2, c2)
	cReconstructed := MultiplyMatrices(a, b)
	for i := 0; i < n-1; i++ {
		AddMatrices(aReconstructed, aShares[i])
		AddMatrices(bReconstructed, bShares[i])
		AddMatrices(cReconstructed, cShares[i])
	}
	if reflect.DeepEqual(c, MultiplyMatrices(aReconstructed, bReconstructed)) {
		panic(fmt.Errorf("c is not the product of a and b"))
	}

	return aShares, bShares, cShares
}

func AddMatrices(a, b [][]byte) {
	for i := range a {
		for j := range a[i] {
			a[i][j] ^= b[i][j]
		}
	}
}

func MultiplyMatricesShares(a, b [][]byte) [][]byte {
	result := make([][]byte, 0)
	return result
}

func MultiplyMatrices(A, B [][]byte) [][]byte {
	rowsA, colsA := len(A), len(A[0])
	rowsB, colsB := len(B), len(B[0])

	if colsA != rowsB {
		panic(fmt.Sprintf("Cannot multiply matrices colsA: '%d', rowsB: '%d'", colsA, rowsB))
	}

	C := make([][]byte, rowsA)
	for i := range C {
		C[i] = make([]byte, colsB)
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				C[i][j] ^= gf16Mul(A[i][k], B[k][j])
			}
		}
	}

	return C
}

func gf16Mul(a, b byte) byte {
	var r byte

	// Multiply each coefficient with y
	r = (a & 0x1) * b
	r ^= (a & 0x2) * b
	r ^= (a & 0x4) * b
	r ^= (a & 0x8) * b

	overFlowBits := r & 0xF0

	// Reduce with respect to x^4 + x + 1
	reducedOverFlowBits := overFlowBits>>4 ^ overFlowBits>>3

	// Subtract and remove overflow bits
	r = (r ^ reducedOverFlowBits) & 0x0F

	return r
}

func MultiplyMatrixWithConstantMatrix(a, b [][]byte) [][]byte {
	result := make([][]byte, 0)
	return result
}

func MatrixTranspose(a [][]byte) [][]byte {
	if len(a) == 0 {
		return [][]byte{}
	}

	rows, cols := len(a), len(a[0])
	result := make([][]byte, cols)
	for i := range result {
		result[i] = make([]byte, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			result[j][i] = a[i][j]
		}
	}

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
