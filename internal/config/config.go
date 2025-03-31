package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	LLM struct {
		Provider  string `mapstructure:"provider"`
		Model     string `mapstructure:"model"`
		APIKey    string `mapstructure:"api_key"`
		Endpoint  string `mapstructure:"endpoint"`
		MaxTokens int    `mapstructure:"max_tokens"`
	} `mapstructure:"llm"`

	UI struct {
		Theme         string `mapstructure:"theme"`
		MaxWidth      int    `mapstructure:"max_width"`
		ShowTimestamp bool   `mapstructure:"show_timestamp"`
	} `mapstructure:"ui"`

	Storage struct {
		ChatsDir string `mapstructure:"chats_dir"`
	} `mapstructure:"storage"`
}

// Load reads the configuration from a file and environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("llm.provider", "openai")
	v.SetDefault("llm.model", "gpt-3.5-turbo")
	v.SetDefault("llm.max_tokens", 2000)
	v.SetDefault("ui.max_width", 100)
	v.SetDefault("ui.show_timestamp", true)
	v.SetDefault("storage.chats_dir", "chats")

	// Config file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("configs/")
	v.AddConfigPath("$HOME/.config/gochat/")

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("GOCHAT")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, using defaults and env vars
	}

	// Special handling for API key from environment
	if key := os.Getenv("GOCHAT_LLM_API_KEY"); key != "" {
		v.Set("llm.api_key", key)
	}

	// Create chat directory if it doesn't exist
	chatsDir := v.GetString("storage.chats_dir")
	if err := os.MkdirAll(chatsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create chats directory: %w", err)
	}

	// Convert absolute paths if needed
	if !filepath.IsAbs(chatsDir) {
		absPath, err := filepath.Abs(chatsDir)
		if err == nil {
			v.Set("storage.chats_dir", absPath)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
