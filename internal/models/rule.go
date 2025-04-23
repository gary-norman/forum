package models

import (
	"time"
)

type Rule struct {
	ID         int64     `db:"id"`
	Rule       string    `db:"rule"`
	Created    time.Time `db:"created"`
	Predefined bool      `db:"predefined"`
}

func (r Rule) TableName() string { return "rules" }
func (r Rule) GetID() int64      { return r.ID }
func (r *Rule) SetID(id int64)   { r.ID = id }

type PostRule struct {
	ID      int64     `db:"id"`
	RuleID  string    `db:"rule"`
	PostID  int64     `db:"postId"`
	Created time.Time `db:"created"`
}

func (pr PostRule) TableName() string { return "post_rules" }
func (pr PostRule) GetID() int64      { return pr.ID }
func (pr *PostRule) SetID(id int64)   { pr.ID = id }

type ChannelRule struct {
	ID        int64     `db:"id"`
	RuleID    string    `db:"rule"`
	ChannelID int64     `db:"channelId"`
	Created   time.Time `db:"created"`
}

func (cr ChannelRule) TableName() string { return "channel_rules" }
func (cr ChannelRule) GetID() int64      { return cr.ID }
func (cr *ChannelRule) SetID(id int64)   { cr.ID = id }
