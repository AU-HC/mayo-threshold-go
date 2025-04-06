//go:build mayo2

package mpc

const n = 81
const m = 64
const o = 17
const k = 4
const v = n - o

const shifts = k * (k + 1) / 2

var tailF = [4]byte{8, 0, 2, 8}
