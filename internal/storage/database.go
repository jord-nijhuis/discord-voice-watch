package storage

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitializeDatabase() error {
	var err error

	db, err = gorm.Open(sqlite.Open("voice-watch.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {

		return fmt.Errorf("failed to connect to database: %w", err)
	}
	err = db.AutoMigrate(&Server{}, &User{}, &Registration{})

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func getDatabase() (*gorm.DB, error) {

	if db == nil {
		return nil, errors.New("database not initialized")
	}

	return db, nil
}
