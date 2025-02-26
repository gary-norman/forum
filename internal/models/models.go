package models

import (
	"fmt"
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
	TimeSince   string
	IsFlagged   bool `json:"isFlagged,omitempty"`
	Followers   int  `json:"followers"`
	Following   int  `json:"following"`
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
	PostID    *int      `json:"postId,omitempty"`
	CommentID *int      `json:"commentId,omitempty"`
	ChannelID *int      `json:"channelId,omitempty"`
	Created   time.Time `json:"created"`
}

type Loyalty struct {
	ID       int `json:"id"`
	Follower int `json:"follower"`
	Followee int `json:"followee"`
}

type Channel struct {
	ID               int       `json:"id"`
	OwnerID          int       `json:"ownerId"`
	Name             string    `json:"name"`
	Avatar           string    `json:"avatar,omitempty"`
	Banner           string    `json:"banner,omitempty"`
	Description      string    `json:"description"`
	Created          time.Time `json:"created"`
	Rules            []Rule
	UnsubmittedRules []string
	Owned            bool
	Joined           bool
	Privacy          bool `json:"privacy"`
	IsMuted          bool `json:"isMuted"`
	IsFlagged        bool `json:"isFlagged,omitempty"`
}

type Rule struct {
	ID         int       `json:"id"`
	Rule       string    `json:"rule"`
	Created    time.Time `json:"created"`
	Predefined bool      `json:"predefined"`
}

type PostRule struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type ChannelWithDaysAgo struct {
	Channel
	TimeSince string
}
type ChannelData struct {
	ChannelName string `json:"channelName"`
	ChannelID   string `json:"channelId"`
}

type ChannelRule struct {
	ID        int    `json:"id"`
	RuleID    string `json:"rule"`
	ChannelID int    `json:"channelId"`
	Created   int    `json:"created"`
}

type MutedChannel struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ChannelID int       `json:"channelId"`
	Created   time.Time `json:"created"`
}
type OwnedAndJoinedChannels struct {
	Owned    bool
	Joined   bool
	Channels []Channel
}

type Membership struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ChannelID int       `json:"channelId"`
	Created   time.Time `json:"created"`
}

type Mod struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ChannelID int       `json:"channelId"`
	Created   time.Time `json:"created"`
}

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
	Comments      []Comment
}

//type reactable interface {
//	*Post | *Comment
//}

func getTimeSince(created time.Time) string {
	now := time.Now()
	hours := now.Sub(created).Hours()
	var timeSince string
	if hours > 24 {
		timeSince = fmt.Sprintf("%.0f days ago", hours/24)
	} else if hours > 1 {
		timeSince = fmt.Sprintf("%.0f hours ago", hours)
	} else if minutes := now.Sub(created).Minutes(); minutes > 1 {
		timeSince = fmt.Sprintf("%.0f minutes ago", minutes)
	} else {
		timeSince = "just now"
	}
	return timeSince
}

type TimeUpdatable interface {
	UpdateTimeSince()
}

func UpdateTimeSince(t TimeUpdatable) {
	t.UpdateTimeSince()
}

func (p *Post) React(likes, dislikes int) {
	p.Likes += likes
	p.Dislikes += dislikes
}
func (c *Comment) React(likes, dislikes int) {
	c.Likes += likes
	c.Dislikes += dislikes
}
func (p *Post) UpdateTimeSince() {
	p.TimeSince = getTimeSince(p.Created)
}
func (c *Comment) UpdateTimeSince() {
	c.TimeSince = getTimeSince(c.Created)
}

func (u *User) UpdateTimeSince() {
	u.TimeSince = getTimeSince(u.Created)
}

type Image struct {
	ID       string    `json:"id"`
	Created  time.Time `json:"created"`
	AuthorID int       `json:"authorId"`
	PostID   int       `json:"postId"`
}

type Comment struct {
	ID                 int       `json:"id"`
	Content            string    `json:"content"`
	Created            time.Time `json:"created"`
	TimeSince          string
	Author             string `json:"author"`
	AuthorID           int    `json:"author_id"`
	AuthorAvatar       string `json:"author_avatar"`
	ChannelID          int    `json:"channel_id"`
	ChannelName        string `json:"channel_name"`
	CommentedPostID    *int   `json:"commented_post_id,omitempty"`
	CommentedCommentID *int   `json:"commented_comment_id,omitempty"`
	IsCommentable      bool   `json:"is_commentable"`
	IsFlagged          bool   `json:"is_flagged,omitempty"`
	Likes              int    `json:"likes"`
	Dislikes           int    `json:"dislikes"`
	Comments           []Comment
	Replies            []Comment
}

type Reaction struct {
	ID               int       `json:"id"`
	Liked            bool      `json:"liked"`
	Disliked         bool      `json:"disliked"`
	Created          time.Time `json:"created"`
	AuthorID         int       `json:"authorId"`
	ReactedPostID    *int      `json:"reactedPostId,omitempty"`
	ReactedCommentID *int      `json:"reactedCommentId,omitempty"`
}

type Flag struct {
	ID               int       `json:"id"`
	FlagType         string    `json:"flagType"`
	Content          string    `json:"content,omitempty"`
	Created          time.Time `json:"created"`
	Approved         bool      `json:"approved"`
	AuthorID         int       `json:"authorId"`
	ChannelID        int       `json:"channelId"`
	FlaggedUserID    *int      `json:"flaggedUserId,omitempty"`
	FlaggedPostID    *int      `json:"flaggedPostId,omitempty"`
	FlaggedCommentID *int      `json:"flaggedCommentId,omitempty"`
}

type Notifications struct {
	ID           int       `json:"id"`
	Notification string    `json:"notification"`
	Created      time.Time `json:"created"`
	Read         bool      `json:"read"`
	Archived     bool      `json:"archived"`
}
type NotificationsUsers struct {
	ID             int `json:"id"`
	UserID         int `json:"userId"`
	NotificationID int `json:"notificationId"`
}

type Notify struct {
	BadPass      string
	RegisterOk   string
	RegisterFail string
	BadLogin     string
	LoginOk      string
	LoginFail    string
}

//Notify := models.Notify{
//BadPass:      "The passwords do not match.",
//RegisterOk:   "Registration successful.",
//RegisterFail: "Registration failed.",
//BadLogin:     "Username or email address not found",
//LoginOk:      "Logged in successfully!",
//LoginFail:    "Unable to log in.",
//}

type TemplateData struct {
	// ---------- users ----------
	AllUsers    []User
	RandomUser  User
	CurrentUser *User
	// ---------- posts ----------
	Posts []Post
	// ---------- channels ----------
	Channels                   []Channel
	AllChannels                []Channel
	ThisChannel                ChannelWithDaysAgo
	ThisChannelOwnerName       string
	ThisChannelIsOwnedOrJoined bool
	ThisChannelIsOwned         bool
	ThisChannelRules           []Rule
	ThisChannelPosts           []Post
	OwnedChannels              []Channel
	JoinedChannels             []Channel
	OwnedAndJoinedChannels     []Channel
	// ---------- misc ----------
	Images    []Image
	Reactions []Reaction
}
type Session struct {
	Username string
	Expires  time.Time
}
