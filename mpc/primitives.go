package mpc

import (
	"fmt"
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
)

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

	a := rand.RandMatrix(r1, c1)
	b := rand.RandMatrix(r2, c2)
	c := MultiplyMatrices(a, b)

	aShares := make([][][]byte, n)
	bShares := make([][][]byte, n)
	cShares := make([][][]byte, n)
	for i := 0; i < n-1; i++ {
		aShares[i] = rand.RandMatrix(r1, c1)
		bShares[i] = rand.RandMatrix(r2, c2)
		cShares[i] = rand.RandMatrix(r1, c2)

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
