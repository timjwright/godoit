package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	// TODO...
	jobExec := JobExecutorFromScript("./test_wrapper.sh")
	jobExec("name","path")
	assert.True(t, true, "Failed to parse job")
}
