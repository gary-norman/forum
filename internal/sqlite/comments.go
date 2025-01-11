package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Insert(content string, channel, author, commentedPostID, commentedCommentID int, isReply, commentable bool) error {
	stmt := "INSERT INTO Comments (Content, Created, ChannelID, AuthorID, Commented_postID, Commented_commentID, Is_reply, Is_flagged) VALUES (?, DateTime('now'), ?, ?, ?, ?, 0, 0)"
	_, err := m.DB.Exec(stmt, content, channel, author, commentedPostID, commentedCommentID, isReply, commentable)
	return err
}

func (m *CommentModel) All() ([]models.Comment, error) {
	stmt := "SELECT ID, Content, Created, ChannelID, AuthorID, Commented_postID, Commented_commentID, Is_reply, Is_flagged FROM Comments ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Comments []models.Comment
	for rows.Next() {
		p := models.Comment{}
		err = rows.Scan(&p.ID, &p.Content, &p.Created, &p.ChannelID, &p.AuthorID, &p.CommentedPostID, &p.CommentedCommentID, &p.IsReply, &p.IsFlagged)
		if err != nil {
			return nil, err
		}
		Comments = append(Comments, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Comments, nil
}
