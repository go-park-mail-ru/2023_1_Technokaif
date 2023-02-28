package logger

type Logger interface {
	Error(msg string)
	Info(msg string)
}

func NewLogger() (Logger, error) {
	return NewFLogger()
}
