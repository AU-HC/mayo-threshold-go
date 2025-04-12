package main

import (
	"encoding/hex"
	"fmt"
	"mayo-threshold-go/flags"
	"mayo-threshold-go/mpc"
	"time"
)

const AmountOfParties = 4
const Threshold = 4

func main() {
	// Get application flags
	arguments := flags.GetApplicationArguments()
	amountOfBenchmarkSamples := arguments.AmountBenchmarkingSamples

	// If amount of samples > 0, then benchmark and write benchmarks to results/
	if amountOfBenchmarkSamples > 0 {
		benchmark(amountOfBenchmarkSamples)
		return
	}

	// Define the message and the context
	context := mpc.CreateContext(AmountOfParties, Threshold)
	message := []byte("Hello, world!")

	// Generate expanded public key, and shares of expanded secret key for the parties
	before := time.Now()
	epk, parties := context.KeyGenAPI(AmountOfParties)
	fmt.Println(fmt.Sprintf("Key generation with %d parties took: %dms", AmountOfParties, time.Since(before).Milliseconds()))

	// Threshold sign message
	before = time.Now()
	sig := context.SignAPI(message, parties)
	fmt.Println(fmt.Sprintf("Signing with %d parties took: %dms", AmountOfParties, time.Since(before).Milliseconds()))

	// Verify message
	before = time.Now()
	valid := context.Verify(epk, message, sig)
	fmt.Println(fmt.Sprintf("Verify took: %dms", time.Since(before).Milliseconds()))
	if valid {
		fmt.Println(fmt.Sprintf("Signature: '%s' is a valid signature on the message: '%s'",
			hex.EncodeToString(sig.Bytes()), message))
	} else {
		fmt.Println(fmt.Sprintf("Signature: '%s' is not a valid signature on the message: '%s'",
			hex.EncodeToString(sig.Bytes()), message))
	}

	mpc.Test()
}

func benchmark(n int) {
	path, err := mpc.Benchmark(n)

	if err != nil {
		fmt.Println("Got error while benchmarking: ", err)
	}

	fmt.Println(fmt.Sprintf("Benchmarking done, see /%s for more information", path))
}
