package logger

import (
	"errors"

	"go.uber.org/zap"
)

// Can't come up with name
// type Logger interface {
// 	Error(msg string)
// }

type Logger struct {
	logger *zap.Logger
}

func NewLogger() (*Logger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, errors.New("Can't initialize logger: " + err.Error())
	}

	return &Logger{logger: logger}, nil
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}
