package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"mayo-threshold-go/mock"
	"mayo-threshold-go/mpc"
)

const AmountOfParties = 2

func main() {
	// Set seed for easier testing
	rand.Seed(99)

	// Get mock esk, epk from json files and define message
	message := []byte("Hello, world!")
	esk, epk := mock.GetExpandedKeyPair()

	// Start the parties, by giving them the epk and shares of the esk
	parties := mock.CreatePartiesAndSharesForEsk(esk, epk, AmountOfParties)

	// Sign message
	signature := mpc.Sign(message, parties)
	fmt.Println(hex.EncodeToString(signature.Encode()))

	// Verify message TODO: make this print prettier
	valid := mpc.ThresholdVerify(parties, signature)
	fmt.Println("Signature was valid:", valid)
}
