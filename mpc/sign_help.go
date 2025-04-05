package mpc

import (
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
	"reflect"
)

func computeM(parties []*model.Party, message []byte) {
	salt := rand.Coin(parties, lambda)
	t := rand.Shake256(m, message, salt)
	triples := GenerateMultiplicationTriples(len(parties), k, v, v, o, m) // V: k x v, Li: v x o

	VReconstructed := generateZeroMatrix(k, v) // FOR CORRECTNESS
	for _, party := range parties {
		V := rand.RandMatrix(k, v)
		AddMatrices(VReconstructed, V) // FOR CORRECTNESS
		party.Salt = salt
		party.V = V
		party.LittleT = t
		party.M = make([][][]byte, m)
		party.Y = make([][][]byte, m)
	}

	for i := 0; i < m; i++ {
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))

		// Compute locally
		for partyNumber, party := range parties {
			ai := triples[i].A[partyNumber]
			bi := triples[i].B[partyNumber]
			di := AddMatricesNew(party.V, ai)
			ei := AddMatricesNew(party.EskShare.L[i], bi)

			party.VReconstructed = VReconstructed

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d, e and compute locally
		zShares := multiplicationProtocol(parties, triples[i], dShares, eShares, k, v, v, o)
		for partyNumber, party := range parties {
			party.M[i] = zShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		MReconstructed := generateZeroMatrix(k, o)
		LReconstructed := generateZeroMatrix(v, o)
		for _, party := range parties {
			AddMatrices(MReconstructed, party.M[i])
			AddMatrices(LReconstructed, party.EskShare.L[i])
		}
		if !reflect.DeepEqual(MReconstructed, MultiplyMatrices(VReconstructed, LReconstructed)) {
			panic("M is not equal to V * L")
		}
		// CHECK FOR CORRECTNESS
	}
}

func computeY(parties []*model.Party) {
	triples := GenerateMultiplicationTriples(len(parties), k, v, v, k, m) // V*P_1: k * v, V^T: v x k
	for i := 0; i < m; i++ {
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))

		// Compute locally
		for partyNumber, party := range parties {
			ai := triples[i].A[partyNumber]
			bi := triples[i].B[partyNumber]
			di := AddMatricesNew(MultiplyMatrices(party.V, party.Epk.P1[i]), ai)
			ei := AddMatricesNew(MatrixTranspose(party.V), bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d, e and compute locally
		zShares := multiplicationProtocol(parties, triples[i], dShares, eShares, k, v, v, k)
		for partyNumber, party := range parties {
			party.Y[i] = zShares[partyNumber]
		}

		// CHECK FOR CORRECTNESS
		YReconstructed := generateZeroMatrix(k, k)
		for _, party := range parties {
			AddMatrices(YReconstructed, party.Y[i])
		}
		if !reflect.DeepEqual(YReconstructed, MultiplyMatrices(MultiplyMatrices(parties[0].VReconstructed,
			parties[0].Epk.P1[i]), MatrixTranspose(parties[0].VReconstructed))) {
			panic("Y is not equal to V * P1 * V^T")
		}
		// CHECK FOR CORRECTNESS
	}
}

func localComputeA(parties []*model.Party) {
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

func localComputeY(parties []*model.Party) {
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
		if partyNumber == 0 {
			t := party.LittleT
			for i := 0; i < m; i++ {
				y[i] ^= t[i]
			}
		}
		party.LittleY = y
	}
}

func computeSignature(parties []*model.Party) model.Signature {
	// [X * O^T] = [X] * [O^t]
	triple := GenerateMultiplicationTriple(len(parties), k, o, o, v)
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.X = matrixify(party.LittleX, k, o)

		ai := triple.A[partyNumber]
		bi := triple.B[partyNumber]
		di := AddMatricesNew(party.X, ai)
		ei := AddMatricesNew(MatrixTranspose(party.EskShare.O), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	xTimesOTransposedShares := multiplicationProtocol(parties, triple, dShares, eShares, k, o, o, v)

	// CHECK FOR CORRECTNESS
	xTimesOTransposedReconstructed := generateZeroMatrix(k, v)
	XReconstructed := generateZeroMatrix(k, o)
	OReconstructed := generateZeroMatrix(v, o)
	VReconstructed := generateZeroMatrix(k, v)
	for i, party := range parties {
		AddMatrices(xTimesOTransposedReconstructed, xTimesOTransposedShares[i])
		AddMatrices(XReconstructed, party.X)
		AddMatrices(OReconstructed, party.EskShare.O)
		AddMatrices(VReconstructed, party.V)
	}
	if !reflect.DeepEqual(xTimesOTransposedReconstructed, MultiplyMatrices(XReconstructed, MatrixTranspose(OReconstructed))) {
		panic("XO^T != X * O^T")
	}
	if !reflect.DeepEqual(xTimesOTransposedReconstructed, MatrixTranspose(MultiplyMatrices(OReconstructed, MatrixTranspose(XReconstructed)))) {
		panic("XO^T != (OX^T)^T")
	}
	// CHECK FOR CORRECTNESS

	// [S'] = [V + (OX^T)^T)]
	for i, party := range parties {
		party.SPrime = AddMatricesNew(party.V, xTimesOTransposedShares[i])
	}

	// Open S' and X
	SPrimeReconstructed := generateZeroMatrix(k, v)
	xReconstructed := generateZeroMatrix(k, o)
	for _, party := range parties {
		AddMatrices(SPrimeReconstructed, party.SPrime)
		AddMatrices(xReconstructed, party.X)
	}

	s := appendMatrixHorizontal(SPrimeReconstructed, xReconstructed)
	for _, party := range parties {
		party.Signature = appendMatrixHorizontal(party.SPrime, party.X)
	}

	// CHECK FOR CORRECTNESS
	if !reflect.DeepEqual(SPrimeReconstructed, AddMatricesNew(VReconstructed, xTimesOTransposedReconstructed)) {
		panic("S' != V + XO^T")
	}
	if (len(s) * len(s[0])) != (k * n) {
		panic("signature invalid size")
	}
	// CHECK FOR CORRECTNESS

	return model.Signature{
		S:    s,
		Salt: parties[0].Salt,
	}
}
