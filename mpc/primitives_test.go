package mpc

import "testing"

func TestMatrixMultiplication(t *testing.T) {
	_, _, _ = GenerateMultiplicationTriples(2, 5, 5, 5, 5)
}
