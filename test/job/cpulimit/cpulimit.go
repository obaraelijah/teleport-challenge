package main

import (
	"fmt"
	"runtime"

	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/cgroupv1"
	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager"
)

func runTest(controllers ...cgroupv1.Controller) {
	job := jobmanager.NewJob("theOwner", "my-test", controllers,
		"/usr/bin/stress-ng",
		"--cpu",
		fmt.Sprintf("%d", runtime.NumCPU()),
		"--timeout",
		"20",
		"--times",
	)
	if err := job.Start(); err != nil {
		panic(err)
	}
	for output := range job.StderrStream().Stream() {
		fmt.Print(string(output))
	}
	fmt.Printf("\n")
}

// Sample run:
//     $ sudo go run test/job/cpulimit/cpulimit.go
//     Password:
//     Running CPU test with no cgroup constraints
//     stress-ng: info:  [1] setting to a 20 second run per stressor
//     stress-ng: info:  [1] dispatching hogs: 12 cpu
//     stress-ng: info:  [1] successful run completed in 20.01s
//     stress-ng: info:  [1] for a 20.01s run time:
//     stress-ng: info:  [1]     240.09s available CPU time
//     stress-ng: info:  [1]     239.37s user time   ( 99.70%)
//     stress-ng: info:  [1]       0.05s system time (  0.02%)
//     stress-ng: info:  [1]     239.42s total time  ( 99.72%)
//     stress-ng: info:  [1] load average: 5.83 2.88 2.22
//
//     Running CPU test with cgroup constraints at 0.5 CPU
//     stress-ng: info:  [1] setting to a 20 second run per stressor
//     stress-ng: info:  [1] dispatching hogs: 12 cpu
//     stress-ng: info:  [1] successful run completed in 20.13s
//     stress-ng: info:  [1] for a 20.13s run time:
//     stress-ng: info:  [1]     241.60s available CPU time
//     stress-ng: info:  [1]      10.08s user time   (  4.17%)
//     stress-ng: info:  [1]       0.02s system time (  0.01%)
//     stress-ng: info:  [1]      10.10s total time  (  4.18%)
//     stress-ng: info:  [1] load average: 5.47 2.99 2.27

func main() {
	fmt.Println("Running CPU test with no cgroup constraints")
	runTest()

	fmt.Println("Running CPU test with cgroup constraints at 0.5 CPU")
	runTest(cgroupv1.NewCpuController().SetCpus(0.5))
}
