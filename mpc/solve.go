package mpc

import (
	"fmt"
	"mayo-threshold-go/model"
	"slices"
)

func ComputeT(parties []*model.Party) bool {
	s := len(parties[0].A)
	t := len(parties[0].A[0])
	triplesStep2 := GenerateMultiplicationTriple(len(parties), s, t, t, t)
	triplesStep3 := GenerateMultiplicationTriple(len(parties), s, s, s, t)

	// Compute [A * S] = [A] * [S]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.S = RandMatrix(t, t)

		ai := triplesStep2.A[partyNumber]
		bi := triplesStep2.B[partyNumber]
		di := AddMatricesNew(party.A, ai)
		ei := AddMatricesNew(party.S, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesSShares := multiplicationProtocol(parties, triplesStep2, dShares, eShares, s, t, t, t)

	// Compute [T] = [R] * [A * S]
	dShares = make([][][]byte, len(parties))
	eShares = make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.R = RandMatrix(s, s)

		ai := triplesStep3.A[partyNumber]
		bi := triplesStep3.B[partyNumber]
		di := AddMatricesNew(party.R, ai)
		ei := AddMatricesNew(ATimesSShares[partyNumber], bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	TShares := multiplicationProtocol(parties, triplesStep3, dShares, eShares, s, s, s, t)

	// Open T and check rank
	T := generateZeroMatrix(s, t)
	for _, tShare := range TShares {
		AddMatrices(T, tShare)
	}

	for _, party := range parties {
		party.T = T
	}

	copyOfT := generateZeroMatrix(s, t)
	for row := 0; row < len(T); row++ {
		copy(copyOfT[row][:], T[row][:])
	}

	return rankOfMatrix(copyOfT) < s
}

func rankOfMatrix(t [][]byte) int {
	if len(t) == 0 || len(t[0]) == 0 {
		return 0
	}

	rows, cols := len(t), len(t[0])
	rank := 0

	for col := 0; col < cols; col++ {
		pivotRow := -1
		for row := rank; row < rows; row++ {
			if t[row][col] != 0 {
				pivotRow = row
				break
			}
		}

		if pivotRow == -1 {
			continue
		}

		t[pivotRow], t[rank] = t[rank], t[pivotRow]

		pivot := t[rank][col]
		for c := col; c < cols; c++ {
			t[rank][c] /= pivot
		}

		for row := 0; row < rows; row++ {
			if row != rank && t[row][col] != 0 {
				factor := t[row][col]
				for c := col; c < cols; c++ {
					t[row][c] -= factor * t[rank][c]
				}
			}
		}

		rank++
	}

	return rank
}

func ComputeAInverse(parties []*model.Party) {
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	triple := GenerateMultiplicationTriple(len(parties), t, s, s, s)
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))

	// Compute locally
	for partyNumber, party := range parties {
		ai := triple.A[partyNumber]
		bi := triple.B[partyNumber]
		TInverse := RightInverse(party.T)

		di := AddMatricesNew(MultiplyMatrices(party.S, TInverse), ai)
		ei := AddMatricesNew(party.R, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	zShares := multiplicationProtocol(parties, triple, dShares, eShares, t, s, s, s)
	for i, party := range parties {
		party.AInverse = zShares[i]
	}

	// TODO: Remove this check when benchmarking
	ARecovered := generateZeroMatrix(s, t)
	AInverseRecovered := generateZeroMatrix(t, s)
	for _, party := range parties {
		AddMatrices(ARecovered, party.A)
		AddMatrices(AInverseRecovered, party.AInverse)
	}

	Identity := MultiplyMatrices(ARecovered, AInverseRecovered)
	for _, row := range Identity {
		fmt.Println(fmt.Sprintf("%2d", row))
	}
}

func RightInverse(t [][]byte) [][]byte {
	_, invTable := generateMulAndInvTable() // TODO: refactor

	M := len(t)    // Rows
	N := len(t[0]) // Columns

	if M > N {
		return nil
	}

	// Augment A with an identity matrix to form (A | I)
	augmented := make([][]byte, M)
	for i := 0; i < M; i++ {
		augmented[i] = make([]byte, N+M)
		copy(augmented[i], t[i])
		for j := 0; j < M; j++ {
			if i == j {
				augmented[i][N+j] = 1
			}
		}
	}

	// Perform Gaussian elimination
	for i := 0; i < M; i++ {
		// Find pivot
		if augmented[i][i] == 0 {
			// Swap with a row below that has a nonzero pivot
			for k := i + 1; k < M; k++ {
				if augmented[k][i] != 0 {
					augmented[i], augmented[k] = augmented[k], augmented[i]
					break
				}
			}
		}

		// Ensure pivot is nonzero
		if augmented[i][i] == 0 {
			return nil
		}

		// Normalize pivot row
		pivotInv := invTable[augmented[i][i]]
		for j := 0; j < N+M; j++ {
			augmented[i][j] = gf16Mul(augmented[i][j], pivotInv)
		}

		// Eliminate other rows
		for k := 0; k < M; k++ {
			if k != i && augmented[k][i] != 0 {
				factor := augmented[k][i]
				for j := 0; j < N+M; j++ {
					augmented[k][j] = augmented[k][j] ^ gf16Mul(factor, augmented[i][j])
				}
			}
		}
	}

	// Extract the right inverse (n x m matrix)
	B := make([][]byte, N)
	for i := 0; i < N; i++ {
		B[i] = make([]byte, M)
		if i < M {
			copy(B[i], augmented[i][N:])
		}
	}

	return B
}

func ComputeLittleX(parties []*model.Party) {
	s := len(parties[0].A)
	t := len(parties[0].A[0])

	basis := Kernel(parties[0].T)

	if len(basis) != t-s {
		panic(fmt.Errorf("length of basis is incorrect"))
	}

	for _, party := range parties {
		z := make([]byte, t)
		zVector := RandVector(t - s)

		for index, elem := range zVector {
			z = AddVec(z, MultiplyVecConstant(elem, basis[index]))
		}

		party.Z = z
	}

	triplesStep7 := GenerateMultiplicationTriple(len(parties), t, s, s, 1)
	triplesStep8 := GenerateMultiplicationTriple(len(parties), t, t, t, 1)

	// Compute [A^-1] * [b]
	dShares := make([][][]byte, len(parties))
	eShares := make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.S = RandMatrix(t, t)

		ai := triplesStep7.A[partyNumber]
		bi := triplesStep7.B[partyNumber]
		di := AddMatricesNew(party.AInverse, ai)
		ei := AddMatricesNew(vectorToMatrix(party.LittleY), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesB := multiplicationProtocol(parties, triplesStep7, dShares, eShares, t, s, s, 1)

	// Compute [S] * [z]
	dShares = make([][][]byte, len(parties))
	eShares = make([][][]byte, len(parties))
	for partyNumber, party := range parties {
		party.R = RandMatrix(s, s)

		ai := triplesStep8.A[partyNumber]
		bi := triplesStep8.B[partyNumber]
		di := AddMatricesNew(party.S, ai)
		ei := AddMatricesNew(vectorToMatrix(party.Z), bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	STimesZ := multiplicationProtocol(parties, triplesStep8, dShares, eShares, t, t, t, 1)

	// [x] = [A^-1] * [b] + [S] * [z]
	for i, _ := range parties {
		AddMatrices(ATimesB[i], STimesZ[i])
	}

	for i, party := range parties {
		party.LittleX = matrixToVec(ATimesB[i])
	}
}

// matrixToVec takes a column vector (as a matrix) and returns a row vector
func matrixToVec(A [][]byte) []byte {
	result := make([]byte, len(A))

	for index, elem := range A {
		result[index] = elem[0]
	}

	return result
}

func Kernel(T [][]byte) [][]byte {
	rows := len(T)    // s
	cols := len(T[0]) // t
	basis := make([][]byte, 0)

	// Add the identity matrix below T
	Identity := generateIdentityMatrix(cols)
	TWithIdentity := appendMatrixBelow(T, Identity)

	// Perform Gaussian elimination
	EchelonForm := echelonForm(MatrixTranspose(TWithIdentity))

	// Determine the basis which are the rows where the first s values are 0
	for i := rows; i < len(EchelonForm); i++ { // Look at the added identity rows
		isKernelVector := true
		for j := 0; j < rows; j++ {
			if EchelonForm[i][j] != 0 {
				isKernelVector = false
				break
			}
		}
		if isKernelVector {
			basis = append(basis, EchelonForm[i][rows:])
		}
	}

	return basis
}

func echelonForm(B [][]byte) [][]byte {
	rows := len(B)
	cols := len(B[0])
	pivotColumn := 0
	pivotRow := 0

	// TODO: Fix
	field := InitField()

	for pivotRow < rows && pivotColumn < cols+1 {
		var possiblePivots []int
		for i := pivotRow; i < rows; i++ {
			if B[i][pivotColumn] != 0 {
				possiblePivots = append(possiblePivots, i)
			}
		}

		if len(possiblePivots) == 0 {
			pivotColumn++
			continue
		}

		nextPivotRow := slices.Min(possiblePivots)
		B[pivotRow], B[nextPivotRow] = B[nextPivotRow], B[pivotRow]

		// Make the leading entry a 1
		B[pivotRow] = MultiplyVecConstant(field.Gf16Inv(B[pivotRow][pivotColumn]), B[pivotRow])

		// Eliminate entries below the pivot
		for row := nextPivotRow + 1; row < rows; row++ {
			B[row] = AddVec(B[row], MultiplyVecConstant(B[row][pivotColumn], B[pivotRow]))
		}

		pivotRow++
		pivotColumn++
	}

	return B
}

func MultiplyVecConstant(b byte, a []byte) []byte {
	C := make([]byte, len(a))
	for i := range C {
		C[i] = gf16Mul(b, a[i])
	}
	return C
}

func appendMatrixBelow(A, B [][]byte) [][]byte {
	if len(A[0]) != len(B[0]) {
		panic(123)
	}

	resultSize := len(A) + len(B)
	result := make([][]byte, resultSize)

	for i := 0; i < len(A); i++ {
		result[i] = A[i]
	}

	for i := 0; i < len(B); i++ {
		result[i+len(A)] = B[i]
	}

	return result
}

type Field struct {
	mulTable [][]byte
	invTable []byte
}

func InitField() *Field {
	mulTable, invTable := generateMulAndInvTable()

	return &Field{
		mulTable: mulTable,
		invTable: invTable,
	}
}

func (f *Field) VectorTransposedMatrixMul(vec []byte, matrix [][]byte) []byte {
	cols := len(matrix)
	if cols == 0 || len(vec) != len(matrix) {
		panic("Vector length must match matrix row count")
	}

	rows := len(matrix[0])
	result := make([]byte, rows)

	for i := 0; i < rows; i++ {
		var sum byte
		for j := 0; j < cols; j++ {
			sum ^= f.Gf16Mul(vec[j], matrix[j][i])
		}
		result[i] = sum
	}

	return result
}

// MatrixVectorMul Takes a matrix and vector, here we assume that the output of this multiplication will be a
// vector, since this is the case in MAYO.
func (f *Field) MatrixVectorMul(matrix [][]byte, vec []byte) []byte {
	rows := len(matrix)
	if rows == 0 || len(vec) != len(matrix[0]) {
		panic("Vector length must match matrix column count")
	}

	cols := len(matrix[0])
	result := make([]byte, rows)

	for i := 0; i < rows; i++ {
		var sum byte
		for j := 0; j < cols; j++ {
			sum ^= f.Gf16Mul(vec[j], matrix[i][j])
		}
		result[i] = sum
	}

	return result
}

// MultiplyMatrices multiplies two matrices
func (f *Field) MultiplyMatrices(A, B [][]byte) [][]byte {
	rowsA, colsA := len(A), len(A[0])
	rowsB, colsB := len(B), len(B[0])

	if colsA != rowsB {
		panic(fmt.Sprintf("Cannot multiply matrices colsA: '%d', rowsB: '%d'", colsA, rowsB))
	}

	C := make([][]byte, rowsA)
	for i := range C {
		C[i] = make([]byte, colsB)
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				C[i][j] ^= f.Gf16Mul(A[i][k], B[k][j])
			}
		}
	}

	return C
}

// MultiplyVecConstant multiplies a vector by a constant element-wise
func (f *Field) MultiplyVecConstant(b byte, a []byte) []byte {
	C := make([]byte, len(a))
	for i := range C {
		C[i] = f.Gf16Mul(b, a[i])
	}
	return C
}

// Gf16Mul multiplies two elements in GF(16)
func (f *Field) Gf16Mul(a, b byte) byte {
	return f.mulTable[a][b]
}

// Gf16Inv calculates the inverse of an element in GF(16)
func (f *Field) Gf16Inv(a byte) byte {
	return f.invTable[a]
}

func (f *Field) VecInnerProduct(vec1Transposed []byte, vec2 []byte) byte {
	if len(vec1Transposed) != len(vec2) {
		panic("Vectors must have the same length")
	}

	var result byte = 0
	for i := 0; i < len(vec1Transposed); i++ {
		result ^= f.Gf16Mul(vec1Transposed[i], vec2[i])
	}

	return result
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

func AddVec(A, B []byte) []byte {
	if len(A) != len(B) {
		panic("Cannot add vectors of different lengths")
	}

	C := make([]byte, len(A))
	for i := range C {
		C[i] = A[i] ^ B[i]
	}

	return C
}
