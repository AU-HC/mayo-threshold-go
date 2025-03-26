package mock

import (
	"encoding/json"
	"io"
	"mayo-threshold-go/rand"
	"os"
)

const eskFileName = "mock/resources/mock_esk.json"
const epkFileName = "mock/resources/mock_epk.json"

func GetExpandedKeyPair() (ExpandedSecretKey, ExpandedPublicKey) {
	var esk ExpandedSecretKey
	eskBytes := getBytesFromFile(eskFileName)
	if err := json.Unmarshal(eskBytes, &esk); err != nil {
		panic(err)
	}

	var epk ExpandedPublicKey
	epkBytes := getBytesFromFile(epkFileName)
	if err := json.Unmarshal(epkBytes, &epk); err != nil {
		panic(err)
	}

	return esk, epk
}

func getBytesFromFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return bytes
}

func CreatePartiesAndAddShares(esk ExpandedSecretKey, epk ExpandedPublicKey, n int) []Party {
	// First create the empty structs
	eskShares := make([]ExpandedSecretKey, n)
	for i := 0; i < n; i++ {
		eskShares[i] = ExpandedSecretKey{
			P1: esk.P1,
			L:  esk.L,
			O:  esk.O,
		}
	}

	// Create the shares of the esk L
	for i, matrix := range esk.L {
		for j, row := range matrix {
			for k, element := range row {
				shares := generateSharesForElement(n, element)
				for p, share := range shares {
					eskShares[p].L[i][j][k] = share
				}
			}
		}
	}

	// Create the shares of the esk O
	for i, row := range esk.O {
		for j, element := range row {
			shares := generateSharesForElement(n, element)
			for p, share := range shares {
				eskShares[p].O[i][j] = share
			}
		}
	}

	// Set the epk and the esk shares
	parties := make([]Party, n)
	for i := 0; i < n; i++ {
		party := Party{EskShare: eskShares[i], Epk: epk}
		parties[i] = party
	}
	return parties
}

func generateSharesForElement(n int, element byte) []byte {
	shares := make([]byte, n)
	var sharesSum byte

	for l := 0; l < n-1; l++ { // sample shares for n-1 parties
		share := rand.SampleFieldElement()
		shares[l] = share
		sharesSum ^= share
	}
	shares[n-1] = element ^ sharesSum // then the last share is given by the n-1 shares

	return shares
}
