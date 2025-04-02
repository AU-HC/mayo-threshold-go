package model

type Signature struct {
	S    [][]byte
	Salt []byte
}

func (sig *Signature) Encode() []byte {
	encodedSig := make([]byte, 0)

	for _, row := range sig.S {
		encodedVec := encodeVec(row)
		encodedSig = append(encodedSig, encodedVec...)
	}

	return append(encodedSig, sig.Salt...)
}

func encodeVec(bytes []byte) []byte {
	encoded := make([]byte, (len(bytes)+1)/2)

	for i := 0; i < len(bytes)-1; i += 2 {
		encoded[i/2] = bytes[i+1]<<4 | bytes[i]&0xf
	}

	if (len(bytes) % 2) == 1 {
		encoded[(len(bytes)-1)/2] = bytes[len(bytes)-1] //<< 4
	}

	return encoded
}
