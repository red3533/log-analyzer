package sorter

import (
	"fmt"
	"sort"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/models"
)

type LogSorter struct {
	log logger.Logger
}

func (s LogSorter) Sort(logs []models.Log, sortFieldFlag string, sortByFlag string) error {

	if sortByFlag != "asc" && sortByFlag != "desc" {
		return fmt.Errorf("unknown sort direction: %s", sortByFlag)
	}

	comparators := map[string]func(i, j models.Log) bool{
		"ip":        func(i, j models.Log) bool { return i.IP < j.IP },
		"timestamp": func(i, j models.Log) bool { return i.Timestamp.Before(j.Timestamp) },
		"method":    func(i, j models.Log) bool { return i.Method < j.Method },
		"url":       func(i, j models.Log) bool { return i.URL < j.URL },
		"status":    func(i, j models.Log) bool { return i.Status < j.Status },
		"size_byte": func(i, j models.Log) bool { return i.SizeByte < j.SizeByte },
	}

	comparator, ok := comparators[sortFieldFlag]
	if !ok {
		return fmt.Errorf("unknown sort field: %s", sortFieldFlag)
	}

	if sortByFlag == "asc" {
		sort.SliceStable(logs, func(i, j int) bool { return comparator(logs[i], logs[j]) })
	} else {
		sort.SliceStable(logs, func(i, j int) bool { return comparator(logs[j], logs[i]) })

	}

	return nil

}

func NewLogSorter(log logger.Logger) LogSorter {
	return LogSorter{
		log: log,
	}
}
