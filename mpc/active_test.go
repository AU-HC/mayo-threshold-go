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
	secretOpen, _ := openMatrix(shares)

	if !reflect.DeepEqual(secret, secretOpen) {
		t.Errorf("Expected %2d, but got %2d", secret, secretOpen)
	}
}
