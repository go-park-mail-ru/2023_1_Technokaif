package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Fluire + Logger = FLogger :)
// Customized minimalistic zap.Logger
type FLogger struct {
	logger      *zap.Logger
	reqIdGetter ReqIDGetter
}

func NewFLogger(getter ReqIDGetter) (*FLogger, error) {
	logger, err := initZapLogger()

	if err != nil {
		return nil, fmt.Errorf("can't initialize logger: %w", err)
	}

	return &FLogger{
		logger:      logger,
		reqIdGetter: getter,
	}, nil
}

// Error is used to log error-sort events
func (l *FLogger) Error(msg string) {
	l.logger.Error(msg)
}

// Errorf is used to log error-sort events with fromat string
func (l *FLogger) Errorf(format string, a ...any) {
	l.logger.Error(fmt.Sprintf(format, a...))
}

// Info is used to log informational messages
func (l *FLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *FLogger) Infof(format string, a ...any) {
	l.logger.Info(fmt.Sprintf(format, a...))
}

func (l *FLogger) ErrorReqID(ctx context.Context, msg string) {
	reqId, err := l.reqIdGetter(ctx)
	if err != nil {
		l.Error(msg)
		return
	}

	l.Errorf("ReqID:%d %s", reqId, msg)
}

func (l *FLogger) ErrorfReqID(ctx context.Context, format string, a ...any) {
	reqId, err := l.reqIdGetter(ctx)
	if err != nil {
		l.Errorf(format, a...)
		return
	}

	l.Errorf(fmt.Sprintf("ReID:%d ", reqId)+format, a...)
}

func (l *FLogger) InfoReqID(ctx context.Context, msg string) {
	reqId, err := l.reqIdGetter(ctx)
	if err != nil {
		l.Info(msg)
		return
	}

	l.Infof("ReqID:%d %s", reqId, msg)
}

func (l *FLogger) InfofReqID(ctx context.Context, format string, a ...any) {
	reqId, err := l.reqIdGetter(ctx)
	if err != nil {
		l.Infof(format, a...)
		return
	}

	l.Infof(fmt.Sprintf("ReqID:%d ", reqId)+format, a...)
}

// initZapLogger customizes zap.Logger and returns, generally, FLogger
func initZapLogger() (*zap.Logger, error) {
	configConsole := zap.NewProductionEncoderConfig()
	configConsole.EncodeTime = consoleTimeEncoder
	configConsole.EncodeCaller = zapcore.ShortCallerEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(configConsole)

	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}

func consoleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("15:04:05") + "]")
}
