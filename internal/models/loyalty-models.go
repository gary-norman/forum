package models

type Loyalty struct {
	ID       int64 `db:"id"`
	Follower int64 `db:"follower"`
	Followee int64 `db:"followee"`
}

func (l Loyalty) TableName() string { return "loyalty" }
func (l Loyalty) GetID() int64      { return l.ID }
func (l *Loyalty) SetID(id int64)   { l.ID = id }
