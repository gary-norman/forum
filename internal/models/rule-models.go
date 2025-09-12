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

// TODO figure out the use of strings/int64s for this

type PostRule struct {
	ID      string    `db:"id"`
	Rule    string    `db:"rule"`
	RuleID  string    `db:"rule"`
	PostID  string    `db:"postId"`
	Created time.Time `db:"created"`
}

func (pr PostRule) TableName() string { return "postRules" }
func (pr PostRule) GetID() string     { return pr.ID }
func (pr *PostRule) SetID(id string)  { pr.ID = id }

type ChannelRule struct {
	ID        int64     `db:"id"`
	RuleID    string    `db:"rule"`
	ChannelID int64     `db:"channelId"`
	Created   time.Time `db:"created"`
}

func (cr ChannelRule) TableName() string { return "channel_rules" }
func (cr ChannelRule) GetID() int64      { return cr.ID }
func (cr *ChannelRule) SetID(id int64)   { cr.ID = id }
