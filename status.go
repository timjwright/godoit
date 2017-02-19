package main

import (
	"os"
	"os/exec"
	"log"
	"io"
	"encoding/json"
	"time"
)

type StatusReporter func(jobSets map[string]*JobSet)

type GodoitInfo struct {
	Time string				   `json:"time"`
	Hostname string 		   `json:"hostname"`
	JobInfo []JobCollection	   `json:"jobInfo"`
	Environment map[string]string `json:"environment"`
}

type JobCollection struct {
	Path string		`json:"path"`
	Jobs []JobInfo	`json:"jobs"`
}


type JobInfo struct {
	Name string	`json:"name"`
	Spec string `json:"spec"`
	Timezone string `json:"timezone"`
	Path string `json:"path"`
	Timeout int `json:"timeout"`
	Enabled bool `json:"enabled"`
	Errors []string `json:"errors"`
}

func ToJson(jobSets map[string]*JobSet, statusEnvironment []string) []byte {
	jobCollections := make([]JobCollection, len(jobSets))
	i := 0
	for _, jobSet := range jobSets {
		jobs := make([]JobInfo, len(jobSet.jobs))
		j := 0
		for _, job := range jobSet.jobs {
			jobs[j] =
				JobInfo{
					job.Name,
					job.Spec,
					job.Timezone.String(),
					job.Filepath,
					int(job.Timeout.Seconds()),
					job.Enabled,
					job.Errors}
			j++

		}
		jobCollections[i] = JobCollection{jobSet.directory, jobs}
		i++
	}

	hostname, _ := os.Hostname()
	time := time.Now().UTC().Format("20060102T15:04:05Z")

	// Setup environment variables
	environment := make(map[string]string)
	for _, environmentVariable := range statusEnvironment {
		environment[environmentVariable] = os.Getenv(environmentVariable)
	}

	godoitInfo := &GodoitInfo{time,hostname,jobCollections, environment}
	info, _ := json.Marshal(godoitInfo)
	return info
}

func StatusReporterFromScript(statusScript string, statusEnvironment []string, output io.Writer) StatusReporter {
	if len(statusScript) == 0 {
		log.Fatalf("Status script is not defined")
	}
	statusScript = os.ExpandEnv(statusScript)
	return func(jobSets map[string]*JobSet) {
		cmd := exec.Command(statusScript)
		log.Printf("Running status script: %s", statusScript)
		cmd.Stdout = output
		cmd.Stderr = output
		pipe, _ := cmd.StdinPipe()
		err := cmd.Start()
		pipe.Write(ToJson(jobSets, statusEnvironment))
		pipe.Close()
		if err != nil {
			log.Printf("ERROR: Failed to execute status script %s", statusScript)
		}

		cmd.Wait()
		if err != nil {
			log.Printf("ERROR: Status script completed with error code %s", statusScript)
		}
	}

}
