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

	a := rand.Matrix(r1, c1)
	b := rand.Matrix(r2, c2)
	c := MultiplyMatrices(a, b)

	aShares := algo.createSharesForMatrix(a)
	bShares := algo.createSharesForMatrix(b)
	cShares := algo.createSharesForMatrix(c)

	// Reconstruct a, b, c
	aReconstructed := algo.openMatrix(aShares)
	bReconstructed := algo.openMatrix(bShares)
	cReconstructed := algo.openMatrix(cShares)
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

	d := algo.openMatrix(dShares)
	e := algo.openMatrix(eShares)

	for partyNumber := range parties {
		a := triple.A[partyNumber]
		b := triple.B[partyNumber]
		c := triple.C[partyNumber]

		db := MultiplyMatrices(d, b) // d * [b]
		de := MultiplyMatrices(d, e) // d * e
		ae := MultiplyMatrices(a, e) // [a] * e
		AddMatrices(db, ae)          // d * [b] + [a] * e
		AddMatrices(db, c)           // d * [b] + [a] * e + [c]

		if partyNumber == partyNumber { // TODO: variable point
			AddMatrices(db, de) // d * [b] + [a] * e + [c] + d * e
		}

		zShares[partyNumber] = db
	}

	return zShares
}
