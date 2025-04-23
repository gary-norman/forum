package sqlite

import (
	"database/sql"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type SavedModel struct {
	DB *sql.DB
}

func (m *SavedModel) Insert(postID, commentID, channelID int64) error {
	stmt := "INSERT INTO Bookmarks (PostID, CommentID, ChannelID, Created) VALUES (?, ?, ?, DateTime('now'))"
	_, err := m.DB.Exec(stmt, postID, commentID, channelID)
	return err
}

func (m *SavedModel) All() ([]models.Bookmark, error) {
	stmt := "SELECT ID, PostID, CommentID, ChannelID, Created FROM Bookmarks ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Bookmarks []models.Bookmark
	for rows.Next() {
		p := models.Bookmark{}
		err = rows.Scan(&p.ID, &p.PostID, &p.CommentID, &p.ChannelID, &p.Created)
		if err != nil {
			return nil, err
		}
		Bookmarks = append(Bookmarks, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Bookmarks, nil
}
