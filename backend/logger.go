package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func logger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0750); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logName := time.Now().Format("2006-01-01") + "-logs.log"

	if strings.Contains(logName, string(filepath.Separator)) {
		log.Fatalf("Invalid log filename")
	}

	logFilePath := filepath.Join(logDir, logName)
	if !strings.HasPrefix(filepath.Clean(logFilePath), filepath.Clean(logDir)) {
		log.Fatalf("Invalid log file path")
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
