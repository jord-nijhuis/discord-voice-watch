package events

import (
	"discord-voice-watch/internal/notifications"
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"time"
)

// OnGuildCreate is called when the bot joins a new guild
// This function gets called on two occassions: on startup and when the bot joins a new guild
// We use this event to send a welcome message to the guild if we join a new guild.
// We also use this event to count the number of members in each voice channel
func OnGuildCreate(session *discordgo.Session, guild *discordgo.GuildCreate) {

	slog.Debug("Guild create event", "guild", guild.Guild.ID, "name", guild.Guild.Name)

	setOccupancy(guild.Guild)

	if time.Since(guild.JoinedAt) < time.Minute {
		notifications.SendWelcomeMessage(session, guild.Guild)
	}
}

// setOccupancy sets the occupancy for the guild
func setOccupancy(guild *discordgo.Guild) {
	occpuancy := 0

	for _, vs := range guild.VoiceStates {
		if vs.ChannelID != "" {
			occpuancy += 1
		}
	}

	storage.SetOccupancy(guild.ID, occpuancy)

	slog.Info("Setting the occupancy for the guild", "guild", guild.ID, "occupancy", occpuancy)

}
