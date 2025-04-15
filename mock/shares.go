package mock

import (
	"mayo-threshold-go/mpc"
	"mayo-threshold-go/rand"
	"reflect"
)

func CreatePartiesAndSharesForEsk(esk mpc.ExpandedSecretKey, epk mpc.ExpandedPublicKey, n int) []*mpc.Party {
	// First create the empty structs
	eskShares := make([]mpc.ExpandedSecretKey, n)
	for i := 0; i < n; i++ {
		eskShares[i] = getNewExpandedSecretKey()
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
	parties := make([]*mpc.Party, n)
	for i := 0; i < n; i++ {
		party := &mpc.Party{EskShare: eskShares[i], Epk: epk}
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

func VerifyShares(esk mpc.ExpandedSecretKey, parties []*mpc.Party) bool {
	n := len(parties)
	if n == 0 {
		return false
	}

	reconstructedL := make([][][]byte, len(esk.L))
	for i := range esk.L {
		reconstructedL[i] = make([][]byte, len(esk.L[i]))
		for j := range esk.L[i] {
			reconstructedL[i][j] = make([]byte, len(esk.L[i][j]))
		}
	}

	reconstructedO := make([][]byte, len(esk.O))
	for i := range esk.O {
		reconstructedO[i] = make([]byte, len(esk.O[i]))
	}

	// Reconstruct esk.L
	for i, matrix := range esk.L {
		for j, row := range matrix {
			for k := range row {
				for _, party := range parties {
					reconstructedL[i][j][k] ^= party.EskShare.L[i][j][k] // XOR to reconstruct
				}
			}
		}
	}

	// Reconstruct esk.O
	for i, row := range esk.O {
		for j := range row {
			for _, party := range parties {
				reconstructedO[i][j] ^= party.EskShare.O[i][j] // XOR to reconstruct
			}
		}
	}

	// Compare reconstructed values with original esk
	return reflect.DeepEqual(esk.L, reconstructedL) && reflect.DeepEqual(esk.O, reconstructedO)
}
