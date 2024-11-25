package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type ReactionModel struct {
	DB *sql.DB
}

func (m *ReactionModel) Insert(liked, disliked bool, authorID, channelID, reactedPostID, reactedCommentID int) error {
	stmt := "INSERT INTO Reactions (Liked, Disliked, AuthorID, ChannelID, Created, Reacted_postID, Reacted_commentID) VALUES (?, ?, ?, ?, DateTime('now'), ?, ?)"
	_, err := m.DB.Exec(stmt, liked, disliked, authorID, channelID, reactedPostID, reactedCommentID)
	return err
}

func (m *ReactionModel) All() ([]models.Reaction, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Liked, Disliked, AuthorID, ChannelID, Created, Reacted_postID, Reacted_commentID FROM Reactions ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Reactions []models.Reaction
	for rows.Next() {
		p := models.Reaction{}
		err = rows.Scan(&p.ID, &p.Liked, &p.Disliked, &p.AuthorID, &p.ChannelID, &p.Created, &p.ReactedPostID, &p.ReactedCommentID)
		if err != nil {
			return nil, err
		}
		Reactions = append(Reactions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Reactions, nil
}
