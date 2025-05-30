# Threshold MAYO in Go
*Andreas Skriver Nielsen, Markus Valdemar Grønkjær Jensen, and Hans-Christian Kjeldsen*

## Overview
This project provides a **threshold implementation** of the [MAYO digital signature scheme](https://pqmayo.org/assets/specs/mayo-round2.pdf), following the design outlined in the [threshold variant proposal](https://eprint.iacr.org/2024/1960.pdf). MAYO is a multivariate signature scheme based on the Oil and Vinegar (O&V) framework, optimized for compact public keys and post-quantum security.

This threshold construction enables a set of \( n \) parties to jointly sign a message such that any subset of \( t \) or more parties can produce a valid signature, while fewer than \( t \) learn nothing about the secret key.

### Key Features
- **Threshold signing:** Distributed key generation and signing without reconstructing the full private key.
- **Post-quantum security:** Inherits the security properties of MAYO under multivariate assumptions.
- **Modular implementation:** Built in Go for portability and clarity.
- **Benchmarking support:** Evaluate performance for various \( n \), \( t \), and system sizes.

## Installation
Ensure you have Go installed: [Download Go](https://go.dev/doc/install)
Clone the repository:
```bash
$ git clone https://github.com/AU-HC/mayo-threshold-go
$ cd mayo-threshold-go
```

## Usage
The implementation is currently a command-line tool. To run key generation, threshold signing, and verification, for 5 parties with a threshold of 3 execute:

```bash
go run -tags=mayo1 main.go -n=5 -t=3
```
Note that the following flags can be set:
- `-tags` (string): The level of MAYO.
- `-n` (int): Total number of parties participating in the threshold signature scheme.
- `-t` (int): Threshold number of parties required to generate a valid signature.
- `-b` (int): Optional. Number of benchmark iterations to run. This flag will ignore the n and t flags.

## Remarks
- This is a prototype implementation of the threshold construction based on the referenced paper.
- The current focus is correctness and clarity, not optimization.
- For a version that includes active security protections, see the active-security branch.