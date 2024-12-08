package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"strings"
)

// Config The configuration for this bot
type Config struct {
	Discord DiscordConfig
	Logging LoggingConfig
	Bot     BotConfig
}

// DiscordConfig The configuration related to discord
type DiscordConfig struct {
	// Token The discord token of the bot
	Token string
}

type LoggingConfig struct {
	Level slog.Leveler `mapstructure:"-"`
	File  string
}

type BotConfig struct {
	Delay int
}

// LoadConfig Loads the configuration from the config file
// The config file should be named config.yaml and should be in the working directory
func LoadConfig() (Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("Discord.Token", "")
	viper.SetDefault("Logging.Level", "info")
	viper.SetDefault("Bot.Delay", "60")

	err := viper.ReadInConfig() // Find and read the config file

	var configFileNotFoundError viper.ConfigFileNotFoundError

	if errors.As(err, &configFileNotFoundError) {
		err := viper.SafeWriteConfig()

		if err != nil {
			return Config{}, fmt.Errorf("Could not store default config: %s \n", err)
		}
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return config, fmt.Errorf("Fatal error config file: %s \n", err)
	}

	config.Logging.Level = parseLogLevel(viper.GetString("Logging.Level"))

	return config, nil
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		log.Printf("Unknown log level: %s, defaulting to INFO", level)
		return slog.LevelInfo
	}
}
