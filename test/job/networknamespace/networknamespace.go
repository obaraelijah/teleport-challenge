package main

import (
	"fmt"

	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager"
)

func runTest() {
	job := jobmanager.NewJob("theOwner", "my-test", nil,
		"/bin/ip",
		"link",
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
//
//	$ sudo go run networknamespace.go
//	Running test to list all network interfaces avaialble to a job
//	1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN mode DEFAULT group default qlen 1000
//	    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
//	2: sit0@NONE: <NOARP> mtu 1480 qdisc noop state DOWN mode DEFAULT group default qlen 1000
//	    link/sit 0.0.0.0 brd 0.0.0.0
func main() {
	fmt.Println("Running test to list all network interfaces avaialble to a job")
	runTest()
}
