package main

import (
	"bufio"
	"bytes"
	"sync"

	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager"
)

func runTest() {

	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/bash",
		"-c",
		"for ((i = 0; i < 100; ++i)); do for((j = 0; j < 1000; ++j)); do echo $RANDOM; done; sleep 0.25; done",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	const numGoroutines = 100
	var buckets [numGoroutines][]byte

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineNum int) {
			for output := range job.StdoutStream().Stream() {
				buckets[goroutineNum] = append(buckets[goroutineNum], output...)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	var readers [numGoroutines]*bufio.Reader
	for i := 0; i < numGoroutines-1; i++ {
		readers[i] = bufio.NewReader(bytes.NewReader(buckets[i]))
	}

	for i := 0; i < numGoroutines-1; i++ {
	}
}

// The job generates 10000 random numbers and prints them to standard output
// The program starts 100 goroutines to consume that output.  Each goroutine
// counts and prints the number of bytes that it receives.

func main() {
	runTest()
}
