package storage

var (
	// channelOccupancy is a map of guild IDs to the number of users in voice channels
	channelOccupancy = make(map[string]int)
)

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
