package main

import (
	"path/filepath"
	"strings"
	"regexp"
)

type Job struct {
	Filepath string
	Spec string
	Name string
}

var cronSpecRegex,_ = regexp.Compile(`\s*($|#|\w+\s*=|(x|\*|(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?(?:,(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?)*)\s+(x|\*|(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?(?:,(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?)*)\s+(x|\*|(?:[01]?\d|2[0-3])(?:(?:-|%|\,)(?:[01]?\d|2[0-3]))?(?:,(?:[01]?\d|2[0-3])(?:(?:-|%|\,)(?:[01]?\d|2[0-3]))?)*)\s+(x|\*|(?:0?[1-9]|[12]\d|3[01])(?:(?:-|%|\,)(?:0?[1-9]|[12]\d|3[01]))?(?:,(?:0?[1-9]|[12]\d|3[01])(?:(?:-|%|\,)(?:0?[1-9]|[12]\d|3[01]))?)*)\s+(x|\*|(?:[1-9]|1[012])(?:(?:-|%|\,)(?:[1-9]|1[012]))?(?:L|W)?(?:,(?:[1-9]|1[012])(?:(?:-|%|\,)(?:[1-9]|1[012]))?(?:L|W)?)*|x|\*|(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)(?:(?:-)(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC))?(?:,(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)(?:(?:-)(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC))?)*)\s+(x|\*|(?:[0-6])(?:(?:-|%|\,|#)(?:[0-6]))?(?:L)?(?:,(?:[0-6])(?:(?:-|%|\,|#)(?:[0-6]))?(?:L)?)*|x|\*|(?:MON|TUE|WED|THU|FRI|SAT|SUN)(?:(?:-)(?:MON|TUE|WED|THU|FRI|SAT|SUN))?(?:,(?:MON|TUE|WED|THU|FRI|SAT|SUN)(?:(?:-)(?:MON|TUE|WED|THU|FRI|SAT|SUN))?)*)(|\s)+(x|\*|(?:|\d{4})(?:(?:-|%|\,)(?:|\d{4}))?(?:,(?:|\d{4})(?:(?:-|%|\,)(?:|\d{4}))?)*)) (.*)\.godoit`)

func ParseJobFile(path, filename string) *Job {
	// Return commented out files
	if strings.HasPrefix(filename, "--") {
		return nil
	}
	if result := cronSpecRegex.FindStringSubmatch(filename); result != nil {
		cronspec := strings.Replace(result[1], "x", "*", -1)
		cronspec = strings.Replace(cronspec, "%", "/", -1)
		return &Job{filepath.Join(path, filename), cronspec, strings.TrimSpace(result[10])}
	} else {
		return nil
	}
}