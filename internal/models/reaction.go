package models

import "time"

type Reaction struct {
	ID               int64     `db:"id"`
	Liked            bool      `db:"liked"`
	Disliked         bool      `db:"disliked"`
	Created          time.Time `db:"created"`
	AuthorID         int64     `db:"authorId"`
	ReactedPostID    *int64    `db:"reactedPostId,omitempty"`
	ReactedCommentID *int64    `db:"reactedCommentId,omitempty"`
}

func (r Reaction) TableName() string { return "reactions" }
func (r Reaction) GetID() int64      { return r.ID }
func (r *Reaction) SetID(id int64)   { r.ID = id }
