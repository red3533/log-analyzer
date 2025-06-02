package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/models"
	"github.com/red3533/log-analyzer/internal/parser"
)

func main() {
	// TODO: load logger config from file
	loggerConfig := models.LoggerConfig{
		LogFile:    "logs/app.log",
		LogLevel:   "debug",
		MaxSizeMB:  20,
		MaxBackups: 100,
		MaxAgeDays: 30,
	}

	log := logger.NewLogger(loggerConfig)

	filepath := flag.String("file", "", "Path to log file (required)")

	flag.Parse()

	if *filepath == "" {
		fmt.Println("Flag -file not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	nginxParser := parser.NewNginxParser(log)
	logParsed, err := parser.NginxParser.Parse(nginxParser, *filepath)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("logParsed: %v\n", logParsed)

}
