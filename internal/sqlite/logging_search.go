package sqlite

import (
	"database/sql"
	"time"

	"github.com/gary-norman/forum/internal/models"
)

// SearchableEvent represents a structured log event with extracted JSON fields
type SearchableEvent struct {
	ID          int64
	Timestamp   time.Time
	MetricName  string
	UserID      string
	PostID      int64
	ChannelID   int64
	CommentID   int64
	Action      string
	Keywords    string // JSON array as string
	Metadata    string // Full JSON
}

// UserActivityEvent represents a user action with context
type UserActivityEvent struct {
	Timestamp   time.Time
	UserID      string
	Action      string
	EntityType  string
	EntityID    int64
	Description string
}

// GetUserActivity returns all actions by a specific user
func (m *LoggingModel) GetUserActivity(userID string, limit int) ([]UserActivityEvent, error) {
	query := `
		SELECT
			Timestamp,
			json_extract(Details, '$.user_id') as user_id,
			json_extract(Details, '$.action') as action,
			json_extract(Details, '$.metadata.entity_type') as entity_type,
			COALESCE(
				json_extract(Details, '$.post_id'),
				json_extract(Details, '$.comment_id'),
				json_extract(Details, '$.channel_id'),
				0
			) as entity_id,
			MetricName as description
		FROM SystemMetrics
		WHERE MetricType = 'user_activity'
			AND json_extract(Details, '$.user_id') = ?
		ORDER BY Timestamp DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []UserActivityEvent
	for rows.Next() {
		var e UserActivityEvent
		err := rows.Scan(
			&e.Timestamp,
			&e.UserID,
			&e.Action,
			&e.EntityType,
			&e.EntityID,
			&e.Description,
		)
		if err != nil {
			continue
		}
		events = append(events, e)
	}

	return events, nil
}

// GetPostInteractions returns all user interactions with a specific post
func (m *LoggingModel) GetPostInteractions(postID int64) ([]UserActivityEvent, error) {
	query := `
		SELECT
			Timestamp,
			json_extract(Details, '$.user_id') as user_id,
			json_extract(Details, '$.action') as action,
			json_extract(Details, '$.metadata.entity_type') as entity_type,
			json_extract(Details, '$.post_id') as entity_id,
			MetricName as description
		FROM SystemMetrics
		WHERE MetricType = 'user_activity'
			AND json_extract(Details, '$.post_id') = ?
		ORDER BY Timestamp DESC
	`

	rows, err := m.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []UserActivityEvent
	for rows.Next() {
		var e UserActivityEvent
		err := rows.Scan(
			&e.Timestamp,
			&e.UserID,
			&e.Action,
			&e.EntityType,
			&e.EntityID,
			&e.Description,
		)
		if err != nil {
			continue
		}
		events = append(events, e)
	}

	return events, nil
}

// GetChannelActivity returns all activity in a specific channel
func (m *LoggingModel) GetChannelActivity(channelID int64, limit int) ([]UserActivityEvent, error) {
	query := `
		SELECT
			Timestamp,
			json_extract(Details, '$.user_id') as user_id,
			json_extract(Details, '$.action') as action,
			json_extract(Details, '$.metadata.entity_type') as entity_type,
			json_extract(Details, '$.channel_id') as entity_id,
			MetricName as description
		FROM SystemMetrics
		WHERE MetricType = 'user_activity'
			AND json_extract(Details, '$.channel_id') = ?
		ORDER BY Timestamp DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, channelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []UserActivityEvent
	for rows.Next() {
		var e UserActivityEvent
		err := rows.Scan(
			&e.Timestamp,
			&e.UserID,
			&e.Action,
			&e.EntityType,
			&e.EntityID,
			&e.Description,
		)
		if err != nil {
			continue
		}
		events = append(events, e)
	}

	return events, nil
}

// SearchByKeyword searches logs for events matching a keyword
func (m *LoggingModel) SearchByKeyword(keyword string, limit int) ([]SearchableEvent, error) {
	query := `
		SELECT
			ID,
			Timestamp,
			MetricName,
			json_extract(Details, '$.user_id') as user_id,
			COALESCE(json_extract(Details, '$.post_id'), 0) as post_id,
			COALESCE(json_extract(Details, '$.channel_id'), 0) as channel_id,
			COALESCE(json_extract(Details, '$.comment_id'), 0) as comment_id,
			json_extract(Details, '$.action') as action,
			json_extract(Details, '$.keywords') as keywords,
			Details as metadata
		FROM SystemMetrics
		WHERE MetricType = 'user_activity'
			AND (
				Details LIKE '%' || ? || '%'
				OR json_extract(Details, '$.keywords') LIKE '%' || ? || '%'
			)
		ORDER BY Timestamp DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, keyword, keyword, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []SearchableEvent
	for rows.Next() {
		var e SearchableEvent
		var userID, action sql.NullString
		err := rows.Scan(
			&e.ID,
			&e.Timestamp,
			&e.MetricName,
			&userID,
			&e.PostID,
			&e.ChannelID,
			&e.CommentID,
			&action,
			&e.Keywords,
			&e.Metadata,
		)
		if err != nil {
			continue
		}
		if userID.Valid {
			e.UserID = userID.String
		}
		if action.Valid {
			e.Action = action.String
		}
		events = append(events, e)
	}

	return events, nil
}

// GetUserErrors returns all errors encountered by a specific user
func (m *LoggingModel) GetUserErrors(userID string, limit int) ([]models.ErrorLog, error) {
	query := `
		SELECT
			ID,
			Timestamp,
			Level,
			Message,
			COALESCE(StackTrace, '') as stack_trace,
			COALESCE(RequestPath, '') as request_path,
			UserID,
			COALESCE(Context, '') as context
		FROM ErrorLogs
		WHERE json_extract(Context, '$.user_id') = ?
			OR hex(UserID) = upper(replace(?, '-', ''))
		ORDER BY Timestamp DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, userID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var errors []models.ErrorLog
	for rows.Next() {
		var e models.ErrorLog
		var userIDBytes []byte
		err := rows.Scan(
			&e.ID,
			&e.Timestamp,
			&e.Level,
			&e.Message,
			&e.StackTrace,
			&e.RequestPath,
			&userIDBytes,
			&e.Context,
		)
		if err != nil {
			continue
		}
		errors = append(errors, e)
	}

	return errors, nil
}

// GetMostActiveUsers returns users ranked by activity count
func (m *LoggingModel) GetMostActiveUsers(days int, limit int) ([]struct {
	UserID       string
	ActivityCount int
	LastActivity time.Time
}, error) {
	query := `
		SELECT
			json_extract(Details, '$.user_id') as user_id,
			COUNT(*) as activity_count,
			MAX(Timestamp) as last_activity
		FROM SystemMetrics
		WHERE MetricType = 'user_activity'
			AND Timestamp > datetime('now', '-' || ? || ' days')
			AND json_extract(Details, '$.user_id') IS NOT NULL
			AND json_extract(Details, '$.user_id') != ''
		GROUP BY json_extract(Details, '$.user_id')
		ORDER BY activity_count DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		UserID       string
		ActivityCount int
		LastActivity time.Time
	}

	for rows.Next() {
		var r struct {
			UserID       string
			ActivityCount int
			LastActivity time.Time
		}
		err := rows.Scan(&r.UserID, &r.ActivityCount, &r.LastActivity)
		if err != nil {
			continue
		}
		results = append(results, r)
	}

	return results, nil
}

// GetPopularPosts returns posts ranked by view count
func (m *LoggingModel) GetPopularPosts(days int, limit int) ([]struct {
	PostID    int64
	ViewCount int
	LastViewed time.Time
}, error) {
	query := `
		SELECT
			json_extract(Details, '$.post_id') as post_id,
			COUNT(*) as view_count,
			MAX(Timestamp) as last_viewed
		FROM SystemMetrics
		WHERE MetricType = 'user_activity'
			AND MetricName = 'post_viewed'
			AND Timestamp > datetime('now', '-' || ? || ' days')
			AND json_extract(Details, '$.post_id') IS NOT NULL
		GROUP BY json_extract(Details, '$.post_id')
		ORDER BY view_count DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		PostID    int64
		ViewCount int
		LastViewed time.Time
	}

	for rows.Next() {
		var r struct {
			PostID    int64
			ViewCount int
			LastViewed time.Time
		}
		err := rows.Scan(&r.PostID, &r.ViewCount, &r.LastViewed)
		if err != nil {
			continue
		}
		results = append(results, r)
	}

	return results, nil
}
