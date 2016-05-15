package main

import (
	"github.com/influxdata/config"
	"log"
	"os"
)


type GoDoItConfig struct {
	Include []string `toml:"include" doc:"Paths to scan"`
	JobExecutorScript string`toml:"JobExecutorScript" doc:"Paths for job executor script"`
	ScanTime int `toml:"ScanTime" doc:"Scan time in seconds"`
	LogFile string `toml:"LogFile" doc:"Logfile location"`
	LogMaxSize int `toml:"LogMaxSize" doc:"Log fie max size"`
	LogMaxAge int `toml:"LogMaxAge" doc:"Number of days to keep th log file"`
}

func LoadConfig() *GoDoItConfig {
	if ( len(os.Args) != 2 ) {
		log.Fatalf("Usage: %s <config file>", os.Args[0])
	}
	cfgFile := os.Args[1]
	log.Printf("Loading config file: %s", cfgFile)
	defaults := GoDoItConfig{[]string{}, "", 30, "godoit.log", 100, 14}
	cfg, err := config.NewConfig(cfgFile, defaults)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	var goDoItConfig GoDoItConfig
	if err := cfg.Decode(&goDoItConfig); err != nil {
		log.Fatalf("Error parsing configuration: %s", err.Error())
	}
	log.Printf("Loaded config:\n %s", goDoItConfig)
	return &goDoItConfig
}
