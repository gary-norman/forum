// Package db provides functions to initialize and manage the SQLite database.
// It sets up the database with foreign key support, applies recommended PRAGMAs
package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string, schemaFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enforce foreign keys
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Recommended production PRAGMAs
	pragmas := []string{
		`PRAGMA journal_mode = WAL;`,
		`PRAGMA synchronous = NORMAL;`,
		`PRAGMA ignore_check_constraints = OFF;`,
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to apply pragma '%s': %w", pragma, err)
		}
	}

	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute the schema as a single Exec call
	if _, err := db.Exec(string(schema)); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to execute schema: %w", err)
	}

	return db, nil
}
