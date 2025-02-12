package migrator

import (
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() {
	dbConnString := os.Getenv("POSTGRES_CONN")
	if dbConnString == "" {
		log.Fatal("POSTGRES_CONN is not set in the .env file")
	}

	m, err := migrate.New("file://migrations", dbConnString)
	if err != nil {
		log.Fatalf("error initializing migrations: %v", err)
	}

	migrationErr := m.Up()
	if migrationErr != nil {
		if errors.Is(migrationErr, migrate.ErrNoChange) {
			log.Println("migrations: No changes")
		} else {
			log.Fatalf("migrations error: %v", migrationErr)
		}
	} else {
		log.Println("database migrations completed successfully.")
	}
}
