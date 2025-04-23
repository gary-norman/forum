package models

import "time"

type Flag struct {
	ID               int64     `db:"id"`
	FlagType         string    `db:"flagType"`
	Content          string    `db:"content,omitempty"`
	Created          time.Time `db:"created"`
	Approved         bool      `db:"approved"`
	AuthorID         int64     `db:"authorId"`
	ChannelID        int64     `db:"channelId"`
	FlaggedUserID    *int64    `db:"flaggedUserId,omitempty"`
	FlaggedPostID    *int64    `db:"flaggedPostId,omitempty"`
	FlaggedCommentID *int64    `db:"flaggedCommentId,omitempty"`
}

func (f Flag) TableName() string { return "flags" }
func (f Flag) GetID() int64      { return f.ID }
func (f *Flag) SetID(id int64)   { f.ID = id }
