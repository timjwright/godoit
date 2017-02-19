package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"io/ioutil"
	"os"
	"path"
)

func TestPathMatching(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "0 30 * * * *  job name .godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "0 30 * * * *", job.Spec)
		assert.Equal(t, time.UTC, job.Timezone)
		assert.Equal(t, 0 * time.Second, job.Timeout)
		assert.Equal(t, "job name", job.Name)
		assert.Equal(t, dir + "/0 30 * * * *  job name .godoit", job.Filepath)
		assert.Equal(t, true, job.Enabled)
	})
}

func TestWithCronSpecInFile(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "job name.godoit","#:godoit cronspec 0 30 * * * *")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "0 30 * * * *", job.Spec)
		assert.Equal(t, time.UTC, job.Timezone)
		assert.Equal(t, 0 * time.Second, job.Timeout)
		assert.Equal(t, "job name", job.Name)
		assert.Equal(t, dir + "/job name.godoit", job.Filepath)
		assert.Equal(t, true, job.Enabled)
	})
}

func TestWithParamsInFile(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(
			dir,
			"job name.godoit",
			"#:godoit cronspec 0 30 * * * *",
			"#:godoit timeout 1h30m",
			"#:godoit timezone Europe/London")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "0 30 * * * *", job.Spec)
		assert.Equal(t, "Europe/London", job.Timezone.String())
		assert.Equal(t, 90 * time.Minute, job.Timeout)
		assert.Equal(t, "job name", job.Name)
		assert.Equal(t, dir + "/job name.godoit", job.Filepath)
		assert.Equal(t, true, job.Enabled)
	})
}

func TestPathMatchingWithRepalce(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "0 0%5 x x x x  job name.godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "0 0/5 * * * *", job.Spec)
		assert.Equal(t, "job name", job.Name)
		assert.Equal(t, dir + "/0 0%5 x x x x  job name.godoit", job.Filepath)
		assert.Equal(t, true, job.Enabled)
	})
}

func TestPathMatching2(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "* * * * * * TestScanRemoveJob.godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, true, job.Enabled)
	})
}

func TestPathAll(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "* * * * * * job name.godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, true, job.Enabled)
	})
}

func TestHashCommentedOutJobShouldBeDisabled(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "#0 30 * * * * job name.godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "0 30 * * * *", job.Spec)
		assert.Equal(t, time.UTC, job.Timezone)
		assert.Equal(t, 0 * time.Second, job.Timeout)
		assert.Equal(t, "job name", job.Name)
		assert.Equal(t, dir + "/#0 30 * * * * job name.godoit", job.Filepath)
		assert.Equal(t, false, job.Enabled)
	})
}

func TestCommentedOutJobShouldBeDisabled(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "--0 30 * * * * job name.godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "job name", job.Name)
		assert.Equal(t, dir + "/--0 30 * * * * job name.godoit", job.Filepath)
		assert.Equal(t, false, job.Enabled)
	})
}

func TestIncompleteTaskShouldBeNil(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "0 30 * * *.godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, false, job.Enabled)
		assert.Equal(t, "Missing cronspec", job.Errors[0])
	})
}

func TestWithoutNameShouldBeNil(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "0 30 * * * *.godoit")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "Missing cronspec", job.Errors[0])
	})
}

func TestInvalidTParams(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(
			dir,
			"test.godoit",
			"#:godoit timeout ghgh",
			"#:godoit timezone ohh",
			"#:godoit cronspec ahh")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "Invalid timeout: 'ghgh'", job.Errors[0])
		assert.Equal(t, "Invalid timezone: 'ohh'", job.Errors[1])
		assert.Equal(t, "Invalid cronspec: 'ahh'", job.Errors[2])
	})
}

func TestDuplicateCronSpec(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(
			dir,
			"0 30 * * * * test.godoit",
			"#:godoit cronspec 2 45 * * * *")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "Cronspec in filename and as comment", job.Errors[0])
	})
}

func TestUnknownParam(t *testing.T) {
	withDir(func(dir string) {
		job := createTestJob(dir, "0 30 * * * * test.godoit","#:godoit blahh")
		assert.NotNil(t, job, "Failed to parse job")
		assert.Equal(t, "Invalid parameter 'blahh'", job.Errors[0])
	})
}

type withDirFunc func(dir string)

func withDir(aFunc withDirFunc) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	aFunc(dir)
}

func createTestJob(dir string, file string, lines ...string) *Job {
	f, _ := os.Create(path.Join(dir, file))
	for _, line := range lines {
		f.WriteString(line)
		f.WriteString("\n")
	}
	f.Close()
	return ParseJobFile(dir, file)
}

