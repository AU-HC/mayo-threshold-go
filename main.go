package main

import (
	"fmt"
	"mayo-threshold-go/mock"
)

const n = 2
const lambda = 2

func main() {
	// Get mock esk, epk and define message
	//message := []byte("Hello, world!")
	esk, epk := mock.GetExpandedKeyPair()

	// Start the parties, by giving them the epk and shares of the esk
	parties := mock.CreatePartiesAndSharesForEsk(esk, epk, n)
	//alice := parties[0]
	//bob := parties[1]

	// Begin signing message
	//salt := mpc.Coin(parties, lambda)

	// Do stuff
	fmt.Println(mock.VerifyShares(esk, parties))
}
