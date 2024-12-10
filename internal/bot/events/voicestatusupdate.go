package events

import (
	"discord-voice-watch/internal/notifications"
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

// OnVoiceStateUpdate is called when a user joins or leaves a voice channel
func OnVoiceStateUpdate(session *discordgo.Session, voiceStateUpdate *discordgo.VoiceStateUpdate) {
	guildID := voiceStateUpdate.GuildID
	channelID := voiceStateUpdate.ChannelID

	slog.Debug("Voice state update event", "user", voiceStateUpdate.UserID, "guild", guildID, "channel", channelID)

	// The new channel id is empty. This means that the user has left a voice channel
	if channelID == "" {
		occupancy := storage.DecrementOccupancy(guildID)

		slog.Info("User left voice channel for guild", "user", voiceStateUpdate.UserID, "guild", guildID, "occupancy", storage.GetOccupancy(guildID))

		if occupancy == 0 {
			go notifications.RemovePreviousNotifications(session, guildID)
		}

		return
	}

	// If the before update is not nil, the user was previously in a voice channel and is now in another voice channel
	if voiceStateUpdate.BeforeUpdate != nil {

		// The before and after guild is the same, meaning that the user switched voice channels within the guild
		if voiceStateUpdate.GuildID == voiceStateUpdate.BeforeUpdate.GuildID {
			// User changed voice channel in the same guild
			slog.Info("User switched voice chanel in the guild, ignoring", "user", voiceStateUpdate.UserID, "guild", guildID, "occupancy", storage.GetOccupancy(guildID))
			return
		}

		// The before channel id is not empty, meaning that the user was previously in a voice channel on another guild
		if voiceStateUpdate.BeforeUpdate.ChannelID != "" {
			occupancy := storage.DecrementOccupancy(voiceStateUpdate.BeforeUpdate.GuildID)

			slog.Info("User switched between guilds",
				"user", voiceStateUpdate.UserID,
				"oldGuild", voiceStateUpdate.BeforeUpdate.GuildID,
				"oldOccupancy", storage.GetOccupancy(voiceStateUpdate.BeforeUpdate.GuildID),
				"newGuild", voiceStateUpdate.GuildID,
				"newOccupancy", storage.GetOccupancy(voiceStateUpdate.GuildID),
			)

			// If we sent notifications for the previous guild and the occupancy is now 0, remove the notifications
			if occupancy == 0 {
				go notifications.RemovePreviousNotifications(session, voiceStateUpdate.BeforeUpdate.GuildID)
			}
		}
	}

	// User joined a voice channel
	occupancy := storage.IncrementOccupancy(guildID)
	slog.Info(
		"User joined a voice channel of a guild",
		"user", voiceStateUpdate.UserID,
		"guild", guildID,
		"occupancy", occupancy,
	)

	hasUsersToNotify, err := storage.HasUsersToNotify(guildID)

	if err != nil {
		slog.Error("Could not check if there are users to notify", "guild", guildID, "error", err)
		return
	}

	if occupancy == 1 && hasUsersToNotify {
		// Start notification process
		go notifications.NotifyActivity(session, guildID)
	}
}
