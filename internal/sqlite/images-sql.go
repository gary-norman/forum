package sqlite

import (
	"database/sql"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type ImageModel struct {
	DB *sql.DB
}

func (m *ImageModel) Insert(authorID models.UUIDField, postID int64, path string) (int64, error) {
	stmt := "INSERT INTO Images (Created, Updated, AuthorID, PostID, Path) VALUES (DateTime('now'), DateTime('now'), ?, ?, ?)"

	// Convert UUIDField to driver.Value ([]byte) for database storage
	authorIDBytes, err := authorID.Value()
	if err != nil {
		return 0, err
	}

	result, err := m.DB.Exec(stmt, authorIDBytes, postID, path)
	if err != nil {
		return 0, err
	}

	// Return the ID of the newly inserted image
	imageID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return imageID, nil
}

func (m *ImageModel) All() ([]models.Image, error) {
	stmt := "SELECT ID, Created, Updated, AuthorID, PostID, Path FROM Images ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, rows, "All", closeErr)
		}
	}()

	var Images []models.Image
	for rows.Next() {
		p := models.Image{}
		err = rows.Scan(&p.ID, &p.Created, &p.Updated, &p.AuthorID, &p.PostID, &p.Path)
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
