package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO-refactor: move to internal/models
type LogParsed struct {
	Status int
}

type LogParser interface {
	Parse(filepath string) []LogParsed
}

type NginxParser struct {
}

func (p NginxParser) Parse(filepath string) ([]LogParsed, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s: %w", filepath, err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed open file by path: %s: %w", filepath, err)
	}

	var parsed []LogParsed
	var total int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		logLine := scanner.Text()

		status, err := extractStatus(logLine)
		if err != nil {
			return nil, fmt.Errorf("failed to get status code from line: %s: %w", logLine, err)
		}

		parsed = append(parsed, LogParsed{Status: status})

		total++
	}

	fmt.Println("total", total)

	return parsed, nil
}

func extractStatus(logLine string) (int, error) {
	logParts := strings.Split(logLine, " ")

	for _, part := range logParts {
		if len(part) == 3 {
			status, err := strconv.Atoi(part)
			if err != nil {
				return -1, fmt.Errorf("failed convert to int: %s: %w", part, err)
			}
			return status, nil
		}
	}

	return -1, fmt.Errorf("status code not found")
}
