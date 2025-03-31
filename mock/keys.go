package mock

import (
	"encoding/json"
	"mayo-threshold-go/model"
)

const eskFileName = "mock/resources/mock_esk.json"
const epkFileName = "mock/resources/mock_epk.json"

func GetExpandedKeyPair() (model.ExpandedSecretKey, model.ExpandedPublicKey) {
	var esk model.ExpandedSecretKey
	eskBytes := getBytesFromFile(eskFileName)
	if err := json.Unmarshal(eskBytes, &esk); err != nil {
		panic(err)
	}

	var epk model.ExpandedPublicKey
	epkBytes := getBytesFromFile(epkFileName)
	if err := json.Unmarshal(epkBytes, &epk); err != nil {
		panic(err)
	}

	return esk, epk
}

func getNewExpandedSecretKey() model.ExpandedSecretKey {
	var eskCopy model.ExpandedSecretKey
	eskBytes := getBytesFromFile(eskFileName)
	if err := json.Unmarshal(eskBytes, &eskCopy); err != nil {
		panic(err)
	}
	return eskCopy
}
