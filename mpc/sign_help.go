package mpc

import (
	"mayo-threshold-go/rand"
	"reflect"
)

func (c *Context) computeM(parties []*Party, message []byte, iteration int) {
	salt := rand.Coin(len(parties), lambda)
	t := rand.Shake256(m, message, salt)

	VShares := c.algo.createSharesForRandomMatrix(k, v)
	for i, party := range parties {
		party.Salt = salt
		party.V = VShares[i]
		party.LittleT = t
		party.M = make([]MatrixShare, m)
		party.Y = make([]MatrixShare, m)
	}

	VOpen, err := c.algo.authenticatedOpenMatrix(VShares)
	if err != nil {
		panic(err)
	}

	for i := 0; i < m; i++ {
		dShares := make([]MatrixShare, len(parties))
		eShares := make([]MatrixShare, len(parties))

		// Compute locally
		for partyNumber, party := range parties {
			ai := c.signTriples.ComputeM[iteration][i].A[partyNumber]
			bi := c.signTriples.ComputeM[iteration][i].B[partyNumber]
			di := AddMatrixShares(party.V, ai)
			ei := AddMatrixShares(party.EskShare.L[i], bi)

			party.VReconstructed = VOpen

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d, e and compute locally
		zShares := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeM[iteration][i], dShares, eShares)
		for partyNumber, party := range parties {
			party.M[i] = zShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		MShares := make([]MatrixShare, len(parties))
		LShares := make([]MatrixShare, len(parties))
		for partyNumber, party := range parties {
			MShares[partyNumber] = party.M[i]
			LShares[partyNumber] = party.EskShare.L[i]
		}
		MOpen, err := c.algo.authenticatedOpenMatrix(MShares)
		if err != nil {
			panic(err)
		}
		LOpen, err := c.algo.authenticatedOpenMatrix(LShares)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(MOpen, MultiplyMatrices(VOpen, LOpen)) {
			panic("M is not equal to V * L")
		}
		// CHECK FOR CORRECTNESS
	}
}

func (c *Context) computeY(parties []*Party, iteration int) {
	for i := 0; i < m; i++ {
		dShares := make([]MatrixShare, len(parties))
		eShares := make([]MatrixShare, len(parties))

		// Compute locally
		for partyNumber, party := range parties {
			ai := c.signTriples.ComputeY[iteration][i].A[partyNumber]
			bi := c.signTriples.ComputeY[iteration][i].B[partyNumber]
			di := AddMatrixShares(MulPublicRight(party.V, party.Epk.P1[i]), ai)
			ei := AddMatrixShares(MatrixShareTranspose(party.V), bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d, e and compute locally
		zShares := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeY[iteration][i], dShares, eShares)
		for partyNumber, party := range parties {
			party.Y[i] = zShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		YShares := make([]MatrixShare, len(parties))
		for partyNumber, party := range parties {
			YShares[partyNumber] = party.Y[i]
		}
		YOpen, err := c.algo.authenticatedOpenMatrix(YShares)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(YOpen, MultiplyMatrices(MultiplyMatrices(parties[0].VReconstructed,
			parties[0].Epk.P1[i]), MatrixTranspose(parties[0].VReconstructed))) {
			panic("Y is not equal to V * P1 * V^T")
		}
		// CHECK FOR CORRECTNESS
	}
}

func (c *Context) localComputeA(parties []*Party) {
	for _, party := range parties {
		A := createEmptyMatrixShare(m+shifts, k*o)
		ell := 0
		MHatShares := make([][][]byte, k)
		MHatGammas := make([][][]byte, k)
		for index := 0; index < k; index++ {
			MHatShares[index] = generateZeroMatrix(m, o)
			MHatGammas[index] = generateZeroMatrix(m, o)
		}

		for t := 0; t < k; t++ {
			for j := 0; j < m; j++ {
				copy(MHatShares[t][j][:], party.M[j].shares[t][:])
				copy(MHatGammas[t][j][:], party.M[j].gammas[t][:])
			}
		}

		for t := 0; t < k; t++ {
			for j := k - 1; j >= t; j-- {
				for row := 0; row < m; row++ {
					for column := t * o; column < (t+1)*o; column++ {
						A.shares[row+ell][column] ^= MHatShares[j][row][column%o]
						A.gammas[row+ell][column] ^= MHatGammas[j][row][column%o]
					}

					if t != j {
						for column := j * o; column < (j+1)*o; column++ {
							A.shares[row+ell][column] ^= MHatShares[t][row][column%o]
							A.gammas[row+ell][column] ^= MHatGammas[t][row][column%o]
						}
					}
				}

				ell++
			}
		}

		A.shares = reduceAModF(A.shares)
		A.gammas = reduceAModF(A.gammas)
		A.alpha = party.M[0].alpha
		party.A = A
	}
}

func (c *Context) localComputeY(parties []*Party) {
	for partyNumber, party := range parties {
		y := createEmptyMatrixShare(m+shifts, 1)
		ell := 0

		for j := 0; j < k; j++ {
			for t := k - 1; t >= j; t-- {
				uShares := make([]byte, m)
				uGammas := make([]byte, m)
				if j == t {
					for a := 0; a < m; a++ {
						uShares[a] = party.Y[a].shares[j][j]
						uGammas[a] = party.Y[a].gammas[j][j]
					}
				} else {
					for a := 0; a < m; a++ {
						uShares[a] = party.Y[a].shares[j][t] ^ party.Y[a].shares[t][j]
						uGammas[a] = party.Y[a].gammas[j][t] ^ party.Y[a].gammas[t][j]
					}
				}

				for d := 0; d < m; d++ {
					y.shares[d+ell][0] ^= uShares[d]
					y.gammas[d+ell][0] ^= uGammas[d]
				}

				ell++
			}
		}

		y.shares = vectorToMatrix(reduceVecModF(matrixToVec(y.shares)))
		y.gammas = vectorToMatrix(reduceVecModF(matrixToVec(y.gammas)))
		y.alpha = party.Y[0].alpha

		t := party.LittleT
		y = AddPublicLeft(vectorToMatrix(t), y, partyNumber)
		party.LittleY = y
	}
}

func (c *Context) computeSignature(parties []*Party) ThresholdSignature {
	// [X * O^T] = [X] * [O^t]
	dShares := make([]MatrixShare, len(parties))
	eShares := make([]MatrixShare, len(parties))
	for partyNumber, party := range parties {
		party.X = matrixifyActive(party.LittleX, k, o)

		ai := c.signTriples.ComputeSignature.A[partyNumber]
		bi := c.signTriples.ComputeSignature.B[partyNumber]
		di := AddMatrixShares(party.X, ai)
		ei := AddMatrixShares(MatrixShareTranspose(party.EskShare.O), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	xTimesOTransposedShares := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeSignature, dShares, eShares)

	// CHECK FOR CORRECTNESS
	xTimesOTransposedOpen, err := c.algo.authenticatedOpenMatrix(xTimesOTransposedShares)
	if err != nil {
		panic(err)
	}

	XShares := make([]MatrixShare, len(parties))
	OShares := make([]MatrixShare, len(parties))
	VShares := make([]MatrixShare, len(parties))
	for partyNumber, party := range parties {
		XShares[partyNumber] = party.X
		OShares[partyNumber] = party.EskShare.O
		VShares[partyNumber] = party.V
	}
	XOpen, err := c.algo.authenticatedOpenMatrix(XShares)
	if err != nil {
		panic(err)
	}
	OOpen, err := c.algo.authenticatedOpenMatrix(OShares)
	if err != nil {
		panic(err)
	}
	VOpen, err := c.algo.authenticatedOpenMatrix(VShares)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(xTimesOTransposedOpen, MultiplyMatrices(XOpen, MatrixTranspose(OOpen))) {
		panic("XO^T != X * O^T")
	}
	if !reflect.DeepEqual(xTimesOTransposedOpen, MatrixTranspose(MultiplyMatrices(OOpen, MatrixTranspose(XOpen)))) {
		panic("XO^T != (OX^T)^T")
	}
	// CHECK FOR CORRECTNESS

	// [SShares'] = [V + (OX^T)^T)]
	SPrimeShares := make([]MatrixShare, len(parties))
	for partyNumber, party := range parties {
		SPrimeShares[partyNumber] = AddMatrixShares(party.V, xTimesOTransposedShares[partyNumber])
		party.SPrime = SPrimeShares[partyNumber]
	}

	// Open SShares' and X
	SPrimeOpen, err := c.algo.authenticatedOpenMatrix(SPrimeShares)
	if err != nil {
		panic(err)
	}

	s := appendMatrixHorizontal(SPrimeOpen, XOpen)
	for _, party := range parties {
		party.Signature = appendMatrixShareHorizontal(party.SPrime, party.X)
	}

	// CHECK FOR CORRECTNESS
	if !reflect.DeepEqual(SPrimeOpen, AddMatricesNew(VOpen, xTimesOTransposedOpen)) {
		panic("SShares' != V + XO^T")
	}
	if (len(s) * len(s[0])) != (k * n) {
		panic("signature invalid size")
	}
	// CHECK FOR CORRECTNESS

	signatureShares := make([]MatrixShare, len(parties))
	for i, party := range parties {
		signatureShares[i] = party.Signature
	}

	return ThresholdSignature{
		SShares: signatureShares,
		Salt:    parties[0].Salt,
	}
}
