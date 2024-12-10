package storage

import "sync"

var (
	// channelOccupancy is a map of guild IDs to the number of users in voice channels
	channelOccupancy = make(map[string]int)
	lock             = sync.RWMutex{}
)

// GetOccupancy returns the occupancy for the guild
func GetOccupancy(guildID string) int {
	lock.RLock()
	val, ok := channelOccupancy[guildID]
	lock.RUnlock()

	if ok {
		return val
	}

	return 0
}

// SetOccupancy sets the occupancy for the guild
func SetOccupancy(guildID string, occupancy int) {
	lock.Lock()
	channelOccupancy[guildID] = occupancy
	lock.Unlock()
}

// IncrementOccupancy increments the occupancy for the guild
func IncrementOccupancy(guildID string) int {

	lock.Lock()
	channelOccupancy[guildID]++
	defer lock.Unlock()

	return channelOccupancy[guildID]
}

// DecrementOccupancy decrements the occupancy for the guild
func DecrementOccupancy(guildID string) int {
	// Decrement but not below 0
	lock.Lock()
	channelOccupancy[guildID] = max(channelOccupancy[guildID]-1, 0)
	defer lock.Unlock()

	return channelOccupancy[guildID]
}
