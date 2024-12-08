package commands

import (
	"discord-voice-watch/internal/storage"
	"fmt"
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
	storage.RegisterUser(i.Member.User.ID, guild)

	slog.Info("User registered for guild", "user", i.Member.User.ID, "guild", guild)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("As soon as someone starts voice chatting in this server, you will be notified."),
		},
	})

	if err != nil {
		slog.Error("Could not respond to interaction", "error", err)
	}
}
