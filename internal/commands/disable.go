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

	storage.UnregisterUser(i.Member.User.ID, guild)

	slog.Info("User unregistered for guild", "user", i.Member.User.ID, "guild", guild)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You will no longer be notified.",
		},
	})

	if err != nil {
		slog.Error("Could not respond to interaction", "error", err)
	}
}
