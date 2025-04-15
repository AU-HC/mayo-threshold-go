package mpc

type Triple struct {
	A, B, C [][][]byte
}

type ActiveTriple struct {
	A, B, C []MatrixShare
}
