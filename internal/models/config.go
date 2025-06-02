package models

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

//TODO: add validation method for app config
