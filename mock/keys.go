package mock

import (
	"encoding/json"
	"mayo-threshold-go/mpc"
)

const eskFileName = "mock/resources/mock_esk.json"
const epkFileName = "mock/resources/mock_epk.json"

func GetExpandedKeyPair() (mpc.ExpandedSecretKey, mpc.ExpandedPublicKey) {
	var esk mpc.ExpandedSecretKey
	eskBytes := getBytesFromFile(eskFileName)
	if err := json.Unmarshal(eskBytes, &esk); err != nil {
		panic(err)
	}

	var epk mpc.ExpandedPublicKey
	epkBytes := getBytesFromFile(epkFileName)
	if err := json.Unmarshal(epkBytes, &epk); err != nil {
		panic(err)
	}

	return esk, epk
}

func getNewExpandedSecretKey() mpc.ExpandedSecretKey {
	var eskCopy mpc.ExpandedSecretKey
	eskBytes := getBytesFromFile(eskFileName)
	if err := json.Unmarshal(eskBytes, &eskCopy); err != nil {
		panic(err)
	}
	return eskCopy
}
