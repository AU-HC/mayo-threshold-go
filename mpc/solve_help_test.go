package mpc

import (
	"reflect"
	"testing"
)

func TestRightInverseOfMatrixForSquareMatrix(t *testing.T) {
	matrix := [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 7},
	}
	rightInverse := computeRightInverse(matrix)

	identity := make([][]byte, len(matrix))
	for i := range identity {
		identity[i] = make([]byte, len(matrix))
		identity[i][i] = 1
	}

	actual := MultiplyMatrices(matrix, rightInverse)
	if !reflect.DeepEqual(actual, identity) {
		t.Errorf("Expected identity matrix, but got %v", actual)
		return
	}
}

func TestRightInverseOfMatrixForNonSquareMatrix(t *testing.T) {
	matrix := [][]byte{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	}
	rightInverse := computeRightInverse(matrix)

	identity := make([][]byte, len(matrix))
	for i := range identity {
		identity[i] = make([]byte, len(matrix))
		identity[i][i] = 1
	}

	actual := MultiplyMatrices(matrix, rightInverse)
	if !reflect.DeepEqual(actual, identity) {
		t.Errorf("Expected identity matrix, but got %v", actual)
		return
	}
}

func TestKernel(t *testing.T) {
	matrix := [][]byte{
		{1, 2, 3},
		{2, 4, 6},
	}

	// Expected kernel basis (rows as basis vectors)
	expectedKernelBasis := [][]byte{
		{0, 1, 15},
	}

	// Compute actual kernel basis
	actualKernelBasis := computeBasisOfKernel(matrix)

	// Compare expected vs actual
	if !reflect.DeepEqual(actualKernelBasis, expectedKernelBasis) {
		t.Errorf("computeBasisOfKernel basis mismatch. Expected %v but got %v", expectedKernelBasis, actualKernelBasis)
	}
}

func TestAppendMatrixHorizontal(t *testing.T) {
	A := [][]byte{
		{1, 2},
		{3, 4},
	}

	B := [][]byte{
		{5, 6},
		{7, 8},
	}

	// Expected result after horizontal append
	expected := [][]byte{
		{1, 2, 5, 6},
		{3, 4, 7, 8},
	}

	// Compute actual result
	actual := appendMatrixHorizontal(A, B)

	// Compare expected vs actual
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Matrix append mismatch. Expected %v but got %v", expected, actual)
	}
}
