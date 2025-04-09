package mpc

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

func createShares(secret byte, n, t int) []byte {
	coefficients := generateCoefficients(secret, t)
	shares := make([]byte, n)

	for x := 1; x <= n; x++ {
		y := coefficients[len(coefficients)-1]
		for i := len(coefficients) - 2; i >= 0; i-- {
			y = field.Gf16Mul(y, byte(x)) ^ coefficients[i]
		}
		shares[x-1] = y
	}

	return shares
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
