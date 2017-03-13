package main

import (
	"log"
	"github.com/robfig/cron"
	"path/filepath"
	"io/ioutil"
	"time"
	"os"
	"strings"
)

type JobSet struct {
	executor JobExecutor
	directory string
	jobs map [string]Job
	crons map [string]*cron.Cron
}

func NewJobSet(executor JobExecutor, directory string) *JobSet {
	return &JobSet{executor, directory, make(map[string]Job), make(map[string]*cron.Cron)}
}

func (jobSet *JobSet) Stop() {
	for key, cron := range jobSet.crons {
		log.Printf("  Stopping jobs in directory, %s", jobSet.directory)
		cron.Stop()
		delete(jobSet.crons, key)
	}
}

func (jobSet *JobSet) Scan() bool {
	updated := false

	// Scan for any new jobs
	files, _ := ioutil.ReadDir(jobSet.directory)
	foundFiles := make(map[string]bool)
	for _,file := range files {
		filename := file.Name()
		foundFiles[filename] = true
		if isGodoitFile(file) && shouldParseJob(file, jobSet) {
			job := ParseJobFile(jobSet.directory, filename)
			if job != nil {
				updated = true
				jobSet.jobs[filename] = *job
			}
		}
	}

	// Remove any old jobs
	for filename,_ := range jobSet.jobs {
		if _,ok := foundFiles[filename]; ! ok {
			updated = true
			delete(jobSet.jobs,filename)
		}
	}

	// Setup the cron
	if updated {
		jobSet.setupCron()
	}
	return updated
}

func isGodoitFile(file os.FileInfo) bool {
	return ! file.IsDir() && strings.HasSuffix(file.Name(), GodoitFileSuffix)
}

func shouldParseJob(file os.FileInfo, jobSet *JobSet) bool {
	job,hasJob := jobSet.jobs[file.Name()]
	if hasJob {
		// If the job is know - parse if the modification time has changed
		return job.UpdateTime != file.ModTime()
	} else {
		// If the job is not known - parse it
		return true
	}

}

func (jobSet *JobSet) setupCron() {
	jobSet.Stop()


	log.Printf("  Starting crons for %s", jobSet.directory)
	for _,job := range jobSet.jobs  {
		if job.Enabled {
			addJob(jobSet.cronForLocation(job.Timezone), jobSet.executor, job)
		}
	}
	for timezone, cron := range jobSet.crons {
		log.Printf("    Timezone: %s", timezone)
		cron.Start()
	}
}

func (jobSet *JobSet) cronForLocation(location *time.Location) *cron.Cron {
	timezoneName := location.String()
	if locationCron, ok := jobSet.crons[timezoneName]; ok {
		return locationCron
	} else {
		var newCron = cron.NewWithLocation(location)
		jobSet.crons[timezoneName] = newCron
		return newCron
	}
}


func addJob(cron *cron.Cron, executor JobExecutor, job Job) {
	cron.AddFunc(job.Spec, func() {runJob(executor, job)})
}

func runJob(executor JobExecutor, job Job) {
	log.Printf("Running job %s (%s) Timeout: %s", job.Name, filepath.Dir(job.Filepath), timeoutString(job.Timeout))
	executor(job.Name, job.Filepath, job.Timeout)
}

func (jobSet *JobSet) printJobs() {
	log.Printf("Jobs in %s", jobSet.directory)
	if len(jobSet.jobs) == 0 {
		log.Printf("  -- no jobs --")
	}
	for _,job := range jobSet.jobs {
		log.Printf(
			"  %s (%s): %s (Timeout: %s, Enabled: %t)",
			job.Spec,
			job.Timezone.String(),
			job.Name,
			timeoutString(job.Timeout),
			job.Enabled)
	}
	log.Printf("")
}

func timeoutString(timeout time.Duration) string {
	if timeout == 0 {
		return "None"
	} else {
		return timeout.String()
	}
}
