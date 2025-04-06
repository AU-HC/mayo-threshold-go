//go:build mayo3

package mpc

const n = 118
const m = 108
const o = 10
const k = 11
const v = n - o

const shifts = k * (k + 1) / 2

var tailF = [4]byte{8, 0, 1, 7}
