package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"time"
)

func TestExecutor(t *testing.T) {
	// TODO...
	jobExec := JobExecutorFromScript("./test_wrapper.sh", os.Stdout)
	jobExec("my job","/path/to/@ 1 @ @ @ @ my job.godoit", noTimeout)
	assert.True(t, true, "Failed to parse job")
}

func TestExecutorWithTimeout(t *testing.T) {
	jobExec := JobExecutorFromScript("./test_wrapper_sleep.sh", os.Stdout)
	start := time.Now()
	jobExec("my job","/path/to/@ 1 @ @ @ @ my job.godoit", time.Second * 3)
	duration := time.Since(start)
	assert.True(t, duration.Seconds() < 4.0, "Job took to long")
}
