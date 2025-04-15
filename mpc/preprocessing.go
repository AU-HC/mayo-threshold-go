package mpc

type PreprocessedMultiplicationSignTriples struct {
	ComputeM         [][]ActiveTriple
	ComputeY         [][]ActiveTriple
	ComputeT1        []ActiveTriple
	ComputeT2        []ActiveTriple
	ComputeAInverse  ActiveTriple
	ComputeX1        ActiveTriple
	ComputeX2        ActiveTriple
	ComputeSignature ActiveTriple
}

func (c *Context) PreprocessMultiplicationSignTriples(amountOfTries int) {
	mTriples := make([][]ActiveTriple, amountOfTries)
	yTriples := make([][]ActiveTriple, amountOfTries)
	for i := 0; i < amountOfTries; i++ {
		mTriples[i] = c.GenerateMultiplicationActiveTriples(k, v, v, o, m)
		yTriples[i] = c.GenerateMultiplicationActiveTriples(k, v, v, k, m)
	}

	s, t := m, k*o
	c.signTriples = PreprocessedMultiplicationSignTriples{
		ComputeM:         mTriples,
		ComputeY:         yTriples,
		ComputeT1:        c.GenerateMultiplicationActiveTriples(s, t, t, t, amountOfTries),
		ComputeT2:        c.GenerateMultiplicationActiveTriples(s, s, s, t, amountOfTries),
		ComputeAInverse:  c.GenerateMultiplicationActiveTriple(t, s, s, s),
		ComputeX1:        c.GenerateMultiplicationActiveTriple(t, s, s, 1),
		ComputeX2:        c.GenerateMultiplicationActiveTriple(t, t, t, 1),
		ComputeSignature: c.GenerateMultiplicationActiveTriple(k, o, o, v),
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
