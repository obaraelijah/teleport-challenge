package jobmanager_test

import (
	"testing"

	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager"
	"github.com/obaraelijah/teleport-challenge/pkg/jobmanager/jobmanagertest"
	"github.com/stretchr/testify/assert"
)

func Test_JobManager_Start(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, err := jm.Start(userName1, jobName, programPath, nil)

	assert.Nil(t, err)
	assert.NotNil(t, job)
}

func Test_JobManager_DuplicateJobName_Error(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)
	_, _ = jm.Start(userName1, jobName, programPath, nil)
	job, err := jm.Start(userName1, jobName, programPath, nil)

	assert.Error(t, err)
	assert.Nil(t, job)
}

func Test_JobManager_Status_MatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	status, err := jm.Status(userName1, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, true, status.Running)
	assert.Equal(t, jobName, status.Name)
	assert.Equal(t, userName1, status.Owner)
}

func Test_JobManager_Status_NonMatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_, err := jm.Status("someOtherUser", job.Id().String())

	assert.Error(t, err)
}

func Test_JobManager_Status_Superuser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	status, err := jm.Status(jobmanager.Superuser, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, true, status.Running)
	assert.Equal(t, jobName, status.Name)
	assert.Equal(t, userName1, status.Owner)
}

func Test_JobManager_Stop_MatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = jm.Stop(userName1, job.Id().String())
	status, err := jm.Status(userName1, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, false, status.Running)
	assert.Equal(t, jobName, status.Name)
}

func Test_JobManager_Stop_NonmatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	err := jm.Stop("someOtherUser", job.Id().String())

	assert.Error(t, err)
}
func Test_JobManager_Stop_Superuser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = jm.Stop(jobmanager.Superuser, job.Id().String())

	status, err := jm.Status(userName1, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, false, status.Running)
	assert.Equal(t, jobName, status.Name)
}

func Test_JobManager_List_MatchingUser(t *testing.T) {
	const userName1 = "user1"
	const userName2 = "user2"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	_, _ = jm.Start(userName1, jobName, programPath, nil)
	_, _ = jm.Start(userName2, jobName, programPath, nil)
	jobList := jm.List(userName1)

	assert.Equal(t, 1, len(jobList))
	assert.Equal(t, true, jobList[0].Running)
	assert.Equal(t, jobName, jobList[0].Name)
}

func Test_JobManager_List_NonmatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	_, _ = jm.Start(userName1, jobName, programPath, nil)
	jobList := jm.List("someOtherUser")

	assert.Equal(t, 0, len(jobList))
}

func Test_JobManager_List_Superuser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	_, _ = jm.Start(userName1, jobName, programPath, nil)
	jobList := jm.List(jobmanager.Superuser)

	assert.Equal(t, 1, len(jobList))
	assert.Equal(t, true, jobList[0].Running)
	assert.Equal(t, jobName, jobList[0].Name)
}

func Test_JobManager_List_Superuser_MultipleUsersJobs(t *testing.T) {
	const userName1 = "user1"
	const userName2 = "user2"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	_, _ = jm.Start(userName1, jobName, programPath, nil)
	_, _ = jm.Start(userName2, jobName, programPath, nil)

	jobList := jm.List(jobmanager.Superuser)

	assert.Equal(t, 2, len(jobList))
}

func Test_JobManager_StdoutStream_MatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = job.Stop()

	stream, err := jm.StdoutStream(userName1, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, "this is standard output", string(<-stream.Stream()))
}

func Test_JobManager_StdoutStream_NonmatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = job.Stop()

	_, err := jm.StdoutStream("someOtherUser", job.Id().String())

	assert.Error(t, err)
}

func Test_JobManager_StdoutStream_Superuser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = job.Stop()

	stream, err := jm.StdoutStream(jobmanager.Superuser, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, "this is standard output", string(<-stream.Stream()))
}

func Test_JobManager_StderrStream_MatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = job.Stop()

	stream, err := jm.StderrStream(userName1, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, "this is standard error", string(<-stream.Stream()))
}

func Test_JobManager_StderrStream_NonmatchingUser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = job.Stop()

	_, err := jm.StderrStream("someOtherUser", job.Id().String())

	assert.Error(t, err)
}

func Test_JobManager_StderrStream_Superuser(t *testing.T) {
	const userName1 = "user1"
	const jobName = "user1-job"
	const programPath = "/bin/true"

	jm := jobmanager.NewManagerDetailed(jobmanagertest.NewMockJob, nil)

	job, _ := jm.Start(userName1, jobName, programPath, nil)
	_ = job.Stop()

	stream, err := jm.StderrStream(jobmanager.Superuser, job.Id().String())

	assert.Nil(t, err)
	assert.Equal(t, "this is standard error", string(<-stream.Stream()))
}
