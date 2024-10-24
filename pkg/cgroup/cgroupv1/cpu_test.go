package cgroupv1_test

import (
	"fmt"
	"testing"

	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os"
	"github.com/obaraelijah/teleport-challenge/pkg/adaptation/os/ostest"
	"github.com/obaraelijah/teleport-challenge/pkg/cgroup/cgroupv1"
	"github.com/stretchr/testify/assert"
)

func Test_cpu_Apply(t *testing.T) {
	path := "/sys/fs/cgroup/jobs/889f7cc2-9935-4773-aaa1-b94478abc923"
	writeRecorder := ostest.WriteFileMock{}
	adapter := &os.Adapter{
		WriteFileFn: writeRecorder.WriteFile,
	}

	cpu := cgroupv1.NewCpuControllerDetailed(adapter).SetCpus(2.0)
	cpu.Apply(path)

	assert.Equal(t, 2, len(writeRecorder.Events))
	assert.Equal(t, fmt.Sprintf("%s/%s", path, cgroupv1.CpuPeriodFilename), writeRecorder.Events[0].Name)
	assert.Equal(t, []byte("100000"), writeRecorder.Events[0].Data)

	assert.Equal(t, fmt.Sprintf("%s/%s", path, cgroupv1.CpuQuotaFilename), writeRecorder.Events[1].Name)
	assert.Equal(t, []byte("200000"), writeRecorder.Events[1].Data)
}
