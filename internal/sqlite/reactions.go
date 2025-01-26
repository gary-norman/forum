package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type ReactionModel struct {
	DB *sql.DB
}

type ReactionStatus struct {
	Liked    bool
	Disliked bool
}

// WhichReaction checks if a reaction exists and returns whether it is Liked or Disliked.
func (m *ReactionModel) WhichReaction(authorID, channelID, reactedPostID, reactedCommentID int) (bool, bool, error) {
	//fmt.Printf("WhichReaction:\nChecking reaction (reactions.go :22 -> WhichReaction) for\nauthorID: %v,\nchannelID: %v,\nreactedPostID: %v,\nreactedCommentID: %v\n", authorID, channelID, reactedPostID, reactedCommentID)

	stmt := `SELECT Liked, Disliked FROM Reactions
             WHERE AuthorID = ?  AND ReactedPostID = ? AND ReactedCommentID = ?`

	var liked, disliked bool
	err := m.DB.QueryRow(stmt, authorID, channelID, reactedPostID, reactedCommentID).Scan(&liked, &disliked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No matching reaction found, return default values
			return false, false, nil
		}
		// Return other errors as is
		return false, false, err
	}

	return liked, disliked, nil
}

// TODO Might have issues cuz it's not in the database

func (m *ReactionModel) GetReactionStatus(authorID, channelID, reactedPostID, reactedCommentID int) (ReactionStatus, error) {
	if m == nil || m.DB == nil {
		return ReactionStatus{}, fmt.Errorf("reaction model or database is nil")
	}
	liked, disliked, err := m.WhichReaction(authorID, channelID, reactedPostID, reactedCommentID)
	if err != nil {
		return ReactionStatus{}, err
	}
	return ReactionStatus{Liked: liked, Disliked: disliked}, nil
}

//func (m *ReactionModel) Insert(liked, disliked bool, authorID int, reactedPostID, reactedCommentID *int) error {
//	if !isValidParent(*reactedPostID, *reactedCommentID) {
//		return fmt.Errorf("only one of ReactedPostID, or ReactedCommentID must be non-zero")
//	}
//
//	tx, err := m.DB.Begin() // Begin transaction
//	if err != nil {
//		return log.Printf(ErrorMsgs().Open, "Reactions", "Insert")
//	}
//
//	stmt := `INSERT INTO Reactions
//			(Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID)
//			VALUES (?, ?, ?, DateTime('now'), ?, ?)`
//
//	// Ensure rollback is called if there's a failure
//	defer func() {
//		if p := recover(); p != nil {
//			_ = tx.Rollback()
//			panic(p)
//		} else if err != nil {
//			_ = tx.Rollback()
//		}
//	}()
//
//	_, err = tx.Exec(stmt, liked, disliked, authorID, reactedPostID, reactedCommentID)
//	//_, err := m.DB.Exec(stmt, liked, disliked, authorID, reactedPostID, reactedCommentID)
//	err = tx.Commit() // Commit the transaction
//	return err
//}

func (m *ReactionModel) Insert(liked, disliked bool, authorID int, reactedPostID, reactedCommentID *int) error {
	// Validate that only one of reactedPostID or reactedCommentID is non-zero
	if !isValidParent(*reactedPostID, *reactedCommentID) {
		return fmt.Errorf("only one of ReactedPostID or ReactedCommentID must be non-zero")
	}

	// Begin the transaction
	tx, err := m.DB.Begin()
	fmt.Println("Beginning transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Insert in Reactions: %w", err)
	}

	// Ensure rollback on failure
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Rolling back transaction")

			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Define the SQL statement
	stmt := `INSERT INTO Reactions 
		(Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID)
		VALUES (?, ?, ?, DateTime('now'), ?, ?)`

	// Execute the query, dereferencing the pointers for ID values
	_, err = tx.Exec(stmt, liked, disliked, authorID,
		dereferenceInt(reactedPostID), dereferenceInt(reactedCommentID))
	if err != nil {
		return fmt.Errorf("failed to execute Insert query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	fmt.Println("Committing transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Insert in Reactions: %w", err)
	}

	return nil
}

// Helper function to safely dereference an integer pointer
func dereferenceInt(value *int) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func (m *ReactionModel) Update(liked, disliked bool, authorID, reactedPostID, reactedCommentID int) error {
	if !isValidParent(reactedPostID, reactedCommentID) {
		return fmt.Errorf("only one of ReactedPostID, or ReactedCommentID must be non-zero")
	}

	stmt := `UPDATE Reactions 
             SET Liked = ?, Disliked = ?, Created = DateTime('now') 
             WHERE AuthorID = ? AND ReactedPostID = ? AND ReactedCommentID = ?`
	_, err := m.DB.Exec(stmt, liked, disliked, authorID, reactedPostID, reactedCommentID)
	return err
}

// Upsert inserts or updates a reaction for a specific combination of AuthorID and the parent fields (ChannelID, ReactedPostID, ReactedCommentID). It uses Exists to determine if the reaction already exists.
func (m *ReactionModel) Upsert(liked, disliked bool, authorID, reactedPostID, reactedCommentID int) error {
	if !isValidParent(reactedPostID, reactedCommentID) {
		return fmt.Errorf("only one of ReactedPostID, or ReactedCommentID must be non-zero")
	}

	// Check if the reaction exists
	exists, err := m.Exists(authorID, reactedPostID, reactedCommentID)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if exists {
		// If the reaction exists, update it
		fmt.Println("Updating a reaction which already exists (reactions.go :53 -> Upsert)")
		return m.Update(liked, disliked, authorID, reactedPostID, reactedCommentID)
	}
	fmt.Println("Inserting a reaction (reactions.go :56 -> Upsert)")

	return m.Insert(liked, disliked, authorID, &reactedPostID, &reactedCommentID)
}

// Exists helps avoid creating duplicate reactions by determining whether a reaction for the specific combination of AuthorID and the parent fields (ChannelID, ReactedPostID, ReactedCommentID).
func (m *ReactionModel) Exists(authorID, reactedPostID, reactedCommentID int) (bool, error) {
	fmt.Printf("Reaction already exists (reactions.go :63 -> Exists) for\nauthorID: %v,\nreactedPostID: %v,\nreactedCommentID: %v\n", authorID, reactedPostID, reactedCommentID)

	stmt := `SELECT EXISTS(
                SELECT 1 FROM Reactions
                WHERE AuthorID = ?  AND ReactedPostID = ? AND ReactedCommentID = ?
             )`
	var exists bool
	err := m.DB.QueryRow(stmt, authorID, reactedPostID, reactedCommentID).Scan(&exists)
	return exists, err
}

// CheckExistingReaction checks if the user has already reacted to a post or comment. For purposes of reactions, the user can only react once to a post or comment.
func (m *ReactionModel) CheckExistingReaction(reactionAuthorID, reactedPostID, reactedCommentID int) (*models.Reaction, error) {
	var reaction models.Reaction

	// Query to find if there's a reaction by the same user for the same post or comment
	stmt := `SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID
			FROM Reactions 
			WHERE AuthorID = ?
			AND (ReactedPostID = ? OR ReactedCommentID = ?)`

	err := m.DB.QueryRow(stmt, reactionAuthorID, reactedPostID, reactedCommentID).Scan(
		&reaction.ID, &reaction.Liked, &reaction.Disliked, &reaction.AuthorID, &reaction.Created, &reaction.ReactedPostID, &reaction.ReactedCommentID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		// No existing reaction found
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &reaction, nil
}

func (m *ReactionModel) CountReactions(channelID, reactedPostID, reactedCommentID int) (likes, dislikes int, err error) {
	if !isValidParent(reactedPostID, reactedCommentID) {
		return 0, 0, fmt.Errorf("only one of  ReactedPostID, or ReactedCommentID must be non-zero")
	}

	stmt := `SELECT 
                 SUM(CASE WHEN Liked = 1 THEN 1 ELSE 0 END) AS Likes,
                 SUM(CASE WHEN Disliked = 1 THEN 1 ELSE 0 END) AS Dislikes
             FROM Reactions
             WHERE ReactedPostID = ? AND 
                   ReactedCommentID = ?`
	var likesNull, dislikesNull sql.NullInt64

	// Run the query
	err = m.DB.QueryRow(stmt, channelID, reactedPostID, reactedCommentID).Scan(&likesNull, &dislikesNull)
	if err != nil {
		return 0, 0, err
	}

	// Convert NULL to 0 for likes and dislikes
	if likesNull.Valid {
		likes = int(likesNull.Int64)
	} else {
		likes = 0
	}

	if dislikesNull.Valid {
		dislikes = int(dislikesNull.Int64)
	} else {
		dislikes = 0
	}

	return
}

// GetReaction checks if a user has already reacted to a post or comment. It retrieves already existing reactions.
func (m *ReactionModel) GetReaction(authorID int, reactedPostID *int, reactedCommentID *int) (*models.Reaction, error) {
	var reaction models.Reaction
	var stmt string

	// Build the SQL query depending on whether the reaction is to a post or comment
	if reactedPostID != nil {
		stmt = `SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID 
				FROM Reactions 
				WHERE AuthorID = ? AND 
				      ReactedPostID = ?`
	} else if reactedCommentID != nil {
		stmt = `SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID 
				FROM Reactions 
				WHERE AuthorID = ? AND 
				      ReactedCommentID = ?`
	} else {
		return nil, nil
	}

	// Query the database
	row := m.DB.QueryRow(stmt, authorID, reactedPostID)
	if reactedCommentID != nil {
		row = m.DB.QueryRow(stmt, authorID, reactedCommentID)
	}

	err := row.Scan(&reaction.ID, &reaction.Liked, &reaction.Disliked, &reaction.AuthorID, &reaction.Created, &reaction.ReactedPostID, &reaction.ReactedCommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No reaction found
			return nil, nil
		}
		// Other errors
		log.Printf("Error fetching reaction: %v", err)
		return nil, err
	}

	// Return the existing reaction
	return &reaction, nil
}

// Delete removes a reaction from the database by ID
func (m *ReactionModel) Delete(reactionID int) error {
	stmt := `DELETE FROM Reactions WHERE ID = ?`
	_, err := m.DB.Exec(stmt, reactionID)
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
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
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

// Ensure only one parent ID is present when inserting a reaction
func isValidParent(reactedPostID, reactedCommentID int) bool {
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
