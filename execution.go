package main

import (
	"os"
	"os/exec"
	"log"
)

type JobExecutor func(jobName, jobPath string)

func JobExecutorFromScript(jobExecutorScript string) JobExecutor {
	expandedJobExecutorScript := os.ExpandEnv(jobExecutorScript)
	return func(jobName, jobPath string) {
		cmd := exec.Command(expandedJobExecutorScript, jobName, jobPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("ERROR: Failed to execute executor script %s", expandedJobExecutorScript)
		}
	}
}
