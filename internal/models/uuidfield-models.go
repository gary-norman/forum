package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type AnyID interface {
	UUIDField | int64
}

type UUIDField struct {
	UUID uuid.UUID
}

// Automatically generate a new UUID if it's not already set
func NewUUIDField() UUIDField {
	return UUIDField{UUID: uuid.New()}
}

// Implement fmt.Stringer
func (u UUIDField) String() string {
	return u.UUID.String()
}

// JSON marshalling
func (u UUIDField) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.UUID.String())
}

func (u *UUIDField) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsed, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	u.UUID = parsed
	return nil
}

// SQL driver interfaces
func (u *UUIDField) copyFromBytes(src any) error {
	switch v := src.(type) {
	case []byte:
		copy(u.UUID[:], v)
		return nil
	default:
		return fmt.Errorf("UUIDField: cannot scan type %T", v)
	}
}

func (u *UUIDField) Exec(value any) error {
	return u.copyFromBytes(value)
}

func (u *UUIDField) Scan(value any) error {
	return u.copyFromBytes(value)
}

func (u *UUIDField) Begin(value any) error {
	return u.copyFromBytes(value)
}

func (u *UUIDField) Commit(value any) error {
	return u.copyFromBytes(value)
}

func (u UUIDField) Value() (driver.Value, error) {
	return u.UUID[:], nil // store as []byte
}

func UUIDFieldFromString(s string) (UUIDField, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return UUIDField{}, err
	}
	return UUIDField{UUID: parsed}, nil
}
