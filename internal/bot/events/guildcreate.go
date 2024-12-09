package events

import (
	"discord-voice-watch/internal/notifications"
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"time"
)

func OnGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {

	slog.Debug("Guild create event", "guild", g.Guild.ID, "name", g.Guild.Name)

	// Count members in each voice channel
	for _, vs := range g.VoiceStates {
		if vs.ChannelID != "" {
			storage.IncrementOccupancy(g.ID)
		}
	}

	slog.Info("Setting the occupancy for the guild", "guild", g.Guild.ID, "occupancy", storage.GetOccupancy(g.Guild.ID))

	if time.Since(g.JoinedAt) < time.Minute {
		notifications.SendWelcomeMessage(s, g.Guild)
	}
}
