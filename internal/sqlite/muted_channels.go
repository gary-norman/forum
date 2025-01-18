package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type MutedChannelModel struct {
	DB *sql.DB
}

func (m *MutedChannelModel) Insert(authorID, postID int) error {
	stmt := "INSERT INTO Muted_channels (UserID, ChannelID, Created) VALUES (?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, authorID, postID)
	return err
}

func (m *MutedChannelModel) All() ([]models.MutedChannel, error) {
	stmt := "SELECT ID, UserID, ChannelID, Created FROM Muted_channels ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var MutedChannels []models.MutedChannel
	for rows.Next() {
		p := models.MutedChannel{}
		err = rows.Scan(&p.ID, &p.UserID, &p.ChannelID)
		if err != nil {
			return nil, err
		}
		MutedChannels = append(MutedChannels, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return MutedChannels, nil
}
