package main

import (
	"encoding/hex"
	"fmt"
	"mayo-threshold-go/mpc"
	"time"
)

const AmountOfParties = 3

func main() {
	// Define the message
	message := []byte("Hello, world!")

	// Generate expanded public key, and shares of expanded secret key for the parties
	before := time.Now()
	epk, parties := mpc.KeyGen(AmountOfParties)
	fmt.Println(fmt.Sprintf("Key generation with %d parties took: %dms", AmountOfParties, time.Since(before).Milliseconds()))

	// Threshold sign message
	before = time.Now()
	sig := mpc.Sign(message, parties)
	fmt.Println(fmt.Sprintf("Signing with %d parties took: %dms", AmountOfParties, time.Since(before).Milliseconds()))

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
