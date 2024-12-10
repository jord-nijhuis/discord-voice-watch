package commands

import (
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// EnableCommand is the command to enable notifications
var enableCommand = &discordgo.ApplicationCommandOption{
	Name:        "enable",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Enable notifications when people start voice chatting in this server",
}

// handleEnableCommand handles the enable command
func handleEnableCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	guild := interaction.GuildID

	var err error
	message := ":white_check_mark: As soon as someone starts voice chatting in this server, I will send you a message!"

	if interaction.Member == nil {
		message = ":exclamation: You can only use this command in a server"
	} else {

		err = storage.CreateUser(interaction.Member.User.ID)

		if err == nil {
			err = storage.RegisterUser(interaction.Member.User.ID, guild)
		}

		if err != nil {
			slog.Error("Could not register user", "user", interaction.Member.User.ID, "guild", guild, "error", err)
			message = ":exclamation: Could not register you. Please try again later."
		} else {
			slog.Info("User registered for guild", "user", interaction.Member.User.ID, "guild", guild)
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
