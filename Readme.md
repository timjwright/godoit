Godoit
======

Godoit is a scheduler which locates and schedules jobs which are
deployed with applications. Godoit will scan a set of
directories for files ending in `.godoit` and schedule them.

The set of directories can be specified to include the wildcard `*` and
include environment variables.

The `.godoit` filename contains the cronspec description of when to run
the job and the scheduled job name.

The `.godoit` file can also include job parameters as comment in the script (including specifying the cronspec
as comments of the form `#:godoit <param> <value>`
e.g.

Comment            | Detail
-------------------|-----------
`#:godoit cronspec ...`| The cron spec (see https://godoc.org/github.com/robfig/cron) 
`#:godoit timeout ...` | Time as a duration after which SIGTERM is sent e.g. `1h30m`, `15s`
`#:godoit timezone ...`| The timezone for the job e.g. `Europe/London`

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

The `scanTime` and `statusInterval` are in seconds. The `logMaxSize` is in megabytes.

###Job Scripts
The filename of the executor script is of the form:

    <cronspec> <jobname>.godoit
    e.g.
    0 0 18 x x SUN weekend restart.godoit

In the cron spec alternative charecters can be used to
make file names eaier to manage:
* `*` can be replaced with `x`
* `/` can be replaced with `%`

###Job Executor

The job executor script will be passed two arguments:
* the job name
* the path to the godoit job whch is to be run

###Status Script
The status script is passed a JSON payload to stdin describing all the jobs.
This can be used to push the set of jobs to a central monitor.

###Logging

Godoit writes to a rotating logfile. The logfile includes the output
of the job executor script, and the status reportng script (`stdout` and `stderr`)