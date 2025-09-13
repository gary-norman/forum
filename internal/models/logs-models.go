package models

import (
	"context"
	"io"
	"log"
	"os"
)

type Logs struct {
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
}

// LogContextKey is the type used for context keys
type LogContextKey string

const RequestIDKey LogContextKey = "requestID"

// CreateLogs creates a new Logs instance with optional file logging
func CreateLogs(logToFile bool, filename string) (*Logs, error) {
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if logToFile {
		logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, err
		}
		writers = append(writers, logFile)
	}

	multi := io.MultiWriter(writers...)

	logs := &Logs{
		infoLog:    log.New(multi, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warningLog: log.New(multi, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog:   log.New(io.MultiWriter(os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return logs, nil
}

// WithRequestID is a helper function that adds a Request ID to a context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// LogWithContext logs a message with a context (if request ID exists)
func (l *Logs) LogWithContext(ctx context.Context, level string, msg string) {
	requestID, _ := ctx.Value(RequestIDKey).(string)

	formatted := msg
	if requestID != "" {
		formatted = "[RequestID: " + requestID + "] " + msg
	}

	switch level {
	case "info":
		l.infoLog.Output(2, formatted)
	case "warn":
		l.warningLog.Output(2, formatted)
	case "error":
		l.errorLog.Output(2, formatted)
	default:
		l.infoLog.Output(2, formatted)
	}
}
