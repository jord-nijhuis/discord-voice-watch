package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

func JoinWithAnd(items []string) string {
	n := len(items)
	if n == 0 {
		return ""
	}
	if n == 1 {
		return items[0]
	}
	if n == 2 {
		return items[0] + " and " + items[1]
	}
	return strings.Join(items[:n-1], ", ") + " and " + items[n-1]
}

func SendDirectMessage(s *discordgo.Session, userID string, message string) (*discordgo.Message, error) {
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		return nil, fmt.Errorf("could not create a DM channel: %w", err)
	}

	m, err := s.ChannelMessageSend(channel.ID, message)

	if err != nil {
		return nil, fmt.Errorf("could not send a message to the user: %w", err)
	}

	return m, nil
}

func GetAllUsersInVoiceChannel(guild *discordgo.Guild) []string {
	// Get all the users that are currently in a voice channel in the guild
	var users []string

	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == "" {
			continue
		}

		users = append(users, vs.UserID)
	}

	return users
}

func GetDisplayName(s *discordgo.Session, guildID string, userID string) string {

	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		slog.Warn("Could not get guild member", "userID", userID, "guildID", guildID, "error", err)

		return ""
	}

	return member.DisplayName()
}

func GetGuild(s *discordgo.Session, guildId string) (*discordgo.Guild, error) {
	// Fetch the guild's voice states
	guild, err := s.State.Guild(guildId)

	if err != nil {
		// If not in cache, fetch from the API
		guild, err = s.Guild(guildId)

		return guild, err
	}

	return guild, nil
}
