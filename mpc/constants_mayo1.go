//go:build mayo1

package mpc

const n = 86
const m = 78
const o = 8
const k = 10
const v = n - o

const shifts = k * (k + 1) / 2

var tailF = [4]byte{8, 1, 1, 0}
