package mpc

import (
	"mayo-threshold-go/model"
)

func Sign(message []byte, parties []*model.Party) model.Signature {
	for true {
		// Steps 1-3 of sign
		computeM(parties, message)
		// Step 4 of sign
		computeY(parties)
		// Step 5 of sign
		localComputeA(parties)
		localComputeY(parties)

		// Step 6 of sign
		// ** Algorithm solve **
		// Steps 1-4 of solve
		isRankDefect := computeT(parties)
		if !isRankDefect {
			break
		}
	}
	// Step 5 of solve
	computeAInverse(parties)
	// Steps 6-9 of solve
	computeLittleX(parties)
	// ** Algorithm solve **

	// Step 7-9 of sign
	signature := computeSignature(parties)
	return signature
}
