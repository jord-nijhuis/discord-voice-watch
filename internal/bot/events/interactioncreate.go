package events

import (
	"discord-voice-watch/internal/commands"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Debug("Interaction create event", "interaction", i.ID, "type", i.Type, "guild", i.GuildID, "channel", i.ChannelID, "member", i.Member.User.ID)

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		commands.HandleCommand(s, i)
	}
}
