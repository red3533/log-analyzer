package models

type LoggerConfig struct {
	LogFile    string
	LogLevel   string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
}
