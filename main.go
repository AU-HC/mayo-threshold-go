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

	// Steps 1-3 of sign
	mpc.ComputeM(parties)
	// Step 4 of sign
	mpc.ComputeY(parties)
	// Step 5 of sign
	mpc.LocalComputeA(parties)
	mpc.LocalComputeY(parties)

	// Steps 1-4 of solve
	mpc.ComputeT(parties)
	// TODO: Remaining steps of solve

	// Step 6 of sign
	// TODO: Remaining steps of sign

	// Do stuff
	fmt.Println(mock.VerifyShares(esk, parties))
}
