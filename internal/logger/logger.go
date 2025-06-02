package logger

import (
	"os"
	"time"

	"github.com/red3533/log-analyzer/internal/models"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger(cfg models.LoggerConfig) Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
		FormatLevel: func(i interface{}) string {
			return "[" + i.(string) + "]"
		},
	}

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   true,
	}

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	baseLogger := zerolog.New(multiWriter).
		Level(level).
		With().
		Timestamp().
		Str("app", "log-analyzer").
		Int("pid", os.Getpid()).
		Logger()

	globalLogger := Logger{baseLogger}

	return globalLogger
}
