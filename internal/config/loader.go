package config

import (
	"fmt"
	"os"

	"github.com/red3533/log-analyzer/internal/models"
	"gopkg.in/yaml.v3"
)

func MustLoadConfig(filepath string) *models.AppConfig {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("file not found: %s: error: %s\n", filepath, err.Error())
		panic(errMsg)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read file: %s: error: %s\n", filepath, err.Error())
		panic(errMsg)
	}

	var cfg models.AppConfig

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		errMsg := fmt.Sprintf("failed to unmarshall data: %s: error: %s\n", data, err.Error())
		panic(errMsg)
	}

	return &cfg
}
