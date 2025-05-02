package models

type DBModel interface {
	TableName() string
	GetID() int64
	SetID(id int64)
}
