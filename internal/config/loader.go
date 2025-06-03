package config

import (
	"fmt"
	"os"

	"github.com/red3533/log-analyzer/internal/models"
	"gopkg.in/yaml.v3"
)

func MustLoadConfig(filepath string) *models.AppConfig {
	cfg, err := LoadConfig(filepath)
	if err != nil {
		panic(err)
	}
	return cfg
}

func LoadConfig(filepath string) (*models.AppConfig, error) {
	// TODO: fix error return

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Printf("file not found: %s: %s\n", filepath, err.Error())
		return nil, err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("failed to read file: %s: %s\n", filepath, err.Error())
		return nil, err
	}

	var cfg models.AppConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Printf("failed to unmarshall data: %s: %s\n", data, err.Error())
		return nil, err
	}

	err = cfg.LoggerConfig.Validate()
	if err != nil {
		return nil, err
	}

	// TODO: add validate for others config

	return &cfg, nil
}
