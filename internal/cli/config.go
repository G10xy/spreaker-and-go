/*
config.go - Configuration management commands

GO/COBRA PATTERN: Command Groups
When you have a command with subcommands (like "spreaker config show"),
you create:
  1. A parent command ("config") that has no RunE of its own
  2. Child commands that do the actual work

The parent command just groups related functionality.
*/
package cli

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	"github.com/G10xy/spreaker-and-go/internal/config"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long: `View and modify CLI configuration settings.

Configuration is stored in a YAML file at:
  Linux:   ~/.config/spreaker-cli/config.yaml
  macOS:   ~/Library/Application Support/spreaker-cli/config.yaml
  Windows: %APPDATA%\spreaker-cli\config.yaml`,
		// No RunE here - this is a parent command.
	}

	cmd.AddCommand(
		newConfigShowCmd(),
		newConfigSetCmd(),
		newConfigPathCmd(),
	)

	return cmd
}

// newConfigShowCmd creates the "config show" subcommand.
func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Display current configuration",
		RunE:  runConfigShow,
	}
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Config file: %s\n\n", config.ConfigFilePath())

	// Mask the token for security.
	tokenDisplay := "(not set)"
	if cfg.Token != "" {
		if len(cfg.Token) > 4 {
			tokenDisplay = "****" + cfg.Token[len(cfg.Token)-4:]
		} else {
			tokenDisplay = "****"
		}
	}

	fmt.Printf("token:           %s\n", tokenDisplay)
	fmt.Printf("default_show_id: %d\n", cfg.DefaultShowID)
	fmt.Printf("output_format:   %s\n", cfg.OutputFormat)
	fmt.Printf("api_url:         %s\n", cfg.APIURL)
	return nil
}

// newConfigSetCmd creates the "config set" subcommand.
func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value. Available keys:

  default_show_id  Your default show ID (used when no show ID is specified)
  output_format    Output format: table, json, plain
  api_url          API base URL (for debugging/testing)

Examples:
  spreaker config set default_show_id 12345
  spreaker config set output_format json`,
		Args: cobra.ExactArgs(2),
		RunE: runConfigSet,
	}
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key, value := args[0], args[1]

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Validate and set the value based on key.
	switch key {
	case "default_show_id":
		var id int
		if _, err := fmt.Sscanf(value, "%d", &id); err != nil {
			return fmt.Errorf("invalid show ID: %s", value)
		}
		cfg.DefaultShowID = id

	case "output_format":
		if value != "table" && value != "json" && value != "plain" {
			return fmt.Errorf("invalid format: %s (must be table, json, or plain)", value)
		}
		cfg.OutputFormat = value

	case "api_url":
		cfg.APIURL = value

	default:
		return fmt.Errorf("unknown key: %s", key)
	}

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("✓ Set %s = %s\n", key, value)
	return nil
}

// newConfigPathCmd creates the "config path" subcommand.
func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show config file path",
		// Using Run instead of RunE because this can't fail.
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.ConfigFilePath())
		},
	}
}
