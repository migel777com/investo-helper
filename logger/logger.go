package logger

import (
	"os"
	"time"
)

type Logger struct {
	Path string
	File os.File
}

func RegisterLogger(path string) Logger {
	f, _ := os.Create(path)
	logger := Logger{Path: path, File: *f}

	return logger
}

func (l *Logger) MakeLog(message string){
	t := time.Now().String()
	_, _ = l.File.WriteString(t + ": " + message + "\r\n")
}
