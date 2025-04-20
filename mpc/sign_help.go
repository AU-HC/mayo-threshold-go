package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
)

func (c *Context) computeM(parties []*model.Party, message []byte, iteration int) {
	salt := rand.Coin(parties, lambda)
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
			ei := AddMatricesNew(party.EskShare.L[i], bi)

			party.VReconstructed = VOpen

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d, e and compute locally
		zShares := c.multiplicationProtocol(parties, c.signTriples.ComputeM[iteration][i], dShares, eShares)
		for partyNumber, party := range parties {
			party.M[i] = zShares[partyNumber]
		}
	}
}

func (c *Context) computeY(parties []*model.Party, iteration int) {
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
	}
}

func (c *Context) localComputeA(parties []*model.Party) {
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

func (c *Context) localComputeY(parties []*model.Party) {
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

func (c *Context) computeSignature(parties []*model.Party) model.ThresholdSignature {
	// [X * O^T] = [X] * [O^t]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.X = matrixify(party.LittleX, k, o)

		ai := c.signTriples.ComputeSignature.A[partyNumber]
		bi := c.signTriples.ComputeSignature.B[partyNumber]
		di := AddMatricesNew(party.X, ai)
		ei := AddMatricesNew(MatrixTranspose(party.EskShare.O), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	xTimesOTransposedShares := c.multiplicationProtocol(parties, c.signTriples.ComputeSignature, dShares, eShares)

	// [S'] = [V + (OX^T)^T)]
	SPrimeShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		SPrimeShares[partyNumber] = AddMatricesNew(party.V, xTimesOTransposedShares[partyNumber])
		party.SPrime = SPrimeShares[partyNumber]
	}

	signatureShares := make([][][]byte, len(parties))
	for i, party := range parties {
		signatureShares[i] = appendMatrixHorizontal(party.SPrime, party.X)
	}

	return model.ThresholdSignature{
		S:    signatureShares,
		Salt: parties[0].Salt,
	}
}
