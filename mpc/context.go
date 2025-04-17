package mpc

import "go/types"

type Context struct {
	algo          SecretSharingAlgo
	f             *Field
	signTriples   PreprocessedMultiplicationSignTriples
	keygenTriples PreprocessedMultiplicationKeyGenTriples
	verifyTriples types.Nil // TODO: Add
}

func CreateContext(n, t int) *Context {
	if t == n {
		return &Context{
			algo: &Additive{n: n},
			f:    InitField(),
		}
	} else {
		return &Context{
			algo: &Shamir{n: n, t: t},
			f:    InitField(),
		}
	}
}
