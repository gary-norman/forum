package models

import "time"

type Posts struct {
	ID                *int      `db:"Posts.post_id"`
	Title             string    `db:"Posts.title"`
	Content           string    `db:"Posts.content"`
	Images            []string  `db:"Posts.images"`
	Created           time.Time `db:"Posts.created_at"`
	Commentable       bool      `db:"Posts.Commentable"`
	AuthorID          *int
	Reactions         string
	Comments          string
	ChannelID         *int
	IsFlagged         string
	IsFlaggedApproved string
}

type Users struct {
	ID                int       `db:"Users.ID"`
	Username          string    `db:"Users.Username"`
	Password          string    `db:"Users.Password"`
	EmailAddress      string    `db:"Users.Email_address"`
	Avatar            string    `db:"Users.Avatar"`
	Banner            string    `db:"Users.Banner"`
	Desc              string    `db:"Users.Description"`
	UserType          string    `db:"Users.Type"`
	Created           time.Time `db:"Users.Created"`
	Membership        string    `db:"Users.Membership"`
	Followers         []string  `db:"Users.Followers"`
	Following         []string  `db:"Users.Following"`
	BookmarksUser     []string  `db:"Users.Bookmarks_user"`
	BookmarksChannel  []string  `db:"Users.Bookmarks_channel"`
	BookmarksPost     []string  `db:"Users.Bookmarks_post"`
	BookmarksComment  []string  `db:"Users.Bookmarks_Comment"`
	Posts             []string  `db:"Users.Posts"`
	Comments          []string  `db:"Users.Comments"`
	FlaggedUsers      []string  `db:"Users.Flagged_users"`
	FlaggedPosts      []string  `db:"Users.Flagged_posts"`
	FlaggedComments   []string  `db:"Users.Flagged_comments"`
	ModOf             []string  `db:"Users.Mod_of"`
	Reactions         []string  `db:"Users.Reactions"`
	IsFlagged         []string  `db:"Users.Is_flagged"`
	IsFlaggedApproved []string  `db:"Users.Is_flagged_approved"`
}

type Comments struct {
	ID                 int       `db:"Comments.ID"`
	Content            string    `db:"Comments.Content"`
	Images             []string  `db:"Comments.Images"`
	Created            time.Time `db:"Comments.Created"`
	AuthorID           int
	Reactions          string
	Replies            string
	ChannelID          int
	IsFlagged          []string `db:"Comments.Is_flagged"`
	IsFlaggedApproved  []string `db:"Comments.Is_flagged_approved"`
	CommentedPostID    int
	CommentedCommentID int
}

type Reactions struct {
	ID               int       `db:"Reactions.ID"`
	Liked            bool      `db:"Reactions.Liked"`
	Disliked         bool      `db:"Reactions.Disliked"`
	Created          time.Time `db:"Reactions.Created"`
	ParentId         int       `db:"Reactions.ParentID"`
	AuthorID         int
	ChannelID        int
	ReactedPostID    int
	ReactedCommentID int
}

type Channels struct {
	ID      int
	Name    string
	Avatar  string
	Desc    string
	Created time.Time
	Rules   string
	Privacy bool
	Members string
	Mods    string
	Posts   string
}

type Flags struct {
	ID               int
	FlagType         string
	Content          string
	Created          time.Time
	Approved         bool
	AuthorID         int
	FlaggedUserID    int
	FlaggedPostID    int
	FlaggedCommentID int
}
