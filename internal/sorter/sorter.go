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

func (s LogSorter) Sort(logs []models.Log, sortField string, sortBy string) error {

	switch sortField {
	case "ip":
		if sortBy == "asc" {
			sort.SliceStable(logs, func(i, j int) bool { return logs[i].IP < logs[j].IP })
			return nil
		}
		sort.SliceStable(logs, func(i, j int) bool { return logs[i].IP > logs[j].IP })

	case "timestamp":
		if sortBy == "asc" {
			sort.SliceStable(logs, func(i, j int) bool { return logs[i].Timestamp.Before(logs[j].Timestamp) })
			return nil
		}
		sort.SliceStable(logs, func(i, j int) bool { return logs[i].Timestamp.After(logs[j].Timestamp) })

	case "method":
		if sortBy == "asc" {
			sort.SliceStable(logs, func(i, j int) bool { return logs[i].Method < logs[j].Method })
			return nil
		}
		sort.SliceStable(logs, func(i, j int) bool { return logs[i].Method > logs[j].Method })

	case "url":
		if sortBy == "asc" {
			sort.SliceStable(logs, func(i, j int) bool { return logs[i].URL < logs[j].URL })
			return nil
		}
		sort.SliceStable(logs, func(i, j int) bool { return logs[i].URL > logs[j].URL })

	case "status":
		if sortBy == "asc" {
			sort.SliceStable(logs, func(i, j int) bool { return logs[i].Status < logs[j].Status })
			return nil
		}
		sort.SliceStable(logs, func(i, j int) bool { return logs[i].Status > logs[j].Status })

	case "size_byte":
		if sortBy == "asc" {
			sort.SliceStable(logs, func(i, j int) bool { return logs[i].SizeByte < logs[j].SizeByte })
			return nil
		}
		sort.SliceStable(logs, func(i, j int) bool { return logs[i].SizeByte > logs[j].SizeByte })

	default:
		return fmt.Errorf("unknown sort field: %s", sortField)
	}

	return nil
}

func NewLogSorter(log logger.Logger) LogSorter {
	return LogSorter{
		log: log,
	}
}
