package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log   *zap.SugaredLogger
	level zapcore.Level
)

// Init configures the global logger based on the LOG_LEVEL environment variable.
func Init() error {
	level = zap.InfoLevel
	if lvl := strings.ToLower(os.Getenv("LOG_LEVEL")); lvl != "" {
		if err := level.UnmarshalText([]byte(lvl)); err != nil {
			// keep default if parsing fails
		}
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	l, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = l.Sugar()
	return nil
}

// Level returns the configured log level.
func Level() zapcore.Level {
	return level
}
