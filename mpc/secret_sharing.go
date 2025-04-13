package mpc

import "mayo-threshold-go/rand"

type SecretSharingAlgo interface {
	openMatrix(shares [][][]byte) [][]byte
	authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error)
	createSharesForMatrix([][]byte) [][][]byte
	createSharesForRandomMatrix(rows, cols int) [][][]byte
	shouldPartyAddConstantShare(partyNumber int) bool
}

type Shamir struct {
	n, t int
}

func (s *Shamir) authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error) {
	//TODO implement me
	panic("implement me")
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
			secretByte := secretMatrix[r][c]
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

func (s *Shamir) shouldPartyAddConstantShare(partyNumber int) bool {
	return true
}

type Additive struct {
	n int
}

func (a *Additive) authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error) {
	return openMatrix(shares)
}

func (a *Additive) openMatrix(shares [][][]byte) [][]byte {
	rows, cols := len(shares[0]), len(shares[0][0])
	result := generateZeroMatrix(rows, cols)

	for _, share := range shares {
		AddMatrices(result, share)
	}

	return result
}

func (a *Additive) createSharesForMatrix(secretMatrix [][]byte) [][][]byte {
	rows, cols := len(secretMatrix), len(secretMatrix[0])
	shares := make([][][]byte, a.n)
	sharesSum := generateZeroMatrix(rows, cols)

	for i := 0; i < a.n-1; i++ { // sample shares for n-1 parties
		share := rand.Matrix(rows, cols)
		shares[i] = share
		AddMatrices(sharesSum, share)
	}
	shares[a.n-1] = AddMatricesNew(sharesSum, secretMatrix)

	return shares
}

func (a *Additive) createSharesForRandomMatrix(rows, cols int) [][][]byte {
	shares := make([][][]byte, a.n)
	for i := 0; i < a.n; i++ {
		shares[i] = rand.Matrix(rows, cols)
	}

	return shares
}

func (a *Additive) shouldPartyAddConstantShare(partyNumber int) bool {
	return partyNumber == 0
}
