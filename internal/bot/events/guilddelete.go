package events

import (
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// OnGuildDelete is called when the bot is removed from a guild
// This function is used to delete the server from the database
func OnGuildDelete(_ *discordgo.Session, guild *discordgo.GuildDelete) {

	slog.Debug("Guild delete event", "guild", guild.ID, "name", guild.Name)

	slog.Info("Deleting server", "guild", guild.ID)

	err := storage.DeleteServer(guild.ID)

	if err != nil {
		slog.Error("Could not delete server", "guild", guild.ID, "error", err)
	}
}
