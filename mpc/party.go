package mpc

type Party struct {
	EskShare       ExpandedSecretKey
	Epk            ExpandedPublicKey
	Salt           []byte
	LittleT        []byte
	LittleX        MatrixShare
	LittleY        MatrixShare
	Z              MatrixShare
	T              [][]byte
	V              MatrixShare
	VReconstructed [][]byte
	A              MatrixShare
	AInverse       MatrixShare
	R              MatrixShare
	S              MatrixShare
	SPrime         MatrixShare
	X              MatrixShare
	Signature      MatrixShare
	M              []MatrixShare
	Y              []MatrixShare
}
