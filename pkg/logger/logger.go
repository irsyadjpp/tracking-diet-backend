package logger

import "log"

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func (l *Logger) Error(msg string, err error) {
	log.Printf("[ERROR] %s: %v\n", msg, err)
}
