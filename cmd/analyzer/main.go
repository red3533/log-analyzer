package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	filepath := flag.String("file", "", "Path to log file (required)")

	flag.Parse()

	if *filepath == "" {
		fmt.Println("Flag -file not set")
		flag.PrintDefaults()
		os.Exit(1)
	}

}
