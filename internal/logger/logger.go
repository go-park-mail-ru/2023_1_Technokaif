package logger

import (
	"errors"

	"go.uber.org/zap"
)

type Logger interface {
	Error(msg string)
	Info(msg string)
}

type FLogger struct {
	logger *zap.Logger
}

func NewFLogger() (*FLogger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, errors.New("Can't initialize logger: " + err.Error())
	}

	return &FLogger{logger: logger}, nil
}

func (l *FLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *FLogger) Info(msg string) {
	l.logger.Info(msg)
}
