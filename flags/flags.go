package flags

import "flag"

type ApplicationArguments struct {
	AmountBenchmarkingSamples int
}

func GetApplicationArguments() ApplicationArguments {
	// Creating struct with empty arguments
	arguments := ApplicationArguments{}

	// Getting arguments from flags
	flag.IntVar(&arguments.AmountBenchmarkingSamples, "b", 0,
		"Decides if the implementation should be benchmarked, and the amount of samples")

	// Parsing flags
	flag.Parse()

	return arguments
}
