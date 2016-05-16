Godoit
======

Godoit is a scheduler which locates and schedules jobs which are
deployed with applications. Godoit will scan a set of
directories for files ending in `.godoit` and schedule them.

The set of directories can be specified to include the wildcard `*` and
include environment variables.

The `.godoit` filename contains the cronspec description of when to run
the job and the scheduled job name.

Godoit jobs are executed by a wrapper script which allows the deployment
to handle specific concerns such as job logging and alerting of failures.

Usage
=============
Usage:

    godoit <godoit.conf>


The configuration file is of the format:

    include = [ '/home/root/systemjobs','$MY_APPS_BASE/*' ]
    scanTime = 60
    logFile = '$LOGDIR/godoit.log'
    logMaxSize = 7
    logMaxAge = 10
    jobExecutorScript = 'job_wrapper.sh'

The `scanTime` is in seconds. The `logMaxSize` is in megabytes.

The job executor script will be passed two arguments:
* the job name
* the path to the godoit job whch is to be run

Godoit writes to a rotating logfile. The logfile includes the output
of the job executor script (`stdout` and `stderr`)