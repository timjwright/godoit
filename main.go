package main

import (
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"log"
)


func main() {
	log.Println("Starting GoDoIt")

	scanner := NewScanner(LoadConfig())

	cron := cron.New()
	cron.AddFunc("@every 5s", func(){scanner.Run()})
	log.Println("Starting scanner")
	cron.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <- c
	log.Println("Shutting down: ", s)
	cron.Stop()
	scanner.Stop()
}

