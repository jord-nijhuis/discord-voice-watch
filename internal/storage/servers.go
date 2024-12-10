package storage

import "fmt"

// CreateServer creates a server in the database
// If the server already exists, this function does nothing
func CreateServer(guildID string) error {

	db, err := Database()

	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	_, err = db.Exec("INSERT INTO servers (id) VALUES (?) ON CONFLICT DO NOTHING", guildID)

	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	return nil
}

func DeleteServer(guildId string) error {
	db, err := Database()

	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	_, err = db.Exec("DELETE FROM servers WHERE id = ?", guildId)

	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	return nil
}
