package mpc

import "testing"

func TestGenerationOfMatrixTriples(t *testing.T) {
	_ = GenerateMultiplicationTriple(2, 5, 7, 7, 5)
}
