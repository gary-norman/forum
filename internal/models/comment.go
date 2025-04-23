package models

import (
	"time"
)

type Comment struct {
	ID                 int64     `db:"id"`
	Content            string    `db:"content"`
	Created            time.Time `db:"created"`
	TimeSince          string
	Author             string `db:"author"`
	AuthorID           int64  `db:"author_id"`
	AuthorAvatar       string `db:"author_avatar"`
	ChannelID          int64  `db:"channel_id"`
	ChannelName        string `db:"channel_name"`
	CommentedPostID    *int64 `db:"commented_post_id,omitempty"`
	CommentedCommentID *int64 `db:"commented_comment_id,omitempty"`
	IsCommentable      bool   `db:"is_commentable"`
	IsReply            bool   `db:"is_reply"`
	IsFlagged          bool   `db:"is_flagged,omitempty"`
	Likes              int    `db:"likes"`
	Dislikes           int    `db:"dislikes"`
	Comments           []Comment
	Replies            []Comment
}

func (c Comment) TableName() string { return "comments" }
func (c Comment) GetID() int64      { return c.ID }
func (c *Comment) SetID(id int64)   { c.ID = id }

func (c *Comment) React(likes, dislikes int) {
	c.Likes += likes
	c.Dislikes += dislikes
}

func (c *Comment) UpdateTimeSince() {
	c.TimeSince = getTimeSince(c.Created)
}
