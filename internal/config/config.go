package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"strings"
	"time"
)

// Config The configuration for this bot
type Config struct {
	Discord       DiscordConfig
	Logging       LoggingConfig
	Notifications NotificationsConfig
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

type NotificationsConfig struct {
	DelayBeforeSending   time.Duration `mapstructure:"delay-before-sending"`
	DelayBetweenMessages time.Duration `mapstructure:"delay-between-messages"`
	NotifySelf           bool          `mapstructure:"notify-self"`
}

var config Config

// LoadConfig Loads the configuration from the config file
// The config file should be named config.yaml and should be in the working directory
func LoadConfig() (Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("Discord.Token", "")
	viper.SetDefault("Logging.Level", "info")
	viper.SetDefault("Notifications.delay-before-sending", time.Minute)
	viper.SetDefault("Notifications.delay-between-messages", time.Hour)
	viper.SetDefault("Notifications.notify-self", false)

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

// parseLogLevel Parses the log level from a string
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

// LoadedConfig Gets the configuration
func LoadedConfig() (Config, error) {
	if config == (Config{}) {
		return LoadConfig()
	}

	return config, nil
}
