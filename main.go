package main

import (
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

	// Steps 1-3 of sign
	mpc.ComputeM(parties)
	// Step 4 of sign
	mpc.ComputeY(parties)
	// Step 5 of sign
	mpc.LocalComputeA(parties)
	mpc.LocalComputeY(parties)

	// Step 6 of sign
	// ** Algorithm solve **
	// Steps 1-4 of solve
	mpc.ComputeT(parties)
	// Step 5 of solve
	mpc.ComputeAInverse(parties)
	// Steps 6-9 of solve
	mpc.Computex(parties)
	// ** Algorithm solve **

	// Step 7 of sign
	mpc.ComputeX(parties)
	// step 8 of sign
	mpc.ComputeS(parties)
	// step 9 of sign
	mpc.ComputeSignature(parties)
}
