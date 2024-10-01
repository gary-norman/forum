package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content, images string, channel, author int, commentable bool) error {
	stmt := "INSERT INTO posts (Title, Content, Images, Created, ChannelID, AuthorID, Commentable, Is_flagged) VALUES (?, ?, ?, DateTime('now'), ?, ?, ?, 0)"
	_, err := m.DB.Exec(stmt, title, content, images, channel, author, commentable)
	return err
}

func (m *PostModel) All() ([]models.Post, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Title, Content, Images, Created, Commentable, AuthorID, ChannelID, Is_flagged FROM posts ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Posts []models.Post
	for rows.Next() {
		p := models.Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Images, &p.Created, &p.Commentable, &p.AuthorID, &p.ChannelID, &p.IsFlagged)
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
