package logger

import (
	"den-den-mushi-Go/internal/proxy/config"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfg *config.Config) *zap.Logger {
	level := getLogLevel(cfg)
	encoder := getLogEncoder(cfg)

	var cores []zapcore.Core

	// stdout writer
	if cfg.Logging.Output == "stdout" || cfg.Logging.Output == "both" {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// file writer
	if cfg.Logging.Output == "file" || cfg.Logging.Output == "both" {
		logPath := cfg.Logging.FilePath
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic("Failed to create log directory: " + err.Error())
		}

		file, err := os.OpenFile(cfg.Logging.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			panic("Failed to open log file: " + err.Error())
		}
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(file),
			level,
		))
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	if cfg.App.Environment == "development" {
		logger = logger.WithOptions(zap.Development())
	}

	return logger
}

func getLogLevel(cfg *config.Config) zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		return zapcore.InfoLevel
	}
	return level
}

func getLogEncoder(cfg *config.Config) zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.CallerKey = "caller"
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "msg"

	if strings.ToLower(cfg.Logging.Format) == "json" {
		return zapcore.NewJSONEncoder(encoderCfg)
	} else if strings.ToLower(cfg.Logging.Format) == "console" {
		if cfg.App.Environment == "development" {
			encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		return zapcore.NewConsoleEncoder(encoderCfg)
	}

	return nil
}
