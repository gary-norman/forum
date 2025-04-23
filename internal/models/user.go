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
	TimeSince   string
	IsFlagged   bool `json:"isFlagged,omitempty"`
	Followers   int  `json:"followers"`
	Following   int  `json:"following"`
}

func (u *User) UpdateTimeSince() {
	u.TimeSince = getTimeSince(u.Created)
}

type UserCheck struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
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
