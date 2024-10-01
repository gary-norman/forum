package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content string) error {
	stmt := "INSERT INTO posts (Title, Content, Created) VALUES (?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, title, content)
	return err
}

func (m *PostModel) All() ([]models.Posts, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Title, Content, Images, Created, Commentable, Author, Reactions, Comments, Channel FROM posts ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Posts []models.Posts
	for rows.Next() {
		p := models.Posts{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Images, &p.Created, &p.Commentable, &p.AuthorID, &p.Reactions, &p.Comments, &p.ChannelID)
		if err != nil {
			return nil, err
		}
		Posts = append(Posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Posts, nil
}
