package mpc

import (
	"fmt"
	"mayo-threshold-go/rand"
	"reflect"
)

// TODO: This should not be a single global alpha
const GlobalAlpha = byte(13)

type Share struct {
	share        byte
	alpha, gamma []byte
}

type MatrixShare struct {
	shares [][]byte
	gammas [][][]byte
	alpha  []byte
}

func createEmptyMatrixShare(rows, cols int) MatrixShare {
	return MatrixShare{
		alpha:  make([]byte, macAmount),
		shares: generateZeroMatrix(rows, cols),
		gammas: generateZeroMatrices(macAmount, rows, cols),
	}
}

func MatrixShareTranspose(m MatrixShare) MatrixShare {
	return MatrixShare{
		alpha:  m.alpha,
		shares: MatrixTranspose(m.shares),
		gammas: MatricesTranspose(m.gammas),
	}
}

func generateSharesForElement(n int, secret byte) []Share {
	shares := make([]byte, n)
	alphas := make([][]byte, n)
	gammas := make([][]byte, n)
	alphaSum := make([]byte, macAmount)

	var sharesSum byte

	for i := 0; i < n-1; i++ {
		alphas[i] = make([]byte, macAmount)
		gammas[i] = make([]byte, macAmount)

		// shares of the secret
		share := rand.SampleFieldElement()
		shares[i] = share
		sharesSum ^= share

		// alpha
		for j := 0; j < macAmount; j++ {
			alpha := rand.SampleFieldElement()
			alphas[i][j] = alpha
			alphaSum[j] ^= alpha
		}
	}
	shares[n-1] = secret ^ sharesSum
	alphas[n-1] = make([]byte, macAmount)
	gammas[n-1] = make([]byte, macAmount)

	for i := 0; i < macAmount; i++ {
		alphas[n-1][i] = alphaSum[i] ^ GlobalAlpha
	}

	// Gamma
	alphaTimesSecret := field.Gf16Mul(GlobalAlpha, secret)
	gammaSum := make([]byte, macAmount)
	for i := 0; i < n-1; i++ {
		for j := 0; j < macAmount; j++ {
			gamma := rand.SampleFieldElement()
			gammas[i][j] = gamma
			gammaSum[j] ^= gamma
		}
	}

	for i := 0; i < macAmount; i++ {
		gammas[n-1][i] = gammaSum[i] ^ alphaTimesSecret
	}

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

func createSharesForMatrix(n int, secretMatrix [][]byte) []MatrixShare {
	rows, cols := len(secretMatrix), len(secretMatrix[0])

	matrixShares := make([]MatrixShare, n)
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
			shareParts := generateSharesForElement(n, secretMatrix[i][j])

			for l := 0; l < n; l++ {
				matrixShares[l].shares[i][j] = shareParts[l].share
				matrixShares[l].alpha = shareParts[l].alpha
				matrixShares[l].gammas[:][i][j] = shareParts[l].gamma
			}
		}
	}

	return matrixShares
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
	result.gammas = make([][][]byte, macAmount)
	result.shares = AddMatricesNew(A.shares, B.shares)

	for i := 0; i < macAmount; i++ {
		result.gammas[i] = AddMatricesNew(A.gammas[i], B.gammas[i])
	}

	result.alpha = A.alpha
	return result
}

func AddPublicLeft(A [][]byte, B MatrixShare, partyNumber int) MatrixShare {
	var result MatrixShare
	if partyNumber == 0 {
		result.shares = AddMatricesNew(A, B.shares)
		for i := 0; i < macAmount; i++ {
			result.gammas[i] = AddMatricesNew(B.gammas[i], MultiplyMatrixWithConstant(A, B.alpha[i]))
		}
		result.alpha = B.alpha
	} else {
		result.shares = B.shares
		for i := 0; i < macAmount; i++ {
			result.gammas[i] = AddMatricesNew(B.gammas[i], MultiplyMatrixWithConstant(A, B.alpha[i]))
		}
		result.alpha = B.alpha
	}
	return result
}

func MulPublicLeft(A [][]byte, B MatrixShare) MatrixShare {
	var result MatrixShare
	result.gammas = make([][][]byte, macAmount)
	result.shares = MultiplyMatrices(A, B.shares)

	for i := 0; i < macAmount; i++ {
		result.gammas[i] = MultiplyMatrices(A, B.gammas[i])
	}

	result.alpha = B.alpha
	return result
}

func MulPublicRight(A MatrixShare, B [][]byte) MatrixShare {
	var result MatrixShare
	result.gammas = make([][][]byte, macAmount)
	result.shares = MultiplyMatrices(A.shares, B)

	for i := 0; i < macAmount; i++ {
		result.gammas[i] = MultiplyMatrices(A.gammas[i], B)
	}

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

	for k := 0; k < macAmount; k++ {
		muShares := make([][][]byte, parties)
		for i, share := range shares {
			muShares[i] = AddMatricesNew(share.gammas[k], MultiplyMatrixWithConstant(sPrime, share.alpha[k]))
		}

		muOpen := generateZeroMatrix(rows, cols)
		for _, share := range muShares {
			AddMatrices(muOpen, share)
		}

		if !reflect.DeepEqual(zero, muOpen) {
			return sPrime, fmt.Errorf("mu was not 0")
		}
	}

	return sPrime, nil
}
