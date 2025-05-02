package mpc

import (
	"fmt"
	"mayo-threshold-go/rand"
	"reflect"
)

type SecretSharingAlgo interface {
	openMatrix(shares [][][]byte) [][]byte
	openMatrixExtension(shares [][][]uint64) [][]uint64
	authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error)
	createSharesForMatrix([][]byte) []MatrixShare
	createSharesForRandomMatrix(rows, cols int) []MatrixShare
	AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare
}

type Shamir struct {
	n, t int
}

/*
func (s *Shamir) AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare {
	var result MatrixShare
	result.shares = AddMatricesNew(A, B.shares)
	result.gammas = AddMatricesNew(B.gammas, MultiplyMatrixWithConstant(A, B.alpha))
	result.alpha = B.alpha
	return result
}

func (s *Shamir) authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error) {
	parties, rows, cols := len(shares), len(shares[0].shares), len(shares[0].shares[0])

	zero := generateZeroMatrix(rows, cols)
	sPrimeShares := make([][][]byte, parties)
	for i, share := range shares {
		sPrimeShares[i] = share.shares
	}
	sPrime := s.openMatrix(sPrimeShares)

	muShares := make([][][]byte, parties)
	for i, share := range shares {
		muShares[i] = AddMatricesNew(share.gammas, MultiplyMatrixWithConstant(sPrime, share.alpha))
	}

	err := commitAndVerify(muShares)
	if err != nil {
		return nil, err
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
	rows, cols := len(secretMatrix), len(secretMatrix[0])
	amountOfParties, threshold := s.n, s.t

	shares := make([]MatrixShare, s.n)
	for i := 0; i < s.n; i++ {
		shares[i] = createEmptyMatrixShare(rows, cols)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			secretByte := secretMatrix[r][c]
			byteShares := createShares(secretByte, amountOfParties, threshold)

			for l := 0; l < amountOfParties; l++ {
				shares[l].shares[r][c] = byteShares[l].share
				shares[l].alpha = byteShares[l].alpha
				shares[l].gammas[r][c] = byteShares[l].gamma
			}
		}
	}

	return shares
}

func (s *Shamir) createSharesForRandomMatrix(rows, cols int) []MatrixShare {
	randomMatrix := rand.Matrix(rows, cols)
	return s.createSharesForMatrix(randomMatrix)
}*/

type Additive struct {
	n int
}

func (a *Additive) AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare {
	var result MatrixShare
	if partyNumber == 0 {
		result.shares = AddMatricesNew(A, B.shares)
		result.gammas = AddMatricesNew(B.gammas, MultiplyMatrixWithConstantExtension(A, B.alpha))
		result.alpha = B.alpha
	} else {
		result.shares = B.shares
		result.gammas = AddMatricesNew(B.gammas, MultiplyMatrixWithConstantExtension(A, B.alpha))
		result.alpha = B.alpha
	}
	return result
}

func (a *Additive) authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error) {
	parties, rows, cols := len(shares), len(shares[0].shares), len(shares[0].shares[0])

	zero := generateZeroMatrix[uint64](rows, cols)
	sPrime := generateZeroMatrix[byte](rows, cols)
	for _, share := range shares {
		AddMatrices(sPrime, share.shares)
	}

	muShares := make([][][]uint64, parties)
	for i, share := range shares {
		muShares[i] = AddMatricesNew(share.gammas, MultiplyMatrixWithConstantExtension(sPrime, share.alpha))
	}

	/*err := commitAndVerify(muShares)
	if err != nil {
		return nil, err
	}*/

	muOpen := a.openMatrixExtension(muShares)

	if !reflect.DeepEqual(zero, muOpen) {
		return sPrime, fmt.Errorf("mu was not 0")
	}
	return sPrime, nil
}

func (a *Additive) openMatrix(shares [][][]byte) [][]byte {
	rows, cols := len(shares[0]), len(shares[0][0])
	result := generateZeroMatrix[byte](rows, cols)

	for _, share := range shares {
		AddMatrices(result, share)
	}

	return result
}

func (a *Additive) openMatrixExtension(shares [][][]uint64) [][]uint64 {
	rows, cols := len(shares[0]), len(shares[0][0])
	result := generateZeroMatrix[uint64](rows, cols)

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
		matrixShares[i].gammas = make([][]uint64, rows)
		for r := 0; r < rows; r++ {
			matrixShares[i].shares[r] = make([]byte, cols)
			matrixShares[i].gammas[r] = make([]uint64, cols)
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

/*func commitAndVerify(shares [][][]uint64) error {
	parties := len(shares)

	// Create commitments
	commitmentRandomness := make([][]uint64, parties)
	commitments := make([][]uint64, parties)
	for p := 0; p < parties; p++ {
		randomVector := make([]uint64, 32)

		for i := 0; i < len(randomVector); i++ {
			randomVector[i] = uint64(random.Int())
		}

		commitmentRandomness[p] = randomVector
		commitments[p] = Commit(shares[p], commitmentRandomness[p])
	}

	// Check commitments
	for p := 0; p < parties; p++ {
		isValid := VerifyCommitment(shares[p], commitmentRandomness[p], commitments[p])
		if !isValid {
			return fmt.Errorf("commitment verification failed")
		}
	}

	return nil
}*/
