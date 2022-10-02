package logger

import (
	"io"
	"log"
)

const (
	logOff   = "OFF"
	logError = "ERROR"
	logWarn  = "WARN"
	logInfo  = "INFO"
	logDebug = "DEBUG"
)

var levels = map[string]int{
	logOff:   0,
	logError: 100,
	logWarn:  200,
	logInfo:  300,
	logDebug: 400,
}

type Logger struct {
	level int
	ex    *log.Logger
}

func New(level string) *Logger {
	logVal, ok := levels[level]
	if !ok {
		logVal = levels[logOff]
	}

	return &Logger{
		level: logVal,
		ex:    log.Default(),
	}
}

func (l *Logger) Info(msg string) {
	if l.IfLevelEnabled(logInfo) {
		l.preparation("log.INFO ")
		l.ex.Println(msg)
	}
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	if l.IfLevelEnabled(logInfo) {
		l.ex.SetFlags(0)
		l.ex.SetPrefix("")
		l.ex.Printf(msg, args...)
	}
}

func (l *Logger) Error(msg string) {
	if l.IfLevelEnabled(logError) {
		l.preparation("log.ERROR ")
		l.ex.Println(msg)
	}
}

func (l *Logger) Warn(msg string) {
	if l.IfLevelEnabled(logWarn) {
		l.preparation("log.WARN ")
		l.ex.Println(msg)
	}
}

func (l *Logger) Debug(msg string) {
	if l.IfLevelEnabled(logDebug) {
		l.preparation("log.DEBUG ")
		l.ex.Println(msg)
	}
}

func (l *Logger) IfLevelEnabled(level string) bool {
	return l.level >= levels[level]
}

func (l *Logger) preparation(str string) {
	l.ex.SetFlags(log.LstdFlags | log.Lmsgprefix)
	l.ex.SetPrefix(str)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.ex.SetOutput(w)
}
