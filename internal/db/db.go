// Package db provides functions to initialize and manage the SQLite database.
// It sets up the database with foreign key support, applies recommended PRAGMAs
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string, schemaFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enforce foreign keys
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
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
			log.Printf("Warning: failed to apply pragma '%s': %v", pragma, err)
		}
	}

	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Split and execute each statement individually
	statements := strings.SplitSeq(string(schema), ";")
	for stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed != "" {
			if _, err := db.Exec(trimmed + ";"); err != nil {
				log.Printf("Error executing statement: %q\nError: %v", trimmed, err)
			}
		}
	}

	return db, nil
}
