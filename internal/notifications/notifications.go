package notifications

import (
	"discord-voice-watch/internal/config"
	"discord-voice-watch/internal/storage"
	"discord-voice-watch/internal/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"slices"
	"time"
)

func NotifyForGuild(s *discordgo.Session, guildID string) {

	cfg, err := config.GetConfig()

	if err != nil {
		slog.Error("Could not get config", "error", err)
		return
	}

	time.Sleep(cfg.Notifications.DelayBeforeSending)

	guild, err := utils.GetGuild(s, guildID)

	if err != nil {

		slog.Error("Could not retrieve guild for id",
			"id", guildID,
			"error", err,
		)
		return
	}

	usersInVoice := utils.GetAllUsersInVoiceChannel(guild)

	if len(usersInVoice) == 0 {
		return
	}

	var names []string

	for _, userID := range usersInVoice {
		names = append(names, utils.GetDisplayName(s, guildID, userID))
	}

	usersToNotify, err := storage.GetUsersToNotify(guildID)

	if err != nil {
		slog.Error("Could not retrieve users to notify for guild",
			"id", guildID,
			"error", err,
		)
		return
	}

	for _, userID := range usersToNotify {
		// Don't send a message to the user if they are already in the voice chat
		if slices.Contains(usersInVoice, userID) {
			continue
		}

		slog.Info("Sending notification for guild to user", "guild", guildID, "user", userID)

		verb := "is"

		if len(names) > 1 {
			verb = "are"
		}

		message, err := utils.SendDirectMessage(
			s,
			userID,
			fmt.Sprintf("%s %s now in voice chat of %s",
				utils.JoinWithAnd(names),
				verb,
				guild.Name),
		)

		if err != nil {
			slog.Error("Could not send notification to user",
				"user", userID,
				"guild", guildID,
				"error", err,
			)
			continue
		}

		err = storage.UpdateNotification(userID, guildID, time.Now(), &message.ChannelID, &message.ID)

		if err != nil {
			slog.Error("Could not update last notification time",
				"user", userID,
				"guild", guildID,
				"error", err,
			)
		}
	}
}

func RemovePreviousNotifications(s *discordgo.Session, guildID string) {
	cfg, err := config.GetConfig()

	if err != nil {
		slog.Error("Could not get config", "error", err)
		return
	}

	time.Sleep(cfg.Notifications.DelayBeforeSending)

	if storage.GetOccupancy(guildID) > 0 {
		return
	}

	registrations, err := storage.GetPreviouslyNotifiedRegistrations(guildID)

	if err != nil {
		slog.Error("Could not retrieve notifications to remove",
			"id", guildID,
			"error", err,
		)
		return
	}

	slog.Info("Removing previous notifications for guild", "guild", guildID)

	for _, registration := range registrations {

		if registration.ChannelID == nil || registration.MessageID == nil {
			continue
		}

		slog.Info("Deleting previous notification for guild to user", "guild", guildID, "user", registration.UserID)

		err = s.ChannelMessageDelete(*registration.ChannelID, *registration.MessageID)

		if err != nil {
			slog.Error("Could not delete previous notification",
				"channel", *registration.ChannelID,
				"message", *registration.MessageID,
				"error", err,
			)
		}

		err = storage.UpdateNotification(registration.UserID, guildID, *registration.LastNotifiedAt, nil, nil)

		if err != nil {
			slog.Error("Could not update last notification time",
				"user", registration,
				"guild", guildID,
				"error", err,
			)
		}
	}
}

func SendWelcomeMessage(session *discordgo.Session, guild *discordgo.Guild) {

	slog.Info("Sending a welcome message", "guild", guild.ID)
	// Choose a channel to send the message
	// Typically, the first available text channel is used
	var defaultChannel *discordgo.Channel
	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			defaultChannel = channel
			break
		}
	}

	if defaultChannel == nil {
		slog.Info("No suitable channel found for welcome message", "guild", guild.ID)
		return
	}
	// If a suitable channel is found, send a message
	_, err := session.ChannelMessageSend(defaultChannel.ID, "Hello there 👋! You can use me to get notified when someone starts voice chatting in this server."+
		" Use `/voice-watch enable` to enable notifications and as soon as someone joins a voice channel, I'll let you "+
		"know in a direct message. You can always disable me by using `/voice-watch disable`.")

	if err != nil {
		slog.Error("Could not send welcome message", "guild", guild.ID, "channel", defaultChannel.ID, "error", err)
	}
}
