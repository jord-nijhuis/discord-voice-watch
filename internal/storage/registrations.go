package storage

import (
	"database/sql"
	"discord-voice-watch/internal/config"
	"errors"
	"fmt"
	"time"
)

type Registration struct {
	// The user ID
	UserID string
	// The server ID the user is registered for
	ServerID string
	// The time the user was last notified
	LastNotifiedAt *time.Time
	// The channel ID of the last notification sent to the user for this server
	ChannelID *string
	// The message ID of the last notification sent to the user for this server
	MessageID *string
}

func RegisterUser(userID string, serverID string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO registrations (user_id, server_id) VALUES (?, ?) ON CONFLICT DO NOTHING", userID, serverID)

	if err != nil {
		return fmt.Errorf("failed to create registration: %w", err)
	}

	return nil
}

func UnregisterUser(userID string, serverID string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	result, err := db.Exec("DELETE FROM registrations WHERE user_id = ? AND server_id = ?", userID, serverID)

	if err != nil {
		return fmt.Errorf("failed to delete registration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return errors.New("registration not found")
	}

	return nil
}

func GetUsersToNotify(serverID string) ([]string, error) {
	db, err := getDatabase()

	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	cfg, err := config.GetConfig()

	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	lastNotifiedThreshold := time.Now().Add(-cfg.Notifications.DelayBetweenMessages)

	rows, err := db.Query("SELECT user_id FROM registrations WHERE server_id = ? AND (registrations.last_notified_at IS NULL OR last_notified_at < ?)", serverID, lastNotifiedThreshold)

	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("failed to close rows", err)
		}
	}(rows)

	var userIDs []string

	for rows.Next() {
		var userID string

		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}

		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func HasUsersToNotify(serverID string) (bool, error) {
	db, err := getDatabase()

	if err != nil {
		return false, fmt.Errorf("failed to get database: %w", err)
	}

	row := db.QueryRow("SELECT COUNT(*) FROM registrations WHERE server_id = ?", serverID)

	var count int

	err = row.Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to count registrations: %w", err)
	}

	return count > 0, nil
}

func UpdateNotification(userID string, serverID string, notifiedAt time.Time, channelId *string, messageId *string) error {
	db, err := getDatabase()

	if err != nil {
		return err
	}

	result, err := db.Exec("UPDATE registrations SET last_notified_at = ?, channel_id = ?, message_id = ? WHERE user_id = ? AND server_id = ?", notifiedAt, channelId, messageId, userID, serverID)

	if err != nil {
		return fmt.Errorf("failed to update registration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return errors.New("registration not found")
	}

	return nil
}

func GetPreviouslyNotifiedRegistrations(serverID string) ([]Registration, error) {
	db, err := getDatabase()

	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	rows, err := db.Query("SELECT user_id, server_id, last_notified_at, channel_id, message_id FROM registrations WHERE server_id = ? AND last_notified_at IS NOT NULL", serverID)

	if err != nil {
		return nil, fmt.Errorf("failed to find registrations: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("failed to close rows", err)
		}
	}(rows)

	var registrations []Registration

	for rows.Next() {
		var registration Registration

		err := rows.Scan(&registration.UserID, &registration.ServerID, &registration.LastNotifiedAt, &registration.ChannelID, &registration.MessageID)
		if err != nil {
			return nil, err
		}

		registrations = append(registrations, registration)
	}

	return registrations, nil
}
