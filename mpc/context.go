package mpc

type Context struct {
	algo          SecretSharingAlgo
	f             *Field
	signTriples   PreprocessedMultiplicationSignTriples
	keygenTriples PreprocessedMultiplicationKeyGenTriples
}

func CreateContext(n, t int) *Context {
	if t == n {
		return &Context{
			algo: CreateAdditive(n),
			f:    InitField(),
		}
	} else {
		return &Context{
			algo: CreateShamir(n, t),
			f:    InitField(),
		}
	}
}
