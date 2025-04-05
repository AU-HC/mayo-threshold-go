package mpc

import (
	"fmt"
	"slices"
)

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

func vectorToMatrix(x []byte) [][]byte {
	result := make([][]byte, len(x))

	for i, elem := range x {
		result[i] = []byte{elem}
	}

	return result
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

// matrixToVec takes a column vector (as a matrix) and returns a row vector
func matrixToVec(A [][]byte) []byte {
	result := make([]byte, len(A))

	for index, elem := range A {
		result[index] = elem[0]
	}

	return result
}

func matrixify(v []byte, rows, cols int) [][]byte {
	if len(v) != rows*cols {
		panic(fmt.Errorf("a does not have the correct dimensions for matrixify"))
	}

	matrix := make([][]byte, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]byte, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j] = v[i*cols+j]
		}
	}
	return matrix
}

func generateZeroMatrix(rows, columns int) [][]byte {
	matrix := make([][]byte, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]byte, columns)
	}

	return matrix
}

func generateIdentityMatrix(dimension int) [][]byte {
	matrix := make([][]byte, dimension)

	for i := 0; i < dimension; i++ {
		matrix[i] = make([]byte, dimension)
		matrix[i][i] = 1
	}

	return matrix
}

func AddMatrices(a, b [][]byte) {
	if len(a) != len(b) || len(a[0]) != len(b[0]) {
		panic(fmt.Errorf("a and b do not have the same dimensions "))
	}

	for i := range a {
		for j := range a[i] {
			a[i][j] ^= b[i][j]
		}
	}
}

func AddMatricesNew(a, b [][]byte) [][]byte {
	if len(a) != len(b) || len(a[0]) != len(b[0]) {
		panic(fmt.Errorf("a and b do not have the same dimensions"))
	}

	c := generateZeroMatrix(len(a), len(a[0]))

	for i := range a {
		for j := range a[i] {
			c[i][j] = a[i][j] ^ b[i][j]
		}
	}

	return c
}

func MultiplyMatrices(A, B [][]byte) [][]byte {
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
				C[i][j] ^= gf16Mul(A[i][k], B[k][j])
			}
		}
	}

	return C
}

func MatrixTranspose(a [][]byte) [][]byte {
	if len(a) == 0 {
		return [][]byte{}
	}

	rows, cols := len(a), len(a[0])
	result := make([][]byte, cols)
	for i := range result {
		result[i] = make([]byte, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			result[j][i] = a[i][j]
		}
	}

	return result
}

func appendMatrixVertical(A, B [][]byte) [][]byte {
	if len(A[0]) != len(B[0]) {
		panic("Cannot append matrices of different column count")
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

func appendMatrixHorizontal(A, B [][]byte) [][]byte {
	if len(A) != len(B) {
		panic("Cannot append matrices of different row count")
	}

	result := make([][]byte, len(A))

	for i := 0; i < len(A); i++ {
		result[i] = append(A[i], B[i]...)
	}

	return result
}

func computeBasisOfKernel(T [][]byte) [][]byte {
	rows := len(T)    // s
	cols := len(T[0]) // t
	basis := make([][]byte, 0)

	// Add the identity matrix below T
	Identity := generateIdentityMatrix(cols)
	TWithIdentity := appendMatrixVertical(T, Identity)

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

func computeRightInverse(t [][]byte) [][]byte {
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

func MultiplyVecConstant(b byte, a []byte) []byte {
	C := make([]byte, len(a))
	for i := range C {
		C[i] = gf16Mul(b, a[i])
	}
	return C
}
