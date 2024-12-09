package bot

import (
	"discord-voice-watch/internal/storage"
	"discord-voice-watch/internal/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"slices"
	"time"
)

func notifyForGuild(s *discordgo.Session, guildID string) {
	time.Sleep(cfg.Notifications.DelayBeforeSending)

	var guild, err = utils.GetGuild(s, guildID)

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
			//continue
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

func removePreviousNotifications(s *discordgo.Session, guildID string) {
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
