package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeUpdatable interface {
	UpdateTimeSince()
}
type Reactable interface {
	React(likes, dislikes int)
}

func UpdateTimeSince(t TimeUpdatable) {
	t.UpdateTimeSince()
}

func React(r Reactable, likes, dislikes int) {
	r.React(likes, dislikes)
}

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

func GetIntFromPathValue(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func ToHTMLVar(s string) string {
	parts := strings.Split(s, "-")
	if len(parts) == 0 {
		return ""
	}
	// Add "s" for plural and append "HTML"
	return parts[0] + "sHTML"
}
