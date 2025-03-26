package main

import (
	"fmt"
	"mayo-threshold-go/mock"
)

const n = 2

func main() {
	// Get mock esk, epk
	esk, epk := mock.GetExpandedKeyPair()

	// Start the parties, by giving them the epk and shares of the esk
	parties := mock.CreatePartiesAndSharesForEsk(esk, epk, n)
	alice := parties[0]
	bob := parties[1]

	// Do stuff
	if alice.Epk.P1[0][0][0] == bob.Epk.P1[1][0][0] {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}
