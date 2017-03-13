package main

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
	"os"
	"io/ioutil"
	"path"
	"time"
	"sync"
)

var executions = make(map [string]int)
var lock sync.RWMutex

var executor = func(name, path string, timeout time.Duration) {
	lock.Lock()
	defer  lock.Unlock()

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

func TestScanDirWithCommentedOutFile(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "--* * * * * * TestScanDirWithCommentedOutFile.godoit")
		assertRescanUpdates(t, jobSet, true)
		time.Sleep(time.Second * 3)
		assertRescanUpdates(t, jobSet, false)
		assertNoExecutions(t, "TestScanDirWithCommentedOutFile")
	})
}

func TestScanDirWithCommentedOutFile2(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "#* * * * * * TestScanDirWithCommentedOutFile2.godoit")
		assertRescanUpdates(t, jobSet, true)
		time.Sleep(time.Second * 3)
		assertRescanUpdates(t, jobSet, false)
		assertNoExecutions(t, "TestScanDirWithCommentedOutFile2")
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

func TestScanCreateJobWithTimeZone(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "TestScanCreateJobWithTimeZone.godoit","#:godoit cronspec * * * * *","#:godoit timezone Europe/Paris")
		assertRescanUpdates(t, jobSet, true)
		time.Sleep(time.Second * 6)
		assertRescanUpdates(t, jobSet, false)
		assertExecutions(t, "TestScanCreateJobWithTimeZone", 5)
	})
}

func TestScanCreateJobWithTimeZoneAndUTCJob(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "aaTestScanCreateJobWithTimeZone_utc.godoit","#:godoit cronspec * * * * *")
		createJob(jobSet, "zzTestScanCreateJobWithTimeZone.godoit","#:godoit cronspec * * * * *","#:godoit timezone Europe/Paris")
		createJob(jobSet, "bbTestScanCreateJobWithTimeZone_utc.godoit","#:godoit cronspec * * * * *")
		createJob(jobSet, "yyTestScanCreateJobWithTimeZone.godoit","#:godoit cronspec * * * * *","#:godoit timezone Europe/Paris")
		assertRescanUpdates(t, jobSet, true)
		time.Sleep(time.Second * 6)
		assertRescanUpdates(t, jobSet, false)
		assertExecutions(t, "zzTestScanCreateJobWithTimeZone", 5)
		assertExecutions(t, "aaTestScanCreateJobWithTimeZone_utc", 5)
		assertExecutions(t, "yyTestScanCreateJobWithTimeZone", 5)
		assertExecutions(t, "bbTestScanCreateJobWithTimeZone_utc", 5)
	})
}

func TestScanCreateJobWithTimeout(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "TestScanCreateJobWithTimeout.godoit","#:godoit cronspec * * * * *","#:godoit timeout 3s")
		assertRescanUpdates(t, jobSet, true)
		time.Sleep(time.Second * 6)
		assertRescanUpdates(t, jobSet, false)
		assertExecutions(t, "TestScanCreateJobWithTimeout", 5)
	})
}

func TestScanCreateJobWithSpecInFile(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "TestScanCreateJobWithSpecInFile.godoit","#:godoit cronspec * * * * * *")
		assertRescanUpdates(t, jobSet, true)
		assertJobCount(t, jobSet, 1)
		time.Sleep(time.Second * 6)
		assertRescanUpdates(t, jobSet, false)
		assertExecutions(t, "TestScanCreateJobWithSpecInFile", 5)
	})
}

func TestReScanJobWithUpdatedSpecInFile(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "TestReScanJobWithUpdatedSpecInFile.godoit","#:godoit cronspec 0 0 12 * * *")
		assertRescanUpdates(t, jobSet, true)
		assertJobCount(t, jobSet, 1)
		time.Sleep(time.Second * 3)
		// Should generate no executions - or 1 if we run right on 12:00!
		assertRescanUpdates(t, jobSet, false)
		// Now update the file
		createJob(jobSet, "TestReScanJobWithUpdatedSpecInFile.godoit","#:godoit cronspec * * * * * *")
		assertRescanUpdates(t, jobSet, true)
		time.Sleep(time.Second * 6)
		assertExecutions(t, "TestReScanJobWithUpdatedSpecInFile", 6)
	})
}

func TestScanCreateJobWithInvalidSpecInFile(t *testing.T) {
	withJobSet(func(jobSet *JobSet) {
		// Commened out job does not create a job
		createJob(jobSet, "TestScanCreateJobWithInvalidSpecInFile.godoit","#:godoit cronspec * * * *")
		assertRescanUpdates(t, jobSet, true)
		time.Sleep(time.Second * 3)
		assertRescanUpdates(t, jobSet, false)
		assertNoExecutions(t, "TestScanCreateJobWithInvalidSpecInFile")
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
		createJob(jobSet1, "0 2 * * * * Job 2.godoit", "#:godoit timeout 10s", "#:godoit timezone Europe/London")
		createJob(jobSet1, "broken job.godoit", "#:godoit wrongparam")
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

func createJob(jobSet *JobSet, script string, lines ...string) {
	f,_ := os.Create(path.Join(jobSet.directory,script))
	for _,line := range lines {
		f.WriteString(line)
		f.WriteString("\n")
	}
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
	lock.RLock()
	defer  lock.RUnlock()

	if actualCount,ok := executions[name]; ok {
		assert.True(
			t,actualCount >= minCount,
			"Number of executions should be at least %n but was %n", minCount, actualCount)
	} else {
		assert.Fail(t,"No executions found for job %s", name)
	}
}

func assertNoExecutions(t *testing.T, name string) {
	lock.RLock()
	defer  lock.RUnlock()

	if actualCount,ok := executions[name]; ok {
		assert.Fail(t, "Expected no executions of %s but found %n", name, actualCount)
	}
}


