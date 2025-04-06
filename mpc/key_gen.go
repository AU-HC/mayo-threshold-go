package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
)

func KeyGen(AmountOfParties int) []*model.Party {
	parties := make([]*model.Party, AmountOfParties)
	P1 := make([][][]byte, m)
	P2 := make([][][]byte, m)

	// Generate P1i and P2i
	for i := 0; i < m; i++ {
		P1[i] = rand.CoinMatrix(parties, v, v)
		P2[i] = rand.CoinMatrix(parties, v, o)
	}

	triplesStep2 := GenerateMultiplicationTriple(len(parties), s, t, t, t)

	// Compute [A * S] = [A] * [S]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.S = rand.Matrix(t, t)

		ai := triplesStep2.A[partyNumber]
		bi := triplesStep2.B[partyNumber]
		di := AddMatricesNew(party.A, ai)
		ei := AddMatricesNew(party.S, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesSShares := multiplicationProtocol(parties, triplesStep2, dShares, eShares, s, t, t, t)

	//for partyNumber, party := range parties {
	//	O := rand.Matrix(v, o)
	//}
	return nil
}
