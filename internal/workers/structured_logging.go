package workers

import (
	"encoding/json"
	"time"

	"github.com/gary-norman/forum/internal/models"
)

// EventContext represents structured context for searchable logging
type EventContext struct {
	UserID     string                 `json:"user_id,omitempty"`     // UUID as string for JSON
	PostID     int64                  `json:"post_id,omitempty"`
	CommentID  int64                  `json:"comment_id,omitempty"`
	ChannelID  int64                  `json:"channel_id,omitempty"`
	Action     string                 `json:"action,omitempty"`      // "create", "update", "delete", "view"
	Keywords   []string               `json:"keywords,omitempty"`    // Searchable terms
	Metadata   map[string]interface{} `json:"metadata,omitempty"`    // Additional context
	IPAddress  string                 `json:"ip_address,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	Path       string                 `json:"path,omitempty"`
}

// ToJSON converts EventContext to JSON string for database storage
func (ec *EventContext) ToJSON() string {
	data, err := json.Marshal(ec)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// Structured logging helpers for common application events

// LogPostCreated logs when a user creates a post
func (pool *LoggerPool) LogPostCreated(userID models.UUIDField, postID int64, channelID int64, title string, tags []string) error {
	ctx := EventContext{
		UserID:    userID.String(),
		PostID:    postID,
		ChannelID: channelID,
		Action:    "create",
		Keywords:  tags,
		Metadata: map[string]interface{}{
			"title":       title,
			"title_len":   len(title),
			"tag_count":   len(tags),
			"entity_type": "post",
		},
	}

	return pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeUserActivity,
		MetricName:  "post_created",
		MetricValue: 1,
		Unit:        "count",
		Details:     ctx.ToJSON(),
	})
}

// LogPostViewed logs when a user views a post
func (pool *LoggerPool) LogPostViewed(userID models.UUIDField, postID int64, channelID int64) error {
	ctx := EventContext{
		UserID:    userID.String(),
		PostID:    postID,
		ChannelID: channelID,
		Action:    "view",
		Metadata: map[string]interface{}{
			"entity_type": "post",
		},
	}

	return pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeUserActivity,
		MetricName:  "post_viewed",
		MetricValue: 1,
		Unit:        "count",
		Details:     ctx.ToJSON(),
	})
}

// LogCommentCreated logs when a user creates a comment
func (pool *LoggerPool) LogCommentCreated(userID models.UUIDField, commentID int64, postID int64, contentLength int) error {
	ctx := EventContext{
		UserID:    userID.String(),
		CommentID: commentID,
		PostID:    postID,
		Action:    "create",
		Metadata: map[string]interface{}{
			"entity_type":    "comment",
			"content_length": contentLength,
		},
	}

	return pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeUserActivity,
		MetricName:  "comment_created",
		MetricValue: 1,
		Unit:        "count",
		Details:     ctx.ToJSON(),
	})
}

// LogChannelJoined logs when a user joins a channel
func (pool *LoggerPool) LogChannelJoined(userID models.UUIDField, channelID int64, channelName string) error {
	ctx := EventContext{
		UserID:    userID.String(),
		ChannelID: channelID,
		Action:    "join",
		Keywords:  []string{channelName},
		Metadata: map[string]interface{}{
			"entity_type":  "channel",
			"channel_name": channelName,
		},
	}

	return pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeUserActivity,
		MetricName:  "channel_joined",
		MetricValue: 1,
		Unit:        "count",
		Details:     ctx.ToJSON(),
	})
}

// LogUserLogin logs successful user authentication
func (pool *LoggerPool) LogUserLogin(userID models.UUIDField, username string, ipAddress string, userAgent string) error {
	ctx := EventContext{
		UserID:    userID.String(),
		Action:    "login",
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Keywords:  []string{username},
		Metadata: map[string]interface{}{
			"entity_type": "user",
			"username":    username,
		},
	}

	return pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeUserActivity,
		MetricName:  "user_login",
		MetricValue: 1,
		Unit:        "count",
		Details:     ctx.ToJSON(),
	})
}

// LogSearchQuery logs search queries for analytics
func (pool *LoggerPool) LogSearchQuery(userID models.UUIDField, query string, resultCount int) error {
	ctx := EventContext{
		UserID:   userID.String(),
		Action:   "search",
		Keywords: []string{query},
		Metadata: map[string]interface{}{
			"query":        query,
			"result_count": resultCount,
			"query_length": len(query),
		},
	}

	return pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeUserActivity,
		MetricName:  "search_query",
		MetricValue: float64(resultCount),
		Unit:        "results",
		Details:     ctx.ToJSON(),
	})
}

// LogReaction logs when a user reacts to content
func (pool *LoggerPool) LogReaction(userID models.UUIDField, postID int64, reactionType string) error {
	ctx := EventContext{
		UserID: userID.String(),
		PostID: postID,
		Action: "react",
		Metadata: map[string]interface{}{
			"entity_type":   "reaction",
			"reaction_type": reactionType,
		},
	}

	return pool.LogMetric(models.SystemMetric{
		Timestamp:   time.Now(),
		MetricType:  models.MetricTypeUserActivity,
		MetricName:  "reaction_" + reactionType,
		MetricValue: 1,
		Unit:        "count",
		Details:     ctx.ToJSON(),
	})
}

// LogApplicationError logs structured error with full context
func (pool *LoggerPool) LogApplicationError(err error, userID models.UUIDField, path string, operation string, entityIDs map[string]interface{}) error {
	ctx := EventContext{
		UserID:   userID.String(),
		Path:     path,
		Action:   operation,
		Metadata: entityIDs,
	}

	return pool.LogError(models.ErrorLog{
		Timestamp:   time.Now(),
		Level:       models.LogLevelError,
		Message:     err.Error(),
		RequestPath: path,
		UserID:      userID,
		Context:     ctx.ToJSON(),
	})
}
