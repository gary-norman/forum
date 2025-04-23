package dao

import (
	"database/sql"

	"github.com/gary-norman/forum/internal/models"
)

type DAO[T models.DBModel] struct {
	DB *sql.DB
}
