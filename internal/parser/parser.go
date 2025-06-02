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

type LogParser interface {
	Parse(filepath string) []models.LogParsed
}

type NginxParser struct {
	log logger.Logger
}

func NewNginxParser(log logger.Logger) NginxParser {
	return NginxParser{log: log}
}

func (p NginxParser) Parse(filepath string) ([]models.LogParsed, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		p.log.Error().Err(err).Str("filepath", filepath).Msg("file not found")
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		p.log.Error().Err(err).Str("filepath", filepath).Msg("failed to open file")
		return nil, err
	}
	// TODO: add file close

	var parsed []models.LogParsed
	var successCount, errorCount int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		logLine := scanner.Text()

		status, err := extractStatus(logLine)
		if err != nil {
			p.log.Warn().Err(err).Str("line", logLine).Msg("failed to parse log line")
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
