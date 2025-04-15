package mpc

type Party struct {
	EskShare       ExpandedSecretKey
	Epk            ExpandedPublicKey
	Salt           []byte
	LittleT        []byte
	LittleX        []byte
	LittleY        []byte
	Z              []byte
	T              [][]byte
	V              MatrixShare
	VReconstructed [][]byte
	A              MatrixShare
	AInverse       MatrixShare
	R              MatrixShare
	S              MatrixShare
	SPrime         MatrixShare
	X              MatrixShare
	Signature      [][]byte
	M              []MatrixShare
	Y              []MatrixShare
}
