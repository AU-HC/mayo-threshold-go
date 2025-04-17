package mpc

import (
	"fmt"
	"mayo-threshold-go/rand"
	"reflect"
)

func (c *Context) GenerateMultiplicationActiveTriples(r1, c1, r2, c2, amount int) []ActiveTriple {
	triples := make([]ActiveTriple, amount)
	for i := 0; i < amount; i++ {
		triples[i] = c.GenerateMultiplicationActiveTriple(r1, c1, r2, c2)
	}
	return triples
}

func (c *Context) GenerateMultiplicationActiveTriple(r1, c1, r2, c2 int) ActiveTriple {
	if c1 != r2 {
		panic(fmt.Errorf("dimensions not suitable for matrix multiplication"))
	}

	aMatrix := rand.Matrix(r1, c1)
	bMatrix := rand.Matrix(r2, c2)
	cMatrix := MultiplyMatrices(aMatrix, bMatrix)

	aShares := createSharesForMatrix(c.n, aMatrix)
	bShares := createSharesForMatrix(c.n, bMatrix)
	cShares := createSharesForMatrix(c.n, cMatrix)

	// Reconstruct aMatrix, bMatrix, c
	aReconstructed, _ := c.algo.authenticatedOpenMatrix(aShares)
	bReconstructed, _ := c.algo.authenticatedOpenMatrix(bShares)
	cReconstructed, _ := c.algo.authenticatedOpenMatrix(cShares)
	if !reflect.DeepEqual(cReconstructed, MultiplyMatrices(aReconstructed, bReconstructed)) {
		panic(fmt.Errorf("c is not the product of aMatrix and bMatrix"))
	}

	return ActiveTriple{
		A: aShares,
		B: bShares,
		C: cShares,
	}
}

func (c *Context) activeMultiplicationProtocol(parties []*Party, triple ActiveTriple, dShares, eShares []MatrixShare) []MatrixShare {
	zShares := make([]MatrixShare, len(parties))

	d, err := c.algo.authenticatedOpenMatrix(dShares)
	if err != nil {
		panic(err)
	}

	e, err := c.algo.authenticatedOpenMatrix(eShares)
	if err != nil {
		panic(err)
	}

	for partyNumber := range parties {
		aTriple := triple.A[partyNumber]
		bTriple := triple.B[partyNumber]
		cTriple := triple.C[partyNumber]

		db := MulPublicLeft(d, bTriple)   // d * [bTriple]
		de := MultiplyMatrices(d, e)      // d * e
		ae := MulPublicRight(aTriple, e)  // [aTriple] * e
		db = AddMatrixShares(db, ae)      // d * [bTriple] + [aTriple] * e
		db = AddMatrixShares(db, cTriple) // d * [bTriple] + [aTriple] * e + [cTriple]

		// TODO: Only one party should add this share for passive, however the variability is encapsulated in the method
		// as of right now
		db = AddPublicLeft(de, db, partyNumber) // d * [bTriple] + [aTriple] * e + [cTriple] + d * e
		zShares[partyNumber] = db
	}

	return zShares
}
