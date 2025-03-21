package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
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
	stmt := "SELECT * FROM posts WHERE ChannelID = ? ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt, channel)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error", "select")
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
			&p.IsCommentable,
			&p.Author,
			&p.AuthorID,
			&p.AuthorAvatar,
			&p.ChannelName,
			&p.ChannelID,
			&p.IsFlagged)
		if scanErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "Error", "scan")
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

// Search queries the database for any post column that contains the values and returns that post
func (m *PostModel) FindCurrentPost(column string, value interface{}) ([]models.Post, error) {
	// Validate column name to prevent SQL injection
	validColumns := map[string]bool{
		"id":            true,
		"title":         true,
		"content":       true,
		"images":        true,
		"created":       true,
		"isCommentable": true,
		"author":        true,
		"authorID":      true,
		"authorAvatar":  true,
		"channelName":   true,
		"channelID":     true,
		"isFlagged":     true,
	}

	if !validColumns[column] {
		return nil, fmt.Errorf("invalid column name: %s", column)
	}

	// Base query
	query := fmt.Sprintf("SELECT id, title, content, images, created, isCommentable, author, authorID, authorAvatar, channelName, channelID, isFlagged FROM posts WHERE %s = ? LIMIT 1", column)

	row := m.DB.QueryRow(query, value)

	// Parse result into a single post
	var posts []models.Post
	var post models.Post
	var avatar, images sql.NullString

	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&images,
		&post.Created,
		&post.IsCommentable,
		&post.Author,
		&post.AuthorID,
		&avatar,
		&post.ChannelName,
		&post.ChannelID,
		&post.IsFlagged,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No post found
		}
		return nil, err
	}

	post.AuthorAvatar = avatar.String
	post.Images = images.String

	posts = append(posts, post)

	return posts, nil
}
