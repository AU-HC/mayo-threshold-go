package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"mayo-threshold-go/mock"
	"mayo-threshold-go/mpc"
)

const n = 1
const lambda = 2

const m = 64
const k = 4
const N = 81
const o = 17
const v = n - o

func main() {
	rand.Seed(99)

	// Get mock esk, epk and define message
	//message := []byte("Hello, world!")
	esk, epk := mock.GetExpandedKeyPair()

	// Start the parties, by giving them the epk and shares of the esk
	parties := mock.CreatePartiesAndSharesForEsk(esk, epk, n)

	// Begin signing the message
	for true {
		// Steps 1-3 of sign
		mpc.ComputeM(parties, []byte("Hello, world!"))
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
	mpc.ComputeLittleX(parties) // TODO: figure out if the spec is correct / test
	// ** Algorithm solve **

	// Step 7-9 of sign
	signature := mpc.ComputeSPrime(parties)
	encodedSignature := signature.Encode()
	fmt.Println(hex.EncodeToString(encodedSignature))

	xd := "0351b3bdd3777207e089b260fc048987ab1b5ff4bb8b93b5c5029d32ad38a8f2494cbe2433c665d3d140199a91903dd9fec45616eaf06553328a26ac4572d9c4511f0036c0cf8356e612bdfbc29a90cea11a62390b17260d0c902a12702de25855913c0e9fa3e819e7162bd799b85326c4aec547e31aedbc064116396b61c73a32ad499a097d806fe2ccaf634fd93be5df400096dc8d8412bcccd01ae90887bb63e34a57137052dfef6649b2a59d6a8b137fd941a62f062e002b48656c6c6f2c20776f726c6421"
	xdBytes, _ := hex.DecodeString(xd)
	s := decodeVec(k*N, xdBytes)
	sVector := make([][]byte, k)
	for i := 0; i < k; i++ {
		sVector[i] = make([]byte, N)
		copy(sVector[i], s[i*N:(i+1)*N])
	}

	//for _, party := range parties {
	//	party.Signature = sVector
	//}

	valid := mpc.Verify(parties, signature)

	fmt.Println("salt    : ", parties[0].Salt)
	fmt.Println("t       : ", parties[0].LittleT)
	fmt.Println("y XOR t : ", parties[0].LittleY)
	fmt.Println("Signature was valid:", valid)
}

// decodeVec decodes a byte slice into a byte slice of length n
// where n is the length of the original byte slice (to accommodate for odd n)
func decodeVec(n int, byteString []byte) []byte {
	decoded := make([]byte, n)

	for i := 0; i < n/2; i++ {
		firstNibble := byteString[i] & 0xf
		secondNibble := byteString[i] >> 4

		decoded[i*2] = firstNibble
		decoded[i*2+1] = secondNibble
	}

	// if 'n' is odd, then fix last nibble. Not second nibble present in the last byte
	if n%2 == 1 {
		decoded[n-1] = byteString[n/2] & 0xf
	}

	return decoded
}
