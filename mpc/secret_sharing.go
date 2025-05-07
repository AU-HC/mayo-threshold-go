package mpc

import (
	"fmt"
	random "math/rand"
	"mayo-threshold-go/rand"
	"reflect"
)

type SecretSharingAlgo interface {
	openMatrix(shares [][][]byte) [][]byte
	authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error)
	createSharesForMatrix(matrix [][]byte) []MatrixShare
	createSharesForRandomMatrix(rows, cols int) []MatrixShare
	AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare
}

type Shamir struct {
	n, t        int
	alphaShares [][]byte
}

func CreateShamir(n, t int) *Shamir {
	return &Shamir{
		n:           n,
		t:           t,
		alphaShares: generateAlphaSharesShamir(n, t),
	}
}

func (s *Shamir) AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare {
	var result MatrixShare
	result.gammas = make([][][]byte, macAmount)
	result.shares = AddMatricesNew(A, B.shares)

	for i := 0; i < macAmount; i++ {
		result.gammas[i] = AddMatricesNew(B.gammas[i], MultiplyMatrixWithConstant(A, s.alphaShares[partyNumber][i]))
	}

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

	for k := 0; k < macAmount; k++ {
		muShares := make([][][]byte, parties)
		for i, share := range shares {
			muShares[i] = AddMatricesNew(share.gammas[k], MultiplyMatrixWithConstant(sPrime, s.alphaShares[i][k]))
		}

		err := commitAndVerify(muShares)
		if err != nil {
			return nil, err
		}

		muOpen := s.openMatrix(muShares)

		if !reflect.DeepEqual(zero, muOpen) {
			return sPrime, fmt.Errorf("mu was not 0")
		}
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
				for i := 0; i < macAmount; i++ {
					shares[l].gammas[i][r][c] = byteShares[l].gamma[i]
				}
			}
		}
	}

	return shares
}

func (s *Shamir) createSharesForRandomMatrix(rows, cols int) []MatrixShare {
	randomMatrix := rand.Matrix(rows, cols)
	return s.createSharesForMatrix(randomMatrix)
}

type Additive struct {
	n           int
	alphaShares [][]byte
}

func CreateAdditive(n int) *Additive {
	return &Additive{
		n:           n,
		alphaShares: generateAlphaSharesAdditive(n),
	}
}

func (a *Additive) AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare {
	var result MatrixShare
	result.gammas = make([][][]byte, macAmount)

	if partyNumber == 0 {
		result.shares = AddMatricesNew(A, B.shares)
	} else {
		result.shares = B.shares
	}

	for i := 0; i < macAmount; i++ {
		result.gammas[i] = AddMatricesNew(B.gammas[i], MultiplyMatrixWithConstant(A, a.alphaShares[partyNumber][i]))
	}

	return result
}

func (a *Additive) authenticatedOpenMatrix(shares []MatrixShare) ([][]byte, error) {
	parties, rows, cols := len(shares), len(shares[0].shares), len(shares[0].shares[0])

	zero := generateZeroMatrix(rows, cols)
	sPrime := generateZeroMatrix(rows, cols)
	for _, share := range shares {
		AddMatrices(sPrime, share.shares)
	}

	for k := 0; k < macAmount; k++ {
		muShares := make([][][]byte, parties)
		for i, share := range shares {
			muShares[i] = AddMatricesNew(share.gammas[k], MultiplyMatrixWithConstant(sPrime, a.alphaShares[i][k]))
		}

		err := commitAndVerify(muShares)
		if err != nil {
			return nil, err
		}

		muOpen := a.openMatrix(muShares)

		if !reflect.DeepEqual(zero, muOpen) {
			return sPrime, fmt.Errorf("mu was not 0")
		}
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
		matrixShares[i].gammas = make([][][]byte, macAmount)
		for r := 0; r < rows; r++ {
			matrixShares[i].shares[r] = make([]byte, cols)
		}

		for k := 0; k < macAmount; k++ {
			matrixShares[i].gammas[k] = make([][]byte, rows)
			for r := 0; r < rows; r++ {
				matrixShares[i].gammas[k][r] = make([]byte, cols)
			}
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			shareParts := a.generateSharesForElement(secretMatrix[i][j])

			for l := 0; l < amountOfParties; l++ {
				for k := 0; k < macAmount; k++ {
					matrixShares[l].shares[i][j] = shareParts[l].share
					matrixShares[l].gammas[k][i][j] = shareParts[l].gamma[k]
				}
			}
		}
	}

	return matrixShares
}

func (a *Additive) generateSharesForElement(secret byte) []Share {
	amountOfParties := len(a.alphaShares)

	shares := make([]byte, amountOfParties)
	gammas := make([][]byte, amountOfParties)
	var sharesSum byte

	for i := 0; i < amountOfParties-1; i++ {
		gammas[i] = make([]byte, macAmount)

		// shares of the secret
		share := rand.SampleFieldElement()
		shares[i] = share
		sharesSum ^= share
	}
	shares[amountOfParties-1] = secret ^ sharesSum
	gammas[amountOfParties-1] = make([]byte, macAmount)

	// Gamma
	gammaSum := make([]byte, macAmount)
	for i := 0; i < amountOfParties-1; i++ {
		for j := 0; j < macAmount; j++ {
			gamma := rand.SampleFieldElement()
			gammas[i][j] = gamma
			gammaSum[j] ^= gamma
		}
	}

	for i := 0; i < macAmount; i++ {
		alphaTimesSecret := field.Gf16Mul(GlobalAlphas[i], secret)
		gammas[amountOfParties-1][i] = gammaSum[i] ^ alphaTimesSecret
	}

	result := make([]Share, amountOfParties)
	for i := 0; i < amountOfParties; i++ {
		result[i] = Share{
			share: shares[i],
			gamma: gammas[i],
		}
	}

	return result
}

func (a *Additive) createSharesForRandomMatrix(rows, cols int) []MatrixShare {
	secret := rand.Matrix(rows, cols)
	return a.createSharesForMatrix(secret)
}

func commitAndVerify(shares [][][]byte) error {
	parties := len(shares)

	// Create commitments
	commitmentRandomness := make([][]byte, parties)
	commitments := make([][]byte, parties)
	for p := 0; p < parties; p++ {
		randomVector := make([]byte, 32)

		for i := 0; i < len(randomVector); i++ {
			randomVector[i] = byte(random.Int())
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
}
