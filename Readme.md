Godoit
======

Godoit is a scheduler which locates and schedules jobs which are
deployed with applications. Godoit will scan a set of
directories for files ending in `.godoit` and schedule them.

The set of directories can be specified to include the wildcard `*` and
include environment variables.

The `.godoit` filename or file contains the cronspec description of when to run
the job and the scheduled job name.

Godoit jobs are executed by a wrapper script which allows the deployment
to handle specific concerns such as job logging and alerting of failures.

Usage
=============
Usage:

    godoit <godoit.conf>


The configuration file is of the format:

    // Paths to scan for jobs
    include = [ '/home/root/systemjobs','$MY_APPS_BASE/*' ]
    // Scan period in seconds
    scanTime = 60
    // Log file
    logFile = '$LOGDIR/godoit.log'
    // Max log file size in MB
    logMaxSize = 7
    // Max log file age in days
    logMaxAge = 10
    // Max number of log files to keep
    logMaxBackups = 5
    // Job executor script
    jobExecutorScript = 'job_wrapper.sh'
    // Godoit status script
    statusScript = 'report_status.sh'
    // Status reporing interval in seconds
    statusInterval = 60
    // Environment variables tp be included on the status json
    statusEnvironment = ['MY_ENV']

The `scanTime` and `statusInterval` are in seconds. The `logMaxSize` is in megabytes.

###Job Scripts
Godoit scripts are named with a `.godoit` suffix. The cronspec can be specified in the 
filename or inside the file as a parameter. When specifying in the filename the form would be:

    <cronspec> <jobname>.godoit
    e.g.
    0 0 18 x x SUN weekend restart.godoit

For the cronspec included in the filename alternative characters can be used to
make file names eaier to manage:
* `*` can be replaced with `x`
* `/` can be replaced with `%`

The `.godoit` file can also include job parameters as comment in the script of 
the form `#:godoit <param> <value>`. This includes specifying the cronspec. 

Supported parameters are:

Comment            | Detail
-------------------|-----------
`#:godoit cronspec ...`| The cron spec (see https://godoc.org/github.com/robfig/cron) 
`#:godoit timeout ...` | Time as a duration after which SIGTERM is sent e.g. `1h30m`, `15s`
`#:godoit timezone ...`| The timezone for the job e.g. `Europe/London`

If the cronspec is specified in both places this is an error and the job will be disabled.
Errors parsing the parameters above will also disable the job.
Godoit will check the modification time of files to detect changes.

*NOTE: Parameters must be specified in the first 10 lines of the file.*

If the `.godoit` filename starts with either `#` or `--` the job will be considered disabled.

###Job Executor

The job executor script will be passed two arguments:
* the job name
* the path to the godoit job whch is to be run

###Status Script
The status script is passed a JSON payload to stdin describing all the jobs.
This can be used to push the set of jobs to a central monitor.

Environment variables can be included which may be useful to add information 
such as server location, environment type etc.

The set of jobs will include disabled jobs and jobs with parameter errors.

###Logging

Godoit writes to a rotating logfile. The logfile includes the output
of the job executor script, and the status reportng script (`stdout` and `stderr`)