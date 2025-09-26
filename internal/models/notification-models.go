package models

import (
	"time"
)

type Notification struct {
	ID           int64     `db:"id"`
	Notification string    `db:"notification"`
	Created      time.Time `db:"created"`
	Updated      time.Time `db:"updated"`
	Read         bool      `db:"read"`
	Archived     bool      `db:"archived"`
}

func (n Notification) TableName() string { return "notifications" }
func (n Notification) GetID() int64      { return n.ID }
func (n *Notification) SetID(id int64)   { n.ID = id }

type NotificationUsers struct {
	ID             int64 `db:"id"`
	UserID         int64 `db:"userId"`
	NotificationID int64 `db:"notificationId"`
}

func (nu NotificationUsers) TableName() string { return "notificationUsers" }
func (nu NotificationUsers) GetID() int64      { return nu.ID }
func (nu *NotificationUsers) SetID(id int64)   { nu.ID = id }

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
