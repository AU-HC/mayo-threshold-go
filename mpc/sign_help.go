package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
)

func (c *Context) computeM(parties []*Party, message []byte, iteration int) {
	salt := rand.Coin(len(parties), lambda)
	t := rand.Shake256(m, message, salt)

	for _, party := range parties {
		V := rand.Matrix(k, v)
		party.Salt = salt
		party.V = V
		party.LittleT = t
		party.M = make([][][]byte, m)
		party.Y = make([][][]byte, m)
	}

	VShares := make([][][]byte, len(parties))
	for i, party := range parties {
		VShares[i] = party.V
	}
	VOpen := c.algo.openMatrix(VShares)

	for i := 0; i < m; i++ {
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))

		// Compute locally
		for partyNumber, party := range parties {
			ai := c.signTriples.ComputeM[iteration][i].A[partyNumber]
			bi := c.signTriples.ComputeM[iteration][i].B[partyNumber]
			di := AddMatricesNew(party.V, ai)
			ei := AddMatricesNew(party.EskShare.L[i].shares, bi)

			party.VReconstructed = VOpen

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d, e and compute locally
		zShares := c.multiplicationProtocol(parties, c.signTriples.ComputeM[iteration][i], dShares, eShares)
		for partyNumber, party := range parties {
			party.M[i] = zShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		MShares := make([][][]byte, len(parties))
		LShares := make([][][]byte, len(parties))
		for partyNumber, party := range parties {
			MShares[partyNumber] = party.M[i]
			LShares[partyNumber] = party.EskShare.L[i].shares
		}
		MOpen := c.algo.openMatrix(MShares)
		LOpen := c.algo.openMatrix(LShares)

		if !reflect.DeepEqual(MOpen, MultiplyMatrices(VOpen, LOpen)) {
			panic("M is not equal to V * L")
		}
		// CHECK FOR CORRECTNESS
	}
}

func (c *Context) computeY(parties []*Party, iteration int) {
	for i := 0; i < m; i++ {
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))

		// Compute locally
		for partyNumber, party := range parties {
			ai := c.signTriples.ComputeY[iteration][i].A[partyNumber]
			bi := c.signTriples.ComputeY[iteration][i].B[partyNumber]
			di := AddMatricesNew(MultiplyMatrices(party.V, party.Epk.P1[i]), ai)
			ei := AddMatricesNew(MatrixTranspose(party.V), bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d, e and compute locally
		zShares := c.multiplicationProtocol(parties, c.signTriples.ComputeY[iteration][i], dShares, eShares)
		for partyNumber, party := range parties {
			party.Y[i] = zShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		YShares := make([][][]byte, len(parties))
		for partyNumber, party := range parties {
			YShares[partyNumber] = party.Y[i]
		}
		YOpen := c.algo.openMatrix(YShares)
		if !reflect.DeepEqual(YOpen, MultiplyMatrices(MultiplyMatrices(parties[0].VReconstructed,
			parties[0].Epk.P1[i]), MatrixTranspose(parties[0].VReconstructed))) {
			panic("Y is not equal to V * P1 * V^T")
		}
		// CHECK FOR CORRECTNESS
	}
}

func (c *Context) localComputeA(parties []*Party) {
	for _, party := range parties {
		A := generateZeroMatrix(m+shifts, k*o)
		ell := 0
		MHat := make([][][]byte, k)
		for index := 0; index < k; index++ {
			MHat[index] = generateZeroMatrix(m, o)
		}

		for t := 0; t < k; t++ {
			for j := 0; j < m; j++ {
				copy(MHat[t][j][:], party.M[j][t][:])
			}
		}

		for t := 0; t < k; t++ {
			for j := k - 1; j >= t; j-- {
				for row := 0; row < m; row++ {
					for column := t * o; column < (t+1)*o; column++ {
						A[row+ell][column] ^= MHat[j][row][column%o]
					}

					if t != j {
						for column := j * o; column < (j+1)*o; column++ {
							A[row+ell][column] ^= MHat[t][row][column%o]
						}
					}
				}

				ell++
			}
		}

		A = reduceAModF(A)
		party.A = A
	}
}

func (c *Context) localComputeY(parties []*Party) {
	for partyNumber, party := range parties {
		y := make([]byte, m+shifts)
		ell := 0

		for j := 0; j < k; j++ {
			for t := k - 1; t >= j; t-- {
				u := make([]byte, m)
				if j == t {
					for a := 0; a < m; a++ {
						u[a] = party.Y[a][j][j]
					}
				} else {
					for a := 0; a < m; a++ {
						u[a] = party.Y[a][j][t] ^ party.Y[a][t][j]
					}
				}

				for d := 0; d < m; d++ {
					y[d+ell] ^= u[d]
				}

				ell++
			}
		}

		y = reduceVecModF(y)
		if c.algo.shouldPartyAddConstantShare(partyNumber) {
			t := party.LittleT
			for i := 0; i < m; i++ {
				y[i] ^= t[i]
			}
		}
		party.LittleY = y
	}
}

func (c *Context) computeSignature(parties []*Party) model.ThresholdSignature {
	// [X * O^T] = [X] * [O^t]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.X = matrixify(party.LittleX, k, o)

		ai := c.signTriples.ComputeSignature.A[partyNumber]
		bi := c.signTriples.ComputeSignature.B[partyNumber]
		di := AddMatricesNew(party.X, ai)
		ei := AddMatricesNew(MatrixTranspose(party.EskShare.O.shares), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	xTimesOTransposedShares := c.multiplicationProtocol(parties, c.signTriples.ComputeSignature, dShares, eShares)

	// CHECK FOR CORRECTNESS
	xTimesOTransposedOpen := c.algo.openMatrix(xTimesOTransposedShares)
	XShares := make([][][]byte, len(parties))
	OShares := make([][][]byte, len(parties))
	VShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		XShares[partyNumber] = party.X
		OShares[partyNumber] = party.EskShare.O.shares
		VShares[partyNumber] = party.V
	}
	XOpen := c.algo.openMatrix(XShares)
	OOpen := c.algo.openMatrix(OShares)
	VOpen := c.algo.openMatrix(VShares)
	if !reflect.DeepEqual(xTimesOTransposedOpen, MultiplyMatrices(XOpen, MatrixTranspose(OOpen))) {
		panic("XO^T != X * O^T")
	}
	if !reflect.DeepEqual(xTimesOTransposedOpen, MatrixTranspose(MultiplyMatrices(OOpen, MatrixTranspose(XOpen)))) {
		panic("XO^T != (OX^T)^T")
	}
	// CHECK FOR CORRECTNESS

	// [S'] = [V + (OX^T)^T)]
	SPrimeShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		SPrimeShares[partyNumber] = AddMatricesNew(party.V, xTimesOTransposedShares[partyNumber])
		party.SPrime = SPrimeShares[partyNumber]
	}

	// Open S' and X
	SPrimeOpen := c.algo.openMatrix(SPrimeShares)

	s := appendMatrixHorizontal(SPrimeOpen, XOpen)
	for _, party := range parties {
		party.Signature = appendMatrixHorizontal(party.SPrime, party.X)
	}

	// CHECK FOR CORRECTNESS
	if !reflect.DeepEqual(SPrimeOpen, AddMatricesNew(VOpen, xTimesOTransposedOpen)) {
		panic("S' != V + XO^T")
	}
	if (len(s) * len(s[0])) != (k * n) {
		panic("signature invalid size")
	}
	// CHECK FOR CORRECTNESS

	signatureShares := make([][][]byte, len(parties))
	for i, party := range parties {
		signatureShares[i] = party.Signature
	}

	return model.ThresholdSignature{
		S:    signatureShares,
		Salt: parties[0].Salt,
	}
}
