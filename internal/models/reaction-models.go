package models

import "time"

type Reaction struct {
	ID               int64     `db:"id"`
	Liked            bool      `db:"liked"`
	Disliked         bool      `db:"disliked"`
	Created          time.Time `db:"created"`
	Updated          time.Time `db:"updated"`
	AuthorID         UUIDField `db:"authorId"`
	PostID           int64
	CommentID        int64
	ReactedPostID    *int64 `db:"reactedPostId,omitempty"`
	ReactedCommentID *int64 `db:"reactedCommentId,omitempty"`
}

type ReactionInput struct {
	Liked            bool   `json:"liked"`
	Disliked         bool   `json:"disliked"`
	AuthorID         string `json:"authorId"` // Convert manually
	ReactedPostID    *int64 `json:"reactedPostId,omitempty"`
	ReactedCommentID *int64 `json:"reactedCommentId,omitempty"`
}

func (r Reaction) TableName() string { return "reactions" }
func (r Reaction) GetID() int64      { return r.ID }
func (r *Reaction) SetID(id int64)   { r.ID = id }
