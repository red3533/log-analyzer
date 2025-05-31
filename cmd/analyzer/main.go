package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/red3533/log-analyzer/internal/parser"
)

func main() {
	filepath := flag.String("file", "", "Path to log file (required)")

	flag.Parse()

	if *filepath == "" {
		fmt.Println("Flag -file not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	logParsed, err := parser.NginxParser.Parse(parser.NginxParser{}, *filepath)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("logParsed: %v\n", logParsed)

}
