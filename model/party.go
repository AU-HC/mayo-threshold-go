package model

type Party struct {
	EskShare          ExpandedSecretKey
	Epk               ExpandedPublicKey
	Salt              []byte
	LittleT           []byte
	LittleX           []byte
	LittleY           []byte
	Z                 []byte
	T                 [][]byte
	V, VReconstructed [][]byte
	M                 [][][]byte
	Y                 [][][]byte
	A                 [][]byte
	AInverse          [][]byte
	R                 [][]byte
	S                 [][]byte
	SPrime            [][]byte
	X                 [][]byte
	Signature         [][]byte
}
