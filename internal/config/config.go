package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Observability ObservabilityConfig
	LogLevel      string
	LogFormat     string
}

// ServerConfig holds HTTP server configuration

// ObservabilityConfig holds observability server configuration
type ObservabilityConfig struct {
	Enabled bool
	Addr    string
}

// Load reads configuration from file and environment variables
func Load(cfgFile string) (*Config, error) {
	// Set defaults
	setDefaults()

	// Config file setup
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".broken-hexagon")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
	}

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; using defaults and environment variables
	}

	// Environment variables override config file
	viper.SetEnvPrefix("BROKEN-HEXAGON")
	viper.AutomaticEnv()

	// Unmarshal config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	// Server defaults

	// Logger defaults
	viper.SetDefault("loglevel", "info")
	viper.SetDefault("logformat", "json")

	// Observability defaults
	viper.SetDefault("observability.enabled", false)
	viper.SetDefault("observability.addr", ":8080")
}
