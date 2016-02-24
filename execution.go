package main

import (
	"os"
	"os/exec"
	"log"
)

type JobExecutor func(jobName, jobPath string)

func JobExecutorFromScript(jobExecutorScript string) JobExecutor {
	if len(jobExecutorScript) == 0 {
		log.Fatalf("Job executor is not defined")
	}
	jobExecutorScript = os.ExpandEnv(jobExecutorScript)
	return func(jobName, jobPath string) {
		cmd := exec.Command(jobExecutorScript, jobName, jobPath)
		log.Printf("Running comand line: %s '%s' '%s'", jobExecutorScript, jobName, jobPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("ERROR: Failed to execute executor script %s %s %s", jobExecutorScript, jobName, jobPath)
		}
	}
}
