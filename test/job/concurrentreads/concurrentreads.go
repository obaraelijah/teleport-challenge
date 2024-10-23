package main

import (
	"fmt"
	"sync"

	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager"
)

func runTest() {
	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/bash",
		"-c",
		"for ((i = 0; i < 10000; ++i)); do echo $RANDOM; done",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			count := 0
			for output := range job.StdoutStream().Stream() {
				count += len(output)
			}
			fmt.Printf("%d\n", count)
			wg.Done()
		}()
	}

	wg.Wait()
}

// The job generates 10000 random numbers and prints them to standard output
// The program starts 100 goroutines to consume that output.  Each goroutine
// counts and prints the number of bytes that it receives.
func main() {
	runTest()
}
