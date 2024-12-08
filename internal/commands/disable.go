package commands

import (
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

var disableCommand = &discordgo.ApplicationCommandOption{
	Name:        "disable",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Disable notifications when people start voice chatting in this server",
}

func handleDisableCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild := i.GuildID

	message := "You will no longer be notified"
	var err error

	if i.Member == nil {
		message = "You can only use this command in a server"
	} else {
		err = storage.UnregisterUser(i.Member.User.ID, guild)

		if err != nil {
			slog.Error("Could not unregister user", "user", i.Member.User.ID, "guild", guild, "error", err)
			message = "Could not unregister you. Please try again later."
		} else {
			slog.Info("User unregistered for guild", "user", i.Member.User.ID, "guild", guild)
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
