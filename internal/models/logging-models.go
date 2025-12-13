package models

import (
	"time"
)

// RequestLog represents an HTTP request log entry
type RequestLog struct {
	ID         int64     `db:"id"`
	Timestamp  time.Time `db:"timestamp"`
	Method     string    `db:"method"`
	Path       string    `db:"path"`
	StatusCode int       `db:"statusCode"`
	Duration   int64     `db:"duration"` // Milliseconds
	UserID     UUIDField `db:"userId"`   // NULL for anonymous users
	IPAddress  string    `db:"ipAddress"`
	UserAgent  string    `db:"userAgent"`
	Referer    string    `db:"referer"`
	BytesSent  int64     `db:"bytesSent"`
}

func (r RequestLog) TableName() string { return "requestLogs" }
func (r RequestLog) GetID() int64      { return r.ID }
func (r *RequestLog) SetID(id int64)   { r.ID = id }

// ErrorLog represents an application error log entry
type ErrorLog struct {
	ID          int64     `db:"id"`
	Timestamp   time.Time `db:"timestamp"`
	Level       string    `db:"level"` // ERROR, WARN, FATAL
	Message     string    `db:"message"`
	StackTrace  string    `db:"stackTrace"`
	RequestPath string    `db:"requestPath"`
	UserID      UUIDField `db:"userId"`
	Context     string    `db:"context"` // JSON
}

func (e ErrorLog) TableName() string { return "errorLogs" }
func (e ErrorLog) GetID() int64      { return e.ID }
func (e *ErrorLog) SetID(id int64)   { e.ID = id }

// SystemMetric represents a system performance/health metric
type SystemMetric struct {
	ID          int64     `db:"id"`
	Timestamp   time.Time `db:"timestamp"`
	MetricType  string    `db:"metricType"` // startup, memory, concurrent_users, health_check
	MetricName  string    `db:"metricName"`
	MetricValue float64   `db:"metricValue"`
	Unit        string    `db:"unit"` // ms, bytes, count, etc.
	Details     string    `db:"details"` // JSON
}

func (s SystemMetric) TableName() string { return "systemMetrics" }
func (s SystemMetric) GetID() int64      { return s.ID }
func (s *SystemMetric) SetID(id int64)   { s.ID = id }

// LogLevel constants
const (
	LogLevelError = "ERROR"
	LogLevelWarn  = "WARN"
	LogLevelFatal = "FATAL"
)

// MetricType constants
const (
	MetricTypeStartup        = "startup"
	MetricTypeMemory         = "memory"
	MetricTypeConcurrentUser = "concurrent_users"
	MetricTypeHealthCheck    = "health_check"
)
