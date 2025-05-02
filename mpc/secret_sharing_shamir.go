package mpc

/*
import (
	"mayo-threshold-go/rand"
)

func generateCoefficients(secret byte, t int) []byte {
	coefficients := make([]byte, t)
	coefficients[0] = secret

	for i := 1; i < t; i++ {
		coefficients[i] = rand.SampleFieldElement()
	}

	return coefficients
}

func createShares(secret byte, n, t int) []Share {
	shareCoefficients := generateCoefficients(secret, t)
	alphaCoefficients := generateCoefficients(GlobalAlpha, t)
	gammaCoefficients := generateCoefficients(field.Gf16Mul(secret, GlobalAlpha), t)

	shares := make([]byte, n)
	alphaShares := make([]uint64, n)
	gammaShares := make([]uint64, n)

	for x := 1; x <= n; x++ {
		y := shareCoefficients[len(shareCoefficients)-1]
		for i := len(shareCoefficients) - 2; i >= 0; i-- {
			y = field.Gf16Mul(y, byte(x)) ^ shareCoefficients[i]
		}
		shares[x-1] = y
	}

	for x := 1; x <= n; x++ {
		y := alphaCoefficients[len(alphaCoefficients)-1]
		for i := len(alphaCoefficients) - 2; i >= 0; i-- {
			y = field.Gf16Mul(y, byte(x)) ^ alphaCoefficients[i]
		}
		alphaShares[x-1] = y
	}

	for x := 1; x <= n; x++ {
		y := gammaCoefficients[len(gammaCoefficients)-1]
		for i := len(gammaCoefficients) - 2; i >= 0; i-- {
			y = field.Gf16Mul(y, byte(x)) ^ gammaCoefficients[i]
		}
		gammaShares[x-1] = y
	}

	// Zip the shares into []Share
	result := make([]Share, n)
	for i := 0; i < n; i++ {
		result[i] = Share{
			share: shares[i],
			alpha: alphaShares[i],
			gamma: gammaShares[i],
		}
	}

	return result
}

func reconstructSecret(shares []byte, t int) byte {
	var secret byte = 0

	for j := 0; j < t; j++ {
		xj := byte(j + 1)
		lj := byte(1)

		for m := 0; m < t; m++ {
			if m == j {
				continue
			}
			xm := byte(m + 1)

			num := xm
			den := xj ^ xm // xj - xm

			// Lagrange coefficient: num / den
			frac := field.Gf16Mul(num, field.Gf16Inv(den))
			lj = field.Gf16Mul(lj, frac)
		}

		secret ^= field.Gf16Mul(shares[j], lj)
	}

	return secret
}
*/
