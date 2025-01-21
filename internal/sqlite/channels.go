package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type ChannelModel struct {
	DB *sql.DB
}

func (m *ChannelModel) Insert(ownerID int, name, avatar, banner, description string, privacy, isFlagged, isMuted bool) error {
	stmt := "INSERT INTO Channels (OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsFlagged, IsMuted) VALUES (?, ?, ?, ?, ?, DateTime('now'),?,  ?, ?)"
	_, err := m.DB.Exec(stmt, ownerID, name, avatar, banner, description, privacy, isFlagged, isMuted)
	return err
}

func (m *ChannelModel) All() ([]models.Channel, error) {
	stmt := "SELECT ID, OwnerID, Name, Avatar, Banner, Description, Created, Privacy, IsFlagged, IsMuted FROM Channels ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Channels []models.Channel
	for rows.Next() {
		p := models.Channel{}
		err = rows.Scan(&p.ID, &p.OwnerID, &p.Name, &p.Avatar, &p.Banner, &p.Description, &p.Created, &p.Privacy, &p.IsFLagged, &p.IsMuted)
		if err != nil {
			return nil, err
		}
		Channels = append(Channels, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "Channels", len(Channels))
	return Channels, nil
}
