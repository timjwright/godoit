package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPathMatching(t *testing.T) {
	job := ParseJobFile("a_path", "0 30 * * * *  job name ")
	assert.NotNil(t, job, "Failed to parse job")
	assert.Equal(t, "0 30 * * * *", job.Spec)
	assert.Equal(t, "job name", job.Name)
	assert.Equal(t, "a_path/0 30 * * * *  job name ", job.Filepath)
}

func TestPathMatching2(t *testing.T) {
	job := ParseJobFile("a_path", "* * * * * * TestScanRemoveJob")
	assert.NotNil(t, job, "Failed to parse job")
}

func TestPathAll(t *testing.T) {
	job := ParseJobFile("a_path", "* * * * * * job name")
	assert.NotNil(t, job, "Failed to parse job")
}

func TestCommentedOutJobShouldBeNil(t *testing.T) {
	job := ParseJobFile("a_path", "#0 30 * * * * job name")
	assert.Nil(t, job, "Should be nil")
}

func TestIncompleteTaskShouldBeNil(t *testing.T) {
	job := ParseJobFile("a_path", "0 30 * * *")
	assert.Nil(t, job, "Should be nil")
}

func TestWithoutNameShouldBeNil(t *testing.T) {
	job := ParseJobFile("a_path", "0 30 * * * *")
	assert.Nil(t, job, "Should be nil")
}
