package logger

import (
	"errors"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOG_FILENAME = "log.log"

// Fluire + Logger = FLogger :)
// Customized zap.Logger
type FLogger struct {
	logger *zap.Logger
}

func NewFLogger() (*FLogger, error) {
	logger, err := initZapLogger()

	if err != nil {
		return nil, errors.New("Can't initialize logger: " + err.Error())
	}

	return &FLogger{logger: logger}, nil
}

// Error is used to log error-sort events
func (l *FLogger) Error(msg string) {
	l.logger.Error(msg)
}

// Info is used to log informational messages
func (l *FLogger) Info(msg string) {
	l.logger.Info(msg)
}

// initZapLogger customizes zap.Logger and returns, generally, FLogger
func initZapLogger() (*zap.Logger, error) {
	configConsole := zap.NewProductionEncoderConfig()
	configConsole.EncodeTime = ConsoleTimeEncoder
	configConsole.EncodeCaller = zapcore.ShortCallerEncoder

	configFile := zap.NewProductionEncoderConfig()
	configFile.EncodeTime = FileTimeEncoder
	configFile.EncodeCaller = zapcore.ShortCallerEncoder

	fileEncoder := zapcore.NewJSONEncoder(configFile)
	consoleEncoder := zapcore.NewConsoleEncoder(configConsole)

	logFile, err := os.OpenFile(LOG_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}

func ConsoleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("15:04:05") + "]")
}

func FileTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan 01, 2006  15:04:05"))
}
