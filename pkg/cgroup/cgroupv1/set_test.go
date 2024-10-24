package cgroupv1_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os/ostest"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/cgroupv1"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/cgroupv1/cgroupv1test"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Create_Success(t *testing.T) {
	jobID := uuid.MustParse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	mkdirAllRecorder := ostest.MkdirAllMock{}
	removeRecorder := ostest.RemoveMock{}

	adapter := &os.Adapter{
		MkdirAllFn: mkdirAllRecorder.MkdirAll,
		RemoveFn:   removeRecorder.Remove,
	}

	controller := &cgroupv1test.ControllerMock{ControllerName: "nil"}
	set := cgroupv1.NewSetDetailed(adapter, cgroupv1.DefaultBasePath, jobID, controller)

	err := set.Create()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(mkdirAllRecorder.Events))
	assert.Equal(t, 0, len(removeRecorder.Events))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s",
			cgroupv1.DefaultBasePath,
			controller.Name(),
			jobID.String(),
		),
		mkdirAllRecorder.Events[0].Path)
}

func Test_Set_Create_Failure(t *testing.T) {
	jobID := uuid.MustParse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	mkdirAllRecorder := ostest.MkdirAllMock{}
	removeRecorder := ostest.RemoveMock{}

	adapter := &os.Adapter{
		MkdirAllFn: mkdirAllRecorder.MkdirAll,
		RemoveFn:   removeRecorder.Remove,
	}

	expectedError := fmt.Errorf("injected error")
	controller := &cgroupv1test.ControllerMock{
		ControllerName:   "nil",
		ApplyReturnValue: expectedError,
	}
	set := cgroupv1.NewSetDetailed(adapter, cgroupv1.DefaultBasePath, jobID, controller)

	err := set.Create()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, len(removeRecorder.Events))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s",
			cgroupv1.DefaultBasePath,
			controller.Name(),
			jobID.String(),
		),
		removeRecorder.Events[0].Path)
}

func Test_Set_Destroy_Success(t *testing.T) {
	jobID := uuid.MustParse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	removeRecorder := ostest.RemoveMock{}

	adapter := &os.Adapter{
		RemoveFn: removeRecorder.Remove,
	}

	controller := &cgroupv1test.ControllerMock{ControllerName: "nil"}
	set := cgroupv1.NewSetDetailed(adapter, cgroupv1.DefaultBasePath, jobID, controller)

	err := set.Destroy()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(removeRecorder.Events))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s",
			cgroupv1.DefaultBasePath,
			controller.Name(),
			jobID.String(),
		),
		removeRecorder.Events[0].Path)
}

func Test_Set_Destroy_Failure(t *testing.T) {
	jobID := uuid.MustParse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")
	injectedError := fmt.Errorf("injected error")
	removeRecorder := ostest.RemoveMock{
		NextError: injectedError,
	}

	adapter := &os.Adapter{
		RemoveFn: removeRecorder.Remove,
	}

	controller := &cgroupv1test.ControllerMock{ControllerName: "nil"}
	set := cgroupv1.NewSetDetailed(adapter, cgroupv1.DefaultBasePath, jobID, controller)

	err := set.Destroy()

	assert.Error(t, err)
	assert.Equal(t, 1, len(removeRecorder.Events))
}

func Test_Set_TaskFiles(t *testing.T) {
	jobID := uuid.MustParse("0b5183b8-b572-49c7-90c4-fffc775b7d7b")

	controller := &cgroupv1test.ControllerMock{ControllerName: "nil"}
	set := cgroupv1.NewSet(jobID, controller)

	taskFiles := set.TaskFiles()

	assert.Equal(t, 1, len(taskFiles))
	assert.Equal(t,
		fmt.Sprintf("%s/%s/jobs/%s/tasks",
			cgroupv1.DefaultBasePath,
			controller.Name(),
			jobID.String(),
		),
		taskFiles[0])
}
