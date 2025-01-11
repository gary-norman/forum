package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type LoyaltyModel struct {
	DB *sql.DB
}

func (m *LoyaltyModel) Insert(follower, folowee int) error {
	stmt := "INSERT INTO Loyalty (Follower, Followee) VALUES (?, ?)"
	_, err := m.DB.Exec(stmt, follower, folowee)
	return err
}

func (m *LoyaltyModel) All() ([]models.Loyalty, error) {
	stmt := "SELECT ID, Follower, Followee FROM Loyalty ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Loyalty []models.Loyalty
	for rows.Next() {
		p := models.Loyalty{}
		err = rows.Scan(&p.ID, &p.Follower, &p.Followee)
		if err != nil {
			return nil, err
		}
		Loyalty = append(Loyalty, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Loyalty, nil
}
