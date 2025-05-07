//go:build mayo3

package mpc

const paramName = "MAYO3"

const n = 118
const m = 108
const o = 10
const k = 11
const v = n - o

const shifts = k * (k + 1) / 2

const macAmount = 16

var tailF = [4]byte{8, 0, 1, 7}
