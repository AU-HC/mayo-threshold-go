package mpc

type Party struct {
	EskShare       ExpandedSecretKey
	Epk            ExpandedPublicKey
	Salt           []byte
	LittleT        []byte
	T              [][]byte
	VReconstructed [][]byte
	LittleX        MatrixShare
	LittleY        MatrixShare
	Z              MatrixShare
	V              MatrixShare
	A              MatrixShare
	AInverse       MatrixShare
	R              MatrixShare
	S              MatrixShare
	SPrime         MatrixShare
	X              MatrixShare
	M              []MatrixShare
	Y              []MatrixShare
}
