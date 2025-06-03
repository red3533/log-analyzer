package models

import (
	"errors"
)

type AppConfig struct {
	LoggerConfig `yaml:"logger_config"`
}

type LoggerConfig struct {
	LogFile    string `yaml:"log_file"`
	LogLevel   string `yaml:"log_level"`
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age_days"`
}

func (c LoggerConfig) Validate() error {
	if c.MaxSizeMB <= 0 {
		return errors.New("MaxSizeMB must be positive")
	}

	if c.MaxAgeDays <= 0 {
		return errors.New("MaxAgeDays must be positive")
	}

	return nil
}
