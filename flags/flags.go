package flags

import "flag"

type ApplicationArguments struct {
	AmountBenchmarkingSamples int
	AmountOfParties           int
	Threshold                 int
}

func GetApplicationArguments() ApplicationArguments {
	// Creating struct with empty arguments
	arguments := ApplicationArguments{}

	// Getting arguments from flags
	flag.IntVar(&arguments.AmountBenchmarkingSamples, "b", 0,
		"Number of samples for benchmarking (0 = no benchmarking)")

	flag.IntVar(&arguments.AmountOfParties, "n", 4,
		"Number of parties")

	flag.IntVar(&arguments.Threshold, "t", 4,
		"Threshold value")

	// Parsing flags
	flag.Parse()

	return arguments
}
