package models

import (
	"time"

	dbutil "github.com/gary-norman/forum/internal/dbutils"
)

type User struct {
	ID       UUIDField `db:"id"`
	SortID   int       `db:"sortId"`
	Username string    `db:"username"`
	Login
	Avatar      string    `db:"avatar,omitempty"`
	Banner      string    `db:"banner,omitempty"`
	Description string    `db:"description,omitempty"`
	Usertype    string    `db:"usertype"`
	Created     time.Time `db:"created"`
	TimeSince   string
	IsFlagged   bool `db:"isFlagged,omitempty"`
	Followers   int  `db:"followers"`
	Following   int  `db:"following"`
}

func (u User) TableName() string   { return "users" }
func (u User) GetID() UUIDField    { return u.ID }
func (u *User) SetID(id UUIDField) { u.ID = id }

func (u *User) UpdateTimeSince() {
	u.TimeSince = getTimeSince(u.Created)
}

type UserCheck struct {
	ID             dbutil.UUID `db:"id"`
	Username       string      `db:"username"`
	Email          string
	HashedPassword string
}

type UserPage struct {
	UserID      UUIDField
	CurrentUser *User
	Instance    string
	ThisUser    User
	OwnerName   string
	Posts       []Post
	ImagePaths
}

func (p UserPage) GetInstance() string { return p.Instance }
