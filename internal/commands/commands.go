package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func RegisterCommands(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "voice-watch",
			Description: "All commands related to the voice watch bot",
			Options:     []*discordgo.ApplicationCommandOption{enableCommand, disableCommand},
		},
	}

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", commands)

	if err != nil {
		slog.Error("Could not create commands", "commands", commands, "error", err)
	}
}

func HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.ApplicationCommandData().Name != "voice-watch" {
		slog.Warn("Unknown command", "command", i.ApplicationCommandData().Name)

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Unknown command %s", i.ApplicationCommandData().Name),
			},
		})

		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "enable":
		handleEnableCommand(s, i)
	case "disable":
		handleDisableCommand(s, i)
	default:
		slog.Warn("Unknown subcommand", "subcommand", i.ApplicationCommandData().Options[0].Name)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Unknown subcommand %s", i.ApplicationCommandData().Options[0].Name),
			},
		})
	}
}
