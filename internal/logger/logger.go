package logger

//go:generate mockgen -source=logger.go -destination=mocks/mock.go

type Logger interface {
	Error(msg string)
	Info(msg string)
}

func NewLogger() (Logger, error) {
	return NewFLogger()
}
