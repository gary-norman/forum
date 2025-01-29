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
func (m *ReactionModel) WhichReaction(authorID, reactedPostID, reactedCommentID int) (bool, bool, error) {
	//fmt.Printf("WhichReaction:\nChecking reaction (reactions.go :22 -> WhichReaction) for\nauthorID: %v,\nchannelID: %v,\nreactedPostID: %v,\nreactedCommentID: %v\n", authorID, channelID, reactedPostID, reactedCommentID)

	stmt := `SELECT Liked, Disliked FROM Reactions
             WHERE AuthorID = ?  AND ReactedPostID = ? AND ReactedCommentID = ?`

	var liked, disliked bool
	err := m.DB.QueryRow(stmt, authorID, reactedPostID, reactedCommentID).Scan(&liked, &disliked)
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

func (m *ReactionModel) GetReactionStatus(authorID, reactedPostID, reactedCommentID int) (ReactionStatus, error) {
	if m == nil || m.DB == nil {
		return ReactionStatus{}, fmt.Errorf("reaction model or database is nil")
	}
	liked, disliked, err := m.WhichReaction(authorID, reactedPostID, reactedCommentID)
	if err != nil {
		return ReactionStatus{}, err
	}
	return ReactionStatus{Liked: liked, Disliked: disliked}, nil
}

func (m *ReactionModel) Insert(liked, disliked bool, authorID, reactedPostID, reactedCommentID int) error {
	// Validate that only one of reactedPostID or reactedCommentID is non-zero
	if !isValidParent(reactedPostID, reactedCommentID) {
		return fmt.Errorf("only one of ReactedPostID or ReactedCommentID must be non-zero")
	}

	// Begin the transaction
	tx, err := m.DB.Begin()
	//fmt.Println("Beginning INSERT INTO transaction")
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
	stmt1 := `INSERT INTO Reactions 
		(Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID)
		VALUES (?, ?, ?, DateTime('now'), ?, ?)`

	// Execute the query, dereferencing the pointers for reactionID values
	_, err = tx.Exec(stmt1, liked, disliked, authorID,
		reactedPostID, reactedCommentID)
	//fmt.Printf("Inserting row:\nLiked: %v, Disliked: %v, userID: %v, PostID: %v\n", liked, disliked, authorID, reactedPostID)
	if err != nil {
		return fmt.Errorf("failed to execute Insert query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	//fmt.Println("Committing INSERT INTO transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Insert in Reactions: %w", err)
	}

	return err
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

	// Begin the transaction
	tx, err := m.DB.Begin()
	//fmt.Println("Beginning UPDATE transaction")
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

	stmt1 := `UPDATE Reactions 
             SET Liked = ?, Disliked = ?, Created = DateTime('now') 
             WHERE AuthorID = ? AND ReactedPostID = ? AND ReactedCommentID = ?`
	//fmt.Printf("Updating Reactions, where reactionID: %v, PostID: %v and UserID: %v with Liked: %v, Disliked: %v\n", reactionID, reactedPostID, authorID, liked, disliked)

	// Execute the query
	_, err = tx.Exec(stmt1, liked, disliked, authorID, reactedPostID, reactedCommentID)
	if err != nil {
		return fmt.Errorf("failed to execute Update query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	//fmt.Println("Committing UPDATE transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Update in Reactions: %w", err)
	}

	return err
}

// Upsert inserts or updates a reaction for a specific combination of AuthorID and the parent fields (ChannelID, ReactedPostID, ReactedCommentID). It uses Exists to determine if the reaction already exists.
func (m *ReactionModel) Upsert(liked, disliked bool, reactionID, authorID, reactedPostID, reactedCommentID int) error {
	if !isValidParent(reactedPostID, reactedCommentID) {
		return fmt.Errorf("only one of ReactedPostID, or ReactedCommentID must be non-zero")
	}

	// Check if the reaction exists
	exists, err := m.Exists(liked, disliked, authorID, reactedPostID)
	if err != nil {
		fmt.Println("Upsert > Exists error")
		return fmt.Errorf(err.Error())
	}

	if exists {
		// If the reaction exists, update it
		//fmt.Println("Updating a reaction which already exists (reactions.go :53)")
		return m.Update(liked, disliked, authorID, reactedPostID, reactedCommentID)
	}
	//fmt.Println("Inserting a reaction (reactions.go :56)")

	return m.Insert(liked, disliked, authorID, reactedPostID, reactedCommentID)
}

// Exists helps avoid creating duplicate reactions by determining whether a reaction for the specific combination of AuthorID, PostID and the reaction itself - liked/disliked
func (m *ReactionModel) Exists(liked, disliked bool, authorID, reactedPostID int) (bool, error) {
	//fmt.Printf("Reaction already exists (reactions.go :63 -> Exists) for\nauthorID: %v,\nreactedPostID: %v,\nLiked: %v\nDisliked: %v", authorID, reactedPostID, liked, disliked)
	stmt := `SELECT EXISTS(
                SELECT 1 FROM Reactions
                WHERE AuthorID = ? AND ReactedPostID = ?)`
	var exists bool
	err := m.DB.QueryRow(stmt, authorID, reactedPostID).Scan(&exists)
	return exists, err
}

// CheckExistingReaction checks if the user has already reacted to a post or comment. For purposes of reactions, the user can only react once to a post or comment.
func (m *ReactionModel) CheckExistingReaction(liked, disliked bool, reactionAuthorID, reactedPostID int) (*models.Reaction, error) {
	var reaction models.Reaction
	// Query to find if there's a reaction by the same user for the same post or comment
	stmt := `SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID
			FROM Reactions 
			WHERE AuthorID = ? 
			AND ReactedPostID = ?
			AND (Liked = ? OR Disliked = ?)`
	err := m.DB.QueryRow(stmt, reactionAuthorID, reactedPostID, liked, disliked).Scan(
		&reaction.ID, &reaction.Liked, &reaction.Disliked, &reaction.AuthorID, &reaction.Created, &reaction.ReactedPostID)
	if errors.Is(err, sql.ErrNoRows) {
		// No existing reaction found
		//fmt.Println("No existing reaction found")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &reaction, nil
}

func (m *ReactionModel) CountReactions(reactedPostID, reactedCommentID int) (likes, dislikes int, err error) {
	if !isValidParent(reactedPostID, reactedCommentID) {
		return 0, 0, fmt.Errorf("only one of  ReactedPostID, or ReactedCommentID must be non-zero")
	}

	stmt := `SELECT 
                 SUM(CASE WHEN Liked = 1 THEN 1 ELSE 0 END) AS Likes,
                 SUM(CASE WHEN Disliked = 1 THEN 1 ELSE 0 END) AS Dislikes
             FROM Reactions
             WHERE ReactedPostID = ? AND 
                   ReactedCommentID = ?`
	var likesSum, dislikesSum sql.NullInt64

	// Run the query
	err = m.DB.QueryRow(stmt, reactedPostID, reactedCommentID).Scan(&likesSum, &dislikesSum)
	if err != nil {
		return 0, 0, err
	}

	likes = int(likesSum.Int64)
	dislikes = int(dislikesSum.Int64)
	//fmt.Println("likes:", likes)
	//fmt.Println("dislikes:", dislikes)

	return likes, dislikes, err
}

// GetReaction checks if a user has already reacted to a post or comment. It retrieves already existing reactions.
func (m *ReactionModel) GetReaction(authorID, reactedPostID, reactedCommentID int) (*models.Reaction, error) {
	var reaction models.Reaction
	var stmt string

	// Build the SQL query depending on whether the reaction is to a post or comment
	if reactedPostID != 0 {
		stmt = `SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID 
				FROM Reactions 
				WHERE AuthorID = ? AND 
				      ReactedPostID = ?`

	} else if reactedCommentID != 0 {
		stmt = `SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID 
				FROM Reactions 
				WHERE AuthorID = ? AND 
				      ReactedCommentID = ?`
	} else {
		log.Printf("Couldn't find the reaction with AuthorID: %v, reactedPostID: %v, reactedCommentID: %v\n", authorID, &reactedPostID, &reactedCommentID)
		return nil, nil
	}

	// Query the database
	row := m.DB.QueryRow(stmt, authorID, reactedPostID)
	if reactedCommentID != 0 {
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
	// Begin the transaction
	tx, err := m.DB.Begin()
	//fmt.Println("Beginning DELETE transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Delete in Reactions: %w", err)
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

	stmt1 := `DELETE FROM Reactions WHERE ID = ?`
	// Execute the query, dereferencing the pointers for ID values
	_, err = m.DB.Exec(stmt1, reactionID)
	//fmt.Printf("Deleting from Reactions where reactionID: %v\n", reactionID)
	if err != nil {
		return fmt.Errorf("failed to execute Delete query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	//fmt.Println("Committing DELETE transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Delete in Reactions: %w", err)
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
