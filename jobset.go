package main

import (
	"log"
	"github.com/robfig/cron"
	"path/filepath"
	"io/ioutil"
)

type JobSet struct {
	executor JobExecutor
	directory string
	jobs map [string]Job
	cron *cron.Cron
}

func NewJobSet(executor JobExecutor, directory string) *JobSet {
	return &JobSet{executor, directory, make(map[string]Job), nil}
}

func (jobSet *JobSet) Stop() {
	if ( jobSet.cron != nil ) {
		log.Printf("  Stopping jobs in directory, %s", jobSet.directory)
		jobSet.cron.Stop()
	}
}

func (jobSet *JobSet) Scan() bool {
	updated := false

	// Scan for any new jobs
	files, _ := ioutil.ReadDir(jobSet.directory)
	filenames := make(map[string]bool)
	for _,file := range files {
		if ( ! file.IsDir() ) {
			filename := file.Name()
			job := ParseJobFile(jobSet.directory, filename)
			if job != nil {
				filenames[file.Name()] = true
				if _,ok := jobSet.jobs[filename]; !ok {
					updated = true
					jobSet.jobs[filename] = *job
				}
			}
		}
	}

	// Remove any old jobs
	for filename,_ := range jobSet.jobs {
		if _,ok := filenames[filename]; ! ok {
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

func (jobSet *JobSet) setupCron() {
	if ( jobSet.cron != nil ) {
		jobSet.cron.Stop()
		jobSet.cron = nil
	}

	if len(jobSet.jobs ) != 0 {
		log.Printf("  Starting cron for %s", jobSet.directory)
		jobSet.cron = cron.New()
		for _,job := range jobSet.jobs {
			addJob(jobSet.cron, jobSet.executor, job)
		}
		jobSet.cron.Start()
	}
}

func addJob(cron *cron.Cron, executor JobExecutor, job Job) {
	cron.AddFunc(job.Spec, func() {runJob(executor, job)})
}

func runJob(executor JobExecutor, job Job) {
	log.Printf("Running job %s (%s)", job.Name, filepath.Dir(job.Filepath))
	executor(job.Name, job.Filepath)
}

func (jobSet *JobSet) printJobs() {
	log.Printf("Jobs in %s", jobSet.directory)
	if len(jobSet.jobs) == 0 {
	log.Printf("  -- no jobs --")
	}
	for _,job := range jobSet.jobs {
		log.Printf("  %s : %s", job.Spec, job.Name)
	}
	log.Printf("")
}
