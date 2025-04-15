package mpc

type PreprocessedMultiplicationSignTriples struct {
	ComputeM         [][]Triple
	ComputeY         [][]Triple
	ComputeT1        []Triple
	ComputeT2        []Triple
	ComputeAInverse  Triple
	ComputeX1        Triple
	ComputeX2        Triple
	ComputeSignature Triple
}

func (c *Context) PreprocessMultiplicationSignTriples(amountOfTries int) {
	mTriples := make([][]Triple, amountOfTries)
	yTriples := make([][]Triple, amountOfTries)
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
	TriplesStep4 []ActiveTriple
}

func (c *Context) PreprocessMultiplicationKeyGenTriples() {
	c.keygenTriples = PreprocessedMultiplicationKeyGenTriples{
		TriplesStep4: c.GenerateMultiplicationActiveTriples(o, v, v, o, m),
	}
}
