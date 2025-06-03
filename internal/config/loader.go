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
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg models.AppConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall data: %w", err)
	}

	err = cfg.LoggerConfig.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate: %w", err)
	}

	// TODO: add validate for others config

	return &cfg, nil
}
