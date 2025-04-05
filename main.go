package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"mayo-threshold-go/mock"
	"mayo-threshold-go/mpc"
	"time"
)

const AmountOfParties = 10

func main() {
	// Set seed for easier testing
	rand.Seed(99)

	// Get mock esk, epk from json files and define message
	message := []byte("Hello, world!")
	esk, epk := mock.GetExpandedKeyPair()

	// Start the parties, by giving them the epk and shares of the esk
	parties := mock.CreatePartiesAndSharesForEsk(esk, epk, AmountOfParties)

	// Threshold sign message
	before := time.Now()
	sig := mpc.Sign(message, parties)
	fmt.Println(fmt.Sprintf("Signing with %d parties, took: %dms", AmountOfParties, time.Since(before).Milliseconds()))

	// Verify message
	before = time.Now()
	valid := mpc.Verify(epk, message, sig)
	fmt.Println(fmt.Sprintf("Verify took: %dms", time.Since(before).Milliseconds()))
	if valid {
		fmt.Println(fmt.Sprintf("Signature: '%s' is a valid signature on the message: '%s'",
			hex.EncodeToString(sig.Bytes()), message))
	} else {
		fmt.Println(fmt.Sprintf("Signature: '%s' is not a valid signature on the message: '%s'",
			hex.EncodeToString(sig.Bytes()), message))
	}
}
