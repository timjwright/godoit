package main

import (
	"github.com/robfig/cron"
	"github.com/natefinch/lumberjack"
	"os"
	"os/signal"
	"log"
	"fmt"
)


func main() {
	log.Println("Starting GoDoIt")
	config := LoadConfig()

	logger := &lumberjack.Logger{
		Filename:   os.ExpandEnv(config.LogFile),
		MaxSize:    config.LogMaxSize, // megabytes
		MaxAge:     config.LogMaxAge, //days
	}
	log.SetOutput(logger)
	scanner := NewScanner(config, logger)

	cron := cron.New()
	cron.AddFunc(fmt.Sprintf("@every %ds",config.ScanTime), func(){scanner.Run()})
	log.Println("Starting scanner")
	cron.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <- c
	log.Println("Shutting down: ", s)
	cron.Stop()
	scanner.Stop()
}

