package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestExecutor(t *testing.T) {
	// TODO...
	jobExec := JobExecutorFromScript("./test_wrapper.sh", os.Stdout)
	jobExec("my job","/path/to/@ 1 @ @ @ @ my job.godoit")
	assert.True(t, true, "Failed to parse job")
}
