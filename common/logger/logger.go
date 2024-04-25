package logger

import "fmt"

var logger *Logger

type Logger struct{}

func (l *Logger) New() *Logger {
	if logger != nil {
		return logger
	}

	logger = &Logger{}

	return logger
}

func Log(msg any) {
	fmt.Println(msg)
}

func Debug(msg any) {
	fmt.Println(msg)
}

func Error(msg any) {
	fmt.Println(msg)
}
