package mpc

import (
	"bytes"
	"mayo-threshold-go/model"
)

func Verify(parties []*model.Party, signature model.Signature) bool {
	P := make([][][][]byte, len(parties))

	for partyNumber, party := range parties {
		P[partyNumber] = calculateP(party.Epk.P1, party.Epk.P2, party.Epk.P3)
	}

	for i := 0; i < m; i++ {
		triple := GenerateMultiplicationTriple(len(parties), k, n, n, k)
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))

		for partyNumber, party := range parties {
			STimesP := MultiplyMatrices(party.Signature, P[partyNumber][i])

			ai := triple.A[partyNumber]
			bi := triple.B[partyNumber]
			di := AddMatricesNew(STimesP, ai)
			ei := AddMatricesNew(MatrixTranspose(party.Signature), bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		YShares := multiplicationProtocol(parties, triple, dShares, eShares, k, n, n, k)

		for partyNumber, party := range parties {
			party.Y[i] = YShares[partyNumber]
		}
	}

	LocalComputeY(parties)

	y := make([]byte, m)
	for _, party := range parties {
		y = AddVec(y, party.LittleY)
	}
	zero := make([]byte, m)

	return bytes.Equal(y, zero)
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
