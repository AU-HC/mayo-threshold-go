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

type Party struct {
	EskShare ExpandedSecretKey
	Epk      ExpandedPublicKey
}

/*
	L := make([][][]byte, len(esk.L))
	O := make([][]byte, len(esk.O))

	// Generate the empty L
	for j := 0; j < len(esk.L); j++ {
		matrix := make([][]byte, len(esk.L[j]))
		for k := 0; k < len(esk.L[j]); k++ {
			row := make([]byte, len(esk.L[j][k]))
			for l := 0; l < len(esk.L[j][k]); l++ {
				row[l] = 0
			}
			matrix[k] = row
		}
		L[j] = matrix
	}

	// Generate the empty O
	for j := 0; j < len(esk.O); j++ {
		row := make([]byte, len(esk.O[j]))
		for k := 0; k < len(esk.O[j]); k++ {
			row[k] = 0
		}
		O[j] = row
	}

*/
