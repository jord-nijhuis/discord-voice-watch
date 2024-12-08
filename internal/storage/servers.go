package storage

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Server struct {
	gorm.Model
	DiscordID     string `gorm:"uniqueIndex"`
	Registrations []Registration
}

func CreateServer(discordID string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	server := &Server{DiscordID: discordID}

	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&server)

	if result.Error != nil {
		return fmt.Errorf("failed to create guild: %w", result.Error)
	}

	return nil
}

func GetServer(discordID string) (*Server, error) {
	db, err := getDatabase()

	if err != nil {
		return nil, err
	}

	var server Server
	result := db.Where("discord_id = ?", discordID).First(&server)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find guild: %w", result.Error)
	}

	return &server, nil
}
