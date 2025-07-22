package logger

import (
	"den-den-mushi-Go/pkg/config"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfg config.Logger) *zap.Logger {
	level := getLogLevel(cfg)
	encoder := getLogEncoder(cfg)

	var cores []zapcore.Core

	// stdout writer
	if cfg.Output == "stdout" || cfg.Output == "both" {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// file writer
	if cfg.Output == "file" || cfg.Output == "both" {
		logPath := cfg.FilePath
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic("Failed to create log directory: " + err.Error())
		}

		file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			panic("Failed to open log file: " + err.Error())
		}
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(file),
			level,
		))
	}

	// todo: addSync() to graylog

	if cfg.Environment == "prod" {
		return zap.New(
			zapcore.NewTee(cores...),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel))
	} else if cfg.Environment == "dev" {
		return zap.New(
			zapcore.NewTee(cores...),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.Development())
	} else {
		return nil
	}
}

func getLogLevel(cfg config.Logger) zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return zapcore.InfoLevel
	}
	return level
}

func getLogEncoder(cfg config.Logger) zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.CallerKey = "caller"
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "msg"

	if strings.ToLower(cfg.Format) == "json" {
		return zapcore.NewJSONEncoder(encoderCfg)
	} else if strings.ToLower(cfg.Format) == "console" {
		if cfg.Environment == "dev" {
			encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		return zapcore.NewConsoleEncoder(encoderCfg)
	}

	return nil
}
