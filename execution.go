package main

import (
	"os"
	"os/exec"
	"log"
	"io"
	"time"
	"syscall"
)

type JobExecutor func(jobName, jobPath string, timeout time.Duration)

func JobExecutorFromScript(jobExecutorScript string, output io.Writer) JobExecutor {
	if len(jobExecutorScript) == 0 {
		log.Fatal("Job executor is not defined")
	}
	jobExecutorScript = os.ExpandEnv(jobExecutorScript)
	return func(jobName, jobPath string, timeout time.Duration) {
		cmd := exec.Command(jobExecutorScript, jobName, jobPath)
		log.Printf("Running comand line: %s '%s' '%s' Timeout: %s", jobExecutorScript, jobName, jobPath, timeout)
		cmd.Stdout = output
		cmd.Stderr = output
		err := runWithTimout(cmd, timeout)
		if err != nil {
			log.Printf("ERROR: Failed to execute executor script %s %s %s", jobExecutorScript, jobName, jobPath)
		}
	}
}

func runWithTimout(cmd *exec.Cmd, timeout time.Duration) error {
	if timeout.Seconds() <= 0 {
		return cmd.Run()
	} else {
		done := make(chan error, 1)
		go func() {
			cmd.Start()
			done <- cmd.Wait()
		}()
		select {
		case <-time.After(timeout):
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				log.Printf("ERROR: Failed to terminate job %s error: %s", cmd.Path, err)
			}
			log.Printf("Job %s timed out", cmd.Path)
			return nil
		case err := <-done:
			return err
		}
	}
}
