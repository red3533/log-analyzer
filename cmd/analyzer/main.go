package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/red3533/log-analyzer/internal/filter"
	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/parser"

	"github.com/red3533/log-analyzer/internal/config"
)

func main() {

	configFlag := flag.String("config", "", "Path to config file (required)")
	logFilepathFlag := flag.String("file", "", "Path to log file (required)")
	logTypeFlag := flag.String("type", "", "Type of logs to analyze (required)")

	filterIPFlag := flag.String("filter-ip", "", "")
	filterURLFlag := flag.String("filter-url", "/api/*", "")
	filterStatusFlag := flag.Int("filter-status", 200, "")

	flag.Parse()

	if *configFlag == "" {
		fmt.Println("Flag -config not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *logFilepathFlag == "" {
		fmt.Println("Flag -file not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *logTypeFlag == "" {
		fmt.Println("Flag -type not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg := config.MustLoadConfig(*configFlag)
	log := logger.NewLogger(cfg.LoggerConfig)

	var logParser parser.LogParser

	switch *logTypeFlag {
	case "json":
		// implement
		log.Fatal().Str("type", *logTypeFlag).Msg("unknown log type")

	case "nginx":
		logParser = parser.NewNginxParser(log, parser.DefaultFileReader{})

	default:
		log.Fatal().Str("type", *logTypeFlag).Msg("unknown log type")
	}

	parsedLogs, err := logParser.Parse(*logFilepathFlag)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse logs")
	}

	// filter
	logFilter := sorter.NewLogFilter()
	filters := []interface{}{*filterIPFlag, *filterURLFlag, *filterStatusFlag}

	filteredLogs, err := logFilter.Filter(parsedLogs, filters)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to sort logs")
	}

	// debug
	fmt.Println("----- SORTED LOGS -----")
	for _, fl := range filteredLogs {
		fmt.Println(fl)
	}

}
