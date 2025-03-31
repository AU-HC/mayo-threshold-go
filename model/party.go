package model

type Party struct {
	EskShare ExpandedSecretKey
	Epk      ExpandedPublicKey
	Salt     []byte
	T        []byte
	V        [][]byte
	M        [][][]byte
	Y        [][][]byte
	A        [][]byte
	Shares   map[string][][][]byte
}
