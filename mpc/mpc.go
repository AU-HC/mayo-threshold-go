package mpc

import (
	"mayo-threshold-go/model"
	"reflect"
)

const m = 64
const k = 4
const n = 81
const o = 17
const v = n - o

const lambda = 2

func ComputeM(parties []*model.Party) {
	salt := Coin(parties, lambda)
	t := make([]byte, 0)                                                  // TODO: Call hash function
	triples := GenerateMultiplicationTriples(len(parties), k, v, v, o, m) // V: k x v, Li: v x o

	VReconstructed := generateZeroMatrix(k, v)
	for _, party := range parties {
		V := RandMatrix(k, v)
		AddMatrices(VReconstructed, V)
		party.Salt = salt
		party.T = t
		party.V = V
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
		for j, party := range parties {
			party.M[i] = zShares[j]
		}

		MReconstructed := generateZeroMatrix(k, o)
		LReconstructed := generateZeroMatrix(v, o)
		for _, party := range parties {
			AddMatrices(MReconstructed, party.M[i])
			AddMatrices(LReconstructed, party.EskShare.L[i])
		}
		if !reflect.DeepEqual(MReconstructed, MultiplyMatrices(VReconstructed, LReconstructed)) {
			panic("M is not equal to V * L")
		}
	}
}

func ComputeY(parties []*model.Party) {
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
		for j, party := range parties {
			party.Y[i] = zShares[j]
		}

		YReconstructed := generateZeroMatrix(k, k)
		for _, party := range parties {
			AddMatrices(YReconstructed, party.Y[i])
		}
		if !reflect.DeepEqual(YReconstructed, MultiplyMatrices(MultiplyMatrices(parties[0].VReconstructed,
			parties[0].Epk.P1[i]), MatrixTranspose(parties[0].VReconstructed))) {
			panic("Y is not equal to V * P1 * V^T")
		}
	}
}

func LocalComputeA(parties []*model.Party) {
	for _, party := range parties {
		A := generateZeroMatrix(m, k*o)
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
			for j := t; j < k; j++ {
				elmjhat := MultiplyMatrixWithConstant(MHat[j], byte(ell))
				for i := 0; i < m; i++ {
					for y := t * o; y < (t+1)*o; y++ {
						A[i][y] ^= elmjhat[i][y-t*o]
					}
				}

				if t != j {
					elmthat := MultiplyMatrixWithConstant(MHat[t], byte(ell))
					for i := 0; i < m; i++ {
						for y := j * o; y < (j+1)*o; y++ {
							A[i][y] ^= elmthat[i][y-j*o]
						}
					}
				}

				ell++
			}
		}

		party.A = A
	}
}

func LocalComputeY(parties []*model.Party) {
	for _, party := range parties {
		y := make([]byte, m)
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

				for i := 0; i < len(y); i++ {
					y[i] ^= gf16Mul(byte(ell), u[i])
				}

				ell++
			}
		}

		t := make([]byte, m) // TODO: This should be H(msg || salt)
		for i := 0; i < m; i++ {
			y[i] ^= t[i]
		}
		party.LittleY = y
	}
}

func ComputeX(parties []*model.Party) {

}

func ComputeS(parties []*model.Party) {

}

func ComputeSignature(parties []*model.Party) {

}

func ComputeT(parties []*model.Party) bool {
	s := len(parties[0].A)
	t := len(parties[0].A[0])
	triplesStep2 := GenerateMultiplicationTriple(len(parties), s, t, t, t)
	triplesStep3 := GenerateMultiplicationTriple(len(parties), s, s, s, t)

	// Compute [A * S] = [A] * [S]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		S := RandMatrix(t, t)

		ai := triplesStep2.A[partyNumber]
		bi := triplesStep2.B[partyNumber]
		di := AddMatricesNew(party.A, ai)
		ei := AddMatricesNew(S, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesSShares := multiplicationProtocol(parties, triplesStep2, dShares, eShares, s, t, t, t)

	// Compute [T] = [R] * [A * S]
	dShares = make([][][]byte, len(parties))
	eShares = make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.R = RandMatrix(s, s)

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

	return rankOfMatrix(T) < s
}

func rankOfMatrix(t [][]byte) int {
	if len(t) == 0 || len(t[0]) == 0 {
		return 0
	}

	rows, cols := len(t), len(t[0])
	rank := 0

	for col := 0; col < cols; col++ {
		pivotRow := -1
		for row := rank; row < rows; row++ {
			if t[row][col] != 0 {
				pivotRow = row
				break
			}
		}

		if pivotRow == -1 {
			continue
		}

		t[pivotRow], t[rank] = t[rank], t[pivotRow]

		pivot := t[rank][col]
		for c := col; c < cols; c++ {
			t[rank][c] /= pivot
		}

		for row := 0; row < rows; row++ {
			if row != rank && t[row][col] != 0 {
				factor := t[row][col]
				for c := col; c < cols; c++ {
					t[row][c] -= factor * t[rank][c]
				}
			}
		}

		rank++
	}

	return rank
}

func ComputeAInverse(parties []*model.Party) {

}

func Computex(parties []*model.Party) {

}
