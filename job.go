package main

import (
	"path"
	"path/filepath"
	"strings"
	"regexp"
	"time"
	"os"
	"bufio"
	"log"
	"github.com/robfig/cron"
	"fmt"
)

type Job struct {
	Filepath string
	Spec string
	Timezone *time.Location
	Name string
	Timeout time.Duration
	Enabled bool
	Errors []string
	UpdateTime time.Time
}

var cronSpecRegex,_ = regexp.Compile(`\s*($|#|\w+\s*=|(x|\*|(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?(?:,(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?)*)\s+(x|\*|(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?(?:,(?:[0-5]?\d)(?:(?:-|%|\,)(?:[0-5]?\d))?)*)\s+(x|\*|(?:[01]?\d|2[0-3])(?:(?:-|%|\,)(?:[01]?\d|2[0-3]))?(?:,(?:[01]?\d|2[0-3])(?:(?:-|%|\,)(?:[01]?\d|2[0-3]))?)*)\s+(x|\*|(?:0?[1-9]|[12]\d|3[01])(?:(?:-|%|\,)(?:0?[1-9]|[12]\d|3[01]))?(?:,(?:0?[1-9]|[12]\d|3[01])(?:(?:-|%|\,)(?:0?[1-9]|[12]\d|3[01]))?)*)\s+(x|\*|(?:[1-9]|1[012])(?:(?:-|%|\,)(?:[1-9]|1[012]))?(?:L|W)?(?:,(?:[1-9]|1[012])(?:(?:-|%|\,)(?:[1-9]|1[012]))?(?:L|W)?)*|x|\*|(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)(?:(?:-)(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC))?(?:,(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)(?:(?:-)(?:JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC))?)*)\s+(x|\*|(?:[0-6])(?:(?:-|%|\,|#)(?:[0-6]))?(?:L)?(?:,(?:[0-6])(?:(?:-|%|\,|#)(?:[0-6]))?(?:L)?)*|x|\*|(?:MON|TUE|WED|THU|FRI|SAT|SUN)(?:(?:-)(?:MON|TUE|WED|THU|FRI|SAT|SUN))?(?:,(?:MON|TUE|WED|THU|FRI|SAT|SUN)(?:(?:-)(?:MON|TUE|WED|THU|FRI|SAT|SUN))?)*)(|\s)+(x|\*|(?:|\d{4})(?:(?:-|%|\,)(?:|\d{4}))?(?:,(?:|\d{4})(?:(?:-|%|\,)(?:|\d{4}))?)*)) (.*)\.godoit`)
var noTimeout = time.Second * 0
var GodoitFileSuffix = ".godoit"
var godoitCommentPrefix = "#:godoit "


func ParseJobFile(directory, filename string) *Job {
	jobPath := path.Join(directory, filename)
	var cronspec string
	var name string
	enabled :=
		!strings.HasPrefix(filename, "--") &&
		!strings.HasPrefix(filename, "#")

	if result := cronSpecRegex.FindStringSubmatch(filename); result != nil {
		cronspec = strings.Replace(result[1], "x", "*", -1)
		cronspec = strings.Replace(cronspec, "%", "/", -1)
		name = strings.TrimSpace(result[10])
	} else if strings.HasSuffix(filename, GodoitFileSuffix) {
		name = strings.TrimSuffix(filename, GodoitFileSuffix)
	} else {
		return nil
	}

	cronspec, timeout, timezone, errors, updateTime := parseJobParameters(jobPath, cronspec)

	if cronspec == "" {
		errors = append(errors, "Missing cronspec")
	}

	if len(errors) > 0 {
		enabled = false
		log.Printf("Errors parsing job %s: %v", jobPath, errors)
	}

	return &Job{
		filepath.Join(directory, filename),
		cronspec,
		timezone,
		name,
		timeout,
		enabled,
		errors,
		updateTime}
}

func parseJobParameters(jobPath, cronspec string) (string, time.Duration, *time.Location, []string, time.Time) {
	timeout := 0 * time.Second
	timezone := time.UTC
	var updateTime time.Time
	errors := make ([]string,0,10)

	if file, err := os.Open(jobPath); err == nil {
		defer file.Close()
		if info, err := file.Stat() ; err == nil {
			updateTime = info.ModTime()
		}

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		i := 0
		for scanner.Scan() {
			// Only scan start of file for params...
			i++
			if i > 10 {
				break
			}

			line := scanner.Text()
			if strings.HasPrefix(line, godoitCommentPrefix) {
				line = strings.TrimPrefix(line, godoitCommentPrefix)
				line = strings.TrimSpace(line)
				parts := strings.SplitN(line," ",2)
				if len(parts) == 2 {
					param := parts[0]
					value := parts[1]
					if param == "cronspec" {
						if cronspec != "" {
							errors = append(errors, "Cronspec in filename and as comment")
						}
						if _, err := cron.Parse(value); err == nil {
							cronspec = value
						} else {
							errors = append(errors, fmt.Sprintf("Invalid cronspec: '%s'", value))
						}
					} else if param == "timeout" {
						if d, err := time.ParseDuration(value); err == nil {
							timeout = d
						} else {
							errors = append(errors, fmt.Sprintf("Invalid timeout: '%s'", value))
						}
					} else if param == "timezone" {
						if l, err := time.LoadLocation(value); err == nil {
							timezone = l
						} else {
							errors = append(errors, fmt.Sprintf("Invalid timezone: '%s'", value))
						}
					}

				} else {
					errors = append(errors, fmt.Sprintf("Invalid parameter '%s'", line))
				}
			}
		}
	} else {
		errors = append(errors, "Unable to open file to parse parameters")
	}
	return cronspec, timeout, timezone, errors, updateTime
}