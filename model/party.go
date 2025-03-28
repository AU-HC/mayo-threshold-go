package model

import "mayo-threshold-go/mock"

type Party struct {
	EskShare mock.ExpandedSecretKey
	Epk      mock.ExpandedPublicKey
	V        [][]byte
	M        [][][]byte
	Y        [][][]byte
	A        [][]byte
	Shares   map[string][][][]byte
}
