package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	JobExecutor exec = JobExecutorFromScript("/bin/bash/echo Hello")
	exec.
	assert.True(t, true, "Failed to parse job")
}
