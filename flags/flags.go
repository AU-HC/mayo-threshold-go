package flags

import "flag"

type ApplicationArguments struct {
	AmountBenchmarkingSamples, NumberOfParties, Threshold int
}

func GetApplicationArguments() ApplicationArguments {
	// Creating struct with empty arguments
	arguments := ApplicationArguments{}

	// Getting arguments from flags
	flag.IntVar(&arguments.AmountBenchmarkingSamples, "b", 0,
		"Decides if the implementation should be benchmarked, and the amount of samples")
	flag.IntVar(&arguments.NumberOfParties, "NumberOfParties", 3,
		"Decides how many parties should participate in generating the signature")
	flag.IntVar(&arguments.Threshold, "Threshold", 3,
		"Decides the Threshold of the secret sharing algorithm. If Threshold < NumberOfParties Shamir will be used")

	// Parsing flags
	flag.Parse()

	return arguments
}
