package mpc

import (
	"fmt"
	"mayo-threshold-go/rand"
	"reflect"
)

type Share struct {
	share, alpha, gamma byte
}

type MatrixShare struct {
	shares, gammas [][]byte
	alpha          byte
}

func createEmptyMatrixShare(rows, cols int) MatrixShare {
	return MatrixShare{
		shares: generateZeroMatrix(rows, cols),
		alpha:  0,
		gammas: generateZeroMatrix(rows, cols),
	}
}

var globalAlpha = byte(13)

func MatrixShareTranspose(m MatrixShare) MatrixShare {
	return MatrixShare{
		shares: MatrixTranspose(m.shares),
		alpha:  m.alpha,
		gammas: MatrixTranspose(m.gammas),
	}
}

func generateSharesForElement(n int, secret byte) []Share {
	shares := make([]byte, n)
	alphas := make([]byte, n)
	gammas := make([]byte, n)

	var sharesSum byte
	var alphaSum byte
	for i := 0; i < n-1; i++ {
		// shares of the secret
		share := rand.SampleFieldElement()
		shares[i] = share
		sharesSum ^= share

		// alpha
		alpha := rand.SampleFieldElement()
		alphas[i] = alpha
		alphaSum ^= alpha
	}
	shares[n-1] = secret ^ sharesSum
	alphas[n-1] = alphaSum ^ globalAlpha

	// Gamma
	alphaTimesSecret := field.Gf16Mul(globalAlpha, secret)
	var gammaSum byte
	for i := 0; i < n-1; i++ {
		gamma := rand.SampleFieldElement()
		gammas[i] = gamma
		gammaSum ^= gamma
	}
	gammas[n-1] = gammaSum ^ alphaTimesSecret

	result := make([]Share, n)
	for i := 0; i < n; i++ {
		result[i] = Share{
			share: shares[i],
			alpha: alphas[i],
			gamma: gammas[i],
		}
	}

	return result
}

func generateSharesForRandomElement(n int) []Share {
	secret := rand.SampleFieldElement()
	return generateSharesForElement(n, secret)
}

func createSharesForMatrix(n int, secretMatrix [][]byte) []MatrixShare {
	rows, cols := len(secretMatrix), len(secretMatrix[0])

	matrixShares := make([]MatrixShare, n)
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
			shareParts := generateSharesForElement(n, secretMatrix[i][j])

			for l := 0; l < n; l++ {
				matrixShares[l].shares[i][j] = shareParts[l].share
				matrixShares[l].alpha = shareParts[l].alpha
				matrixShares[l].gammas[i][j] = shareParts[l].gamma
			}
		}
	}

	return matrixShares
}

func createSharesForRandomMatrix(n, rows, cols int) []MatrixShare {
	secret := rand.Matrix(rows, cols)
	return createSharesForMatrix(n, secret)
}

func MultiplyMatrixWithConstant(a [][]byte, b byte) [][]byte {
	rows, cols := len(a), len(a[0])
	out := make([][]byte, rows)
	for i := range out {
		out[i] = make([]byte, cols)
		for j := 0; j < cols; j++ {
			out[i][j] = gf16Mul(a[i][j], b)
		}
	}
	return out
}

func AddMatrixShares(A, B MatrixShare) MatrixShare {
	var result MatrixShare
	result.shares = AddMatricesNew(A.shares, B.shares)
	result.gammas = AddMatricesNew(A.gammas, B.gammas)
	result.alpha = A.alpha
	return result
}

func AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare {
	var result MatrixShare
	if partyNumber == 0 {
		result.shares = AddMatricesNew(A, B.shares)
		result.gammas = AddMatricesNew(B.gammas, MultiplyMatrixWithConstant(A, B.alpha))
		result.alpha = B.alpha
	} else {
		result.shares = B.shares
		result.gammas = AddMatricesNew(B.gammas, MultiplyMatrixWithConstant(A, B.alpha))
		result.alpha = B.alpha
	}

	return result
}

func MulPublicLeft(A [][]byte, B MatrixShare) MatrixShare {
	var result MatrixShare
	result.shares = MultiplyMatrices(A, B.shares)
	result.gammas = MultiplyMatrices(A, B.gammas)
	result.alpha = B.alpha
	return result
}

func MulPublicRight(A MatrixShare, B [][]byte) MatrixShare {
	var result MatrixShare
	result.shares = MultiplyMatrices(A.shares, B)
	result.gammas = MultiplyMatrices(A.gammas, B)
	result.alpha = A.alpha
	return result
}

func openMatrix(shares []MatrixShare) ([][]byte, error) {
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

	muOpen := generateZeroMatrix(rows, cols)
	for _, share := range muShares {
		AddMatrices(muOpen, share)
	}

	if !reflect.DeepEqual(zero, muOpen) {
		return sPrime, fmt.Errorf("mu was not 0")
	}
	return sPrime, nil
}
