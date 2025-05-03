package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
	"log"
	"os"
	"path/filepath"
)

func ResetDb() {
    os.Create("example.db")

	// SQLite DB connection
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Directory containing DDL SQL files
	ddlDir := "/Users/giuseppe/Documents/personal/gapt-atc-2025/backend/internal/db/ddl"

	// Apply all DDL files in the directory
	err = applyDDLsFromDirectory(db, ddlDir)
	if err != nil {
		log.Fatalf("Error applying DDLs: %v", err)
	}

	fmt.Println("All DDLs applied successfully!")
}

// applyDDLsFromDirectory reads all SQL files in a directory and applies them to the SQLite database
func applyDDLsFromDirectory(db *sql.DB, ddlDir string) error {
	// Read all files in the directory
	files, err := os.ReadDir(ddlDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", ddlDir, err)
	}

	// Loop over all files in the directory
	for _, file := range files {
		// Only process .sql files
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(ddlDir, file.Name())
			err := applyDDL(db, filePath)
			if err != nil {
				return fmt.Errorf("failed to apply DDL from file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// applyDDL reads the DDL from a file and executes it on the SQLite database
func applyDDL(db *sql.DB, ddlFile string) error {
	// Read DDL file
	ddlContent, err := os.ReadFile(ddlFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", ddlFile, err)
	}

	// Execute the DDL query
	_, err = db.Exec(string(ddlContent))
	if err != nil {
		return fmt.Errorf("failed to execute DDL from file %s: %w", ddlFile, err)
	}

	fmt.Printf("Applied DDL from %s\n", ddlFile)
	return nil
}

func SeedDb() {
	// SQLite DB connection
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Directory containing DDL SQL files
	ddlDir := "/Users/giuseppe/Documents/personal/gapt-atc-2025/backend/internal/db/seeds"

	// Apply all DDL files in the directory
	err = applyDDLsFromDirectory(db, ddlDir)
	if err != nil {
		log.Fatalf("Error seeding DDLs: %v", err)
	}

	fmt.Println("All seeds applied successfully!")
}

func Refresh() {
	ResetDb()
	SeedDb()
}
