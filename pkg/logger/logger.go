package logger

import "net/http"

//go:generate mockgen -source=logger.go -destination=mocks/mock.go

type ReqIDGetter func(r *http.Request) (uint32, error)

type Logger interface {
	Error(msg string)
	Errorf(format string, a ...any)
	Info(msg string)
	Infof(format string, a ...any)

	ErrorReqID(r *http.Request, msg string)
	ErrorfReqID(r *http.Request, format string, a ...any)
	InfoReqID(r *http.Request, msg string)
	InfofReqID(r *http.Request, format string, a ...any)
}

func NewLogger(reqIdGetter ReqIDGetter) (Logger, error) {
	return NewFLogger(reqIdGetter)
}
