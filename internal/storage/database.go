package storage

import (
	"database/sql"
	projectRoot "discord-voice-watch"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"
	"log/slog"
	_ "modernc.org/sqlite"
)

var params = "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-16000)"
var db *sql.DB

// InitializeDatabase initializes the database
func InitializeDatabase() (*sql.DB, error) {
	var err error

	db, err = sql.Open("sqlite", "file:voice-watch.db"+params)

	if err != nil {

		return nil, fmt.Errorf("failed to open the database: %w", err)
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})

	migrationFileSystem, err := iofs.New(projectRoot.MigrationFiles, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		migrationFileSystem,
		"sqlite",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Up()

	if err != nil {

		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Database already up to date")
		} else {
			return nil, fmt.Errorf("failed to apply migrations: %w", err)
		}
	}

	slog.Info("Migrations applied successfully!")

	return db, nil
}

// Database returns the database connection
func Database() (*sql.DB, error) {

	if db == nil {
		return nil, errors.New("database not initialized")
	}

	return db, nil
}
