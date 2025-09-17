package migrations

import (
	"database/sql"
	"log"
	"os"
)

func Migrate(db *sql.DB, migrationPath string) {
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration applied successfully")
}
