package mpc

import "mayo-threshold-go/model"

type PreprocessedMultiplicationSignTriples struct {
	ComputeM         [][]model.Triple
	ComputeY         [][]model.Triple
	ComputeT1        []model.Triple
	ComputeT2        []model.Triple
	ComputeAInverse  model.Triple
	ComputeX1        model.Triple
	ComputeX2        model.Triple
	ComputeSignature model.Triple
}

func (c *Context) PreprocessMultiplicationSignTriples(amountOfTries int) {
	mTriples := make([][]model.Triple, amountOfTries)
	yTriples := make([][]model.Triple, amountOfTries)
	for i := 0; i < amountOfTries; i++ {
		mTriples[i] = c.GenerateMultiplicationTriples(k, v, v, o, m)
		yTriples[i] = c.GenerateMultiplicationTriples(k, v, v, k, m)
	}

	s, t := m, k*o
	c.signTriples = PreprocessedMultiplicationSignTriples{
		ComputeM:         mTriples,
		ComputeY:         yTriples,
		ComputeT1:        c.GenerateMultiplicationTriples(s, t, t, t, amountOfTries),
		ComputeT2:        c.GenerateMultiplicationTriples(s, s, s, t, amountOfTries),
		ComputeAInverse:  c.GenerateMultiplicationTriple(t, s, s, s),
		ComputeX1:        c.GenerateMultiplicationTriple(t, s, s, 1),
		ComputeX2:        c.GenerateMultiplicationTriple(t, t, t, 1),
		ComputeSignature: c.GenerateMultiplicationTriple(k, o, o, v),
	}
}

type PreprocessedMultiplicationKeyGenTriples struct {
	TriplesStep4 []model.Triple
}

func (c *Context) PreprocessMultiplicationKeyGenTriples() {
	c.keygenTriples = PreprocessedMultiplicationKeyGenTriples{
		TriplesStep4: c.GenerateMultiplicationTriples(o, v, v, o, m),
	}
}
