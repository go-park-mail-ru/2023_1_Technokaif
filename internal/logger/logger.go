package logger

//go:generate mockgen -source=logger.go -destination=mocks/mock.go

type Logger interface {
	Error(msg string)
	Errorf(format string, a ...any)
	Info(msg string)
	Infof(format string, a ...any)
}

func NewLogger() (Logger, error) {
	return NewFLogger()
}
