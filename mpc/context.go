package mpc

import "go/types"

type Context struct {
	algo          SecretSharingAlgo
	n             int
	f             *Field
	signTriples   PreprocessedMultiplicationSignTriples
	keygenTriples PreprocessedMultiplicationKeyGenTriples
	verifyTriples types.Nil // TODO: Add
}

func CreateContext(amountOfParties, t int) *Context {
	if t == amountOfParties {
		return &Context{
			n:    amountOfParties,
			algo: &Additive{n: amountOfParties},
			f:    InitField(),
		}
	} else {
		return &Context{
			n:    amountOfParties,
			algo: &Shamir{n: amountOfParties, t: t},
			f:    InitField(),
		}
	}
}
