package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeUpdatable is an interface for types that can update their TimeSince field.
type TimeUpdatable interface {
	UpdateTimeSince()
}

// Reactable is an interface for types that can be reacted to with likes and dislikes.
type Reactable interface {
	React(likes, dislikes int)
}

// UpdateTimeSince updates the TimeSince field of any struct that implements the TimeUpdatable interface.
func UpdateTimeSince(t TimeUpdatable) {
	t.UpdateTimeSince()
}

// React updates the likes and dislikes of any struct that implements the Reactable interface.
func React(r Reactable, likes, dislikes int) {
	r.React(likes, dislikes)
}

// getTimeSince returns a human-readable string representing the time elapsed since the provided creation time.
func getTimeSince(created time.Time) string {
	now := time.Now()
	hours := now.Sub(created).Hours()
	var timeSince string
	if hours > 24 {
		timeSince = fmt.Sprintf("%.0f days ago", hours/24)
	} else if hours > 1 {
		timeSince = fmt.Sprintf("%.0f hours ago", hours)
	} else if minutes := now.Sub(created).Minutes(); minutes > 1 {
		timeSince = fmt.Sprintf("%.0f minutes ago", minutes)
	} else {
		timeSince = "just now"
	}
	return timeSince
}

// GetIntFromPathValue parses a string value from the URL path and converts it to int64, returning an error if conversion fails.
func GetIntFromPathValue(value string) (int64, error) {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

// ToHTMLVar converts a string like "post-page" to "postsHTML"
func ToHTMLVar(s string) string {
	parts := strings.Split(s, "-")
	if len(parts) == 0 {
		return ""
	}
	// Add "s" for plural and append "HTML"

	return parts[0] + "sHTML"
}

func NotFoundLocation(location string) NotFound {
	data := NotFound{
		Instance: location + "-page",
		Location: "not-found",
		Message:  location,
	}
	return data
}
