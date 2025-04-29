package mpc

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const directory = "results"
const fileName = "results.json"

var amountOfPartiesToBenchmark = []int{3, 5, 10, 15}

type Result struct {
	AmountOfParties           int                       `json:"AmountOfParties"`
	ResultsForAmountOfParties ResultsForAmountOfParties `json:"ResultsForAmountOfParties"`
}

type ResultsForAmountOfParties struct {
	KeyGen          []int64 `json:"KeyGen"`
	Sign            []int64 `json:"Sign"`
	ThresholdVerify []int64 `json:"ThresholdVerify"`
}

func Benchmark(n int) (string, error) {
	message := make([]byte, 32)

	results := make([]Result, len(amountOfPartiesToBenchmark))
	for i, amountOfParties := range amountOfPartiesToBenchmark {
		context := CreateContext(amountOfParties, amountOfParties)

		// Benchmark KeyGen
		keyGenResults := make([]int64, n)
		for j := 0; j < n; j++ {
			context.PreprocessMultiplicationKeyGenTriples()
			before := time.Now()
			context.KeyGen(amountOfParties)
			duration := time.Since(before)
			keyGenResults[j] = duration.Nanoseconds()
		}

		// Benchmark Sign
		_, parties := context.KeyGen(amountOfParties)
		signResults := make([]int64, n)
		for j := 0; j < n; j++ {
			context.PreprocessMultiplicationSignTriples(AmountOfMultiplicationTriples)
			before := time.Now()
			context.Sign(message, parties)
			duration := time.Since(before)
			signResults[j] = duration.Nanoseconds()
		}

		// Benchmark ThresholdVerify
		verifyResults := make([]int64, n)
		//_ = ThresholdVerifiableSign(message, parties)
		for j := 0; j < n; j++ {
			before := time.Now()
			//ThresholdVerify(parties, sig)
			duration := time.Since(before)
			verifyResults[j] = duration.Nanoseconds()
		}

		// Create struct to contain data-points
		results[i] = Result{
			AmountOfParties: amountOfParties,
			ResultsForAmountOfParties: ResultsForAmountOfParties{
				KeyGen:          keyGenResults,
				Sign:            signResults,
				ThresholdVerify: verifyResults,
			},
		}
	}

	// Write the results to JSON in results directory
	resultsJson, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return "", err
	}
	pathToResults := fmt.Sprintf("%s/%s-%s-%s", directory, paramName, time.Now().Format("2006-01-02-15-04-05"), fileName)
	fmt.Println(pathToResults)
	err = os.WriteFile(pathToResults, resultsJson, 0644)
	if err != nil {
		return "", err
	}

	return pathToResults, nil
}
