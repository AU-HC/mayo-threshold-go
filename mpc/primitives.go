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

func GenerateMultiplicationTriples(n, r1, c1, r2, c2, amount int) []model.Triple {
	triples := make([]model.Triple, amount)
	for i := 0; i < amount; i++ {
		triples[i] = GenerateMultiplicationTriple(n, r1, c1, r2, c2)
	}
	return triples
}

func GenerateMultiplicationTriple(n, r1, c1, r2, c2 int) model.Triple {
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
	aReconstructed := generateZeroMatrix(r1, c1)
	bReconstructed := generateZeroMatrix(r2, c2)
	cReconstructed := generateZeroMatrix(r1, c2)
	for i := 0; i < n; i++ {
		AddMatrices(aReconstructed, aShares[i])
		AddMatrices(bReconstructed, bShares[i])
		AddMatrices(cReconstructed, cShares[i])
	}
	if !reflect.DeepEqual(cReconstructed, MultiplyMatrices(aReconstructed, bReconstructed)) {
		panic(fmt.Errorf("c is not the product of a and b"))
	}

	return model.Triple{
		A: aShares,
		B: bShares,
		C: cShares,
	}
}

// generateZeroMatrix generates a matrix of bytes with all elements set to zero
func generateZeroMatrix(rows, columns int) [][]byte {
	matrix := make([][]byte, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]byte, columns)
	}

	return matrix
}

func AddMatrices(a, b [][]byte) {
	if len(a) != len(b) && len(a[0]) != len(b[0]) {
		panic(fmt.Errorf("a and b do not have the same dimensions "))
	}

	for i := range a {
		for j := range a[i] {
			a[i][j] ^= b[i][j]
		}
	}
}

func AddMatricesNew(a, b [][]byte) [][]byte {
	if len(a) != len(b) && len(a[0]) != len(b[0]) {
		panic(fmt.Errorf("a and b do not have the same dimensions"))
	}

	c := generateZeroMatrix(len(a), len(a[0]))

	for i := range a {
		for j := range a[i] {
			c[i][j] = a[i][j] ^ b[i][j]
		}
	}

	return c
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

func MultiplyMatrixWithConstant(A [][]byte, c byte) [][]byte {
	rowsA, colsA := len(A), len(A[0])
	B := make([][]byte, rowsA)

	for i := 0; i < rowsA; i++ {
		B[i] = make([]byte, colsA)
		for j := 0; j < colsA; j++ {
			B[i][j] = gf16Mul(A[i][j], c)
		}
	}

	return B
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

func multiplicationProtocol(parties []*model.Party, triple model.Triple, dShares, eShares [][][]byte, dRow, dCol, eRow, eCol int) [][][]byte {
	zShares := make([][][]byte, len(parties))

	d := generateZeroMatrix(dRow, dCol)
	e := generateZeroMatrix(eRow, eCol)
	for j := range parties {
		AddMatrices(d, dShares[j])
		AddMatrices(e, eShares[j])
	}

	for partyNumber := range parties {
		a := triple.A[partyNumber]
		b := triple.B[partyNumber]
		c := triple.C[partyNumber]

		db := MultiplyMatrices(d, b) // d * [b]
		de := MultiplyMatrices(d, e) // d * e
		ae := MultiplyMatrices(a, e) // [a] * e
		AddMatrices(db, ae)          // d * [b] + [a] * e
		AddMatrices(db, c)           // d * [b] + [a] * e + [c]

		if partyNumber == 0 {
			AddMatrices(db, de) // d * [b] + [a] * e + [c] + d * e
		}

		zShares[partyNumber] = db
	}

	return zShares
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
