package logger

import (
    "log"
    "os"
)

type Logger struct {
    *log.Logger
}

var GlobalLogger *Logger

func Init() {
    GlobalLogger = &Logger{
        Logger: log.New(os.Stdout, "[WALLET-API] ", log.LstdFlags|log.Lshortfile),
    }
}

func (l *Logger) Info(format string, v ...interface{}) {
    l.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
    l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
    l.Printf("[DEBUG] "+format, v...)
}

func (l *Logger) Warning(format string, v ...interface{}) {
    l.Printf("[WARNING] "+format, v...)
}
