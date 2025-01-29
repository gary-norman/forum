package models

import (
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Login
	Avatar      string    `json:"avatar,omitempty"` // Store UUID as string
	Banner      string    `json:"banner,omitempty"` // Store UUID as string
	Description string    `json:"description,omitempty"`
	Usertype    string    `json:"usertype"`
	Created     time.Time `json:"created"`
	TimeSince   string    `json:"time_since"`
	IsFlagged   bool      `json:"is_flagged,omitempty"`
}
type UserCheck struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Email          string
	HashedPassword string
}
type Login struct {
	Email          string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
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
	OwnerID     int       `json:"owner_id"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar,omitempty"` // Store UUID as string
	Banner      string    `json:"banner,omitempty"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Rules       []ChannelRule
	Privacy     bool `json:"privacy"`
	IsMuted     bool `json:"is_muted"`
	IsFLagged   bool `json:"is_flagged,omitempty"`
}
type ChannelData struct {
	ChannelName string `json:"channelName"`
	ChannelID   string `json:"channelID"`
}

type ChannelRule struct {
	ID        int    `json:"id"`
	Rule      string `json:"rule"`
	ChannelID int    `json:"channel_id"`
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
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Images        string    `json:"images,omitempty"` // Store JSON as string
	Created       time.Time `json:"created"`
	TimeSince     string    `json:"time_since"`
	IsCommentable bool      `json:"commentable"`
	Author        string    `json:"author"`
	AuthorID      int       `json:"author_id"`
	AuthorAvatar  string    `json:"author_avatar"`
	ChannelID     int       `json:"channel_id"`
	ChannelName   string    `json:"channel_name"`
	IsFlagged     bool      `json:"is_flagged,omitempty"`
	Likes         int       `json:"likes"`
	Dislikes      int       `json:"dislikes"`
}

type PostWithDaysAgo struct {
	Post
	TimeSince string
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
	CommentedPostID    *int      `json:"commented_post_id,omitempty"`
	CommentedCommentID *int      `json:"commented_comment_id,omitempty"`
	IsCommentable      bool      `json:"is_commentable"`
	IsFlagged          bool      `json:"is_flagged,omitempty"`
}

type Reaction struct {
	ID               int       `json:"id"`
	Liked            bool      `json:"liked"`
	Disliked         bool      `json:"disliked"`
	Created          time.Time `json:"created"`
	AuthorID         int       `json:"author_id"`
	ReactedPostID    *int      `json:"reacted_post_id,omitempty"`
	ReactedCommentID *int      `json:"reacted_comment_id,omitempty"`
}

type PostReaction struct {
	ID         int  `json:"id"`
	UserID     *int `json:"user_id"`
	PostID     *int `json:"post_id,omitempty"`
	ReactionID *int `json:"reaction_id,omitempty"`
}

type CommentReaction struct {
	ID         int  `json:"id"`
	UserID     *int `json:"user_id"`
	PostID     *int `json:"comment_id,omitempty"`
	ReactionID *int `json:"reaction_id,omitempty"`
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

type Notify struct {
	BadPass      string
	RegisterOk   string
	RegisterFail string
	BadLogin     string
	LoginOk      string
	LoginFail    string
}

type NotifyPlaceholder struct {
	Register string
	Login    string
}

type TemplateData struct {
	CurrentUser       *User             `json:"user"`
	CurrentUserName   string            `json:"currentUserName"`
	Channels          []Channel         `json:"channels"`
	OwnedChannels     []Channel         `json:"ownedChannels"`
	Posts             []PostWithDaysAgo `json:"posts"`
	Avatar            string            `json:"avatar"`
	Bio               string            `json:"bio"`
	Images            []Image           `json:"images"`
	Comments          []Comment         `json:"comments"`
	Reactions         []Reaction        `json:"reactions"`
	NotifyPlaceholder `json:"notifyPlaceholder"`
}

type Session struct {
	Username string
	Expires  time.Time
}
