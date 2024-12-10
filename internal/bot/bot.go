package bot

import (
	projectRoot "discord-voice-watch"
	"discord-voice-watch/internal/bot/events"
	"discord-voice-watch/internal/commands"
	"discord-voice-watch/internal/config"
	"discord-voice-watch/internal/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var cfg config.Config

// Start starts the bot
func Start() {
	slog.Info("Starting bot", "version", projectRoot.Version)
	// Load the configuration
	var err error
	cfg, err = config.LoadConfig()

	if err != nil {
		slog.Error("Could not load config file", "error", err)
		return
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
		Level: cfg.Logging.Level,
	})))

	db, err := storage.InitializeDatabase()

	if err != nil {
		slog.Error("Could not initialize database", err)
		return
	}

	// Initialize Discord session
	dg, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		slog.Error("Could not create discord session", err)
		return
	}

	// Add handlers
	dg.AddHandler(events.OnInteractionCreate)
	dg.AddHandler(events.OnVoiceStateUpdate)
	dg.AddHandler(events.OnGuildCreate)
	dg.AddHandler(events.OnGuildDelete)

	// Make sure we have tbe right permissions
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	// Open the connection
	err = dg.Open()
	if err != nil {
		slog.Error("Could not open connection", err)
		return
	}

	// Register slash commands
	commands.RegisterCommands(dg)

	slog.Info("Bot is running. Press CTRL+C to exit.")

	// Make sure we correctly close the bot when we receive a signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	slog.Info("Received signal, stopping bot")

	slog.Info("Closing discord session")
	// Cleanly close down the Discord session.
	err = dg.Close()

	if err != nil {
		slog.Error("Could not close discord session", err)
	}

	slog.Info("Closing database connection")
	err = db.Close()

	if err != nil {
		slog.Error("Could not close database connection", err)
	}
}
