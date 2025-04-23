package models

type DBModel interface {
	TableName() string
	GetID() int64
	SetID(int64)
}
