package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[APP] ", log.LstdFlags),
	}
}
func (l *Logger) Println(v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.Println("[INFO]", v)
}

func (l *Logger) Error(v ...interface{}) {
	l.Println("[ERROR]", v)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Println("[DEBUG]", v)
}
