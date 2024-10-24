package command

import (
	"fmt"

	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/syscall"
)

// Cgexec adds the current process to 0 or more specified cgroups, then
// execs the specfied command.  The format of args is:
//
//	args[0:n]   - cgroups files
//	args[n:n+1] - --
//	args[n+2:]  - command to exec and its arguments
//
// It returns an error if it failed to add itself to the requested cgroups
// or if it fails to exec the command.
func Cgexec(args []string) error {
	return CgexecDetailed(args, nil, nil)
}

// CgexecDetailed is wrapped by Cgexec and performs the same operation,
// optionally with concrete os and syscall adapters.
func CgexecDetailed(args []string, osa *os.Adapter, sa *syscall.Adapter) error {
	const DefaultPerms os.FileMode = 0644

	var (
		taskFileList []string
		commandList  []string
	)

	for i := range args {
		if args[i] == "--" {
			taskFileList = args[1:i]
			commandList = args[i+1:]
			break
		}
	}

	// If we never found --, treat all arguments as the command list
	if len(taskFileList) == 0 && len(commandList) == 0 {
		commandList = args[1:]
	}

	if len(commandList) == 0 {
		return fmt.Errorf("cgexec: no command provided")
	}

	pid := fmt.Sprintf("%d", osa.Getpid())
	for _, taskFile := range taskFileList {
		if err := osa.WriteFile(taskFile, []byte(pid), DefaultPerms); err != nil {
			return err
		}
	}

	if err := sa.Exec(commandList[0], commandList, osa.Environ()); err != nil {
		return err
	}

	// This should never happen
	return fmt.Errorf("reached end of Cgexec unexpectedly")
}
