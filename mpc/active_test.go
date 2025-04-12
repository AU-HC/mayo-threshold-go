package mpc

import (
	"mayo-threshold-go/rand"
	"reflect"
	"testing"
)

func TestActiveOpenMatrix(t *testing.T) {
	parties := 2
	rows, cols := 3, 3
	secret := rand.Matrix(rows, cols)

	shares := createSharesForMatrix(parties, secret)
	secretOpen, err := openMatrix(shares)

	if !reflect.DeepEqual(secret, secretOpen) {
		t.Errorf("Expected %2d, but got %2d", secret, secretOpen)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestActiveAddMatrices(t *testing.T) {
	parties := 2
	rows, cols := 3, 3
	matrix1 := rand.Matrix(rows, cols)
	matrix2 := rand.Matrix(rows, cols)
	shares1 := createSharesForMatrix(parties, matrix1)
	shares2 := createSharesForMatrix(parties, matrix2)

	result := make([]MatrixShare, parties)
	for i := 0; i < parties; i++ {
		result[i] = AddMatrixShares(shares1[i], shares2[i])
	}

	actual, err := openMatrix(result)
	expected := AddMatricesNew(matrix1, matrix2)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %2d, but got %2d", expected, actual)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestActiveMulMatrixShareWithConstantLeft(t *testing.T) {
	parties := 2
	matrix1 := rand.Matrix(2, 3)
	matrix2 := rand.Matrix(3, 2)
	shares2 := createSharesForMatrix(parties, matrix2)

	result := make([]MatrixShare, parties)
	for i := 0; i < parties; i++ {
		result[i] = MulMatrixShareWithConstantLeft(matrix1, shares2[i])
	}

	actual, err := openMatrix(result)
	expected := MultiplyMatrices(matrix1, matrix2)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %2d, but got %2d", expected, actual)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestActiveMulMatrixShareWithConstantLeftOtherDimension(t *testing.T) {
	parties := 2
	matrix1 := rand.Matrix(3, 2)
	matrix2 := rand.Matrix(2, 3)
	shares2 := createSharesForMatrix(parties, matrix2)

	result := make([]MatrixShare, parties)
	for i := 0; i < parties; i++ {
		result[i] = MulMatrixShareWithConstantLeft(matrix1, shares2[i])
	}

	actual, err := openMatrix(result)
	expected := MultiplyMatrices(matrix1, matrix2)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %2d, but got %2d", expected, actual)
	}
	if err != nil {
		t.Error(err)
	}
}
