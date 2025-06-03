package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/parser"

	"github.com/red3533/log-analyzer/internal/config"
)

func main() {

	configFlag := flag.String("config", "", "Path to config file (required)")
	filepathFlag := flag.String("file", "", "Path to log file (required)")

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

	cfg := config.MustLoadConfig(*configFlag)
	log := logger.NewLogger(cfg.LoggerConfig)

	nginxParser := parser.NewNginxParser(log, parser.NginxFileReader{})
	logParsed, err := nginxParser.Parse(*filepathFlag)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse logs")
	}

	log.Debug().Msgf("logParsed: %v", logParsed)

}
