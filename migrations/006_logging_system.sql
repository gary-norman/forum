-- Migration: Logging and Metrics System
-- Creates tables for request logs, error logs, and system metrics

BEGIN TRANSACTION;

-- Request Logs: Track all HTTP requests
CREATE TABLE IF NOT EXISTS RequestLogs (
    ID INTEGER PRIMARY KEY,
    Timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Method TEXT NOT NULL,               -- GET, POST, etc.
    Path TEXT NOT NULL,                 -- /post/123
    StatusCode INTEGER NOT NULL,        -- 200, 404, 500, etc.
    Duration INTEGER NOT NULL,          -- Request duration in milliseconds
    UserID BLOB,                        -- UUID of authenticated user (NULL if anonymous)
    IPAddress TEXT,                     -- Client IP address
    UserAgent TEXT,                     -- Browser/client info
    Referer TEXT,                       -- Where request came from
    BytesSent INTEGER DEFAULT 0,        -- Response size in bytes

    -- Indexes for common queries
    INDEX idx_timestamp ON RequestLogs(Timestamp),
    INDEX idx_path ON RequestLogs(Path),
    INDEX idx_status ON RequestLogs(StatusCode),
    INDEX idx_user ON RequestLogs(UserID)
);

-- Error Logs: Track application errors
CREATE TABLE IF NOT EXISTS ErrorLogs (
    ID INTEGER PRIMARY KEY,
    Timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Level TEXT NOT NULL,                -- ERROR, WARN, FATAL
    Message TEXT NOT NULL,              -- Error message
    StackTrace TEXT,                    -- Full stack trace
    RequestPath TEXT,                   -- Path where error occurred
    UserID BLOB,                        -- User who encountered error
    Context TEXT,                       -- Additional context (JSON)

    INDEX idx_timestamp ON ErrorLogs(Timestamp),
    INDEX idx_level ON ErrorLogs(Level)
);

-- System Metrics: Track application health and performance
CREATE TABLE IF NOT EXISTS SystemMetrics (
    ID INTEGER PRIMARY KEY,
    Timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    MetricType TEXT NOT NULL,           -- startup, memory, concurrent_users, health_check
    MetricName TEXT NOT NULL,           -- Specific metric name
    MetricValue REAL NOT NULL,          -- Numeric value (duration, count, bytes, etc.)
    Unit TEXT NOT NULL,                 -- ms, bytes, count, etc.
    Details TEXT,                       -- Additional JSON details

    INDEX idx_timestamp ON SystemMetrics(Timestamp),
    INDEX idx_type ON SystemMetrics(MetricType),
    INDEX idx_name ON SystemMetrics(MetricName)
);

COMMIT;
