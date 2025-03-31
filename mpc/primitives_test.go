package mpc

import "testing"

func TestGenerationOfMatrixTriples(t *testing.T) {
	_, _, _ = GenerateMultiplicationTriples(2, 5, 5, 5, 5)
}
