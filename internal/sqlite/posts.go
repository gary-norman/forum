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
	stmt := "INSERT INTO Posts (Title, Content, Images, Created, Author, ChannelName, AuthorAvatar, ChannelID, AuthorID, IsCommentable, IsFlagged) VALUES (?, ?, ?, DateTime('now'), ?, ?, ?, ?, ?, ?, ?)"
	_, err := m.DB.Exec(stmt, title, content, images, author, channel, authorAvatar, channelID, authorID, commentable, isFlagged)
	return err
}

func (m *PostModel) All() ([]models.Post, error) {
	stmt := "SELECT ID, Title, Content, Images, Created, Author, AuthorAvatar, IsCommentable, AuthorID, ChannelID, ChannelName, IsFlagged FROM Posts ORDER BY ID DESC"
	rows, selectErr := m.DB.Query(stmt)
	if selectErr != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error:", "select")
		log.Printf(ErrorMsgs().Query, stmt, selectErr)
		return nil, selectErr
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Posts []models.Post
	for rows.Next() {
		p := models.Post{}
		scanErr := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.Images,
			&p.Created,
			&p.Author,
			&p.AuthorAvatar,
			&p.IsCommentable,
			&p.AuthorID,
			&p.ChannelID,
			&p.ChannelName,
			&p.IsFlagged)
		if scanErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "Error:", "scan")
			log.Printf(ErrorMsgs().Query, stmt, scanErr)
			return nil, scanErr
		}
		Posts = append(Posts, p)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error:", "rows")
		log.Printf(ErrorMsgs().Query, stmt, rowsErr)
		return nil, rowsErr
	}

	return Posts, nil
}

func (m *PostModel) GetPostsByChannel(channel int) ([]models.Post, error) {
	stmt := "SELECT * from posts where channelID = ? ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt, channel)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error:", "select")
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
		scanErr := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.Images,
			&p.Created,
			&p.Author,
			&p.AuthorAvatar,
			&p.IsCommentable,
			&p.AuthorID,
			&p.ChannelID,
			&p.ChannelName,
			&p.IsFlagged)
		if scanErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "Error:", "scan")
			log.Printf(ErrorMsgs().Query, stmt, scanErr)
			return nil, scanErr
		}
		Posts = append(Posts, p)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error:", "rows")
		log.Printf(ErrorMsgs().Query, stmt, rowsErr)
		return nil, rowsErr
	}

	return Posts, nil
}
