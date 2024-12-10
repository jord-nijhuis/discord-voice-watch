package commands

import (
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// DisableCommand is the command to disable notifications
var disableCommand = &discordgo.ApplicationCommandOption{
	Name:        "disable",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Disable notifications when people start voice chatting in this server",
}

// handleDisableCommand handles the disable command
func handleDisableCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild := interaction.GuildID

	var err error
	message := ":negative_squared_cross_mark: I will no longer send you notifications when someone starts voice chatting in this server"

	if interaction.Member == nil {
		message = ":exclamation: You can only use this command in a server"
	} else {
		err = storage.UnregisterUser(interaction.Member.User.ID, guild)

		if err != nil {
			slog.Error("Could not unregister user", "user", interaction.Member.User.ID, "guild", guild, "error", err)
			message = ":exclamation: Could not unregister you. Please try again later."
		} else {
			slog.Info("User unregistered for guild", "user", interaction.Member.User.ID, "guild", guild)
		}
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		slog.Error("Could not respond to interaction", "error", err)
	}
}
