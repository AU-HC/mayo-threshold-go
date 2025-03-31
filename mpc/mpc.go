package mpc

import "mayo-threshold-go/model"

const m = 64
const k = 4
const n = 0
const o = 0
const v = n - o
const lambda = 2

func ComputeM(parties []*model.Party) {
	salt := Coin(parties, lambda)
	t := make([]byte, 0) // TODO: Call hash function

	for _, party := range parties {
		V := RandMatrix(k, v)
		M := make([][][]byte, m)

		for i, Li := range party.EskShare.L {
			M[i] = MultiplyMatrices(V, Li)
		}

		party.Salt = salt
		party.T = t
		party.M = M
		party.V = V
	}
}

func ComputeY(parties []*model.Party) {
	for _, party := range parties {
		V := party.V
		Y := make([][][]byte, m)

		for i, P1i := range party.Epk.P1 {
			Y[i] = MultiplyMatrices(MultiplyMatrixWithConstantMatrix(V, P1i), MatrixTranspose(V))
		}

		party.Y = Y
	}
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
	for _, party := range parties {
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
}

func Rank(t [][]byte) int {
	return 0 // TODO: Implement rank of matrix
}

func ComputeAInverse(parties []*model.Party) {

}

func Computex(parties []*model.Party) {

}
