package test

import (
	"fmt"
	"time"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

func BenchmarkFunctions(function1 func(), function2 func(), iterations int) {
	startFunction1 := time.Now()
	for x := 0; x < iterations; x++ {
		function1()
	}
	durationFunction1 := time.Since(startFunction1)

	startFunction2 := time.Now()
	for x := 0; x < iterations; x++ {
		function2()
	}
	durationFunction2 := time.Since(startFunction2)

	function1Milli := float64(durationFunction1.Milliseconds()) / float64(iterations)
	function2Milli := float64(durationFunction2.Milliseconds()) / float64(iterations)

	TraceLog(fmt.Sprintf(
		"Benchmark || Function 1 || Total time taken %fs, Call Duration: %fms",
		durationFunction1.Seconds(), function1Milli,
	))
	TraceLog(fmt.Sprintf(
		"Benchmark || Function 2 || Total time taken %fs, Call Duration: %fms",
		durationFunction2.Seconds(), function2Milli,
	))
}
