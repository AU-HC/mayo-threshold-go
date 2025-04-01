package model

type Party struct {
	EskShare          ExpandedSecretKey
	Epk               ExpandedPublicKey
	Salt              []byte
	T                 [][]byte
	V, VReconstructed [][]byte
	M                 [][][]byte
	Y                 [][][]byte
	A                 [][]byte
	AInverse          [][]byte
	R                 [][]byte
	S                 [][]byte
	X                 []byte
	LittleY           []byte
	Z                 byte
	Shares            map[string][][][]byte
}
