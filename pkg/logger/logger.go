package logger

import (
	"context"
)

//go:generate mockgen -source=logger.go -destination=mocks/mock.go

type ReqIDGetter func(ctx context.Context) (uint32, error)

type Logger interface {
	Error(msg string)
	Errorf(format string, a ...any)
	Info(msg string)
	Infof(format string, a ...any)

	ErrorReqID(ctx context.Context, msg string)
	ErrorfReqID(ctx context.Context, format string, a ...any)
	InfoReqID(ctx context.Context, msg string)
	InfofReqID(ctx context.Context, format string, a ...any)
}

func NewLogger(reqIdGetter ReqIDGetter) (Logger, error) {
	return NewFLogger(reqIdGetter)
}
