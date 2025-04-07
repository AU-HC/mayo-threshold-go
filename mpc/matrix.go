package mpc

import (
	"fmt"
	"reflect"
	"slices"
)

// upper transposes the lower triangular part of a matrix to the upper triangular part
func upper(matrix [][]byte) [][]byte {
	n := len(matrix)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			matrix[i][j] = matrix[i][j] ^ matrix[j][i] // Update upper triangular part
			matrix[j][i] = 0
		}
	}

	return matrix
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
		panic(fmt.Errorf("a and b do not have the same dimensions (%d, %d), (%d, %d)", len(a), len(a[0]), len(b), len(b[0])))
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
				C[i][j] ^= field.Gf16Mul(A[i][k], B[k][j])
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

func rowReduceToRREF(B [][]byte) [][]byte {
	rows := len(B)
	cols := len(B[0])
	pivotRow := 0

	for pivotCol := 0; pivotCol < cols && pivotRow < rows; pivotCol++ {
		// Find pivot
		pivotIndex := -1
		for i := pivotRow; i < rows; i++ {
			if B[i][pivotCol] != 0 {
				pivotIndex = i
				break
			}
		}

		if pivotIndex == -1 {
			continue // No pivot in this column
		}

		// Swap to top of remaining submatrix
		if pivotIndex != pivotRow {
			B[pivotRow], B[pivotIndex] = B[pivotIndex], B[pivotRow]
		}

		// Normalize pivot row
		pivotVal := B[pivotRow][pivotCol]
		inv := field.Gf16Inv(pivotVal)
		B[pivotRow] = MultiplyVecConstant(inv, B[pivotRow])

		// Eliminate above and below
		for i := 0; i < rows; i++ {
			if i != pivotRow && B[i][pivotCol] != 0 {
				factor := B[i][pivotCol]
				sub := MultiplyVecConstant(factor, B[pivotRow])
				B[i] = AddVec(B[i], sub)
			}
		}

		pivotRow++
	}

	return B
}

func isFullRank(matrix [][]byte) bool {
	r := len(matrix)
	c := len(matrix[0])
	minRank := min(r, c)
	rank := 0

	// Copy the matrix
	B := make([][]byte, len(matrix))
	for i := range matrix {
		B[i] = make([]byte, len(matrix[i]))
		copy(B[i], matrix[i])
	}

	// Reduce to reduced row echelon form
	B = rowReduceToRREF(B)

	// Count non-zero rows
	for i := 0; i < r; i++ {
		if B[i][i] == 1 {
			rank++
		}
	}

	return rank == minRank
}

func computeRightInverse(A [][]byte) [][]byte {
	M := len(A)    // Rows
	N := len(A[0]) // Columns

	if M > N {
		fmt.Println("Matrix is not full row rank or too tall; no right inverse.")
		return nil
	}

	// Create identity matrix of size M
	identity := generateIdentityMatrix(M)

	// Right inverse will be N x M
	B := make([][]byte, N)
	for i := range B {
		B[i] = make([]byte, M)
	}

	// For each column of the identity matrix
	for col := 0; col < M; col++ {
		// Create a deep copy of A and augment with the col-th column of identity
		aug := make([][]byte, M)
		for i := 0; i < M; i++ {
			aug[i] = make([]byte, N+1)
			copy(aug[i], A[i])
			aug[i][N] = identity[i][col] // Augment with e_col
		}

		// Perform Gaussian elimination to solve A x = e_col
		for row := 0; row < M; row++ {
			// Find pivot
			pivotRow := -1
			for i := row; i < M; i++ {
				if aug[i][row] != 0 {
					pivotRow = i
					break
				}
			}
			if pivotRow == -1 {
				return nil // Not full rank
			}

			if pivotRow != row {
				aug[row], aug[pivotRow] = aug[pivotRow], aug[row]
			}

			// Normalize pivot
			inv := field.invTable[aug[row][row]]
			for j := 0; j <= N; j++ {
				aug[row][j] = field.Gf16Mul(aug[row][j], inv)
			}

			// Eliminate other rows
			for i := 0; i < M; i++ {
				if i != row && aug[i][row] != 0 {
					factor := aug[i][row]
					for j := 0; j <= N; j++ {
						aug[i][j] ^= field.Gf16Mul(factor, aug[row][j])
					}
				}
			}
		}

		// Extract solution (column vector) into B
		for i := 0; i < N; i++ {
			if i < M {
				B[i][col] = aug[i][N]
			}
		}
	}

	// Optional: verify A * B = I_M
	if !reflect.DeepEqual(MultiplyMatrices(A, B), identity) {
		panic("Matrix multiplication check failed; no valid right inverse")
	}

	return B
}

func MultiplyVecConstant(b byte, a []byte) []byte {
	C := make([]byte, len(a))
	for i := range C {
		C[i] = field.Gf16Mul(b, a[i])
	}
	return C
}
