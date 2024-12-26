package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	var zapLevel zapcore.Level

	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	createLogsFolder()

	config := zap.Config{
		Level:            atomicLevel,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout", "logs/app.json"},
		ErrorOutputPaths: []string{"stderr", "logs/error.json"},
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return zapLogger, nil
}

func createLogsFolder() error {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
