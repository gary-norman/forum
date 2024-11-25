package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type ModModel struct {
	DB *sql.DB
}

func (m *ModModel) Insert(userID, channelID int) error {
	stmt := "INSERT INTO Mods (UserID, ChannelID) VALUES (?, ?)"
	_, err := m.DB.Exec(stmt, userID, channelID)
	return err
}

func (m *ModModel) All() ([]models.Mod, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, UserID, ChannelID FROM Mods ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Mods []models.Mod
	for rows.Next() {
		p := models.Mod{}
		err = rows.Scan(&p.ID, &p.UserID, &p.ChannelID)
		if err != nil {
			return nil, err
		}
		Mods = append(Mods, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Mods, nil
}
