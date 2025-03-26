package mock

import (
	"encoding/json"
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
