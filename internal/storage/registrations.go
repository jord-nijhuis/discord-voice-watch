package storage

import (
	"discord-voice-watch/internal/config"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"time"
)

type Registration struct {
	UserID             uint   `gorm:"primaryKey"`
	User               User   `gorm:"foreignKey:UserID"`
	ServerID           uint   `gorm:"primaryKey"`
	Server             Server `gorm:"foreignKey:ServerID"`
	LastNotificationAt *time.Time
	MessageID          *string
	ChannelID          *string
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

func UpdateNotification(userID string, serverID string, notifiedAt time.Time, channelId *string, messageId *string) error {
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

	result := db.Model(&registration).Update("last_notification_at", notifiedAt).Update("channel_id", channelId).Update("message_id", messageId)

	if result.Error != nil {
		return fmt.Errorf("failed to update registration: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("registration not found")
	}

	return nil
}

func GetPreviouslyNotifiedRegistrations(serverID string) ([]Registration, error) {
	db, err := getDatabase()

	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	server, err := GetServer(serverID)

	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	var registrations []Registration

	result := db.Preload("User").Where("server_id = ? AND last_notification_at IS NOT NULL", server.ID).Find(&registrations)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find registrations: %w", result.Error)
	}

	return registrations, nil
}
