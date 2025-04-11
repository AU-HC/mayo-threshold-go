package mpc

type Context struct {
	algo SecretSharingAlgo
	f    *Field
}

func CreateContext(amountOfParties, t int) *Context {
	if t == amountOfParties {
		return &Context{
			algo: &Additive{n: amountOfParties},
			f:    InitField(),
		}
	} else {
		return &Context{
			algo: &Shamir{n: amountOfParties, t: t},
			f:    InitField(),
		}
	}
}
