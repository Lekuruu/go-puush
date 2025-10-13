package app

import (
	"fmt"
	"time"
)

type Logger struct {
	Name string
}

func NewLogger(name string) *Logger {
	return &Logger{Name: name}
}

func (l *Logger) Write(bs []byte) (int, error) {
	now := time.Now().UTC().Format(time.DateTime)
	logMessage := fmt.Sprintf("[%s] <%s> - %s", now, l.Name, bs)
	return fmt.Print(logMessage)
}

func (l *Logger) Log(args ...interface{}) {
	logMessage := fmt.Sprint(args...)
	l.Write([]byte(logMessage + "\n"))
}

func (l *Logger) Logf(format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	l.Write([]byte(logMessage + "\n"))
}
