package bot

import (
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

func Start() {
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

	err = storage.InitializeDatabase()

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
	dg.AddHandler(commands.OnInteractionCreate)
	dg.AddHandler(onVoiceStateUpdate)

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

	// Sync the occupancies
	err = SyncOccupancies(dg)

	if err != nil {
		slog.Error("Could not sync occupancy", err)
		return
	}

	slog.Info("Notifications is running. Press CTRL+C to exit.")

	// Make sure we correctly close the bot when we receive a signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	err = dg.Close()

	if err != nil {
		slog.Error("Could not close discord session", err)
	}
}
