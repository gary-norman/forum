package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type ModModel struct {
	DB *sql.DB
}

func (m *ModModel) Insert(userID, channelID int) error {
	stmt := "INSERT INTO Mods (UserID, ChannelID, Created) VALUES (?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, userID, channelID)
	return err
}

func (m *ModModel) All() ([]models.Mod, error) {
	stmt := "SELECT ID, UserID, ChannelID FROM Mods ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
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

type ModIds struct {
	UserID    int `json:"userId"`
	ChannelID int `json:"channelId"`
}

// GetModOrModdedChannelIDs returns a slice of UserID if the ChannelID is provided, and vice-versa
func (m *ModModel) GetModOrModdedChannelIDs(ID int, column string) ([]int, error) {
	validColumns := map[string]bool{
		"ID":        true,
		"UserID":    true,
		"ChannelID": true,
		"Created":   true,
	}
	stmt := fmt.Sprintf("SELECT UserID, ChannelID FROM Mods WHERE %s = ?", column)

	if !validColumns[column] {
		return nil, fmt.Errorf("invalid column name: %s", column)
	}
	rows, queryErr := m.DB.Query(stmt, ID)
	if queryErr != nil {
		return nil, queryErr
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "ModSearch", closeErr)
		}
	}()

	var ids []ModIds
	for rows.Next() {
		var m ModIds
		if err := rows.Scan(&m.UserID, &m.ChannelID); err != nil {
			return nil, err
		}
		ids = append(ids, m)
	}
	var returnedIds []int
	for _, id := range ids {
		switch column {
		case "UserID":
			returnedIds = append(returnedIds, id.UserID)
		case "ChannelID":
			returnedIds = append(returnedIds, id.ChannelID)
		}
	}
	return returnedIds, nil
}
