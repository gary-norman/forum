package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content, images, author, authorAvatar string, authorID int, commentable, isFlagged bool) (int, error) {
	stmt := "INSERT INTO Posts (Title, Content, Images, Created, Author, AuthorAvatar, AuthorID, IsCommentable, IsFlagged) VALUES (?, ?, ?, DateTime('now'), ?, ?, ?, ?, ?)"
	result, err := m.DB.Exec(stmt, title, content, images, author, authorAvatar, authorID, commentable, isFlagged)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// fmt.Printf(ErrorMsgs().KeyValuePair, "Inserting a new post with ID: ", id)

	return int(id), nil
}

func (m *PostModel) All() ([]models.Post, error) {
	stmt := "SELECT * FROM Posts ORDER BY ID DESC"
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
			&p.IsCommentable,
			&p.Author,
			&p.AuthorID,
			&p.AuthorAvatar,
			&p.IsFlagged)
		if scanErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "Error", "scan")
			log.Printf(ErrorMsgs().Query, stmt, scanErr)
			return nil, scanErr
		}
		Posts = append(Posts, p)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error", "rows")
		log.Printf(ErrorMsgs().Query, stmt, rowsErr)
		return nil, rowsErr
	}

	return Posts, nil
}

func (m *PostModel) GetPostsByUserID(user int64) ([]models.Post, error) {
	stmt := "SELECT * FROM posts WHERE AuthorID = ? ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt, user)
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
			&p.IsFlagged)
		if scanErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "Error", "scan")
			log.Printf(ErrorMsgs().Query, stmt, scanErr)
			return nil, scanErr
		}
		Posts = append(Posts, p)
	}
	return Posts, nil
}

func (m *PostModel) GetPostsByChannel(channel int64) ([]models.Post, error) {
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
			&p.IsFlagged)
		if scanErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "Error", "scan")
			log.Printf(ErrorMsgs().Query, stmt, scanErr)
			return nil, scanErr
		}
		Posts = append(Posts, p)
	}

	// if rowsErr := rows.Err(); rowsErr != nil {
	// 	log.Printf(ErrorMsgs().KeyValuePair, "Error:", "rows")
	// 	log.Printf(ErrorMsgs().Query, stmt, rowsErr)
	// 	return nil, rowsErr
	// }

	return Posts, nil
}

func (m *PostModel) GetPostByID(id int64) (models.Post, error) {
	stmt := "SELECT * FROM Posts WHERE ID = ?"
	row := m.DB.QueryRow(stmt, id)
	p := models.Post{}
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Images,
		&p.Created,
		&p.IsCommentable,
		&p.Author,
		&p.AuthorID,
		&p.AuthorAvatar,
		&p.IsFlagged)
	if err != nil {
		log.Printf(ErrorMsgs().KeyValuePair, "Error", "scan")
		log.Printf(ErrorMsgs().Query, stmt, err)
		return p, err
	}

	return p, nil
}

// Search queries the database for any post column that contains the values and returns that post
func (m *PostModel) FindCurrentPost(column string, value any) ([]models.Post, error) {
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
		"isFlagged":     true,
	}

	if !validColumns[column] {
		return nil, fmt.Errorf("invalid column name: %s", column)
	}

	// Base query
	query := fmt.Sprintf("SELECT id, title, content, images, created, isCommentable, author, authorID, authorAvatar,  isFlagged FROM posts WHERE %s = ? LIMIT 1", column)

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
		&post.IsFlagged,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No post found")
			return nil, nil // No post found
		}
		return nil, err
	}

	post.AuthorAvatar = avatar.String
	post.Images = images.String
	posts = append(posts, post)

	return posts, nil
}
