package commands

import (
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// EnableCommand is the command to enable notifications
var statusCommand = &discordgo.ApplicationCommandOption{
	Name:        "status",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Check whether you are registered for notifications",
}

// handleEnableCommand handles the enable command
func handleStatusCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild := interaction.GuildID

	var (
		err     error
		message string
	)
	if interaction.Member == nil {
		message = ":exclamation: You can only use this command in a server"
	} else {

		exists, err := storage.RegistrationExists(interaction.Member.User.ID, guild)

		if err != nil {
			slog.Error("Could not check status", "user", interaction.Member.User.ID, "guild", guild, "error", err)
			message = ":exclamation: Could not check the status. Please try again later."
		} else {

			if exists {
				message = ":white_check_mark: You are registered for notifications for this server!"
			} else {
				message = ":x: You are not registered for notifications for this server"
			}
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
