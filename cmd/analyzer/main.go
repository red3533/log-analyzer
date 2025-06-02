package main

import (
	"flag"
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
		log.Error().Msg("Flag -file not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	nginxParser := parser.NewNginxParser(log)
	logParsed, err := parser.NginxParser.Parse(nginxParser, *filepath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse logs")
	}

	log.Debug().Msgf("logParsed: %v", logParsed)

}
