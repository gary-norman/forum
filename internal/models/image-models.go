package models

import (
	"time"
)

type Image struct {
	ID       int64     `db:"id"`
	Created  time.Time `db:"created"`
	Updated  time.Time `db:"updated"`
	AuthorID int64     `db:"authorId"`
	PostID   int64     `db:"postId"`
}

func (i Image) TableName() string { return "images" }
func (i Image) GetID() int64      { return i.ID }
func (i *Image) SetID(id int64)   { i.ID = id }

type PostImage struct {
	ID      int64 `db:"id"`
	PostID  int64 `db:"postId"`
	ImageID int64 `db:"imageId"`
}

func (ip PostImage) TableName() string { return "postImages" }
func (ip PostImage) GetID() int64      { return ip.ID }
func (ip *PostImage) SetID(id int64)   { ip.ID = id }

type ImagePaths struct {
	Channel string
	Post    string
	User    string
}
