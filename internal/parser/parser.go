package parser

import (
	"bufio"
	"os"

	"github.com/red3533/log-analyzer/internal/models"
)

type LogParser interface {
	Parse(filepath string) ([]models.Log, error)
}

type FileReader interface {
	ReadLines(filepath string) ([]string, error)
}

type DefaultFileReader struct {
}

func (r DefaultFileReader) ReadLines(filepath string) ([]string, error) {
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
