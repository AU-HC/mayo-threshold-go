package mpc

import "mayo-threshold-go/rand"

type SecretSharingAlgo interface {
	openMatrix(shares [][][]byte) [][]byte
	createSharesForMatrix([][]byte) [][][]byte
	createSharesForRandomMatrix(rows, cols int) [][][]byte
}

type Shamir struct {
	n, t int
}

func (s *Shamir) openMatrix(shares [][][]byte) [][]byte {
	rows := len(shares[0])
	cols := len(shares[0][0])

	secretMatrix := make([][]byte, rows)
	for row := 0; row < rows; row++ {
		secretMatrix[row] = make([]byte, cols)
		for col := 0; col < cols; col++ {
			// Gather the t shares for this matrix element
			elementShares := make([]byte, s.t)
			for partyNumber := 0; partyNumber < s.t; partyNumber++ {
				elementShares[partyNumber] = shares[partyNumber][row][col]
			}

			// Reconstruct the secret byte at position (row, col)
			secretMatrix[row][col] = reconstructSecret(elementShares, s.t)
		}
	}

	return secretMatrix
}

func (s *Shamir) createSharesForMatrix(secretMatrix [][]byte) [][][]byte {
	rows, cols := len(secretMatrix), len(secretMatrix[0])

	shares := make([][][]byte, s.n)
	for i := 0; i < s.n; i++ {
		shares[i] = generateZeroMatrix(rows, cols)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			secretByte := secretMatrix[r][c] //rand.SampleFieldElement()
			byteShares := createShares(secretByte, s.n, s.t)

			for partyNumber := 0; partyNumber < s.n; partyNumber++ {
				shares[partyNumber][r][c] = byteShares[partyNumber]
			}
		}
	}

	return shares
}

func (s *Shamir) createSharesForRandomMatrix(rows, cols int) [][][]byte {
	randomMatrix := rand.Matrix(rows, cols)
	return s.createSharesForMatrix(randomMatrix)
}

type Additive struct {
	n int
}

func (a *Additive) openMatrix(shares [][][]byte) [][]byte {
	rows, cols := len(shares[0]), len(shares[0][0])
	result := generateZeroMatrix(rows, cols)

	for _, share := range shares {
		AddMatrices(result, share)
	}

	return result
}

func (a *Additive) createSharesForMatrix(rows, cols int) [][][]byte {
	shares := make([][][]byte, a.n)
	for i := 0; i < a.n; i++ {
		shares[i] = rand.Matrix(rows, cols)
	}

	return shares
}
