package storage

import "sync"

var (
	// channelOccupancy is a map of guild IDs to the number of users in voice channels
	channelOccupancy = make(map[string]int)
	lock             = sync.RWMutex{}
)

func GetOccupancy(guildID string) int {
	lock.RLock()
	val, ok := channelOccupancy[guildID]
	lock.RUnlock()

	if ok {
		return val
	}

	return 0
}

func IncrementOccupancy(guildID string) {

	lock.Lock()
	channelOccupancy[guildID]++
	lock.Unlock()
}

func DecrementOccupancy(guildID string) {
	// Decrement but not below 0
	lock.Lock()
	channelOccupancy[guildID] = max(channelOccupancy[guildID]-1, 0)

	lock.Unlock()
}
