package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content, images, author, channel, authorAvatar string, channelID, authorID int, commentable, isFlagged bool) error {
	stmt := "INSERT INTO Posts (Title, Content, Images, Created, Author, ChannelName, AUthorAVatar,ChannelID, AuthorID, Commentable, Is_flagged) VALUES (?, ?, ?, DateTime('now'), ?, ?, ?, ?, ?, ?, ?)"
	_, err := m.DB.Exec(stmt, title, content, images, author, channel, authorAvatar, channelID, authorID, commentable, isFlagged)
	return err
}

func (m *PostModel) All() ([]models.Post, error) {
	stmt := "SELECT ID, Title, Content, Images, Created, Author, Commentable, AuthorID, ChannelID, ChannelName, Is_flagged FROM Posts ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		log.Printf("Error called by %v\n", "1")
		log.Printf(ErrorMsgs().Query, stmt, err)
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Posts []models.Post
	for rows.Next() {
		p := models.Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Images, &p.Created, &p.Author, &p.Commentable, &p.AuthorID, &p.ChannelID, &p.ChannelName, &p.IsFlagged)
		if err != nil {
			log.Printf("Error called by %v\n", "2")
			log.Printf(ErrorMsgs().Query, stmt, err)
			return nil, err
		}
		Posts = append(Posts, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error called by %v\n", "3")
		log.Printf(ErrorMsgs().Query, stmt, err)
		return nil, err
	}

	return Posts, nil
}
