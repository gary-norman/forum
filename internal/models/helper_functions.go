package models

import (
	"fmt"
	"time"
)

type TimeUpdatable interface {
	UpdateTimeSince()
}
type Reactable interface {
	React()
}

func UpdateTimeSince(t TimeUpdatable) {
	t.UpdateTimeSince()
}
func React(r Reactable) {
	r.React()
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
