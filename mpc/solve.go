package mpc

import (
	"fmt"
	"mayo-threshold-go/model"
)

func ComputeT(parties []*model.Party) bool {
	s := len(parties[0].A)
	t := len(parties[0].A[0])
	triplesStep2 := GenerateMultiplicationTriple(len(parties), s, t, t, t)
	triplesStep3 := GenerateMultiplicationTriple(len(parties), s, s, s, t)

	// Compute [A * S] = [A] * [S]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.S = RandMatrix(t, t)

		ai := triplesStep2.A[partyNumber]
		bi := triplesStep2.B[partyNumber]
		di := AddMatricesNew(party.A, ai)
		ei := AddMatricesNew(party.S, bi)

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

	for _, party := range parties {
		party.T = T
	}

	copyOfT := generateZeroMatrix(s, t)
	for row := 0; row < len(T); row++ {
		copy(copyOfT[row][:], T[row][:])
	}

	return rankOfMatrix(copyOfT) < s
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
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	triple := GenerateMultiplicationTriple(len(parties), t, s, s, s)
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))

	// Compute locally
	for partyNumber, party := range parties {
		ai := triple.A[partyNumber]
		bi := triple.B[partyNumber]
		TInverse := RightInverse(party.T)

		di := AddMatricesNew(MultiplyMatrices(party.S, TInverse), ai)
		ei := AddMatricesNew(party.R, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	zShares := multiplicationProtocol(parties, triple, dShares, eShares, t, s, s, s)
	for i, party := range parties {
		party.AInverse = zShares[i]
	}

	// TODO: Remove this check when benchmarking
	ARecovered := generateZeroMatrix(s, t)
	AInverseRecovered := generateZeroMatrix(t, s)
	for _, party := range parties {
		AddMatrices(ARecovered, party.A)
		AddMatrices(AInverseRecovered, party.AInverse)
	}

	Identity := MultiplyMatrices(ARecovered, AInverseRecovered)
	for _, row := range Identity {
		fmt.Println(fmt.Sprintf("%2d", row))
	}
}

func RightInverse(t [][]byte) [][]byte {
	_, invTable := generateMulAndInvTable() // TODO: refactor

	M := len(t)    // Rows
	N := len(t[0]) // Columns

	if M > N {
		return nil
	}

	// Augment A with an identity matrix to form (A | I)
	augmented := make([][]byte, M)
	for i := 0; i < M; i++ {
		augmented[i] = make([]byte, N+M)
		copy(augmented[i], t[i])
		for j := 0; j < M; j++ {
			if i == j {
				augmented[i][N+j] = 1
			}
		}
	}

	// Perform Gaussian elimination
	for i := 0; i < M; i++ {
		// Find pivot
		if augmented[i][i] == 0 {
			// Swap with a row below that has a nonzero pivot
			for k := i + 1; k < M; k++ {
				if augmented[k][i] != 0 {
					augmented[i], augmented[k] = augmented[k], augmented[i]
					break
				}
			}
		}

		// Ensure pivot is nonzero
		if augmented[i][i] == 0 {
			return nil
		}

		// Normalize pivot row
		pivotInv := invTable[augmented[i][i]]
		for j := 0; j < N+M; j++ {
			augmented[i][j] = gf16Mul(augmented[i][j], pivotInv)
		}

		// Eliminate other rows
		for k := 0; k < M; k++ {
			if k != i && augmented[k][i] != 0 {
				factor := augmented[k][i]
				for j := 0; j < N+M; j++ {
					augmented[k][j] = augmented[k][j] ^ gf16Mul(factor, augmented[i][j])
				}
			}
		}
	}

	// Extract the right inverse (n x m matrix)
	B := make([][]byte, N)
	for i := 0; i < N; i++ {
		B[i] = make([]byte, M)
		if i < M {
			copy(B[i], augmented[i][N:])
		}
	}

	return B
}

func ComputeLittleX(parties []*model.Party) {
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	basis := RandVector(t - s)

	for _, party := range parties {
		var z byte
		zVector := RandVector(t - s)

		for i := 0; i < t-s; i++ {
			z ^= gf16Mul(zVector[i], basis[i])
		}

		party.Z = z
	}

	triplesStep7 := GenerateMultiplicationTriple(len(parties), t, s, s, 1)
	triplesStep8 := GenerateMultiplicationTriple(len(parties), t, t, t, 1)
	// Compute [A^-1] * [b]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.S = RandMatrix(t, t)

		ai := triplesStep7.A[partyNumber]
		bi := triplesStep7.B[partyNumber]
		di := AddMatricesNew(party.AInverse, ai)
		ei := AddMatricesNew(vectorToMatrix(party.LittleY), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesB := multiplicationProtocol(parties, triplesStep7, dShares, eShares, t, s, s, 1)

	// Compute [S] * [z]
	dShares = make([][][]byte, len(parties))
	eShares = make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.R = RandMatrix(s, s)

		ai := triplesStep8.A[partyNumber]
		bi := triplesStep8.B[partyNumber]
		di := AddMatricesNew(party.S, ai)
		ei := AddMatricesNew(vectorToMatrix(party.Z), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	STimesZ := multiplicationProtocol(parties, triplesStep8, dShares, eShares, t, t, t, 1)

	// [x] = [A^-1] * [b] + [S] * [z]
	for i, _ := range parties {
		AddMatrices(STimesZ[i], ATimesB[i])
	}
}
