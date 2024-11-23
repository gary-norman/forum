package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type FlagModel struct {
	DB *sql.DB
}

func (m *FlagModel) Insert(flagType, content string, approved bool, authorID, channelID, flaggedUserID, flaggedPostID, flaggedCommentID int) error {
	stmt := "INSERT INTO Flags (Flag_type, Content, Created, Approved, AuthorID, ChannelID, Flagged_userID, Flagged_postID, Flagged_commentID) VALUES (?, ?, DateTime('now'), ?, ?, ?, ?, ?, ?)"
	_, err := m.DB.Exec(stmt, flagType, content, approved, authorID, channelID, flaggedUserID, flaggedPostID, flaggedCommentID)
	return err
}

func (m *FlagModel) All() ([]models.Flag, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Flag_type, Content, Created, Approved, AuthorID, ChannelID, Flagged_userID, Flagged_postID, Flagged_commentID FROM Flags ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Flags []models.Flag
	for rows.Next() {
		p := models.Flag{}
		err = rows.Scan(&p.ID, &p.FlagType, &p.Content, &p.Created, &p.Approved, &p.AuthorID, &p.ChannelID, &p.FlaggedUserID, &p.FlaggedPostID, &p.FlaggedCommentID)
		if err != nil {
			return nil, err
		}
		Flags = append(Flags, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Flags, nil
}
