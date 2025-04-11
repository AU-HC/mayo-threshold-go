package mpc

import (
	"fmt"
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
)

func (c *Context) computeT(parties []*model.Party) bool {
	s := len(parties[0].A)
	t := len(parties[0].A[0])
	triplesStep2 := c.GenerateMultiplicationTriple(s, t, t, t)
	triplesStep3 := c.GenerateMultiplicationTriple(s, s, s, t)

	SShares := c.algo.createSharesForRandomMatrix(t, t)
	RShares := c.algo.createSharesForRandomMatrix(s, s)
	for partyNumber, party := range parties {
		party.S = SShares[partyNumber]
		party.R = RShares[partyNumber]
	}

	// Compute [A * S] = [A] * [S]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		ai := triplesStep2.A[partyNumber]
		bi := triplesStep2.B[partyNumber]
		di := AddMatricesNew(party.A, ai)
		ei := AddMatricesNew(party.S, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesSShares := c.multiplicationProtocol(parties, triplesStep2, dShares, eShares)

	// Compute [T] = [R] * [A * S]
	dShares = make([][][]byte, len(parties))
	eShares = make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		ai := triplesStep3.A[partyNumber]
		bi := triplesStep3.B[partyNumber]
		di := AddMatricesNew(party.R, ai)
		ei := AddMatricesNew(ATimesSShares[partyNumber], bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	TShares := c.multiplicationProtocol(parties, triplesStep3, dShares, eShares)

	// Open T and check rank
	T := c.algo.openMatrix(TShares)

	for _, party := range parties {
		party.T = T
	}

	return isFullRank(T)
}

func (c *Context) computeAInverse(parties []*model.Party) {
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	triple := c.GenerateMultiplicationTriple(t, s, s, s)
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))

	// Compute locally
	for partyNumber, party := range parties {
		ai := triple.A[partyNumber]
		bi := triple.B[partyNumber]
		TInverse := computeRightInverse(party.T)

		di := AddMatricesNew(MultiplyMatrices(party.S, TInverse), ai)
		ei := AddMatricesNew(party.R, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	zShares := c.multiplicationProtocol(parties, triple, dShares, eShares)
	for partyNumber, party := range parties {
		party.AInverse = zShares[partyNumber]
	}
}

func (c *Context) computeLittleX(parties []*model.Party) {
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	basis := computeBasisOfKernel(parties[0].T)

	if len(basis) != t-s {
		panic(fmt.Errorf("length of basis is incorrect"))
	}

	for _, party := range parties {
		z := make([]byte, t)
		zVector := rand.Vector(t - s)

		for index, elem := range zVector {
			z = AddVec(z, MultiplyVecConstant(elem, basis[index]))
		}

		party.Z = z
	}

	triplesStep7 := c.GenerateMultiplicationTriple(t, s, s, 1)
	triplesStep8 := c.GenerateMultiplicationTriple(t, t, t, 1)

	// Compute [A^-1] * [b]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		ai := triplesStep7.A[partyNumber]
		bi := triplesStep7.B[partyNumber]
		di := AddMatricesNew(party.AInverse, ai)
		ei := AddMatricesNew(vectorToMatrix(party.LittleY), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	AInvTimesB := c.multiplicationProtocol(parties, triplesStep7, dShares, eShares)

	// Compute [S] * [z]
	dShares = make([][][]byte, len(parties))
	eShares = make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		ai := triplesStep8.A[partyNumber]
		bi := triplesStep8.B[partyNumber]
		di := AddMatricesNew(party.S, ai)
		ei := AddMatricesNew(vectorToMatrix(party.Z), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	STimesZ := c.multiplicationProtocol(parties, triplesStep8, dShares, eShares)

	// [x] = [A^-1] * [b] + [S] * [z]
	for i, party := range parties {
		party.LittleX = matrixToVec(AddMatricesNew(AInvTimesB[i], STimesZ[i]))
	}

	// CHECK FOR CORRECTNESS
	AShares := make([][][]byte, len(parties))
	XShares := make([][][]byte, len(parties))
	YShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		AShares[partyNumber] = party.A
		XShares[partyNumber] = vectorToMatrix(party.LittleX)
		YShares[partyNumber] = vectorToMatrix(party.LittleY)
	}
	AOpen := c.algo.openMatrix(AShares)
	XOpen := c.algo.openMatrix(XShares)
	YOpen := c.algo.openMatrix(YShares)
	// CHECK FOR CORRECTNESS

	//p := parties[0]
	ATimesX := MultiplyMatrices(AOpen, XOpen)
	//xd1 := MultiplyMatrices(p.A, AddMatricesNew(MultiplyMatrices(p.AInverse, vectorToMatrix(p.LittleY)), MultiplyMatrices(p.S, vectorToMatrix(p.Z))))
	//xd2 := AddMatricesNew(MultiplyMatrices(p.A, MultiplyMatrices(p.S, MultiplyMatrices(computeRightInverse(p.T), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY))))), MultiplyMatrices(p.A, MultiplyMatrices(p.S, vectorToMatrix(p.Z))))
	//xd3 := MultiplyMatrices(computeRightInverse(p.R), AddMatricesNew(MultiplyMatrices(p.R, MultiplyMatrices(p.A, MultiplyMatrices(p.S, MultiplyMatrices(computeRightInverse(p.T), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY)))))), MultiplyMatrices(p.R, MultiplyMatrices(p.A, MultiplyMatrices(p.S, vectorToMatrix(p.Z))))))
	//xd4 := MultiplyMatrices(computeRightInverse(p.R), AddMatricesNew(MultiplyMatrices(p.T, MultiplyMatrices(computeRightInverse(p.T), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY)))), MultiplyMatrices(p.T, vectorToMatrix(p.Z))))
	//xd5 := MultiplyMatrices(computeRightInverse(p.R), MultiplyMatrices(Identity(s), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY))))
	xd6 := YOpen
	if !reflect.DeepEqual(ATimesX, xd6) {
		panic("solve did not find a correct solution")
	}
}
