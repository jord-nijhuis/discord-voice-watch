package commands

import (
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

var enableCommand = &discordgo.ApplicationCommandOption{
	Name:        "enable",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Enable notifications when people start voice chatting in this server",
}

func handleEnableCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild := i.GuildID

	var err error
	message := "As soon as someone starts voice chatting in this server, you will be notified"

	if i.Member == nil {
		message = "You can only use this command in a server"
	} else {
		err = storage.RegisterUser(i.Member.User.ID, guild)

		if err != nil {
			slog.Error("Could not register user", "user", i.Member.User.ID, "guild", guild, "error", err)
			message = "Could not register you. Please try again later."
		} else {
			slog.Info("User registered for guild", "user", i.Member.User.ID, "guild", guild)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if err != nil {
		slog.Error("Could not respond to interaction", "error", err)
	}
}
