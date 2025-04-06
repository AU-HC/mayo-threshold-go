package mpc

// Constants from MAYO spec
const m = 64
const k = 4
const n = 81
const o = 17
const v = n - o

// Constants from Threshold MAYO
const lambda = 0

// Constants needed to do operations
const shifts = k * (k + 1) / 2

// Field arithmetic
var field *Field = InitField()
