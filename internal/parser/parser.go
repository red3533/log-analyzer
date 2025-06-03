package parser

import (
	"github.com/red3533/log-analyzer/internal/models"
)

type LogParser interface {
	Parse(filepath string) ([]models.LogParsed, error)
}

type FileReader interface {
	ReadLines(filepath string) ([]string, error)
}
