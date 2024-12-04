package models

import "time"

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email_address"`
	Avatar      string    `json:"avatar,omitempty"` // Store UUID as string
	Banner      string    `json:"banner,omitempty"` // Store UUID as string
	Description string    `json:"description,omitempty"`
	Usertype    string    `json:"usertype"`
	Created     time.Time `json:"created"`
	IsFlagged   bool      `json:"is_flagged,omitempty"`
}

type Bookmark struct {
	ID        int       `json:"id"`
	PostID    *int      `json:"post_id,omitempty"`
	CommentID *int      `json:"comment_id,omitempty"`
	ChannelID *int      `json:"channel_id,omitempty"`
	Created   time.Time `json:"created"`
}

type Loyalty struct {
	ID       int `json:"id"`
	Follower int `json:"follower"`
	Followee int `json:"followee"`
}

type Channel struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar,omitempty"` // Store UUID as string
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Rules       string    `json:"rules,omitempty"`
	Privacy     bool      `json:"privacy"`
}

type MutedChannel struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ChannelID int       `json:"channel_id"`
	Created   time.Time `json:"created"`
}

type Membership struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ChannelID int       `json:"channel_id"`
	Created   time.Time `json:"created"`
}

type Mod struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ChannelID int `json:"channel_id"`
}

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Images      string    `json:"images,omitempty"` // Store JSON as string
	Created     time.Time `json:"created"`
	Commentable bool      `json:"commentable"`
	AuthorID    int       `json:"author_id"`
	ChannelID   int       `json:"channel_id"`
	IsFlagged   bool      `json:"is_flagged,omitempty"`
	Likes       int       `json:"likes"`
	Dislikes    int       `json:"dislikes"`
}

type Image struct {
	ID       string    `json:"id"` // UUID
	Created  time.Time `json:"created"`
	AuthorID int       `json:"author_id"`
	PostID   int       `json:"post_id"`
}

type Comment struct {
	ID                 int       `json:"id"`
	Content            string    `json:"content"`
	Created            time.Time `json:"created"`
	AuthorID           int       `json:"author_id"`
	ChannelID          int       `json:"channel_id"`
	IsReply            bool      `json:"is_reply"`
	CommentedPostID    *int      `json:"commented_post_id,omitempty"`
	CommentedCommentID *int      `json:"commented_comment_id,omitempty"`
	IsFlagged          bool      `json:"is_flagged,omitempty"`
}

type Reaction struct {
	ID               int       `json:"id"`
	Liked            bool      `json:"liked"`
	Disliked         bool      `json:"disliked"`
	Created          time.Time `json:"created"`
	AuthorID         int       `json:"author_id"`
	ChannelID        int       `json:"channel_id"`
	ReactedPostID    *int      `json:"reacted_post_id,omitempty"`
	ReactedCommentID *int      `json:"reacted_comment_id,omitempty"`
}

type Flag struct {
	ID               int       `json:"id"`
	FlagType         string    `json:"flag_type"`
	Content          string    `json:"content,omitempty"`
	Created          time.Time `json:"created"`
	Approved         bool      `json:"approved"`
	AuthorID         int       `json:"author_id"`
	ChannelID        int       `json:"channel_id"`
	FlaggedUserID    *int      `json:"flagged_user_id,omitempty"`
	FlaggedPostID    *int      `json:"flagged_post_id,omitempty"`
	FlaggedCommentID *int      `json:"flagged_comment_id,omitempty"`
}
