package mpc

import (
	"fmt"
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
)

func computeT(parties []*model.Party) bool {
	s := len(parties[0].A)
	t := len(parties[0].A[0])
	triplesStep2 := GenerateMultiplicationTriple(len(parties), s, t, t, t)
	triplesStep3 := GenerateMultiplicationTriple(len(parties), s, s, s, t)

	// Compute [A * S] = [A] * [S]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.S = rand.RandMatrix(t, t)

		ai := triplesStep2.A[partyNumber]
		bi := triplesStep2.B[partyNumber]
		di := AddMatricesNew(party.A, ai)
		ei := AddMatricesNew(party.S, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesSShares := multiplicationProtocol(parties, triplesStep2, dShares, eShares, s, t, t, t)

	// Compute [T] = [R] * [A * S]
	dShares = make([][][]byte, len(parties))
	eShares = make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.R = rand.RandMatrix(s, s)

		ai := triplesStep3.A[partyNumber]
		bi := triplesStep3.B[partyNumber]
		di := AddMatricesNew(party.R, ai)
		ei := AddMatricesNew(ATimesSShares[partyNumber], bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	TShares := multiplicationProtocol(parties, triplesStep3, dShares, eShares, s, s, s, t)

	// Open T and check rank
	T := generateZeroMatrix(s, t)
	for _, tShare := range TShares {
		AddMatrices(T, tShare)
	}

	for _, party := range parties {
		party.T = T
	}

	copyOfT := generateZeroMatrix(s, t)
	for row := 0; row < len(T); row++ {
		copy(copyOfT[row][:], T[row][:])
	}

	return rankOfMatrix(copyOfT) < s
}

func computeAInverse(parties []*model.Party) {
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	triple := GenerateMultiplicationTriple(len(parties), t, s, s, s)
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
	zShares := multiplicationProtocol(parties, triple, dShares, eShares, t, s, s, s)
	for partyNumber, party := range parties {
		party.AInverse = zShares[partyNumber]
	}
}

func computeLittleX(parties []*model.Party) {
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	basis := computeBasisOfKernel(parties[0].T)

	if len(basis) != t-s {
		panic(fmt.Errorf("length of basis is incorrect"))
	}

	for _, party := range parties {
		z := make([]byte, t)
		zVector := rand.RandVector(t - s)

		for index, elem := range zVector {
			z = AddVec(z, MultiplyVecConstant(elem, basis[index]))
		}

		party.Z = z
	}

	triplesStep7 := GenerateMultiplicationTriple(len(parties), t, s, s, 1)
	triplesStep8 := GenerateMultiplicationTriple(len(parties), t, t, t, 1)

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
	AInvTimesB := multiplicationProtocol(parties, triplesStep7, dShares, eShares, t, s, s, 1)

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
	STimesZ := multiplicationProtocol(parties, triplesStep8, dShares, eShares, t, t, t, 1)

	// [x] = [A^-1] * [b] + [S] * [z]
	for i, party := range parties {
		party.LittleX = matrixToVec(AddMatricesNew(AInvTimesB[i], STimesZ[i]))
	}

	// CHECK FOR CORRECTNESS
	AReconstructed := generateZeroMatrix(s, t)
	XReconstructed := generateZeroMatrix(t, 1)
	YReconstructed := generateZeroMatrix(s, 1)
	for _, party := range parties {
		AddMatrices(AReconstructed, party.A)
		AddMatrices(XReconstructed, vectorToMatrix(party.LittleX))
		AddMatrices(YReconstructed, vectorToMatrix(party.LittleY))
	}
	// CHECK FOR CORRECTNESS

	//p := parties[0]
	ATimesX := MultiplyMatrices(AReconstructed, XReconstructed)
	//xd1 := MultiplyMatrices(p.A, AddMatricesNew(MultiplyMatrices(p.AInverse, vectorToMatrix(p.LittleY)), MultiplyMatrices(p.S, vectorToMatrix(p.Z))))
	//xd2 := AddMatricesNew(MultiplyMatrices(p.A, MultiplyMatrices(p.S, MultiplyMatrices(computeRightInverse(p.T), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY))))), MultiplyMatrices(p.A, MultiplyMatrices(p.S, vectorToMatrix(p.Z))))
	//xd3 := MultiplyMatrices(computeRightInverse(p.R), AddMatricesNew(MultiplyMatrices(p.R, MultiplyMatrices(p.A, MultiplyMatrices(p.S, MultiplyMatrices(computeRightInverse(p.T), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY)))))), MultiplyMatrices(p.R, MultiplyMatrices(p.A, MultiplyMatrices(p.S, vectorToMatrix(p.Z))))))
	//xd4 := MultiplyMatrices(computeRightInverse(p.R), AddMatricesNew(MultiplyMatrices(p.T, MultiplyMatrices(computeRightInverse(p.T), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY)))), MultiplyMatrices(p.T, vectorToMatrix(p.Z))))
	//xd5 := MultiplyMatrices(computeRightInverse(p.R), MultiplyMatrices(Identity(s), MultiplyMatrices(p.R, vectorToMatrix(p.LittleY))))
	xd6 := YReconstructed
	if !reflect.DeepEqual(ATimesX, xd6) {
		panic("solve did not find a correct solution")
	}
}

func MultiplyVecConstant(b byte, a []byte) []byte {
	C := make([]byte, len(a))
	for i := range C {
		C[i] = gf16Mul(b, a[i])
	}
	return C
}
