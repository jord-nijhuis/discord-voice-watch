package storage

import (
	"discord-voice-watch/internal/config"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"time"
)

type Registration struct {
	UserID             uint `gorm:"primaryKey"`
	ServerID           uint `gorm:"primaryKey"`
	LastNotificationAt *time.Time
}

func RegisterUser(userID string, serverID string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	err = CreateUser(userID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user, err := GetUser(userID)

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	server, err := GetServer(serverID)

	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	registration := &Registration{UserID: user.ID, ServerID: server.ID}

	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&registration)

	if result.Error != nil {
		return fmt.Errorf("failed to create registration: %w", result.Error)
	}

	return nil
}

func UnregisterUser(userID string, serverID string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	user, err := GetUser(userID)

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	server, err := GetServer(serverID)

	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	registration := &Registration{UserID: user.ID, ServerID: server.ID}

	result := db.Delete(&registration)

	if result.Error != nil {
		return fmt.Errorf("failed to delete registration: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("registration not found")
	}

	return nil
}

func GetUsersToNotify(serverID string) ([]string, error) {
	db, err := getDatabase()

	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	server, err := GetServer(serverID)

	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	cfg, err := config.GetConfig()

	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	var users []User

	lastNotifiedThreshold := time.Now().Add(-cfg.Notifications.DelayBetweenMessages)

	result := db.Joins("INNER JOIN registrations ON registrations.user_id = users.id AND registrations.last_notification_at < ?", lastNotifiedThreshold).Where("server_id = ?", server.ID).Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %w", result.Error)
	}

	var userIDs []string

	for _, user := range users {
		userIDs = append(userIDs, user.DiscordID)
	}

	return userIDs, nil
}

func HasUsersToNotify(serverID string) (bool, error) {
	db, err := getDatabase()

	if err != nil {
		return false, fmt.Errorf("failed to get database: %w", err)
	}

	var server Server
	result := db.Where("discord_id = ?", serverID).First(&server)

	if result.Error != nil {
		return false, fmt.Errorf("failed to find server: %w", result.Error)
	}

	var count int64

	result = db.Model(&Registration{}).Where("server_id = ?", server.ID).Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to count registrations: %w", result.Error)
	}

	return count > 0, nil
}

func UpdateLastNotificationAt(userID string, serverID string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	user, err := GetUser(userID)

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	server, err := GetServer(serverID)

	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	registration := &Registration{UserID: user.ID, ServerID: server.ID}

	result := db.Model(&registration).Update("last_notification_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to update registration: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("registration not found")
	}

	return nil
}
