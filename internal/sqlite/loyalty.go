package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gary-norman/forum/internal/models"
)

type LoyaltyModel struct {
	DB *sql.DB
}

func (m *LoyaltyModel) InsertLoyalty(follower, following models.UUIDField) error {
	err := m.InsertFollowing(follower, following)
	if err != nil {
		fmt.Println("Error adding a following")
		return errors.New(err.Error())
	}

	err = m.InsertFollower(following, follower)
	if err != nil {
		fmt.Println("Error adding a follower")
		return errors.New(err.Error())
	}

	return err
}

// InsertFollower inserts a
func (m *LoyaltyModel) InsertFollower(user, follower models.UUIDField) error {
	// Begin the transaction
	tx, err := m.DB.Begin()
	// fmt.Println("Beginning UPDATE transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Insert Follower: %w", err)
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

	stmt := "INSERT INTO Followers (UserID, FollowerUserID) VALUES (?, ?)"
	_, InsertErr := m.DB.Exec(stmt, user, follower)
	// fmt.Printf("Updating Comments, where reactionID: %v, PostID: %v and UserID: %v with Liked: %v, Disliked: %v\n", reactionID, reactedPostID, authorID, liked, disliked)
	if InsertErr != nil {
		return fmt.Errorf("failed to execute Insert query in Insert Follower: %w", err)
	}

	// Commit the transaction
	commitErr := tx.Commit()
	// fmt.Println("Committing UPDATE transaction")
	if commitErr != nil {
		return fmt.Errorf("failed to commit transaction for Insert query in Insert Follower: %w", err)
	}

	return commitErr
}

func (m *LoyaltyModel) CountUsers(userID models.UUIDField) (followers, following int, err error) {
	stmt1 := `SELECT COUNT(*) AS FollowingCount
             FROM Following
             WHERE UserID = ?`

	stmt2 := `SELECT COUNT(*) AS FollowersCount
             FROM Followers
             WHERE UserID = ?`

	var followingCount, followersCount sql.NullInt64

	// Run the query
	err = m.DB.QueryRow(stmt1, userID).Scan(&followingCount)
	if err != nil {
		return 0, 0, err
	}

	// Run the query
	err = m.DB.QueryRow(stmt2, userID).Scan(&followersCount)
	if err != nil {
		return 0, 0, err
	}

	followers = int(followersCount.Int64)
	following = int(followingCount.Int64)
	// fmt.Println("likes:", likes)
	// fmt.Println("dislikes:", dislikes)

	return followers, following, err
}

// Delete removes an entry in the Following table by ID
func (m *LoyaltyModel) Delete(followingID, followersID models.UUIDField) error {
	// Begin the transaction
	tx, err := m.DB.Begin()
	// fmt.Println("Beginning DELETE transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Delete in Following: %w", err)
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

	stmt1 := `DELETE FROM Following WHERE ID = ?`
	// Execute the query, dereferencing the pointers for ID values
	_, err = m.DB.Exec(stmt1, followingID)
	// fmt.Printf("Deleting from Reactions where commentID: %v\n", commentID)
	if err != nil {
		return fmt.Errorf("failed to execute Delete query: %w", err)
	}

	stmt2 := `DELETE FROM Followers WHERE ID = ?`
	// Execute the query, dereferencing the pointers for ID values
	_, err = m.DB.Exec(stmt2, followersID)
	// fmt.Printf("Deleting from Reactions where commentID: %v\n", commentID)
	if err != nil {
		return fmt.Errorf("failed to execute Delete query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	// fmt.Println("Committing DELETE transaction")
	if err != nil {
		return fmt.Errorf("failed to commit transaction for Delete in Following: %w", err)
	}

	return err
}

// InsertFollowing inserts a new user to the Following list of a target use
func (m *LoyaltyModel) InsertFollowing(user, following models.UUIDField) error {
	// Begin the transaction
	tx, err := m.DB.Begin()
	// fmt.Println("Beginning UPDATE transaction")
	if err != nil {
		return fmt.Errorf("failed to begin transaction for Insert Following: %w", err)
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

	stmt := "INSERT INTO Following (UserID, FollowingUserID) VALUES (?, ?)"
	_, InsertErr := m.DB.Exec(stmt, user, following)
	// fmt.Printf("Updating Comments, where reactionID: %v, PostID: %v and UserID: %v with Liked: %v, Disliked: %v\n", reactionID, reactedPostID, authorID, liked, disliked)
	if InsertErr != nil {
		return fmt.Errorf("failed to execute Insert query in Insert Following: %w", err)
	}

	// Commit the transaction
	commitErr := tx.Commit()
	// fmt.Println("Committing UPDATE transaction")
	if commitErr != nil {
		return fmt.Errorf("failed to commit transaction for Insert query in Insert Follower: %w", err)
	}

	return commitErr
}

func (m *LoyaltyModel) All() ([]models.Loyalty, error) {
	stmt := "SELECT ID, Follower, Followee FROM Loyalty ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, rows, "All", closeErr)
		}
	}()

	var Loyalty []models.Loyalty
	for rows.Next() {
		p := models.Loyalty{}
		err = rows.Scan(&p.ID, &p.Follower, &p.Followee)
		if err != nil {
			return nil, err
		}
		Loyalty = append(Loyalty, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Loyalty, nil
}
