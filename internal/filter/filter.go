package sorter

import (
	"regexp"
	"strconv"

	"github.com/red3533/log-analyzer/internal/models"
)

type LogFilter struct {
}

func (s LogFilter) Filter(logs []models.Log, filters []interface{}) ([]models.Log, error) {

	ipRe := regexp.MustCompile(filters[0].(string))
	urlRe := regexp.MustCompile(filters[1].(string))
	statusRe := regexp.MustCompile(strconv.Itoa(filters[2].(int)))

	var filteredLogs []models.Log
	for _, log := range logs {

		logStatus := strconv.Itoa(log.Status)

		if ipRe.FindStringSubmatch(log.IP) != nil && urlRe.FindStringSubmatch(log.URL) != nil && statusRe.FindStringSubmatch(logStatus) != nil {
			filteredLogs = append(filteredLogs, log)
		}

	}

	return filteredLogs, nil

}

func NewLogFilter() LogFilter {
	return LogFilter{}
}
