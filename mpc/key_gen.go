package mpc

import (
	"mayo-threshold-go/rand"
	"reflect"
)

func (c *Context) KeyGenAPI(amountOfParties int) (ExpandedPublicKey, []*Party) {
	c.PreprocessMultiplicationKeyGenTriples()
	return c.KeyGen(amountOfParties)
}

func (c *Context) KeyGen(amountOfParties int) (ExpandedPublicKey, []*Party) {
	parties := make([]*Party, amountOfParties)
	P1 := make([][][]byte, m)
	P2 := make([][][]byte, m)
	P3 := make([][][]byte, m)
	OShares := c.algo.createSharesForRandomMatrix(v, o)
	LShares := make([][]MatrixShare, amountOfParties)

	OReconstructed, err := c.algo.authenticatedOpenMatrix(OShares) // FOR CORRECTNESS
	if err != nil {
		panic(err)
	}

	for i := 0; i < amountOfParties; i++ {
		LShares[i] = make([]MatrixShare, m)
	}

	// Generate P1i and P2i
	for i := 0; i < m; i++ {
		P1[i] = rand.CoinMatrix(len(parties), v, v)
		P2[i] = rand.CoinMatrix(len(parties), v, o)
	}

	for i := 0; i < m; i++ {
		// Compute [P1i * O]
		P1iTimeOShares := make([]MatrixShare, amountOfParties)
		for partyNumber, _ := range parties {
			P1iTimeOShares[partyNumber] = MulPublicLeft(P1[i], OShares[partyNumber])
		}
		P1iTimeO, err := c.algo.authenticatedOpenMatrix(P1iTimeOShares)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(P1iTimeO, MultiplyMatrices(P1[i], OReconstructed)) {
			panic("incorrect computation")
		}

		// Compute [O^T * (P1i * O - P2i)]
		dShares := make([]MatrixShare, len(parties))
		eShares := make([]MatrixShare, len(parties))
		for partyNumber, _ := range parties {
			ai := c.keygenTriples.TriplesStep4[i].A[partyNumber]
			bi := c.keygenTriples.TriplesStep4[i].B[partyNumber]
			di := AddMatrixShares(MatrixShareTranspose(OShares[partyNumber]), ai)
			ei := AddMatrixShares(AddPublicLeft(P2[i], P1iTimeOShares[partyNumber], partyNumber), bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}
		step4Shares := c.activeMultiplicationProtocol(parties, c.keygenTriples.TriplesStep4[i], dShares, eShares)

		// CHECK FOR CORRECTNESS
		Step4Reconstructed, err := c.algo.authenticatedOpenMatrix(step4Shares)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(Step4Reconstructed, MultiplyMatrices(MatrixTranspose(OReconstructed), AddMatricesNew(MultiplyMatrices(P1[i], OReconstructed), P2[i]))) {
			panic("Step4 is not equal to O^T * (P1i * O - P2i)")
		}
		// CHECK FOR CORRECTNESS

		// Compute Upper of P3i
		P3iShares := make([]MatrixShare, amountOfParties)
		for partyNumber, _ := range parties {
			P3iShares[partyNumber] = upper(step4Shares[partyNumber])
		}

		// Open P3
		p3i, err := c.algo.authenticatedOpenMatrix(P3iShares)
		P3[i] = p3i
		if err != nil {
			panic(err)
		}

		// Compute locally [(P1i + P1i^T) * OShares] + P2i
		LiShares := make([]MatrixShare, amountOfParties)
		for partyNumber, _ := range parties {
			LiShares[partyNumber] = AddPublicLeft(P2[i],
				MulPublicLeft(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OShares[partyNumber]), partyNumber)
		}

		for partyNumber, _ := range parties {
			LShares[partyNumber][i] = LiShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		Li, err := c.algo.authenticatedOpenMatrix(LiShares)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(Li, AddMatricesNew(MultiplyMatrices(AddMatricesNew(P1[i], MatrixTranspose(P1[i])), OReconstructed), P2[i])) {
			panic("Li is not equal to (P1i + P1i^T) * O + P2i")
		}
		// CHECK FOR CORRECTNESS
	}

	// Generate the models for the expanded public key / expanded secret key
	epk := ExpandedPublicKey{
		P1: P1,
		P2: P2,
		P3: P3,
	}

	// Set the epk and the esk shares
	for partyNumber, _ := range parties {
		eskShare := ExpandedSecretKey{P1: P1, L: LShares[partyNumber], O: OShares[partyNumber]}
		party := &Party{EskShare: eskShare, Epk: epk}
		parties[partyNumber] = party
	}

	return epk, parties
}
