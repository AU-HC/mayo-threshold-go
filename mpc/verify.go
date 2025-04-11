package mpc

import (
	"bytes"
	"mayo-threshold-go/model"
	"mayo-threshold-go/rand"
)

const Tau = 1

// ThresholdVerify takes an 'secret shared' signature and checks if it is valid for the message under the public key
func ThresholdVerify(parties []*model.Party, signature model.ThresholdSignature) bool {
	p := parties[0]
	P := calculateP(p.Epk.P1, p.Epk.P2, p.Epk.P3)

	for i := 0; i < m; i++ {
		triple := GenerateMultiplicationTriple(k, n, n, k)
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))

		for partyNumber, _ := range parties {
			STimesP := MultiplyMatrices(signature.S[partyNumber], P[i])

			ai := triple.A[partyNumber]
			bi := triple.B[partyNumber]
			di := AddMatricesNew(STimesP, ai)
			ei := AddMatricesNew(MatrixTranspose(signature.S[partyNumber]), bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		YShares := multiplicationProtocol(parties, triple, dShares, eShares)

		for partyNumber, party := range parties {
			party.Y[i] = YShares[partyNumber]
		}
	}

	localComputeY(parties)
	zShares := make([][]byte, len(parties))
	for i, party := range parties {
		zShares[i] = party.LittleY
	}

	alphaValues := rand.CoinMatrix(parties, m, Tau)
	w := make([][]byte, m)

	for i := 0; i < Tau; i++ {
		for j := 0; j < m; j++ {
			w[i] = matrixToVec(AddMatricesNew(vectorToMatrix(w[i]), vectorToMatrix(MultiplyVecConstant(alphaValues[j][i], zShares[j]))))
		}
	}

	z := make([]byte, m)
	for _, party := range parties {
		z = AddVec(z, party.LittleY)
	}
	zero := make([]byte, m)

	return bytes.Equal(z, zero)
}

// Verify takes an 'opened' signature and checks if it is valid for the message under the public key
func Verify(epk model.ExpandedPublicKey, message []byte, signature model.Signature) bool {
	P := calculateP(epk.P1, epk.P2, epk.P3)
	Y := make([][][]byte, m)
	t := rand.Shake256(m, message, signature.Salt)

	for i := 0; i < m; i++ {
		STimesP := MultiplyMatrices(signature.S, P[i])
		Y[i] = MultiplyMatrices(STimesP, MatrixTranspose(signature.S))
	}

	// Create party, due to how code is structured
	parties := make([]*model.Party, 1)
	parties[0] = &model.Party{Y: Y, LittleT: t}
	localComputeY(parties)

	zero := make([]byte, m)
	return bytes.Equal(parties[0].LittleY, zero)
}

func calculateP(P1, P2, P3 [][][]byte) [][][]byte {
	P := make([][][]byte, m)
	for i := 0; i < m; i++ {
		P[i] = make([][]byte, n)
		for j := 0; j < n; j++ {
			P[i][j] = make([]byte, n)
		}
	}

	for i := 0; i < m; i++ {
		// Set P1
		for row := 0; row < v; row++ {
			for column := 0; column < v; column++ {
				P[i][row][column] = P1[i][row][column]
			}
		}
		// Set P2
		for row := 0; row < v; row++ {
			for column := 0; column < o; column++ {
				P[i][row][column+v] = P2[i][row][column]
			}
		}
		// Set P3
		for row := 0; row < o; row++ {
			for column := 0; column < o; column++ {
				P[i][row+v][column+v] = P3[i][row][column]
			}
		}
	}

	return P
}
