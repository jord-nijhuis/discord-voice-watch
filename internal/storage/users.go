package storage

import "fmt"

// CreateUser creates a user in the database
// If the user already exists, this function does nothing
func CreateUser(userID string) error {

	db, err := Database()

	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	_, err = db.Exec("INSERT INTO users (id) VALUES (?) ON CONFLICT DO NOTHING", userID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
