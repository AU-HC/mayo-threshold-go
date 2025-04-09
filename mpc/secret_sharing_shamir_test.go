package mpc

import (
	"testing"
)

func TestGenerateCoefficientsSecretValueIsCorrect(t *testing.T) {
	secret := byte(2)

	coefficients := generateCoefficients(secret, 5)

	if coefficients[0] != secret {
		t.Errorf("Expected coefficients[0] to be %v, got %v", secret, coefficients[0])
	}
}

func TestGenerateCoefficientsAmountOfCoefficientsIsCorrect(t *testing.T) {
	secret := byte(2)

	coefficients := generateCoefficients(secret, 3)

	if len(coefficients) != 3 {
		t.Errorf("Expected amount of coefficients to be %v, got %v", 3, len(coefficients))
	}
}

func TestCreateSharesOutputsTheCorrectAmountOfShares(t *testing.T) {
	secret := byte(2)

	shares := createShares(secret, 5, 3)

	if len(shares) != 5 {
		t.Errorf("Expected amount of shares to be %v, got %v", 3, len(shares))
	}
}

func TestCreateSharesAndReconstructIsCorrect(t *testing.T) {
	secret := byte(2)

	shares := createShares(secret, 5, 3)
	actual := reconstructSecret(shares, 3)

	if actual != secret {
		t.Errorf("Expected reconstructed secret to be %v, got %v", secret, actual)
	}
}

func TestCreateSharesAndReconstructIsCorrectWhenExactlyTSharesArePresent(t *testing.T) {
	secret := byte(2)

	shares := createShares(secret, 5, 3)
	actual := reconstructSecret(shares[:3], 3)

	if actual != secret {
		t.Errorf("Expected reconstructed secret to be %v, got %v", secret, actual)
	}
}

func TestAdditionOfShares(t *testing.T) {
	secret1 := byte(2)
	secret2 := byte(5)
	shares1 := createShares(secret1, 5, 3)
	shares2 := createShares(secret2, 5, 3)
	additionShares := make([]byte, 5)
	for i := 0; i < 5; i++ {
		additionShares[i] = shares1[i] ^ shares2[i]
	}

	actual := reconstructSecret(additionShares[:5], 3)

	if actual != (secret1 ^ secret2) {
		t.Errorf("Expected reconstructed secret to be %v, got %v", secret1^secret2, actual)
	}
}

func TestMultiplicationWithConstant(t *testing.T) {
	secret1 := byte(2)
	constant := byte(5)
	shares1 := createShares(secret1, 5, 3)
	multiplicationShares := make([]byte, 5)
	for i := 0; i < 5; i++ {
		multiplicationShares[i] = gf16Mul(shares1[i], constant)
	}

	actual := reconstructSecret(multiplicationShares[:5], 3)

	expected := gf16Mul(secret1, constant)
	if actual != expected {
		t.Errorf("Expected reconstructed secret to be %v, got %v", expected, actual)
	}
}
