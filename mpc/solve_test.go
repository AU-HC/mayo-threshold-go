package mpc

import (
	"reflect"
	"testing"
)

func TestRankOfMatrix(t *testing.T) {
	matrix := [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 7},
	}
	expected := 3
	result := rankOfMatrix(matrix)
	if result != expected {
		t.Errorf("Expected rank %d, but got %d", expected, result)
	}

	matrix = [][]byte{
		{1, 2, 1, 2},
		{1, 3, 2, 2},
		{2, 4, 3, 4},
		{3, 7, 4, 6},
	}
	expected = 3
	result = rankOfMatrix(matrix)
	if result != expected {
		t.Errorf("Expected rank %d, but got %d", expected, result)
	}
}

func TestRightInverseOfMatrixForSquareMatrix(t *testing.T) {
	matrix := [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 7},
	}
	rightInverse := RightInverse(matrix)

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
	rightInverse := RightInverse(matrix)

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
	actualKernelBasis := Kernel(matrix)

	// Compare expected vs actual
	if !reflect.DeepEqual(actualKernelBasis, expectedKernelBasis) {
		t.Errorf("Kernel basis mismatch. Expected %v but got %v", expectedKernelBasis, actualKernelBasis)
	}
}
