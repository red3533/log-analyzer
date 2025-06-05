package parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/models"
)

var (
	ipRegexp        = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	timestampRegexp = regexp.MustCompile(`\[.{26}\]`)
	methodRegexp    = regexp.MustCompile(`\"([A-Z]+)`)
	urlRegexp       = regexp.MustCompile(`"[A-Z]+\s+([^ ]+)\s+HTTP/\d\.\d"`)
	sizeByteRegexp  = regexp.MustCompile(`\d+\s(\-*[0-9]*)$`)
)

type NginxParser struct {
	log    logger.Logger
	reader FileReader
}

func (p NginxParser) Parse(filepath string) ([]models.Log, error) {

	var parsed []models.Log
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

		timestamp, err := extractTimestamp(line)
		if err != nil {
			p.log.Warn().Err(err).Str("line", line).Msg("failed to parse log line")
			errorCount++
			continue
		}

		method, err := extractMethod(line)
		if err != nil {
			p.log.Warn().Err(err).Str("line", line).Msg("failed to parse log line")
			errorCount++
			continue
		}

		url, err := extractURL(line)
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

		sizeByte, err := extractSizeByte(line)
		if err != nil {
			p.log.Warn().Err(err).Str("line", line).Msg("failed to parse log line")
			errorCount++
			continue
		}

		parsed = append(parsed, models.Log{
			IP:        ip,
			Timestamp: timestamp,
			Method:    method,
			URL:       url,
			Status:    status,
			SizeByte:  sizeByte,
		})

		successCount++

	}

	p.log.Debug().Int("successCount", successCount).Int("errorCount", errorCount).Msg("Parsed lines")

	return parsed, nil
}

func extractIP(line string) (string, error) {
	ip := ipRegexp.FindString(line)

	if ip == "" {
		return "", fmt.Errorf("ip not found: %s", line)
	}

	return ip, nil
}

func extractTimestamp(line string) (time.Time, error) {
	timestampRaw := timestampRegexp.FindString(line)

	if timestampRaw == "" {
		return time.Time{}, fmt.Errorf("timestamp not found: %s", line)
	}

	// remove []
	timestampRaw = timestampRaw[1:27]

	layout := "02/Jan/2006:15:04:05 -0700"
	timestamp, err := time.Parse(layout, timestampRaw)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to convert timestamp raw: %w", err)
	}

	return timestamp, nil
}

func extractMethod(line string) (string, error) {
	methodGroup := methodRegexp.FindStringSubmatch(line)
	if len(methodGroup) < 1 {
		return "", fmt.Errorf("method not found: %s", line)
	}

	method := methodGroup[1]

	correctMethods := map[string]interface{}{
		"GET":     struct{}{},
		"POST":    struct{}{},
		"PUT":     struct{}{},
		"DELETE":  struct{}{},
		"HEAD":    struct{}{},
		"OPTIONS": struct{}{},
		"PATCH":   struct{}{},
		"TRACE":   struct{}{},
		"CONNECT": struct{}{},
	}

	if _, ok := correctMethods[method]; !ok {
		return "", fmt.Errorf("incorrect method: %s", method)
	}

	return method, nil
}

func extractURL(line string) (string, error) {
	urlGroup := urlRegexp.FindStringSubmatch(line)
	if len(urlGroup) < 1 {
		return "", fmt.Errorf("url not found: %s", line)
	}

	urlPart := urlGroup[1]

	if _, err := url.Parse(urlPart); err != nil {
		return "", fmt.Errorf("incorrect url: %s: %w", urlPart, err)
	}

	return urlPart, nil
}

func extractStatus(logLine string) (int, error) {
	logParts := strings.Split(logLine, " ")

	for _, part := range logParts {
		if len(part) == 3 {
			status, err := strconv.Atoi(part)
			if err != nil {
				return -1, fmt.Errorf("failed convert to int: %s: %w", part, err)
			}

			if status < 100 || status > 599 {
				return -1, fmt.Errorf("status code: %d must be in range 100-599", status)
			}

			return status, nil
		}
	}

	return -1, fmt.Errorf("status code not found: %s", logLine)
}

func extractSizeByte(line string) (int, error) {
	sizeByteGroup := sizeByteRegexp.FindStringSubmatch(line)

	if len(sizeByteGroup) < 1 {
		return -1, fmt.Errorf("size byte not found: %s", line)
	}

	sizeByteStr := string(sizeByteGroup[1])
	sizeByte, err := strconv.Atoi(sizeByteStr)
	if err != nil {
		return -1, fmt.Errorf("failed to convert: %w", err)
	}

	if sizeByte < 0 {
		return -1, fmt.Errorf("incorrect size byte: %d", sizeByte)
	}

	return sizeByte, nil
}

func NewNginxParser(log logger.Logger, reader FileReader) NginxParser {
	return NginxParser{
		log:    log,
		reader: reader,
	}
}
