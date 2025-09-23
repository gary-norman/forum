// Package sqlite contains the implementation of all database operations
package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type ReactionModel struct {
	DB *sql.DB
}

type ReactionStatus struct {
	Liked    bool
	Disliked bool
}

func (m *ReactionModel) GetLastReaction(reactedPostID, reactedCommentID int64) (models.Reaction, error) {
	whereArgs, arg := preparePostChannelDynamicWhere(reactedPostID, reactedCommentID)

	stmt := fmt.Sprintf(`
	SELECT
		ID,
		Liked,
		Disliked,
		Created,
		AuthorID,
		ReactedPostID,
		ReactedCommentID
	FROM Reactions
	WHERE %s
	ORDER BY id DESC
	LIMIT 1`, whereArgs)

	row := m.DB.QueryRow(stmt, arg)

	var reaction models.Reaction

	err := row.Scan(
		&reaction.ID,
		&reaction.Liked,
		&reaction.Disliked,
		&reaction.Created,
		&reaction.AuthorID,
		&reaction.ReactedPostID,
		&reaction.ReactedCommentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no matching reaction
			return models.Reaction{}, nil
		}
		return models.Reaction{}, err
	}

	// fmt.Println("Reaction: ", reaction)

	return reaction, nil
}

func (m *ReactionModel) GetReactionStatus(authorID models.UUIDField, reactedPostID, reactedCommentID int64) (ReactionStatus, error) {
	var liked, disliked int
	var reactions ReactionStatus
	if m == nil || m.DB == nil {
		return reactions, fmt.Errorf("reaction model or database is nil")
	}

	whereArgs, arg := preparePostChannelDynamicWhere(reactedPostID, reactedCommentID)

	stmt := fmt.Sprintf(`
	SELECT
	CASE WHEN (SUM(Liked)) = 1 THEN 1 ELSE 0 END,
	CASE WHEN (SUM(Disliked)) = 1 THEN 1 ELSE 0 END
	FROM Reactions
	WHERE AuthorID = ? AND %s
	`, whereArgs)

	if err := m.DB.QueryRow(stmt, authorID, arg).Scan(&liked, &disliked); err != nil {
		return reactions, err
	}

	reactions.Liked = liked == 1
	reactions.Disliked = disliked == 1

	return reactions, nil
}

func (m *ReactionModel) Upsert(liked, disliked bool, authorID models.UUIDField, reactedPostID, reactedCommentID int64) error {
	if !isValidParent(reactedPostID, reactedCommentID) {
		return fmt.Errorf("only one of ReactedPostID or ReactedCommentID must be non-zero")
	}

	var (
		query string
		args  []any
	)

	// TODO refactor so that query inserts ID/NULL to PostID AND CommentID
	if reactedPostID != 0 {
		query = `
		WITH existing AS (
    SELECT ID,
    COALESCE(Liked, 0) AS existing_liked,
    COALESCE(Disliked, 0) AS existing_disliked
    FROM Reactions
    WHERE AuthorID = ? AND ReactedPostID = ?
		)
		INSERT OR REPLACE INTO Reactions (ID, Liked, Disliked, Created, AuthorID, ReactedPostID)
		VALUES (
			(SELECT ID FROM existing),
			CASE WHEN (SELECT existing_liked FROM existing) + 1 = 2 THEN 0 ELSE ? END,
			CASE WHEN (SELECT existing_disliked FROM existing) + 1 = 2 THEN 0 ELSE ? END,
			CURRENT_TIMESTAMP,
			?,
			?
		);
		`
		args = []any{authorID, reactedPostID, liked, disliked, authorID, reactedPostID}
	} else {
		query = `
		WITH existing AS (
    SELECT ID,
    COALESCE(Liked, 0) AS existing_liked,
    COALESCE(Disliked, 0) AS existing_disliked
    FROM Reactions
    WHERE AuthorID = ? AND ReactedCommentID = ?
		)
		INSERT OR REPLACE INTO Reactions (ID, Liked, Disliked, Created, AuthorID, ReactedCommentID)
		VALUES (
			(SELECT ID FROM existing),
			CASE WHEN (SELECT existing_liked FROM existing) + 1 = 2 THEN 0 ELSE ? END,
			CASE WHEN (SELECT existing_disliked FROM existing) + 1 = 2 THEN 0 ELSE ? END,
			CURRENT_TIMESTAMP,
			?,
			?
		);
		`
		args = []any{authorID, reactedCommentID, liked, disliked, authorID, reactedCommentID}
	}

	_, err := m.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to upsert reaction: %w", err)
	}

	return nil
}

func (m *ReactionModel) CountReactions(reactedPostID, reactedCommentID int64) (likes, dislikes int, err error) {
	if !isValidParent(reactedPostID, reactedCommentID) {
		return 0, 0, fmt.Errorf("only one of  ReactedPostID, or ReactedCommentID must be non-zero")
	}

	whereArgs, arg := preparePostChannelDynamicWhere(reactedPostID, reactedCommentID)

	stmt := fmt.Sprintf(`
		SELECT
		SUM(Liked) AS Likes,
		SUM(Disliked) AS Dislikes
		FROM Reactions
		WHERE %s`, whereArgs)
	var likesSum, dislikesSum sql.NullInt64

	// Run the query
	err = m.DB.QueryRow(stmt, arg).Scan(&likesSum, &dislikesSum)
	if err != nil {
		return 0, 0, err
	}
	likes = int(likesSum.Int64)
	dislikes = int(dislikesSum.Int64)

	return likes, dislikes, err
}

// Delete removes a reaction from the database by ID
func (m *ReactionModel) Delete(reactionID int64) error {
	stmt := `DELETE FROM Reactions WHERE ID = ?`
	// Execute the query, dereferencing the pointers for ID values
	_, err := m.DB.Exec(stmt, reactionID)
	// fmt.Printf("Deleting from Reactions where reactionID: %v\n", reactionID)
	if err != nil {
		return fmt.Errorf("failed to execute Delete query: %w", err)
	}

	return err
}

func (m *ReactionModel) All() ([]models.Reaction, error) {
	stmt := "SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID FROM Reactions ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, rows, "All", closeErr)
		}
	}()

	var Reactions []models.Reaction
	for rows.Next() {
		p := models.Reaction{}
		err = rows.Scan(&p.ID, &p.Liked, &p.Disliked, &p.AuthorID, &p.Created, &p.ReactedPostID, &p.ReactedCommentID)
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

// ***** helper functions *****

// Ensure only one parent ID is present when inserting a reaction
func isValidParent(reactedPostID, reactedCommentID int64) bool {
	// Ensure only one parent ID is non-zero
	nonZeroCount := 0
	if reactedPostID != 0 {
		nonZeroCount++
	}
	if reactedCommentID != 0 {
		nonZeroCount++
	}
	return nonZeroCount == 1
}

// preparePostChannelDynamicWhere prepares the tail of the UPDATE statement
func preparePostChannelDynamicWhere(post, comment int64) (string, int64) {
	if post == 0 {
		return "ReactedPostID IS NULL AND ReactedCommentID = ?", comment
	}
	return "ReactedPostID = ? AND ReactedCommentID IS NULL", post
}

// Helper function to safely dereference an integer pointer
// func dereferenceInt(value *int) any {
// 	if value == nil {
// 		return nil
// 	}
// 	return *value
// }
