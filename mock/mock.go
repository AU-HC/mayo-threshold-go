package mock

import (
	"encoding/json"
	"io"
	"os"
)

const eskFileName = "mock/resources/mock_esk.json"
const epkFileName = "mock/resources/mock_epk.json"

func GetExpandedKeyPair() (ExpandedSecretKey, ExpandedPublicKey) {
	var esk ExpandedSecretKey
	eskBytes := readFileAndReturnBytes(eskFileName)
	if err := json.Unmarshal(eskBytes, &esk); err != nil {
		panic(err)
	}

	var epk ExpandedPublicKey
	epkBytes := readFileAndReturnBytes(epkFileName)
	if err := json.Unmarshal(epkBytes, &esk); err != nil {
		panic(err)
	}

	return esk, epk
}

func readFileAndReturnBytes(filename string) []byte {
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
