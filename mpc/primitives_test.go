package mpc

import (
	"reflect"
	"testing"
)

func TestGenerationOfMatrixTriples(t *testing.T) {
	_ = GenerateMultiplicationTriple(2, 5, 7, 7, 5)
}

func TestMatrixify(t *testing.T) {
	vector := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedMatrix := [][]byte{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10},
	}

	matrix := matrixify(vector, 2, 5)

	if !reflect.DeepEqual(matrix, expectedMatrix) {
		t.Errorf("Expected %v, got %v", expectedMatrix, matrix)
	}
}
