package main

import (
	"fmt"

	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager"
)

func runTest() {
	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/usr/bin/id",
	)

	if err := job.Start(); err != nil {
		panic(err)
	}

	for output := range job.StdoutStream().Stream() {
		fmt.Print(string(output))
	}
	fmt.Printf("\n")
}

// Sample run:
//     Determining the job's PID in its namespace
//     uid=0(root) gid=0(root) groups=0(root)

func main() {
	fmt.Println("Determining the job's PID in its namespace")
	runTest()
}
