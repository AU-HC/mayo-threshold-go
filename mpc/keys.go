package mpc

type ExpandedSecretKey struct {
	P1 [][][]byte
	L  []MatrixShare
	O  MatrixShare
}

type ExpandedPublicKey struct {
	P1 [][][]byte
	P2 [][][]byte
	P3 [][][]byte
}
