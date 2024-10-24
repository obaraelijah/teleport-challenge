package command_test

import (
	"fmt"
	"testing"

	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os/ostest"
	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/syscall"
	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/syscall/syscalltest"
	"github.com/obaraelijah/teleport-challenge/pkg/command"
	"github.com/stretchr/testify/assert"
)

func Test_Cgexec_WriteCgroupFiles_Success(t *testing.T) {
	writeFileRecorder := &ostest.WriteFileMock{}
	pidGenerator := ostest.GetpidMock(1234)

	osa := &os.Adapter{
		WriteFileFn: writeFileRecorder.WriteFile,
		GetpidFn:    pidGenerator.Getpid,
	}

	sc := &syscall.Adapter{
		ExecFn: (&syscalltest.ExecMock{}).Exec,
	}

	cgfile := "/sys/fs/cgroup/cpu/job/1e71d42d-b7e2-4f1c-893f-b16415b96e1a/tasks"

	args := []string{
		"nameOfTheTool",
		cgfile,
		"--",
		"ignored in this test",
	}

	_ = command.CgexecDetailed(args, osa, sc)

	assert.Equal(t, 1, len(writeFileRecorder.Events))
	assert.Equal(t, cgfile, writeFileRecorder.Events[0].Name)
	assert.Equal(t, fmt.Sprintf("%d", pidGenerator), string(writeFileRecorder.Events[0].Data))
}

func Test_Cgexec_WriteCgroupFiles_Failure(t *testing.T) {
	expectedError := fmt.Errorf("injected error")
	writeFileRecorder := &ostest.WriteFileMock{
		NextError: expectedError,
	}
	var pidGenerator ostest.GetpidMock
	osa := &os.Adapter{
		WriteFileFn: writeFileRecorder.WriteFile,
		GetpidFn:    pidGenerator.Getpid,
	}

	sc := &syscall.Adapter{
		ExecFn: (&syscalltest.ExecMock{}).Exec,
	}

	cgfile := "/sys/fs/cgroup/cpu/job/1e71d42d-b7e2-4f1c-893f-b16415b96e1a/tasks"

	args := []string{
		"nameOfTheTool",
		cgfile,
		"--",
		"ignored in this test",
	}

	err := command.CgexecDetailed(args, osa, sc)

	assert.Equal(t, expectedError, err)
}

func Test_Cgexec_Exec(t *testing.T) {
	env := ostest.EnvironMock{"x=y"}
	var pidGenerator ostest.GetpidMock

	osa := &os.Adapter{
		WriteFileFn: (&ostest.WriteFileMock{}).WriteFile,
		GetpidFn:    pidGenerator.Getpid,
		EnvironFn:   env.Environ,
	}

	execRecorder := &syscalltest.ExecMock{}
	sc := &syscall.Adapter{
		ExecFn: execRecorder.Exec,
	}

	commandName := "commandName"
	commandArgs := []string{"arg1", "arg2", "--", "arg3"}
	args := []string{
		"nameOfTheTool",
		"--",
		commandName,
	}
	args = append(args, commandArgs...)

	var argv []string
	argv = append(argv, commandName)
	argv = append(argv, commandArgs...)

	err := command.CgexecDetailed(args, osa, sc)

	assert.Error(t, err)
	assert.Equal(t, commandName, execRecorder.Argv0)
	assert.Equal(t, argv, execRecorder.Argv)
	assert.Equal(t, env, ostest.EnvironMock(execRecorder.Envv))
}
