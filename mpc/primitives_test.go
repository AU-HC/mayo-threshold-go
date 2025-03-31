package mpc

import "testing"

func TestGenerationOfMatrixTriples(t *testing.T) {
	_ = GenerateMultiplicationTriple(2, 5, 5, 5, 5)
}
