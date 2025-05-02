package models

import (
	"time"
)

type Bookmark struct {
	ID        int64     `db:"id"`
	PostID    *int64    `db:"postId,omitempty"`
	CommentID *int64    `db:"commentId,omitempty"`
	ChannelID *int64    `db:"channelId,omitempty"`
	Created   time.Time `db:"created"`
}

func (b Bookmark) TableName() string { return "bookmarks" }
func (b Bookmark) GetID() int64      { return b.ID }
func (b *Bookmark) SetID(id int64)   { b.ID = id }
