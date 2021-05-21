package database

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func DefaultMigrate(log logger) {
	cfg := parseConfig()
	ApplyDefaultMigrations(log, cfg.PostgresURI)
}

func ApplyDefaultMigrations(log logger, postgresURI string) {
	d, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		log.Fatalf("Failed to read embedded migrations. Error: %+v", err)
	}
	ApplyToPostgres(log, postgresURI, d)
}

func ApplyToPostgres(log logger, postgresURI string, driver source.Driver) {
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		log.Fatalf("Failed connecting to postgres database. Error: %v", err)
		return
	}
	ApplyToPostgresInstance(log, db, driver)
}

func ApplyToPostgresInstance(log logger, db *sql.DB, sDriver source.Driver) {
	dDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed instantiating postgres driver. Error: %v", err)
		return
	}
	m, err := migrate.NewWithInstance("iofs", sDriver, "postgres", dDriver)
	if err != nil {
		log.Fatalf("Failed instantiating golang migrate. Error %v", err)
		return
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Info("No changes to the database required.")
	} else if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return
	}
	log.Info("Successfully performed migrate step.")
}

type logger interface {
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
}
