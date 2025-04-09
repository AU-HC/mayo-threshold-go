package mpc

import (
	"mayo-threshold-go/model"
)

// ThresholdVerifiableSign takes a message, parties and outputs an 'opened' signature, which can be verified using
// original MAYO, or using the Verify method.
func ThresholdVerifiableSign(message []byte, parties []*model.Party) model.ThresholdSignature {
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
		isTFullRank := computeT(parties)
		if isTFullRank {
			break
		}
	}
	// Step 5 of solve
	computeAInverse(parties)
	// Steps 6-9 of solve
	computeLittleX(parties)
	// ** Algorithm solve **

	// Step 7-9 of sign
	thresholdSignature := computeSignature(parties)
	return thresholdSignature
}

func Sign(message []byte, parties []*model.Party) model.Signature {
	// Compute signature
	thresholdSignature := ThresholdVerifiableSign(message, parties)

	// Recover the signature
	s := algo.openMatrix(thresholdSignature.S)

	// Return the 'revealed' signature
	return model.Signature{
		S:    s,
		Salt: thresholdSignature.Salt,
	}
}
