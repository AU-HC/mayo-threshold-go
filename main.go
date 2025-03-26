package main

import (
	"fmt"
	"mayo-threshold-go/mock"
)

func main() {
	esk, epk := mock.GetExpandedKeyPair()

	fmt.Println(esk)
	fmt.Println(epk)

	//skShares := mock.GenerateRandomShares(esk)

}
