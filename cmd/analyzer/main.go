package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/parser"
	"github.com/red3533/log-analyzer/internal/sorter"

	"github.com/red3533/log-analyzer/internal/config"
)

func main() {

	configFlag := flag.String("config", "", "Path to config file (required)")
	filepathFlag := flag.String("file", "", "Path to log file (required)")
	logTypeFlag := flag.String("type", "", "Type of logs to analyze (required)")
	sortFieldFlag := flag.String("field", "ip", "Sort field or \"\" for no sort")
	sortByFlag := flag.String("by", "desc", "Direction of sort: asc, desc")

	flag.Parse()

	if *configFlag == "" {
		fmt.Println("Flag -config not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *filepathFlag == "" {
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

	logEntries, err := logParser.Parse(*filepathFlag)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse logs")
	}

	log.Debug().Msgf("logParsed: %v", logEntries)

	logSorter := sorter.NewLogSorter(log)
	logSorter.Sort(logEntries, *sortFieldFlag, *sortByFlag)

	log.Debug().Msgf("sorted logs: %v", logEntries)

}
