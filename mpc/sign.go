package mpc

const AmountOfMultiplicationTriples = 4

// ThresholdVerifiableSignAPI takes a message, parties and outputs an 'opened' signature, which can be verified using
// original MAYO, or using the Verify method. Note that this method also preprocesses the multiplication triples needed
// for an execution of the protocol.
func (c *Context) ThresholdVerifiableSignAPI(message []byte, parties []*Party) ThresholdSignature {
	c.PreprocessMultiplicationSignTriples(AmountOfMultiplicationTriples)
	thresholdSignature := c.ThresholdVerifiableSign(message, parties)
	return thresholdSignature
}

// ThresholdVerifiableSign takes a message, parties and outputs an 'opened' signature, which can be verified using
// original MAYO, or using the Verify method.
func (c *Context) ThresholdVerifiableSign(message []byte, parties []*Party) ThresholdSignature {
	iteration := 0

	for true {
		// Steps 1-3 of sign
		c.computeM(parties, message, iteration)
		// Step 4 of sign
		c.computeY(parties, iteration)
		// Step 5 of sign
		c.localComputeA(parties)
		c.localComputeY(parties)

		// Step 6 of sign
		// ** Algorithm solve **
		// Steps 1-4 of solve
		isTFullRank := c.computeT(parties, iteration)
		if isTFullRank {
			break
		}
		iteration++
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

// SignAPI Also preprocesses the multiplication triples needed for an execution of the protocol.
func (c *Context) SignAPI(message []byte, parties []*Party) Signature {
	// Compute signature
	c.PreprocessMultiplicationSignTriples(AmountOfMultiplicationTriples)
	thresholdSignature := c.ThresholdVerifiableSign(message, parties)

	// Recover the signature
	s, err := c.algo.authenticatedOpenMatrix(thresholdSignature.S)
	if err != nil {
		panic(err)
	}

	// Return the 'revealed' signature
	return Signature{
		S:    s,
		Salt: thresholdSignature.Salt,
	}
}

func (c *Context) Sign(message []byte, parties []*Party) Signature {
	// Compute signature
	thresholdSignature := c.ThresholdVerifiableSign(message, parties)

	// Recover the signature
	s, err := c.algo.authenticatedOpenMatrix(thresholdSignature.S)
	if err != nil {
		panic(err)
	}

	// Return the 'revealed' signature
	return Signature{
		S:    s,
		Salt: thresholdSignature.Salt,
	}
}
