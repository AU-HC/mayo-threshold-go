package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
)

func (c *Context) KeyGenAPI(amountOfParties int) (model.ExpandedPublicKey, []*model.Party) {
	c.PreprocessMultiplicationKeyGenTriples()
	return c.KeyGen(amountOfParties)
}

func (c *Context) KeyGen(amountOfParties int) (model.ExpandedPublicKey, []*model.Party) {
	parties := make([]*model.Party, amountOfParties)
	P1 := make([][][]byte, m)
	P2 := make([][][]byte, m)
	P3 := make([][][]byte, m)
	OShares := c.algo.createSharesForRandomMatrix(v, o)
	LShares := make([][][][]byte, amountOfParties)

	for i := 0; i < amountOfParties; i++ {
		LShares[i] = make([][][]byte, m)
	}

	// Generate P1i and P2i
	for i := 0; i < m; i++ {
		P1[i] = rand.CoinMatrix(parties, v, v)
		P2[i] = rand.CoinMatrix(parties, v, o)
	}

	for i := 0; i < m; i++ {
		// Compute [P1i * O]
		P1iTimeOShares := make([][][]byte, amountOfParties)
		for partyNumber, _ := range parties {
			P1iTimeOShares[partyNumber] = MultiplyMatrices(P1[i], OShares[partyNumber])
		}

		// Compute [O^T * (P1i * O - P2i)]
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))
		for partyNumber, _ := range parties {
			ai := c.keygenTriples.TriplesStep4[i].A[partyNumber]
			bi := c.keygenTriples.TriplesStep4[i].B[partyNumber]
			di := AddMatricesNew(MatrixTranspose(OShares[partyNumber]), ai)
			ei := AddMatricesNew(c.algo.addPublicLeft(P2[i], P1iTimeOShares[partyNumber], partyNumber), bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}
		step4Shares := c.multiplicationProtocol(parties, c.keygenTriples.TriplesStep4[i], dShares, eShares)

		// Compute Upper of P3i
		P3iShares := make([][][]byte, amountOfParties)
		for partyNumber, _ := range parties {
			P3iShares[partyNumber] = upper(step4Shares[partyNumber])
		}

		// Open P3
		p3i := c.algo.openMatrix(P3iShares)
		P3[i] = p3i

		// Compute locally [(P1i + P1i^T) * OShares] + P2i
		LiShares := make([][][]byte, amountOfParties)
		for partyNumber, _ := range parties {
			LiShares[partyNumber] = c.algo.addPublicLeft(P2[i], MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OShares[partyNumber]), partyNumber)
		}

		for partyNumber, _ := range parties {
			LShares[partyNumber][i] = LiShares[partyNumber]
		}
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
