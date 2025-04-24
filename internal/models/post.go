package models

import (
	"time"
)

type Post struct {
	ID            int64     `db:"id,primary"`
	Title         string    `db:"title"`
	Content       string    `db:"content"`
	Images        string    `db:"images,omitempty"`
	Created       time.Time `db:"created"`
	TimeSince     string
	IsCommentable bool   `db:"commentable"`
	Author        string `db:"author"`
	AuthorID      int64  `db:"authorId"`
	AuthorAvatar  string `db:"authorAvatar"`
	ChannelID     int64  `db:"channelId"`
	ChannelName   string `db:"channelName"`
	IsFlagged     bool   `db:"isFlagged,omitempty"`
	Likes         int    `db:"likes"`
	Dislikes      int    `db:"dislikes"`
	CommentsCount int    `db:"commentsCount"`
	Comments      []Comment
}

func (*Post) TableName() string { return "posts" }
func (p Post) GetID() int64     { return p.ID }
func (p *Post) SetID(id int64)  { p.ID = id }

func (p *Post) UpdateTimeSince() {
	p.TimeSince = getTimeSince(p.Created)
}

func (p *Post) React(likes, dislikes int) {
	p.Likes += likes
	p.Dislikes += dislikes
}

type PostPage struct {
	UserID      int64
	CurrentUser *User
	Instance    string
	ThisPost    Post
	OwnerName   string
	ImagePaths
}

func (p PostPage) GetInstance() string { return p.Instance }
