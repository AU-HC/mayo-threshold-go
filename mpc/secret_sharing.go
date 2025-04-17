package mpc

import (
	"fmt"
	"mayo-threshold-go/rand"
	"reflect"
)

type SecretSharingAlgo interface {
	openMatrix(shares [][][]byte) [][]byte
	authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error)
	createSharesForMatrix([][]byte) []MatrixShare
	createSharesForRandomMatrix(rows, cols int) []MatrixShare
}

type Shamir struct {
	n, t int
}

func (s *Shamir) authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error) {
	parties, rows, cols := len(shares), len(shares[0].shares), len(shares[0].shares[0])

	zero := generateZeroMatrix(rows, cols)
	sPrime := generateZeroMatrix(rows, cols)
	for _, share := range shares {
		AddMatrices(sPrime, share.shares)
	}

	muShares := make([][][]byte, parties)
	for i, share := range shares {
		muShares[i] = AddMatricesNew(share.gammas, MultiplyMatrixWithConstant(sPrime, share.alpha))
	}
	muOpen := s.openMatrix(muShares)

	if !reflect.DeepEqual(zero, muOpen) {
		return sPrime, fmt.Errorf("mu was not 0")
	}
	return sPrime, nil
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

func (s *Shamir) createSharesForMatrix(secretMatrix [][]byte) []MatrixShare {
	/*rows, cols := len(secretMatrix), len(secretMatrix[0])

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

	*/
	return nil
}

func (s *Shamir) createSharesForRandomMatrix(rows, cols int) []MatrixShare {
	randomMatrix := rand.Matrix(rows, cols)
	return s.createSharesForMatrix(randomMatrix)
}

type Additive struct {
	n int
}

func (a *Additive) authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error) {
	parties, rows, cols := len(shares), len(shares[0].shares), len(shares[0].shares[0])

	zero := generateZeroMatrix(rows, cols)
	sPrime := generateZeroMatrix(rows, cols)
	for _, share := range shares {
		AddMatrices(sPrime, share.shares)
	}

	muShares := make([][][]byte, parties)
	for i, share := range shares {
		muShares[i] = AddMatricesNew(share.gammas, MultiplyMatrixWithConstant(sPrime, share.alpha))
	}
	muOpen := a.openMatrix(muShares)

	if !reflect.DeepEqual(zero, muOpen) {
		return sPrime, fmt.Errorf("mu was not 0")
	}
	return sPrime, nil
}

func (a *Additive) openMatrix(shares [][][]byte) [][]byte {
	rows, cols := len(shares[0]), len(shares[0][0])
	result := generateZeroMatrix(rows, cols)

	for _, share := range shares {
		AddMatrices(result, share)
	}

	return result
}

func (a *Additive) createSharesForMatrix(secretMatrix [][]byte) []MatrixShare {
	rows, cols := len(secretMatrix), len(secretMatrix[0])
	amountOfParties := a.n

	matrixShares := make([]MatrixShare, amountOfParties)
	for i := range matrixShares {
		matrixShares[i].shares = make([][]byte, rows)
		matrixShares[i].gammas = make([][]byte, rows)
		for r := 0; r < rows; r++ {
			matrixShares[i].shares[r] = make([]byte, cols)
			matrixShares[i].gammas[r] = make([]byte, cols)
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			shareParts := generateSharesForElement(amountOfParties, secretMatrix[i][j])

			for l := 0; l < amountOfParties; l++ {
				matrixShares[l].shares[i][j] = shareParts[l].share
				matrixShares[l].alpha = shareParts[l].alpha
				matrixShares[l].gammas[i][j] = shareParts[l].gamma
			}
		}
	}

	return matrixShares
}

func (a *Additive) createSharesForRandomMatrix(rows, cols int) []MatrixShare {
	secret := rand.Matrix(rows, cols)
	return a.createSharesForMatrix(secret)
}
