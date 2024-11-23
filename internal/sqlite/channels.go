package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type ChannelModel struct {
	DB *sql.DB
}

func (m *ChannelModel) Insert(name, avatar, description, rules string, privacy bool) error {
	stmt := "INSERT INTO Channels (Name, Avatar, Description, Created, Rules, Privacy) VALUES (?, ?, ?, DateTime('now'), ?, ?)"
	_, err := m.DB.Exec(stmt, name, avatar, description, rules, privacy)
	return err
}

func (m *ChannelModel) All() ([]models.Channel, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Name, Avatar, Description, Created, Rules, Privacy FROM Channels ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Channels []models.Channel
	for rows.Next() {
		p := models.Channel{}
		err = rows.Scan(&p.ID, &p.Name, &p.Avatar, &p.Description, &p.Created, &p.Rules, &p.Privacy)
		if err != nil {
			return nil, err
		}
		Channels = append(Channels, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Channels, nil
}
