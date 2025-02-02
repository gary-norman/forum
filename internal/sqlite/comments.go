package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type CommentModel struct {
	DB *sql.DB
}

// Upsert inserts or updates a reaction for a specific combination of AuthorID and the parent fields (ChannelID, ReactedPostID, ReactedCommentID). It uses Exists to determine if the reaction already exists.
func (m *CommentModel) Upsert(content, author, channel, authorAvatar string, channelID, authorID, parentPostID, parentCommentID int, isFlagged, isCommentable bool) error {
	if !isValidParent(parentPostID, parentCommentID) {
		return fmt.Errorf("only one of CommentedPostID, or CommentedCommentID must be non-zero")
	}

	// Check if the reaction exists
	exists, err := m.Exists(content, authorID, parentPostID, parentCommentID)
	if err != nil {
		fmt.Println("Upsert Comment > Exists error")
		return fmt.Errorf(err.Error())
	}

	if exists {
		// If the reaction exists, update it
		//fmt.Println("Updating a reaction which already exists (reactions.go :53)")
		return m.Update(content, author, channel, authorAvatar, channelID, authorID, parentPostID, parentCommentID, isFlagged, isCommentable)
	}
	//fmt.Println("Inserting a reaction (reactions.go :56)")

	return m.Insert(content, author, channel, authorAvatar, channelID, authorID, parentPostID, parentCommentID, isFlagged, isCommentable)
}

func (m *CommentModel) Insert(content, author, channel, authorAvatar string, channelID, authorID, parentPostID, parentCommentID int, isFlagged, isCommentable bool) error {
	// Validate that only one of parentPostID or parentCommentID is non-zero
	if !isValidParent(parentPostID, parentCommentID) {
		return fmt.Errorf("only one of CommentedPostID or CommentedCommentID must be non-zero")
	}

	// Begin the transaction
	tx, err := m.DB.Begin()
	//fmt.Println("Beginning INSERT INTO transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Insert in Comments: %w", err)
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
	stmt1 := `INSERT INTO Comments 
		(Content, Created, Author, AuthorID, AuthorAvatar, ChannelName, ChannelID, CommentedPostID, CommentedCommentID, IsCommentable, IsFlagged)
		VALUES (?, DateTime('now'), ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Execute the query, dereferencing the pointers for reactionID values
	_, err = tx.Exec(stmt1, content, author, authorID, authorAvatar, channel, channelID, parentPostID, parentCommentID, isCommentable, isFlagged)
	//fmt.Printf("Inserting row:\nLiked: %v, Disliked: %v, userID: %v, PostID: %v\n", liked, disliked, authorID, parentPostID)
	if err != nil {
		return fmt.Errorf("failed to execute Insert query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	//fmt.Println("Committing INSERT INTO transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Insert in Comments: %w", err)
	}

	return err
}

func (m *CommentModel) Update(content, author, channel, authorAvatar string, channelID, authorID, parentPostID, parentCommentID int, isFlagged, isCommentable bool) error {
	if !isValidParent(parentPostID, parentCommentID) {
		return fmt.Errorf("only one of CommentedPostID, or CommentedCommentID must be non-zero")
	}

	// Begin the transaction
	tx, err := m.DB.Begin()
	//fmt.Println("Beginning UPDATE transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Insert in Comments: %w", err)
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
	stmt1 := `UPDATE Comments 
		SET Content = ?, Created = DateTime('now'), IsCommentable = ?, IsFlagged = ?, Author = ?, AuthorAvatar = ?, ChannelName = ?, ChannelID = ?
		WHERE AuthorID = ? AND (CommentedPostID = ? OR CommentedCommentID = ?)`

	// Execute the query
	_, err = tx.Exec(stmt1, content, isCommentable, isFlagged, author, authorAvatar, channel, channelID, authorID, parentPostID, parentCommentID)
	//fmt.Printf("Updating Comments, where reactionID: %v, PostID: %v and UserID: %v with Liked: %v, Disliked: %v\n", reactionID, reactedPostID, authorID, liked, disliked)
	if err != nil {
		return fmt.Errorf("failed to execute Update query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	//fmt.Println("Committing UPDATE transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Update in Comments: %w", err)
	}

	return err
}

// Delete removes a reaction from the database by ID
func (m *CommentModel) Delete(commentID int) error {
	// Begin the transaction
	tx, err := m.DB.Begin()
	//fmt.Println("Beginning DELETE transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Delete in Comments: %w", err)
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

	stmt1 := `DELETE FROM Comments WHERE ID = ?`
	// Execute the query, dereferencing the pointers for ID values
	_, err = m.DB.Exec(stmt1, commentID)
	//fmt.Printf("Deleting from Reactions where commentID: %v\n", commentID)
	if err != nil {
		return fmt.Errorf("failed to execute Delete query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	//fmt.Println("Committing DELETE transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Delete in Comments: %w", err)
	}

	return err
}

func (m *CommentModel) All() ([]models.Comment, error) {
	stmt := "SELECT ID, Content, Created, Author, AuthorID, AuthorAvatar, ChannelName, ChannelID, CommentedPostID, CommentedCommentID, IsCommentable, IsFlagged FROM Comments ORDER BY ID DESC"

	if m == nil {
		log.Println("Error: DB is nil")
		return nil, fmt.Errorf("database connection is not initialized")
	}

	if m.DB == nil {
		log.Println("Error: CommentModel is nil")
		return nil, fmt.Errorf("database connection is not initialized")
	}

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

	var Comments []models.Comment
	for rows.Next() {
		p := models.Comment{}

		scanErr := rows.Scan(
			&p.ID,
			&p.Content,
			&p.Created,
			&p.Author,
			&p.AuthorID,
			&p.AuthorAvatar,
			&p.ChannelName,
			&p.ChannelID,
			&p.CommentedPostID,
			&p.CommentedCommentID,
			&p.IsFlagged,
			&p.IsCommentable)
		if scanErr != nil {
			log.Printf(ErrorMsgs().KeyValuePair, "Error:", "scan")
			log.Printf(ErrorMsgs().Query, stmt, scanErr)
			return nil, scanErr
		}
		Comments = append(Comments, p)
	}

	return Comments, nil
}

// Exists helps avoid creating duplicate comments by determining whether a comment for the specific combination of AuthorID, PostID/CommentID and the content of the comment
func (m *CommentModel) Exists(content string, authorID, parentPostID, parentCommentID int) (bool, error) {
	//fmt.Printf("Reaction already exists (reactions.go :63 -> Exists) for\nauthorID: %v,\nparentPostID: %v,\nLiked: %v\nDisliked: %v", authorID, parentPostID, liked, disliked)
	stmt := `SELECT EXISTS(
                SELECT 1 FROM Comments
                WHERE AuthorID = ? AND 
                      (CommentedPostID = ? OR CommentedCommentID = ?) AND 
                      Content = ?`

	var exists bool
	err := m.DB.QueryRow(stmt, authorID, parentPostID, parentCommentID, content).Scan(&exists)
	return exists, err
}

// GetReaction checks if a user has already reacted to a post or comment. It retrieves already existing reactions.
func (m *CommentModel) GetComment(authorID int, reactedPostID int, reactedCommentID int) (*models.Reaction, error) {
	var reaction models.Reaction
	var stmt string

	// Build the SQL query depending on whether the reaction is to a post or comment
	if reactedPostID != 0 {
		stmt = `SELECT ID, Created, AuthorID, CommentedPostID, CommentedCommentID, IsCommentable, IsFlagged 
				FROM Comments 
				WHERE AuthorID = ? AND 
				      ReactedPostID = ?`
	} else if reactedCommentID != 0 {
		stmt = `SELECT ID, Liked, Disliked, AuthorID, Created, ReactedPostID, ReactedCommentID 
				FROM Reactions 
				WHERE AuthorID = ? AND 
				      ReactedCommentID = ?`
	} else {
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
