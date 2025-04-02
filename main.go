package main

import (
	"fmt"
	"mayo-threshold-go/mock"
	"mayo-threshold-go/mpc"
)

const n = 2
const lambda = 2

func main() {
	// Get mock esk, epk and define message
	//message := []byte("Hello, world!")
	esk, epk := mock.GetExpandedKeyPair()

	// Start the parties, by giving them the epk and shares of the esk
	parties := mock.CreatePartiesAndSharesForEsk(esk, epk, n)

	// Begin signing the message
	for true {
		// Steps 1-3 of sign
		mpc.ComputeM(parties)
		// Step 4 of sign
		mpc.ComputeY(parties)
		// Step 5 of sign
		mpc.LocalComputeA(parties) // TODO: test
		mpc.LocalComputeY(parties) // TODO: test

		// Step 6 of sign
		// ** Algorithm solve **
		// Steps 1-4 of solve
		isRankDefect := mpc.ComputeT(parties)
		if !isRankDefect {
			break
		}
		fmt.Println("Matrix was rank-defect")
	}
	// Step 5 of solve
	mpc.ComputeAInverse(parties)
	// Steps 6-9 of solve
	mpc.ComputeLittleX(parties) // TODO: figure out if the spec is correct / implement
	// ** Algorithm solve **

	// Step 7-8 of sign
	mpc.ComputeSPrime(parties) // TODO: figure matrixify and dimensions
	// step 9 of sign
	mpc.ComputeSignature(parties) // TODO: implement
}
