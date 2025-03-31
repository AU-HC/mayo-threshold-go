package mpc

import "mayo-threshold-go/model"

const m = 64
const k = 4
const n = 81
const o = 17
const v = n - o

const lambda = 2

func ComputeM(parties []*model.Party) {
	salt := Coin(parties, lambda)
	t := make([]byte, 0)                                                  // TODO: Call hash function
	triples := GenerateMultiplicationTriples(len(parties), k, v, v, o, m) // V: k x v, Li: v x o

	for _, party := range parties {
		V := RandMatrix(k, v)
		party.Salt = salt
		party.T = t
		party.V = V
	}

	for i := 0; i < m; i++ {
		dShares := make([][][]byte, len(parties))
		eShares := make([][][]byte, len(parties))

		// Compute locally
		for partyNumber, party := range parties {
			ai := triples[i].A[partyNumber]
			bi := triples[i].B[partyNumber]
			di := AddMatricesNew(party.V, ai)
			ei := AddMatricesNew(party.EskShare.L[i], bi)

			dShares[partyNumber] = di
			eShares[partyNumber] = ei
		}

		// Open d and e
		d := generateZeroMatrix(k, v)
		e := generateZeroMatrix(v, o)
		for j := range parties {
			AddMatrices(d, dShares[j])
			AddMatrices(e, eShares[j])
		}

		// Compute locally
		for partyNumber, party := range parties {
			M := make([][][]byte, m)
			for j := range party.EskShare.L {
				a := triples[j].A[partyNumber]
				b := triples[j].B[partyNumber]
				c := triples[j].C[partyNumber]

				db := MultiplyMatrices(d, b) // d * [b]
				de := MultiplyMatrices(d, e) // d * e
				ea := MultiplyMatrices(a, e) // e * [a]
				AddMatrices(db, ea)          // d * [b] + e * [a]
				AddMatrices(db, c)           // d * [b] + e * [a] + [c]
				AddMatrices(db, de)          // d * [b] + e * [a] + [c] + d * e

				M[i] = db // k x o
			}
			party.M = M
		}
	}
}

func ComputeY(parties []*model.Party) {

}

func LocalComputeA(parties []*model.Party) {

}

func LocalComputeY(parties []*model.Party) {

}

func ComputeX(parties []*model.Party) {

}

func ComputeS(parties []*model.Party) {

}

func ComputeSignature(parties []*model.Party) {

}

func ComputeT(parties []*model.Party) {
	/*for _, party := range parties {
		A := party.A
		s := len(A)
		t := len(A[0])

		R := RandMatrix(s, s)
		S := RandMatrix(t, t)
		AS := MultiplyMatrices(A, S)
		T := MultiplyMatrices(R, AS)
		TOpened := make([][]byte, len(T)) // TODO: Open correctly
		if Rank(TOpened) < s {
			panic(1) // TODO
		}
	}

	*/
}

func Rank(t [][]byte) int {
	return 0 // TODO: Implement rank of matrix
}

func ComputeAInverse(parties []*model.Party) {

}

func Computex(parties []*model.Party) {

}
