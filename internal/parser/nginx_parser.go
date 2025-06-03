package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/models"
)

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

		ip, err := extractIP(line)
		if err != nil {
			p.log.Warn().Err(err).Str("line", line).Msg("failed to parse log line")
			errorCount++
			continue
		}

		status, err := extractStatus(line)
		if err != nil {
			p.log.Warn().Err(err).Str("line", line).Msg("failed to parse log line")
			errorCount++
			continue
		}

		parsed = append(parsed, models.LogParsed{
			IP:     ip,
			Status: status,
		})

		successCount++

	}

	p.log.Debug().Int("successCount", successCount).Int("errorCount", errorCount).Msg("Parsed lines")

	return parsed, nil
}

func extractIP(line string) (string, error) {
	re := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)

	ip := re.FindString(line)

	return ip, nil
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
