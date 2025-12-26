package logger

import (
	"log"
	"os"
)

type Logger struct {
    debug *log.Logger
    info  *log.Logger
    warn  *log.Logger
    error *log.Logger
}

func NewLogger(environment string) *Logger {
    flags := log.Ldate | log.Ltime | log.Lshortfile
    
    return &Logger{
        debug: log.New(os.Stdout, "DEBUG: ", flags),
        info:  log.New(os.Stdout, "INFO: ", flags),
        warn:  log.New(os.Stdout, "WARN: ", flags),
        error: log.New(os.Stderr, "ERROR: ", flags),
    }
}

func (l *Logger) Debug(format string, v ...interface{}) {
    l.debug.Printf(format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
    l.info.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
    l.warn.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
    l.error.Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.error.Printf(format, v...)
    os.Exit(1)
}

// Добавляем методы с суффиксом 'f' для совместимости
func (l *Logger) Debugf(format string, v ...interface{}) {
    l.Debug(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
    l.Info(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
    l.Warn(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
    l.Error(format, v...)
}