package mock

type ExpandedSecretKey struct {
	P1 [][][]byte
	L  [][][]byte
	O  [][]byte
}

type ExpandedPublicKey struct {
	P1 [][][]byte
	P2 [][][]byte
	P3 [][][]byte
}
