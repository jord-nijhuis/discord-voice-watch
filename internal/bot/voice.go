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

func onVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	guildID := v.GuildID
	channelID := v.ChannelID

	if channelID == "" {
		storage.DecrementOccupancy(guildID)

		slog.Info("User left voice channel for guild", "user", v.UserID, "guild", guildID, "occupancy", storage.GetOccupancy(guildID))
		return
	}

	if v.BeforeUpdate != nil {
		if v.GuildID == v.BeforeUpdate.GuildID {
			// User changed voice channel in the same guild
			slog.Info("User switched voice chanel in the guild, ignoring", "user", v.UserID, "guild", guildID, "occupancy", storage.GetOccupancy(guildID))
			return
		}

		if v.BeforeUpdate.ChannelID != "" {
			storage.DecrementOccupancy(v.BeforeUpdate.GuildID)

			slog.Info("User switched between guilds",
				"user", v.UserID,
				"oldGuild", v.BeforeUpdate.GuildID,
				"oldOccupancy", storage.GetOccupancy(v.BeforeUpdate.GuildID),
				"newGuild", v.GuildID,
				"newOccupancy", storage.GetOccupancy(v.GuildID),
			)

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
		go notifyIfStillOccupied(s, guildID)
	}
}

func notifyIfStillOccupied(s *discordgo.Session, guildID string) {
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

		utils.SendDirectMessage(
			s,
			userID,
			fmt.Sprintf("%s %s now in voice chat of %s",
				utils.JoinWithAnd(names),
				verb,
				guild.Name),
		)

		err = storage.UpdateLastNotificationAt(userID, guildID)

		if err != nil {
			slog.Error("Could not update last notification time",
				"user", userID,
				"guild", guildID,
				"error", err,
			)
		}
	}
}

func SyncOccupancies(s *discordgo.Session) error {
	var after string

	slog.Info("Synchronizing channel occupancy")

	for {
		// Fetch all guilds the bot is part of
		userGuilds, err := s.UserGuilds(200, "", after, false)
		if err != nil {
			return fmt.Errorf("failed to fetch guilds of the bot: %w", err)
		}

		for _, userGuild := range userGuilds {
			after = userGuild.ID

			err = storage.CreateServer(userGuild.ID)

			if err != nil {
				slog.Warn("Could not create server", "id", userGuild.ID, "error", err)
				continue
			}

			slog.Info("Initializing voice channel occupancy for guild", "name", userGuild.Name, "id", userGuild.ID)

			// Fetch the guild's voice states
			guild, err := utils.GetGuild(s, userGuild.ID)

			if err != nil {
				slog.Warn("Could not fetch voice data for guild", "id", userGuild.ID, "error", err)
				continue
			}

			// Count members in each voice channel
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID != "" {
					storage.IncrementOccupancy(guild.ID)
				}
			}
		}

		// If fewer than 200 guilds were returned, we are done
		if len(userGuilds) < 200 {
			break
		}
	}

	slog.Info("Voice channel occupancy successfully synchronized")
	return nil
}
