package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
)

var algo = Shamir{n: 3, t: 2}

func KeyGen(amountOfParties int) (model.ExpandedPublicKey, []*model.Party) {
	parties := make([]*model.Party, amountOfParties)
	P1 := make([][][]byte, m)
	P2 := make([][][]byte, m)
	P3 := make([][][]byte, m)
	OShares := algo.createSharesForMatrix(v, o)
	LiShares := make([][][][]byte, amountOfParties)

	// Generate OShares
	for partyNumber, _ := range parties {
		LiShares[partyNumber] = make([][][]byte, m)
	}

	// Generate P1i and P2i
	for i := 0; i < m; i++ {
		P1[i] = rand.CoinMatrix(parties, v, v)
		P2[i] = rand.CoinMatrix(parties, v, o)
	}

	triplesStep4 := GenerateMultiplicationTriples(amountOfParties, o, v, v, o, m)
	for i := 0; i < m; i++ {
		// Compute [P1i * OShares]
		P1iTimeOShares := make([][][]byte, amountOfParties)
		for partyNumber, _ := range parties {
			P1iTimeOShares[partyNumber] = MultiplyMatrices(P1[i], OShares[partyNumber])
		}

		// Compute [O^T * (P1i * O - P2i)]
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))
		for partyNumber, _ := range parties {
			ai := triplesStep4[i].A[partyNumber]
			bi := triplesStep4[i].B[partyNumber]
			di := AddMatricesNew(MatrixTranspose(OShares[partyNumber]), ai)
			var ei [][]byte
			if partyNumber == 0 {
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
		OReconstructed := algo.openMatrix(OShares)
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
		p3i := generateZeroMatrix(o, o)
		for partyNumber, _ := range parties {
			AddMatrices(p3i, P3iShares[partyNumber])
		}
		P3[i] = p3i

		// Compute locally [(P1i + P1i^T) * OShares] + P2i
		for partyNumber, _ := range parties {
			if partyNumber == 0 {
				LiShares[partyNumber][i] = AddMatricesNew(MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OShares[partyNumber]), P2[i])
			} else {
				LiShares[partyNumber][i] = MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OShares[partyNumber])
			}
		}

		// CHECK FOR CORRECTNESS
		LiReconstructed := generateZeroMatrix(v, o)
		for partyNumber, _ := range parties {
			AddMatrices(LiReconstructed, LiShares[partyNumber][i])
		}
		if !reflect.DeepEqual(LiReconstructed, AddMatricesNew(MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OReconstructed), P2[i])) {
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
		eskShare := model.ExpandedSecretKey{P1: P1, L: LiShares[partyNumber], O: OShares[partyNumber]}
		party := &model.Party{EskShare: eskShare, Epk: epk}
		parties[partyNumber] = party
	}

	return epk, parties
}
