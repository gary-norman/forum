package models

import (
	"time"
)

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Login
	Avatar      string    `db:"avatar,omitempty"` // Store UUID as string
	Banner      string    `db:"banner,omitempty"` // Store UUID as string
	Description string    `db:"description,omitempty"`
	Usertype    string    `db:"usertype"`
	Created     time.Time `db:"created"`
	TimeSince   string
	IsFlagged   bool `db:"isFlagged,omitempty"`
	Followers   int  `db:"followers"`
	Following   int  `db:"following"`
}

func (u User) TableName() string { return "users" }
func (u User) GetID() int64      { return u.ID }
func (u *User) SetID(id int64)   { u.ID = id }

func (u *User) UpdateTimeSince() {
	u.TimeSince = getTimeSince(u.Created)
}

type UserCheck struct {
	ID             int64  `db:"id"`
	Username       string `db:"username"`
	Email          string
	HashedPassword string
}

type UserPage struct {
	CurrentUser *User
	Instance    string
	ThisUser    User
	OwnerName   string
	Posts       []Post
	ImagePaths
}

func (p UserPage) GetInstance() string { return p.Instance }
