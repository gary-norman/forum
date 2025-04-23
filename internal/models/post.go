package models

import (
	"time"
)

type Post struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Images        string    `json:"images,omitempty"` // Store JSON as string
	Created       time.Time `json:"created"`
	TimeSince     string
	IsCommentable bool   `json:"commentable"`
	Author        string `json:"author"`
	AuthorID      int    `json:"authorId"`
	AuthorAvatar  string `json:"authorAvatar"`
	ChannelID     int    `json:"channelId"`
	ChannelName   string `json:"channelName"`
	IsFlagged     bool   `json:"isFlagged,omitempty"`
	Likes         int    `json:"likes"`
	Dislikes      int    `json:"dislikes"`
	CommentsCount int    `json:"commentsCount"`
	Comments      []Comment
}

func (p *Post) UpdateTimeSince() {
	p.TimeSince = getTimeSince(p.Created)
}

func (p *Post) React(likes, dislikes int) {
	p.Likes += likes
	p.Dislikes += dislikes
}

type PostPage struct {
	CurrentUser *User
	Instance    string
	ThisPost    Post
	OwnerName   string
	ImagePaths
}

func (p PostPage) GetInstance() string { return p.Instance }
