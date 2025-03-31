package mpc

import (
	"mayo-threshold-go/model"
	"reflect"
)

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

	VReconstructed := generateZeroMatrix(k, v)
	for _, party := range parties {
		V := RandMatrix(k, v)
		AddMatrices(VReconstructed, V)
		party.Salt = salt
		party.T = t
		party.V = V
		party.M = make([][][]byte, m)
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
			a := triples[i].A[partyNumber]
			b := triples[i].B[partyNumber]
			c := triples[i].C[partyNumber]

			db := MultiplyMatrices(d, b) // d * [b]
			de := MultiplyMatrices(d, e) // d * e
			ae := MultiplyMatrices(a, e) // [a] * e
			AddMatrices(db, ae)          // d * [b] + [a] * e
			AddMatrices(db, c)           // d * [b] + [a] * e + [c]

			if partyNumber == 0 {
				AddMatrices(db, de) // d * [b] + [a] * e + [c] + d * e
			}

			party.M[i] = db // k x o
		}

		MReconstructed := generateZeroMatrix(k, o)
		LReconstructed := generateZeroMatrix(v, o)
		for _, party := range parties {
			AddMatrices(MReconstructed, party.M[i])
			AddMatrices(LReconstructed, party.EskShare.L[i])
		}
		if !reflect.DeepEqual(MReconstructed, MultiplyMatrices(VReconstructed, LReconstructed)) {
			panic("M is not equal to V * L")
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
