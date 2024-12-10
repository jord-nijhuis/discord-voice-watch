package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// RegisterCommands registers the commands with Discord
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

// HandleCommand handles the command
func HandleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	if interaction.ApplicationCommandData().Name != "voice-watch" {
		slog.Warn("Unknown command", "command", interaction.ApplicationCommandData().Name)

		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Unknown command %session", interaction.ApplicationCommandData().Name),
			},
		})

		return
	}

	switch interaction.ApplicationCommandData().Options[0].Name {
	case "enable":
		handleEnableCommand(session, interaction)
	case "disable":
		handleDisableCommand(session, interaction)
	default:
		slog.Warn("Unknown subcommand", "subcommand", interaction.ApplicationCommandData().Options[0].Name)
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Unknown subcommand %session", interaction.ApplicationCommandData().Options[0].Name),
			},
		})
	}
}
