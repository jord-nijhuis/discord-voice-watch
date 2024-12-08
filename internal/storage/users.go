package storage

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	DiscordID     string `gorm:"uniqueIndex"`
	Registrations []Registration
}

func CreateUser(discordID string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	user := &User{DiscordID: discordID}

	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

func GetUser(discordID string) (*User, error) {
	db, err := getDatabase()

	if err != nil {
		return nil, err
	}

	var user User
	result := db.Where("discord_id = ?", discordID).First(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return &user, nil
}
