package model

type Party struct {
	EskShare          ExpandedSecretKey
	Epk               ExpandedPublicKey
	Salt              []byte
	T                 []byte
	V, VReconstructed [][]byte
	M                 [][][]byte
	Y                 [][][]byte
	A                 [][]byte
	R                 [][]byte
	LittleY           []byte
	Shares            map[string][][][]byte
}
