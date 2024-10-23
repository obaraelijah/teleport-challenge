package main

import (
	"fmt"
	"os"

	"github.com/obaraelijah/teleport-challenge/pkg/command"
)

// main is the entrypoint of the cgexec application.  The application accepts
// arguments in the form:
//
//	cgexec [<cgtskfile> ...] -- <command> [<arg> ...]
//
// Everything before the first "--" is treated as a cgroup task file; this
// command will add itself to those cgroups.
//
// Everything after the first "--" is treated as the command and arguments to
// the program to exec.
//
// If no "--" is found, then all arguments are treated as the command and
// arguments to the program to exec.
func main() {
	if err := command.Cgexec(os.Args); err != nil {
		fmt.Printf("cgexec failed: %v", err)
	}
	// Cgexec shouldn't return in a non-error case
	os.Exit(1)
}
