package common

var logger *Logger

type Logger struct{}

func (l *Logger) New() *Logger {
	if logger != nil {
		return logger
	}

	logger = &Logger{}

	return logger
}
