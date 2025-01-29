package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
)

type CommentModel struct {
	DB *sql.DB
}

// Upsert inserts or updates a reaction for a specific combination of AuthorID and the parent fields (ChannelID, ReactedPostID, ReactedCommentID). It uses Exists to determine if the reaction already exists.
func (m *CommentModel) Upsert(content string, authorID, parentPostID, parentCommentID int, isFlagged, isCommentable bool) error {
	if !isValidParent(parentPostID, parentCommentID) {
		return fmt.Errorf("only one of ReactedPostID, or ReactedCommentID must be non-zero")
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
		return m.Update(liked, disliked, reactionID, authorID, parentPostID, parentCommentID, isFlagged)
	}
	//fmt.Println("Inserting a reaction (reactions.go :56)")

	return m.Insert(content, authorID, &parentPostID, &parentCommentID, isFlagged, isCommentable)
}

func (m *CommentModel) Insert(content string, authorID int, parentPostID, parentCommentID *int, isFlagged, isCommentable bool) error {
	// Validate that only one of parentPostID or parentCommentID is non-zero
	if !isValidParent(*parentPostID, *parentCommentID) {
		return fmt.Errorf("only one of CommentedPostID or CommentedCommentID must be non-zero")
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
	stmt1 := `INSERT INTO Comments 
		(Content, Created, AuthorID, CommentedPostID, CommentedCommentID, IsReply, IsFlagged)
		VALUES (?, ?, ?, DateTime('now'), ?, ?)`

	// Execute the query, dereferencing the pointers for reactionID values
	_, err = tx.Exec(stmt1, content, ,
		dereferenceInt(parentPostID), dereferenceInt(parentCommentID))
	//fmt.Printf("Inserting row:\nLiked: %v, Disliked: %v, userID: %v, PostID: %v\n", liked, disliked, authorID, parentPostID)
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

func (m *CommentModel) All() ([]models.Comment, error) {
	stmt := "SELECT ID, Content, Created, ChannelID, AuthorID, CommentedPostID, CommentedCommentID, IsReply, IsFlagged FROM Comments ORDER BY ID DESC"
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

// Exists helps avoid creating duplicate reactions by determining whether a reaction for the specific combination of AuthorID, PostID and the reaction itself - liked/disliked
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
