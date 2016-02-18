package main
import (
	"log"
	"path/filepath"
	"os"
	"path"
)

type GoDoItScanner struct {
	executor JobExecutor
	config *GoDoItConfig
	jobSets map[string]*JobSet
}

func NewScanner(config *GoDoItConfig) *GoDoItScanner {
	return &GoDoItScanner{
		JobExecutorFromScript(config.JobExecutorScript),
		config,
		make(map[string]*JobSet)}
}

func (scanner *GoDoItScanner) Run() {
	if scanner.Scan() {
		scanner.PrintJobs()
	}
}

func (scanner *GoDoItScanner) Scan() bool {
	log.Println("Scanning for changes...")
	foundDirectories := make(map [string]bool)

	// Scan all the patterns
	addedJobSets := scanPatterns(scanner, foundDirectories)

	// Remove any unwanted directories
	removedJobSets := removeDirectories(scanner, foundDirectories)

	// Scan directory entries
	updatedJobSets := scanDirectory(scanner)

	jobsChanged := addedJobSets || removedJobSets || updatedJobSets
	log.Println("  Scanning for changes...done")
	return jobsChanged
}

func scanPatterns(scanner *GoDoItScanner, foundDirectories map[string]bool) bool {
	updated := false
	for _,element := range scanner.config.Include {
		path := path.Clean(os.ExpandEnv(element))
		log.Printf("  Scanning directories matching %s", path)
		directories,err := filepath.Glob(path)
		if err != nil {
			log.Printf("  Failed to scan, %s", path)
			continue
		}
		thisUpdated := ensureDirectory(scanner, directories, foundDirectories)
		updated = updated || thisUpdated
	}
	return updated
}

func ensureDirectory(scanner *GoDoItScanner, directories []string, foundDirectories map[string]bool) bool {
	updated := false
	for _,directory := range directories {
		if info, err := os.Stat(directory); err != nil || ! info.IsDir() {
			continue
		}
		foundDirectories[directory] = true
		if _,ok := scanner.jobSets[directory]; ! ok {
			log.Printf("  Adding directory, %s", directory)
			jobSet := NewJobSet(scanner.executor, directory)
			scanner.jobSets[directory] = jobSet
			jobSet.Scan()
			updated = true
		}
	}
	return updated
}

func removeDirectories(scanner *GoDoItScanner, foundDirectories map[string]bool) bool {
	updated := false
	for directory,jobSet := range scanner.jobSets {
		if _, ok := foundDirectories[directory]; ! ok {
			jobSet.Stop()
			delete(scanner.jobSets, directory)
			updated = true
		}
	}
	return updated
}

func scanDirectory(scanner *GoDoItScanner) bool {
	updated := false
	for _,jobSet := range scanner.jobSets {
		thisUpdated := jobSet.Scan()
		updated = updated || thisUpdated
	}
	return updated
}

func (scanner *GoDoItScanner) PrintJobs() {
	for _,jobSet := range scanner.jobSets {
		jobSet.printJobs()
	}}

func (scanner *GoDoItScanner) Stop() {
	log.Println("Stopping jobs...")
	for _,jobSet := range scanner.jobSets {
		jobSet.Stop()
	}
	log.Println("  Stopping jobs...done")
}

