package mpc

import (
	"fmt"
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
)

func (c *Context) GenerateMultiplicationTriples(r1, c1, r2, c2, amount int) []model.Triple {
	triples := make([]model.Triple, amount)
	for i := 0; i < amount; i++ {
		triples[i] = c.GenerateMultiplicationTriple(r1, c1, r2, c2)
	}
	return triples
}

func (c *Context) GenerateMultiplicationTriple(r1, c1, r2, c2 int) model.Triple {
	if c1 != r2 {
		panic(fmt.Errorf("dimensions not suitable for matrix multiplication"))
	}

	aMatrix := rand.Matrix(r1, c1)
	bMatrix := rand.Matrix(r2, c2)
	cMatrix := MultiplyMatrices(aMatrix, bMatrix)

	aShares := c.algo.createSharesForMatrix(aMatrix)
	bShares := c.algo.createSharesForMatrix(bMatrix)
	cShares := c.algo.createSharesForMatrix(cMatrix)

	return model.Triple{
		A: aShares,
		B: bShares,
		C: cShares,
	}
}

func (c *Context) multiplicationProtocol(parties []*model.Party, triple model.Triple, dShares, eShares [][][]byte) [][][]byte {
	zShares := make([][][]byte, len(parties))

	d := c.algo.openMatrix(dShares)
	e := c.algo.openMatrix(eShares)

	for partyNumber := range parties {
		aTriple := triple.A[partyNumber]
		bTriple := triple.B[partyNumber]
		cTriple := triple.C[partyNumber]

		db := MultiplyMatrices(d, bTriple) // d * [bTriple]
		de := MultiplyMatrices(d, e)       // d * e
		ae := MultiplyMatrices(aTriple, e) // [aTriple] * e
		AddMatrices(db, ae)                // d * [bTriple] + [aTriple] * e
		AddMatrices(db, cTriple)           // d * [bTriple] + [aTriple] * e + [cTriple]

		db = c.algo.addPublicLeft(de, db, partyNumber) // d * [bTriple] + [aTriple] * e + [cTriple] + d * e

		zShares[partyNumber] = db
	}

	return zShares
}
