package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/models"
)

// TODO: replace to default file reader
type NginxFileReader struct {
}

func (r NginxFileReader) ReadLines(filepath string) ([]string, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

type NginxParser struct {
	log    logger.Logger
	reader FileReader
}

func (p NginxParser) Parse(filepath string) ([]models.LogParsed, error) {

	var parsed []models.LogParsed
	var successCount, errorCount int

	lines, err := p.reader.ReadLines(filepath)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {

		status, err := extractStatus(line)
		if err != nil {
			p.log.Warn().Err(err).Str("line", line).Msg("failed to parse log line")
			errorCount++
			continue
		}

		parsed = append(parsed, models.LogParsed{Status: status})

		successCount++

	}

	p.log.Debug().Int("successCount", successCount).Int("errorCount", errorCount).Msg("Parsed lines")

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

func NewNginxParser(log logger.Logger, reader FileReader) NginxParser {
	return NginxParser{
		log:    log,
		reader: reader,
	}
}
