package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// For Viper map config file keys
type Config struct {
	Token string `mapstructure:"token"`

	DefaultShowID int `mapstructure:"default_show_id"`

	// OutputFormat controls how results are displayed: "table", "json", "plain"
	OutputFormat string `mapstructure:"output_format"`

	APIURL string `mapstructure:"api_url"`
}

func DefaultConfig() *Config {
	return &Config{
		Token:         "",
		DefaultShowID: 0,
		OutputFormat:  "table",
		APIURL:        "https://api.spreaker.com",
	}
}

// configDir returns the directory where config files are stored.
func configDir() (string, error) {
	// First, check if user set a custom config location
	if dir := os.Getenv("SPREAKER_CONFIG_DIR"); dir != "" {
		return dir, nil
	}

	// Use the user's config directory (OS-appropriate)
	// Linux:   ~/.config
	// macOS:   ~/Library/Application Support
	// Windows: %APPDATA%
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine config directory: %w", err)
	}

	return filepath.Join(userConfigDir, "spreaker-cli"), nil
}

func configFilePath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load reads configuration from file, environment, and returns a Config.
func Load() (*Config, error) {
	cfg := DefaultConfig()

	dir, err := configDir()
	if err != nil {
		return cfg, err
	}

	viper.SetConfigName("config") 
	viper.SetConfigType("yaml")   
	viper.AddConfigPath(dir)      

	viper.SetEnvPrefix("SPREAKER")
	viper.AutomaticEnv() 

	viper.SetDefault("token", cfg.Token)
	viper.SetDefault("default_show_id", cfg.DefaultShowID)
	viper.SetDefault("output_format", cfg.OutputFormat)
	viper.SetDefault("api_url", cfg.APIURL)

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Error may be due to the fact the user just hasn't configured yet
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return cfg, fmt.Errorf("error reading config file: %w", err)
		}
		// File not found is fine, continue with defaults + env vars
	}

	// Unmarshal the combined configuration into our struct
	if err := viper.Unmarshal(cfg); err != nil {
		return cfg, fmt.Errorf("error parsing config: %w", err)
	}

	return cfg, nil
}

// Save writes the given configuration to the config file.
func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// 0700 so that owner can read/write/execute while others have no access
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	viper.Set("token", cfg.Token)
	viper.Set("default_show_id", cfg.DefaultShowID)
	viper.Set("output_format", cfg.OutputFormat)
	viper.Set("api_url", cfg.APIURL)

	configPath, err := configFilePath()
	if err != nil {
		return err
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

// SaveToken is a convenience function to save just the API token.
func SaveToken(token string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.Token = token
	return Save(cfg)
}

func GetToken() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	if cfg.Token == "" {
		return "", errors.New("not authenticated. Run 'spreaker login' first")
	}

	return cfg.Token, nil
}

func ConfigFilePath() string {
	path, err := configFilePath()
	if err != nil {
		return "(unknown)"
	}
	return path
}
