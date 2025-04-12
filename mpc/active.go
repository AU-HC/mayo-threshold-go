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
	shares, alphas, gammas [][]byte
}

func Test() {
	parties := 3

	// The dealer shares a value
	secret := byte(2)
	shares := generateSharesForElement(parties, secret)

	// Opening the shares
	var sPrime byte
	for _, share := range shares {
		sPrime ^= share.share
	}

	myShares := make([]byte, parties)
	for i, share := range shares {
		myShares[i] = share.gamma ^ gf16Mul(sPrime, share.alpha)
	}

	var myOpen byte
	for _, share := range myShares {
		myOpen ^= share
	}

	fmt.Println(fmt.Sprintf("secret: %d, my: %d", sPrime, myOpen))
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

	alpha := rand.SampleFieldElement()
	alphas[n-1] = alpha
	alphaSum ^= alpha

	// Gamma
	alphaTimesSecret := field.Gf16Mul(alphaSum, secret)
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
		matrixShares[i].alphas = make([][]byte, rows)
		matrixShares[i].gammas = make([][]byte, rows)
		for r := 0; r < rows; r++ {
			matrixShares[i].shares[r] = make([]byte, cols)
			matrixShares[i].alphas[r] = make([]byte, cols)
			matrixShares[i].gammas[r] = make([]byte, cols)
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			shareParts := generateSharesForElement(n, secretMatrix[i][j])

			for l := 0; l < n; l++ {
				matrixShares[l].shares[i][j] = shareParts[l].share
				matrixShares[l].alphas[i][j] = shareParts[l].alpha
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

func MultiplyMatricesElementWise(a, b [][]byte) [][]byte {
	rows, cols := len(a), len(a[0])
	out := make([][]byte, rows)
	for i := range out {
		out[i] = make([]byte, cols)
		for j := 0; j < cols; j++ {
			out[i][j] = gf16Mul(a[i][j], b[i][j])
		}
	}
	return out
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
		muShares[i] = AddMatricesNew(share.gammas, MultiplyMatricesElementWise(sPrime, share.alphas))
	}

	muOpen := generateZeroMatrix(rows, cols)
	for _, share := range muShares {
		AddMatrices(muOpen, share)
	}

	if !reflect.DeepEqual(zero, muOpen) {
		fmt.Println("mu was not 0")
	}
	fmt.Println(fmt.Sprintf("secret: %d, my: %d", sPrime, muOpen))

	return sPrime, nil // TODO: check if mu is 0
}
