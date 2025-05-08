package dbutil

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func (u *UUID) copyFromBytes(value any) error {
	bytes, ok := value.([]byte)
	if !ok || len(bytes) != 16 {
		return fmt.Errorf("invalid UUID format: %v", value)
	}
	copy((*u)[:], bytes)
	return nil
}

func (u *UUID) Exec(value any) error {
	return u.copyFromBytes(value)
}

func (u *UUID) Scan(value any) error {
	return u.copyFromBytes(value)
}

func (u *UUID) Begin(value any) error {
	return u.copyFromBytes(value)
}

func (u *UUID) Commit(value any) error {
	return u.copyFromBytes(value)
}

func (u UUID) Value() (driver.Value, error) {
	return u[:], nil
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func NewUUID() UUID {
	return UUID(uuid.New())
}
