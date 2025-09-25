package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"time"
)

// ApplyMigrations applies multiple SQL files in order and records them in Migrations table
func runMigrations(db *sql.DB, migrationFiles []string) error {
	// Ensure migrations table exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Migrations (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE,
			AppliedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	// Sort files alphabetically to run in order
	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		var exists int
		err := db.QueryRow("SELECT COUNT(1) FROM Migrations WHERE Name = ?", file).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration record: %w", err)
		}
		if exists > 0 {
			fmt.Printf("Skipping already applied migration: %s\n", file)
			continue
		}

		fmt.Printf("Applying migration: %s\n", file)
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		_, err = tx.Exec(string(sqlBytes))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		_, err = tx.Exec("INSERT INTO Migrations (Name, AppliedAt) VALUES (?, ?)", file, time.Now().UTC())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		fmt.Printf("Successfully applied migration: %s\n", file)
	}

	return nil
}
