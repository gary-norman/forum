package models

import "time"

type Posts struct {
	ID       int       `db:"Posts.post_id"`
	ThreadId *string   `db:"Posts.thread_id"`
	UserId   *string   `db:"Posts.user_id"`
	Content  string    `db:"Posts.content"`
	Created  time.Time `db:"Posts.created_at"`
	Title    string    `db:"Posts.title"`
	//Images      []string  `db:"Posts.images"`
	//Commentable bool      `db:"Posts.Commentable"`
}

type Users struct {
	ID                int       `db:"Users.ID"`
	Username          string    `db:"Users.Username"`
	Password          string    `db:"Users.Password"`
	EmailAddress      string    `db:"Users.Email_address"`
	Created           time.Time `db:"Users.Created"`
	Avatar            string    `db:"Users.Avatar"`
	Banner            string    `db:"Users.Banner"`
	Desc              string    `db:"Users.Desc"`
	UserType          string    `db:"Users.Type"`
	Membership        string    `db:"Users.Membership"`
	Followers         []string  `db:"Users.Followers"`
	Following         []string  `db:"Users.Following"`
	BookmarksUser     []string  `db:"Users.Bookmarks_user"`
	BookmarksChannel  []string  `db:"Users.Bookmarks_channel"`
	BookmarksPost     []string  `db:"Users.Bookmarks_post"`
	BookmarksComment  []string  `db:"Users.Bookmarks_Comment"`
	Posts             []string  `db:"Users.Posts"`
	Comments          []string  `db:"Users.Comments"`
	FlaggedItems      []string  `db:"Users.Flagged_items"`
	ModOf             []string  `db:"Users.Mod_of"`
	Reactions         []string  `db:"Users.Reactions"`
	IsFlagged         []string  `db:"Users.Is_flagged"`
	IsFlaggedApproved []string  `db:"Users.Is_flagged_approved"`
}

type Comments struct {
	ID      int       `db:"Comments.ID"`
	Content string    `db:"Comments.Content"`
	Images  []string  `db:"Comments.Images"`
	Created time.Time `db:"Comments.Created"`
}

type Reactions struct {
	ID       int       `db:"Reactions.ID"`
	Liked    bool      `db:"Reactions.Liked"`
	Disliked bool      `db:"Reactions.Disliked"`
	ParentId int       `db:"Reactions.ParentID"`
	Created  time.Time `db:"Reactions.Created"`
}
