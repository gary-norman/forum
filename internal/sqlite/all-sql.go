package sqlite

import (
	"database/sql"

	"github.com/gary-norman/forum/internal/colors"
	"github.com/gary-norman/forum/internal/models"
)

type AllModel struct {
	DB *sql.DB
}

var (
	Colors, _ = colors.UseFlavor("Mocha")
	ErrorMsgs = models.CreateErrorMessages()
)
