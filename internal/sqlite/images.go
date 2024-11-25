package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type ImageModel struct {
	DB *sql.DB
}

func (m *ImageModel) Insert(authorID, postID int) error {
	stmt := "INSERT INTO Images (Created, AuthorID, PostID) VALUES (DateTime('now'), ?, ?)"
	_, err := m.DB.Exec(stmt, authorID, postID)
	return err
}

func (m *ImageModel) All() ([]models.Image, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Created, AuthorID, PostID FROM Images ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Images []models.Image
	for rows.Next() {
		p := models.Image{}
		err = rows.Scan(&p.ID, &p.Created, &p.AuthorID, &p.PostID)
		if err != nil {
			return nil, err
		}
		Images = append(Images, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Images, nil
}
