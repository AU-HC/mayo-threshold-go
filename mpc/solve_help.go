package mpc

import (
	"fmt"
)

func (c *Context) computeT(parties []*Party, iteration int) bool {
	s := len(parties[0].A.shares)
	t := len(parties[0].A.shares[0])

	SShares := c.algo.createSharesForRandomMatrix(t, t)
	RShares := c.algo.createSharesForRandomMatrix(s, s)
	for partyNumber, party := range parties {
		party.S = SShares[partyNumber]
		party.R = RShares[partyNumber]
	}

	// Compute [A * SShares] = [A] * [SShares]
	dShares := make([]MatrixShare, len(parties))
	eShares := make([]MatrixShare, len(parties))
	for partyNumber, party := range parties {
		ai := c.signTriples.ComputeT1[iteration].A[partyNumber]
		bi := c.signTriples.ComputeT1[iteration].B[partyNumber]
		di := AddMatrixShares(party.A, ai)
		ei := AddMatrixShares(party.S, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	ATimesSShares := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeT1[iteration], dShares, eShares)

	// Compute [T] = [R] * [A * SShares]
	dShares = make([]MatrixShare, len(parties))
	eShares = make([]MatrixShare, len(parties))
	for partyNumber, party := range parties {
		ai := c.signTriples.ComputeT2[iteration].A[partyNumber]
		bi := c.signTriples.ComputeT2[iteration].B[partyNumber]
		di := AddMatrixShares(party.R, ai)
		ei := AddMatrixShares(ATimesSShares[partyNumber], bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	TShares := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeT2[iteration], dShares, eShares)

	// Open T and check rank
	T, err := c.algo.authenticatedOpenMatrix(TShares)
	if err != nil {
		panic(err)
	}

	for _, party := range parties {
		party.T = T
	}

	return isFullRank(T)
}

func (c *Context) computeAInverse(parties []*Party) {
	dShares := make([]MatrixShare, len(parties))
	eShares := make([]MatrixShare, len(parties))

	// Compute locally
	for partyNumber, party := range parties {
		ai := c.signTriples.ComputeAInverse.A[partyNumber]
		bi := c.signTriples.ComputeAInverse.B[partyNumber]
		TInverse := computeRightInverse(party.T)

		di := AddMatrixShares(MulPublicRight(party.S, TInverse), ai)
		ei := AddMatrixShares(party.R, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}

	// Open d, e and compute locally
	zShares := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeAInverse, dShares, eShares)
	for partyNumber, party := range parties {
		party.AInverse = zShares[partyNumber]
	}
}

func (c *Context) computeLittleX(parties []*Party) {
	s := len(parties[0].A.shares)
	t := len(parties[0].A.shares[0])

	basis := computeBasisOfKernel(parties[0].T)

	if len(basis) != t-s {
		panic(fmt.Errorf("length of basis is incorrect"))
	}

	zShares := c.algo.createSharesForRandomMatrix(t-s, 1)
	for partyNumber, party := range parties {
		z := createEmptyMatrixShare(t, 1)

		zVector := zShares[partyNumber]

		for index := 0; index < len(zVector.shares); index++ {
			shareElem := zVector.shares[index][0]
			AddMatrices(z.shares, vectorToMatrix(MultiplyVecConstant(shareElem, basis[index])))

			for i := 0; i < macAmount; i++ {
				gammaElem := zVector.gammas[i][index][0]
				AddMatrices(z.gammas[i], vectorToMatrix(MultiplyVecConstant(gammaElem, basis[index])))
			}
		}

		party.Z = z
	}

	// Compute [A^-1] * [b]
	dShares := make([]MatrixShare, len(parties))
	eShares := make([]MatrixShare, len(parties))
	for partyNumber, party := range parties {
		ai := c.signTriples.ComputeX1.A[partyNumber]
		bi := c.signTriples.ComputeX1.B[partyNumber]
		di := AddMatrixShares(party.AInverse, ai)
		ei := AddMatrixShares(party.LittleY, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	AInvTimesB := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeX1, dShares, eShares)

	// Compute [SShares] * [z]
	dShares = make([]MatrixShare, len(parties))
	eShares = make([]MatrixShare, len(parties))
	for partyNumber, party := range parties {
		ai := c.signTriples.ComputeX2.A[partyNumber]
		bi := c.signTriples.ComputeX2.B[partyNumber]
		di := AddMatrixShares(party.S, ai)
		ei := AddMatrixShares(party.Z, bi)

		dShares[partyNumber] = di
		eShares[partyNumber] = ei
	}
	STimesZ := c.activeMultiplicationProtocol(parties, c.signTriples.ComputeX2, dShares, eShares)

	// [x] = [A^-1] * [b] + [SShares] * [z]
	for i, party := range parties {
		party.LittleX = AddMatrixShares(AInvTimesB[i], STimesZ[i])
	}
}
