package sqlite

import (
	"database/sql"

	"github.com/gary-norman/forum/internal/models"
)

// LoggingModel provides database operations for all logging tables
type LoggingModel struct {
	DB *sql.DB
}

// InsertRequestLog stores an HTTP request log entry
func (m *LoggingModel) InsertRequestLog(log models.RequestLog) error {
	stmt := `INSERT INTO RequestLogs
		(Timestamp, Method, Path, StatusCode, Duration, UserID, IPAddress, UserAgent, Referer, BytesSent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	userIDBytes, err := log.UserID.Value()
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(stmt,
		log.Timestamp,
		log.Method,
		log.Path,
		log.StatusCode,
		log.Duration,
		userIDBytes,
		log.IPAddress,
		log.UserAgent,
		log.Referer,
		log.BytesSent,
	)
	return err
}

// InsertErrorLog stores an application error log entry
func (m *LoggingModel) InsertErrorLog(log models.ErrorLog) error {
	stmt := `INSERT INTO ErrorLogs
		(Timestamp, Level, Message, StackTrace, RequestPath, UserID, Context)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	userIDBytes, err := log.UserID.Value()
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(stmt,
		log.Timestamp,
		log.Level,
		log.Message,
		log.StackTrace,
		log.RequestPath,
		userIDBytes,
		log.Context,
	)
	return err
}

// InsertSystemMetric stores a system performance/health metric
func (m *LoggingModel) InsertSystemMetric(metric models.SystemMetric) error {
	stmt := `INSERT INTO SystemMetrics
		(Timestamp, MetricType, MetricName, MetricValue, Unit, Details)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt,
		metric.Timestamp,
		metric.MetricType,
		metric.MetricName,
		metric.MetricValue,
		metric.Unit,
		metric.Details,
	)
	return err
}

// GetRequestLogsSince retrieves request logs after a given timestamp
func (m *LoggingModel) GetRequestLogsSince(since string, limit int) ([]models.RequestLog, error) {
	stmt := `SELECT ID, Timestamp, Method, Path, StatusCode, Duration, UserID, IPAddress, UserAgent, Referer, BytesSent
		FROM RequestLogs
		WHERE Timestamp >= ?
		ORDER BY Timestamp DESC
		LIMIT ?`

	rows, err := m.DB.Query(stmt, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.RequestLog
	for rows.Next() {
		var log models.RequestLog
		err = rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.Method,
			&log.Path,
			&log.StatusCode,
			&log.Duration,
			&log.UserID,
			&log.IPAddress,
			&log.UserAgent,
			&log.Referer,
			&log.BytesSent,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetErrorLogsSince retrieves error logs after a given timestamp
func (m *LoggingModel) GetErrorLogsSince(since string, limit int) ([]models.ErrorLog, error) {
	stmt := `SELECT ID, Timestamp, Level, Message, StackTrace, RequestPath, UserID, Context
		FROM ErrorLogs
		WHERE Timestamp >= ?
		ORDER BY Timestamp DESC
		LIMIT ?`

	rows, err := m.DB.Query(stmt, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ErrorLog
	for rows.Next() {
		var log models.ErrorLog
		err = rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.Level,
			&log.Message,
			&log.StackTrace,
			&log.RequestPath,
			&log.UserID,
			&log.Context,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetSystemMetricsSince retrieves system metrics after a given timestamp
func (m *LoggingModel) GetSystemMetricsSince(since string, limit int) ([]models.SystemMetric, error) {
	stmt := `SELECT ID, Timestamp, MetricType, MetricName, MetricValue, Unit, Details
		FROM SystemMetrics
		WHERE Timestamp >= ?
		ORDER BY Timestamp DESC
		LIMIT ?`

	rows, err := m.DB.Query(stmt, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.SystemMetric
	for rows.Next() {
		var metric models.SystemMetric
		err = rows.Scan(
			&metric.ID,
			&metric.Timestamp,
			&metric.MetricType,
			&metric.MetricName,
			&metric.MetricValue,
			&metric.Unit,
			&metric.Details,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, rows.Err()
}

// GetRequestStats retrieves aggregated request statistics
type RequestStats struct {
	TotalRequests   int64
	AvgDuration     float64
	ErrorRate       float64 // Percentage of 4xx/5xx responses
	UniqueUsers     int64
	RequestsPerPath map[string]int64
}

func (m *LoggingModel) GetRequestStats(since string) (*RequestStats, error) {
	stats := &RequestStats{
		RequestsPerPath: make(map[string]int64),
	}

	// Total requests and average duration
	err := m.DB.QueryRow(`
		SELECT COUNT(*), AVG(Duration)
		FROM RequestLogs
		WHERE Timestamp >= ?`, since).Scan(&stats.TotalRequests, &stats.AvgDuration)
	if err != nil {
		return nil, err
	}

	// Error rate
	var errorCount int64
	err = m.DB.QueryRow(`
		SELECT COUNT(*)
		FROM RequestLogs
		WHERE Timestamp >= ? AND StatusCode >= 400`, since).Scan(&errorCount)
	if err != nil {
		return nil, err
	}
	if stats.TotalRequests > 0 {
		stats.ErrorRate = float64(errorCount) / float64(stats.TotalRequests) * 100
	}

	// Unique users
	err = m.DB.QueryRow(`
		SELECT COUNT(DISTINCT UserID)
		FROM RequestLogs
		WHERE Timestamp >= ? AND UserID IS NOT NULL`, since).Scan(&stats.UniqueUsers)
	if err != nil {
		return nil, err
	}

	// Requests per path
	rows, err := m.DB.Query(`
		SELECT Path, COUNT(*)
		FROM RequestLogs
		WHERE Timestamp >= ?
		GROUP BY Path
		ORDER BY COUNT(*) DESC
		LIMIT 20`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		var count int64
		if err := rows.Scan(&path, &count); err != nil {
			return nil, err
		}
		stats.RequestsPerPath[path] = count
	}

	return stats, rows.Err()
}

// CleanupOldLogs deletes logs older than the specified number of days
func (m *LoggingModel) CleanupOldLogs(daysToKeep int) error {
	cutoff := `datetime('now', '-` + string(rune(daysToKeep+'0')) + ` days')`

	// Clean request logs
	if _, err := m.DB.Exec(`DELETE FROM RequestLogs WHERE Timestamp < ` + cutoff); err != nil {
		return err
	}

	// Clean error logs
	if _, err := m.DB.Exec(`DELETE FROM ErrorLogs WHERE Timestamp < ` + cutoff); err != nil {
		return err
	}

	// Clean system metrics
	if _, err := m.DB.Exec(`DELETE FROM SystemMetrics WHERE Timestamp < ` + cutoff); err != nil {
		return err
	}

	return nil
}
