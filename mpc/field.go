package mpc

type Field struct {
	mulTable          [][]byte
	invTable          []byte
	mulTableExtension [][]uint64
	invTableExtension []uint64
}

func InitField() *Field {
	mulTable, invTable := generateMulAndInvTable()

	return &Field{
		mulTable: mulTable,
		invTable: invTable,
	}
}

// Gf16Mul multiplies two elements in GF(16)
func (f *Field) Gf16Mul(a, b byte) byte {
	return f.mulTable[a][b]
}

// Gf16Inv calculates the inverse of an element in GF(16)
func (f *Field) Gf16Inv(a byte) byte {
	return f.invTable[a]
}

func (f *Field) Gf64Mul(a, b uint64) uint64 {
	return gf64Mul(a, b)
}

func generateMulAndInvTable() ([][]byte, []byte) {
	mulTable := make([][]byte, 16)
	invTable := make([]byte, 16)

	for i := 0; i < 16; i++ {
		mulTable[i] = make([]byte, 16)
		for j := 0; j < 16; j++ {
			mulTable[i][j] = gf16Mul(byte(i), byte(j))

			if mulTable[i][j] == 1 {
				invTable[i] = byte(j)
			}
		}
	}
	return mulTable, invTable
}

func reduceVecModF(y []byte) []byte {
	for i := m + shifts - 1; i >= m; i-- {
		for shift, coefficient := range tailF {
			y[i-m+shift] ^= field.Gf16Mul(y[i], coefficient)
		}
		y[i] = 0
	}

	y = y[:m]
	return y
}

func reduceAModF(A [][]byte) [][]byte {
	for row := m + shifts - 1; row >= m; row-- {
		for column := 0; column < k*o; column++ {
			for shift, coefficient := range tailF {
				A[row-m+shift][column] ^= field.Gf16Mul(A[row][column], coefficient)
			}
			A[row][column] = 0
		}
	}
	A = A[:m]
	return A
}

func gf64Mul(a, b uint64) uint64 {
	// https://en.wikipedia.org/wiki/Finite_field_arithmetic#Multiplication
	p := uint64(0)

	for {
		if a == 0 || b == 0 {
			break
		}

		// Extract the least significant bit of b
		lsb := b & uint64(1)
		if lsb == 1 {
			p = p ^ a
		}

		b = b >> 1

		// Keep track of whether the first bit of a is set to one
		msb := a & uint64(8000000000000000)
		carry := msb == uint64(8000000000000000)

		a = a << 1
		if carry {
			a = a ^ uint64(0xD)
		}
	}

	return p
}

func gf16Mul(a, b byte) byte {
	var r byte

	// Multiply each coefficient with y
	r = (a & 0x1) * b
	r ^= (a & 0x2) * b
	r ^= (a & 0x4) * b
	r ^= (a & 0x8) * b

	overFlowBits := r & 0xF0

	// Reduce with respect to x^4 + x + 1
	reducedOverFlowBits := overFlowBits>>4 ^ overFlowBits>>3

	// Subtract and remove overflow bits
	r = (r ^ reducedOverFlowBits) & 0x0F

	return r
}
