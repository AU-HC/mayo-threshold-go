package mpc

import "testing"

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
