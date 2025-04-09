package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
)

var algo = Shamir{n: 4, t: 2}

func KeyGen(amountOfParties int) (model.ExpandedPublicKey, []*model.Party) {
	parties := make([]*model.Party, amountOfParties)
	P1 := make([][][]byte, m)
	P2 := make([][][]byte, m)
	P3 := make([][][]byte, m)
	OShares := algo.createSharesForRandomMatrix(v, o)
	LShares := make([][][][]byte, amountOfParties)

	OReconstructed := algo.openMatrix(OShares) // FOR CORRECTNESS

	for i := 0; i < amountOfParties; i++ {
		LShares[i] = make([][][]byte, m)
	}

	// Generate P1i and P2i
	for i := 0; i < m; i++ {
		P1[i] = rand.CoinMatrix(parties, v, v)
		P2[i] = rand.CoinMatrix(parties, v, o)
	}

	triplesStep4 := GenerateMultiplicationTriples(amountOfParties, o, v, v, o, m)
	for i := 0; i < m; i++ {
		// Compute [P1i * O]
		P1iTimeOShares := make([][][]byte, amountOfParties)
		for partyNumber, _ := range parties {
			P1iTimeOShares[partyNumber] = MultiplyMatrices(P1[i], OShares[partyNumber])
		}
		if !reflect.DeepEqual(algo.openMatrix(P1iTimeOShares), MultiplyMatrices(P1[i], OReconstructed)) {
			panic("incorrect computation")
		}

		// Compute [O^T * (P1i * O - P2i)]
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))
		for partyNumber, _ := range parties {
			ai := triplesStep4[i].A[partyNumber]
			bi := triplesStep4[i].B[partyNumber]
			di := AddMatricesNew(MatrixTranspose(OShares[partyNumber]), ai)
			var ei [][]byte
			if partyNumber == partyNumber { // TODO: variable point
				ei = AddMatricesNew(AddMatricesNew(P1iTimeOShares[partyNumber], P2[i]), bi)
			} else {
				ei = AddMatricesNew(P1iTimeOShares[partyNumber], bi)
			}

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}
		step4Shares := multiplicationProtocol(parties, triplesStep4[i], dShares, eShares, o, v, v, o)

		// CHECK FOR CORRECTNESS
		Step4Reconstructed := algo.openMatrix(step4Shares)
		if !reflect.DeepEqual(Step4Reconstructed, MultiplyMatrices(MatrixTranspose(OReconstructed), AddMatricesNew(MultiplyMatrices(P1[i], OReconstructed), P2[i]))) {
			panic("Step4 is not equal to O^T * (P1i * O - P2i)")
		}
		// CHECK FOR CORRECTNESS

		// Compute Upper of P3i
		P3iShares := make([][][]byte, m)
		for partyNumber, _ := range parties {
			P3iShares[partyNumber] = upper(step4Shares[partyNumber])
		}

		// Open P3
		p3i := algo.openMatrix(P3iShares)
		P3[i] = p3i

		// Compute locally [(P1i + P1i^T) * OShares] + P2i
		LiShares := make([][][]byte, amountOfParties)
		for partyNumber, _ := range parties {
			if partyNumber == partyNumber { // TODO: variable point
				LiShares[partyNumber] = AddMatricesNew(MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OShares[partyNumber]), P2[i])
			} else {
				LiShares[partyNumber] = MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OShares[partyNumber])
			}
		}

		for partyNumber, _ := range parties {
			LShares[partyNumber][i] = LiShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		Li := algo.openMatrix(LiShares)
		if !reflect.DeepEqual(Li, AddMatricesNew(MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OReconstructed), P2[i])) {
			panic("Li is not equal to (P1i + P1i^T) * O + P2i")
		}
		// CHECK FOR CORRECTNESS
	}

	// Generate the models for the expanded public key / expanded secret key
	epk := model.ExpandedPublicKey{
		P1: P1,
		P2: P2,
		P3: P3,
	}

	// Set the epk and the esk shares
	for partyNumber, _ := range parties {
		eskShare := model.ExpandedSecretKey{P1: P1, L: LShares[partyNumber], O: OShares[partyNumber]}
		party := &model.Party{EskShare: eskShare, Epk: epk}
		parties[partyNumber] = party
	}

	return epk, parties
}
