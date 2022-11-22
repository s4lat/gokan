package log

import (
	"io"
	"log"
	"os"
)

// Log - used for logging info, warnings and errors from other packages.
type Log struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger
}

// Info - logs s in Info level.
func (l *Log) Info(v ...any) {
	l.InfoLogger.Print(v...)
}

// Warning - logs s in Warning level.
func (l *Log) Warning(v ...any) {
	l.WarningLogger.Print(v...)
}

// Error - logs s in Error level.
func (l *Log) Error(v ...any) {
	l.ErrorLogger.Print(v...)
}

// Fatal - logs s in Fatal level.
func (l *Log) Fatal(v ...any) {
	l.ErrorLogger.Print(v...)
	os.Exit(1)
}

// NewLogger - returns new logger that write to w.
func NewLogger(w io.Writer) Log {
	return Log{
		InfoLogger:    log.New(w, "INFO: ", log.Ldate|log.Ltime),
		WarningLogger: log.New(w, "WARNING: ", log.Ldate|log.Ltime),
		ErrorLogger:   log.New(w, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		FatalLogger:   log.New(w, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
