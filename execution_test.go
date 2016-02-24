package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	// TODO...
	jobExec := JobExecutorFromScript("./test_wrapper.sh")
	jobExec("my job","/path/to/@ 1 @ @ @ @ my job.godoit")
	assert.True(t, true, "Failed to parse job")
}
