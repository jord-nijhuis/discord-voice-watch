package events

import (
	"discord-voice-watch/internal/commands"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// OnInteractionCreate is called when an interaction is created
// We use this to handle slash commands
func OnInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	slog.Debug(
		"Interaction create event",
		"interaction", interaction.ID,
		"type", interaction.Type,
		"guild", interaction.GuildID,
		"channel", interaction.ChannelID,
		"member", interaction.Member.User.ID,
	)

	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		commands.HandleCommand(session, interaction)
	}
}
