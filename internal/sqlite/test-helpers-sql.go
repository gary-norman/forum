package sqlite

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

// setupTestDB initializes a test DB by running the schema.sql file
func setupTestDB(t *testing.T) *sql.DB {
	// Open or create a temporary test database (in-memory for fast testing)
	db, err := sql.Open("sqlite3", ":memory:") // or use a file path if you want persistence
	if err != nil {
		t.Fatalf("failed to open in-memory test database: %v", err)
	}

	// Read the schema.sql file into memory
	schema, err := os.ReadFile("/migrations/001_schema.sql")
	if err != nil {
		t.Fatalf("failed to read schema.sql: %v", err)
	}

	// Execute the schema to create tables in the test database
	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatalf("failed to execute schema.sql: %v", err)
	}

	return db
}

func setupTestTx(t *testing.T) (*sql.DB, *sql.Tx) {
	t.Helper()
	db := setupTestDB(t)
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		t.Fatalf("failed to start transaction: %v", err)
	}
	return db, tx
}

// Fixtures holds IDs of inserted records for easy reference in tests.
type Fixtures struct {
	UserID          []byte
	PostID          int64
	CommentID       int64
	SecondUserID    []byte
	SecondPostID    int64
	SecondCommentID int64
}

// insertFixtures populates basic Users, Posts, Comments, etc.
func insertFixtures(db *sql.DB) (*Fixtures, error) {
	userID := uuid.New()
	secondUserID := uuid.New()

	// Insert User 1
	_, err := db.Exec(`
		INSERT INTO Users (ID, SortID, Username, EmailAddress, Usertype, HashedPassword)
		VALUES (?, 1, 'testuser', 'test@example.com', 'member', 'hashedpass')`, userID[:])
	if err != nil {
		return nil, err
	}

	// Insert User 2
	_, err = db.Exec(`
		INSERT INTO Users (ID, SortID, Username, EmailAddress, Usertype, HashedPassword)
		VALUES (?, 2, 'seconduser', 'second@example.com', 'member', 'hashedpass')`, secondUserID[:])
	if err != nil {
		return nil, err
	}

	// Insert Post by User 1
	res, err := db.Exec(`
		INSERT INTO Posts (Title, Content, IsCommentable, Author, AuthorID, Created)
		VALUES ('Test Post', 'Content here', 1, 'testuser', ?, ?)`, userID[:], time.Now())
	if err != nil {
		return nil, err
	}
	postID, _ := res.LastInsertId()

	// Insert Comment by User 2 on Post 1
	res, err = db.Exec(`
		INSERT INTO Comments (Content, CommentedPostID, IsCommentable, Author, AuthorID, ChannelName, ChannelID, Created)
		VALUES ('Nice post!', ?, 1, 'seconduser', ?, 'general', 1, ?)`, postID, secondUserID[:], time.Now())
	if err != nil {
		return nil, err
	}
	commentID, _ := res.LastInsertId()

	// Insert Post 2 by User 2
	res, err = db.Exec(`
		INSERT INTO Posts (Title, Content, IsCommentable, Author, AuthorID, Created)
		VALUES ('Second Post', 'More content', 1, 'seconduser', ?, ?)`, secondUserID[:], time.Now())
	if err != nil {
		return nil, err
	}
	secondPostID, _ := res.LastInsertId()

	// Insert Comment 2 by User 1 on Post 2
	res, err = db.Exec(`
		INSERT INTO Comments (Content, CommentedPostID, IsCommentable, Author, AuthorID, ChannelName, ChannelID, Created)
		VALUES ('Thanks!', ?, 1, 'testuser', ?, 'general', 1, ?)`, secondPostID, userID[:], time.Now())
	if err != nil {
		return nil, err
	}
	secondCommentID, _ := res.LastInsertId()

	return &Fixtures{
		UserID:          userID[:],
		PostID:          postID,
		CommentID:       commentID,
		SecondUserID:    secondUserID[:],
		SecondPostID:    secondPostID,
		SecondCommentID: secondCommentID,
	}, nil
}
