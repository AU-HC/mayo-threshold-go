package mpc

import (
	"mayo-threshold-go/model"
)

// ThresholdVerifiableSign takes a message, parties and outputs an 'opened' signature, which can be verified using
// original MAYO, or using the Verify method.
func (c *Context) ThresholdVerifiableSign(message []byte, parties []*model.Party) model.ThresholdSignature {
	for true {
		// Steps 1-3 of sign
		c.computeM(parties, message)
		// Step 4 of sign
		c.computeY(parties)
		// Step 5 of sign
		c.localComputeA(parties)
		c.localComputeY(parties)

		// Step 6 of sign
		// ** Algorithm solve **
		// Steps 1-4 of solve
		isTFullRank := c.computeT(parties)
		if isTFullRank {
			break
		}
	}
	// Step 5 of solve
	c.computeAInverse(parties)
	// Steps 6-9 of solve
	c.computeLittleX(parties)
	// ** Algorithm solve **

	// Step 7-9 of sign
	thresholdSignature := c.computeSignature(parties)
	return thresholdSignature
}

func (c *Context) Sign(message []byte, parties []*model.Party) model.Signature {
	// Compute signature
	thresholdSignature := c.ThresholdVerifiableSign(message, parties)

	// Recover the signature
	s := c.algo.openMatrix(thresholdSignature.S)

	// Return the 'revealed' signature
	return model.Signature{
		S:    s,
		Salt: thresholdSignature.Salt,
	}
}
