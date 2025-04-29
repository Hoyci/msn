package main

import (
	"database/sql"
	"log"
	"msn/services/user-service/internal/config"
	"os"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.GetConfig()

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	driver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})
	if err != nil {
		log.Fatalf("Error creating migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://services/user-service/internal/infra/database/migrate/migrations",
		"public",
		driver,
	)
	if err != nil {
		log.Fatalf("Error creating migrate instance: %v", err)
	}

	cmd := os.Args[len(os.Args)-1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations applied successfully.")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migrations reverted successfully.")
	default:
		log.Println("No command provided. Use 'up' or 'down'.")
	}
}
