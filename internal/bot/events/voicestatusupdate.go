package events

import (
	"discord-voice-watch/internal/notifications"
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func OnVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	guildID := v.GuildID
	channelID := v.ChannelID

	slog.Debug("Voice state update event", "user", v.UserID, "guild", guildID, "channel", channelID)

	// The new channel id is empty. This means that the user has left a voice channel
	if channelID == "" {
		storage.DecrementOccupancy(guildID)

		slog.Info("User left voice channel for guild", "user", v.UserID, "guild", guildID, "occupancy", storage.GetOccupancy(guildID))

		if storage.GetOccupancy(guildID) == 0 {
			go notifications.RemovePreviousNotifications(s, guildID)
		}

		return
	}

	// If the before update is not nil, the user was previously in a voice channel and is now in another voice channel
	if v.BeforeUpdate != nil {

		// The before and after guild is the same, meaning that the user switched voice channels within the guild
		if v.GuildID == v.BeforeUpdate.GuildID {
			// User changed voice channel in the same guild
			slog.Info("User switched voice chanel in the guild, ignoring", "user", v.UserID, "guild", guildID, "occupancy", storage.GetOccupancy(guildID))
			return
		}

		// The before channel id is not empty, meaning that the user was previously in a voice channel on another guild
		if v.BeforeUpdate.ChannelID != "" {
			storage.DecrementOccupancy(v.BeforeUpdate.GuildID)

			slog.Info("User switched between guilds",
				"user", v.UserID,
				"oldGuild", v.BeforeUpdate.GuildID,
				"oldOccupancy", storage.GetOccupancy(v.BeforeUpdate.GuildID),
				"newGuild", v.GuildID,
				"newOccupancy", storage.GetOccupancy(v.GuildID),
			)

			// If we sent notifications for the previous guild and the occupancy is now 0, remove the notifications
			if storage.GetOccupancy(v.BeforeUpdate.GuildID) == 0 {
				go notifications.RemovePreviousNotifications(s, v.BeforeUpdate.GuildID)
			}
		}
	}

	// User joined a voice channel
	storage.IncrementOccupancy(guildID)
	slog.Info(
		"User joined a voice channel of a guild",
		"user", v.UserID,
		"guild", guildID,
		"occupancy", storage.GetOccupancy(guildID),
	)

	hasUsersToNotify, err := storage.HasUsersToNotify(guildID)

	if err != nil {
		slog.Error("Could not check if there are users to notify", "guild", guildID, "error", err)
		return
	}

	if storage.GetOccupancy(guildID) == 1 && hasUsersToNotify {
		// Start notification process
		go notifications.NotifyForGuild(s, guildID)
	}
}
