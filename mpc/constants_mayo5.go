//go:build mayo5

package mpc

const n = 154
const m = 142
const o = 12
const k = 12
const v = n - o

const shifts = k * (k + 1) / 2

var tailF = [4]byte{4, 0, 8, 1}
