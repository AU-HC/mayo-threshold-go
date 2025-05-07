package mpc

import (
	"mayo-threshold-go/rand"
)

var GlobalAlphas = [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

type Share struct {
	share byte
	gamma []byte
}

type MatrixShare struct {
	shares [][]byte
	gammas [][][]byte
}

func createEmptyMatrixShare(rows, cols int) MatrixShare {
	return MatrixShare{
		shares: generateZeroMatrix(rows, cols),
		gammas: generateZeroMatrices(macAmount, rows, cols),
	}
}

func MatrixShareTranspose(m MatrixShare) MatrixShare {
	return MatrixShare{
		shares: MatrixTranspose(m.shares),
		gammas: MatricesTranspose(m.gammas),
	}
}

func generateAlphaSharesAdditive(n int) [][]byte {
	alphaShares := make([][]byte, n)
	alphaSum := make([]byte, macAmount)

	for i := 0; i < n-1; i++ {
		alphaShares[i] = make([]byte, macAmount)

		for j := 0; j < macAmount; j++ {
			alphaShare := rand.SampleFieldElement()
			alphaShares[i][j] = alphaShare
			alphaSum[j] ^= alphaShare
		}
	}

	alphaShares[n-1] = make([]byte, macAmount)
	for i := 0; i < macAmount; i++ {
		alphaShares[n-1][i] = alphaSum[i] ^ GlobalAlphas[i]
	}

	return alphaShares
}

func generateAlphaSharesShamir(n, t int) [][]byte {
	alphaShares := make([][]byte, n)

	for i := 0; i < n; i++ {
		alphaShares[i] = make([]byte, macAmount)
	}

	for i, alpha := range GlobalAlphas {
		shareCoefficients := generateCoefficients(alpha, t)

		for x := 1; x <= n; x++ {
			y := shareCoefficients[len(shareCoefficients)-1]
			for i := len(shareCoefficients) - 2; i >= 0; i-- {
				y = field.Gf16Mul(y, byte(x)) ^ shareCoefficients[i]
			}
			alphaShares[x-1][i] = y
		}
	}

	return alphaShares
}

func MultiplyMatrixWithConstant(a [][]byte, b byte) [][]byte {
	rows, cols := len(a), len(a[0])
	out := make([][]byte, rows)
	for i := range out {
		out[i] = make([]byte, cols)
		for j := 0; j < cols; j++ {
			out[i][j] = field.Gf16Mul(a[i][j], b)
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

	return result
}

func MulPublicLeft(A [][]byte, B MatrixShare) MatrixShare {
	var result MatrixShare
	result.gammas = make([][][]byte, macAmount)
	result.shares = MultiplyMatrices(A, B.shares)

	for i := 0; i < macAmount; i++ {
		result.gammas[i] = MultiplyMatrices(A, B.gammas[i])
	}

	return result
}

func MulPublicRight(A MatrixShare, B [][]byte) MatrixShare {
	var result MatrixShare
	result.gammas = make([][][]byte, macAmount)
	result.shares = MultiplyMatrices(A.shares, B)

	for i := 0; i < macAmount; i++ {
		result.gammas[i] = MultiplyMatrices(A.gammas[i], B)
	}

	return result
}
