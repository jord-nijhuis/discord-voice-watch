package storage

var (
	guilds           = make(map[string][]string) // Map[guildID]guildID
	channelOccupancy = make(map[string]int)      // Map[guildID][channelID]count
)

func RegisterUser(userID string, guildID string) {

	// Register the user
	guilds[guildID] = append(guilds[guildID], userID)
}

func UnregisterUser(userID string, guildID string) {

	// Unregister the user
	for i, v := range guilds[guildID] {
		if v == userID {
			guilds[guildID] = append(guilds[guildID][:i], guilds[guildID][i+1:]...)
			break
		}
	}

	// Remove the guild if it has no users
	if len(guilds[guildID]) == 0 {
		delete(guilds, guildID)
	}
}

func GetUsers(guildID string) []string {

	// Get the list of users
	val, ok := guilds[guildID]

	if ok {
		return val
	}

	return []string{}
}

func HasUsers(guildID string) bool {

	val, ok := guilds[guildID]

	if !ok {
		return false
	}

	// Check if the guild has users
	return len(val) > 0
}

func GetOccupancy(guildID string) int {

	val, ok := channelOccupancy[guildID]

	if ok {
		return val
	}

	return 0
}

func IncrementOccupancy(guildID string) {

	channelOccupancy[guildID]++
}

func DecrementOccupancy(guildID string) {
	// Decrement but not below 0
	channelOccupancy[guildID] = max(channelOccupancy[guildID]-1, 0)
}
