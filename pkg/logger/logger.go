package logger

import (
	"den-den-mushi-Go/internal/config"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfg *config.Config) *zap.Logger {
	level := setLogLevel(cfg)
	encoder := setLogFormat(cfg)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		level)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func setLogLevel(cfg *config.Config) zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		level = zapcore.InfoLevel
	}
	return level
}

func setLogFormat(cfg *config.Config) zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.CallerKey = "caller"

	if strings.ToLower(cfg.Logging.Format) == "json" {
		return zapcore.NewJSONEncoder(encoderCfg)
	}

	if cfg.App.Environment == "development" {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return zapcore.NewConsoleEncoder(encoderCfg)
}
