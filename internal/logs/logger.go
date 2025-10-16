package logs

import (
	"go.uber.org/zap"
)

func NewLogger(level string) (*zap.Logger, error) {
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level = zapLevel
	logger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
