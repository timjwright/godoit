package main

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
	"os"
	"io/ioutil"
	"path"
	"time"
)

var executions = make(map [string]int)

var executor = func(name, path string) {
	log.Printf("Executing %s: %s", name, path)
	if _,ok := executions[name]; ok {
		executions[name] = executions[name]+1
	} else {
		executions[name] = 1
	}
}

func TestScanEmptyDir(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Scan an empty directory
		assertRescanUpdates(t, jobSet, false)
		assertRescanUpdates(t, jobSet, false)
	})
}

func TestScanDirWithCommentFile(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "random_file")
		assertRescanUpdates(t, jobSet, false)
	})
}

func TestScanDirWithIgnoredFile(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "#* * * * * * job name")
		assertRescanUpdates(t, jobSet, false)
	})
}

func TestScanCreate1SecondJob(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "* * * * * * TestScanCreate1SecondJob.godoit")
		assertRescanUpdates(t, jobSet, true)
		assertJobCount(t, jobSet, 1)
		time.Sleep(time.Second * 6)
		assertRescanUpdates(t, jobSet, false)
		assertExecutions(t, "TestScanCreate1SecondJob", 5)
	})
}


func TestScanRemoveJob(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		createJob(jobSet, "* * * * * * TestScanRemoveJob.godoit")
		assertRescanUpdates(t, jobSet, true)
		assertJobCount(t, jobSet, 1)
		time.Sleep(time.Second * 2)

		removeJob(t, jobSet, "* * * * * * TestScanRemoveJob.godoit")
		assertRescanUpdates(t, jobSet, true)
		assertJobCount(t, jobSet, 0)
		time.Sleep(time.Second * 4)
		assertExecutions(t, "TestScanRemoveJob", 1)
	})
}

func TestScanSwapJobs(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		createJob(jobSet, "* * * * * * TestScanSwapJobs.godoit")
		createJob(jobSet, "* * * * * * TestScanSwapJobsX.godoit")
		assertRescanUpdates(t, jobSet, true)
		jobSet.printJobs()
		assertJobCount(t, jobSet, 2)
		time.Sleep(time.Second * 2)

		createJob(jobSet, "* * * * * * TestScanSwapJobs_Other.godoit")
		removeJob(t, jobSet, "* * * * * * TestScanSwapJobs.godoit")
		assertRescanUpdates(t, jobSet, true)
		jobSet.printJobs()
		assertJobCount(t, jobSet, 2)
		time.Sleep(time.Second * 5)
		assertExecutions(t, "TestScanSwapJobs", 1)
		assertExecutions(t, "TestScanSwapJobsX", 5)
		assertExecutions(t, "TestScanSwapJobs_Other", 3)
	})
}


func TestStatusScript(t *testing.T) {

	withJobSet(func(jobSet1 *JobSet) {
		createJob(jobSet1, "0 1 * * * * Job 1.godoit")
		createJob(jobSet1, "0 2 * * * * Job 2.godoit")
		jobSet1.Scan()

		jobSetsMap := map[string]*JobSet{
			"test_set": jobSet1,
		}

		statusFunc := StatusReporterFromScript("./test_status.sh", []string{"PATH"}, os.Stdout)
		statusFunc(jobSetsMap)
	})
}


type withJobSetFunc func(jobSet *JobSet)

func withJobSet(aFunc withJobSetFunc) {
	dir,_ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	jobSet := NewJobSet(executor, dir)
	defer jobSet.Stop()
	println(dir)
	aFunc(jobSet)
}

func createJob(jobSet *JobSet, script string) {
	f,_ := os.Create(path.Join(jobSet.directory,script))
	f.Close()
}

func removeJob(t *testing.T, jobSet *JobSet, script string) {
	if err := os.Remove(path.Join(jobSet.directory,script)); err != nil {
		t.Fatalf("Failed to remove job %s, error %s", script, err)
	}
}

func assertRescanUpdates(t *testing.T, jobSet *JobSet, expectUpdate bool) {
	r:= jobSet.Scan();
	assert.Equal(t,expectUpdate, r)
}

func assertJobCount(t *testing.T, jobSet *JobSet, expectedJobs int) {
	assert.Equal(t,expectedJobs, len(jobSet.jobs))
}

func assertExecutions(t *testing.T, name string, minCount int) {
	if actualCount,ok := executions[name]; ok {
		assert.True(
			t,actualCount >= minCount,
			"Number of executions should be at least %n but was %n", minCount, actualCount)
	} else {
		t.Fatalf("No executions found for job %s", name)
	}
}


