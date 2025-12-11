/*
Package main is the entry point for the spreaker-cli application.

This file sets up the CLI command structure using Cobra and integrates:
  - Configuration management (internal/config)
  - API client (internal/api)
  - Output formatting (internal/output)
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/G10xy/spreaker-and-go/internal/api"
	"github.com/G10xy/spreaker-and-go/internal/config"
	"github.com/G10xy/spreaker-and-go/internal/output"
)

// version is set at build time using:
//
//	go build -ldflags "-X main.version=1.0.0"
var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------
// Root Command
// -----------------------------------------------------------------------------

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "spreaker",
		Short:        "A CLI for the Spreaker podcast platform",
		Long:         `spreaker-cli is a command line interface for managing your podcasts on Spreaker.`,
		Version:      version,
		SilenceUsage: true,
	}

	// Global flags (available to all subcommands)
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output format: table, json, plain")
	rootCmd.PersistentFlags().String("token", "", "API token (overrides config)")

	// Add command groups
	rootCmd.AddCommand(
		newLoginCmd(),
		newMeCmd(),
		newShowsCmd(),
		newEpisodesCmd(),
		newConfigCmd(),
	)

	return rootCmd
}

// -----------------------------------------------------------------------------
// Helper Functions
// -----------------------------------------------------------------------------

// getClient creates an API client using token from flag, env, or config
func getClient(cmd *cobra.Command) (*api.Client, error) {
	// Priority: flag > env > config
	token, _ := cmd.Flags().GetString("token")

	if token == "" {
		var err error
		token, err = config.GetToken()
		if err != nil {
			return nil, err
		}
	}

	// Load config for base URL
	cfg, _ := config.Load()

	return api.NewClientWithOptions(token, cfg.APIURL, 0), nil
}

// getFormatter creates an output formatter using format from flag or config
func getFormatter(cmd *cobra.Command) *output.Formatter {
	format, _ := cmd.Flags().GetString("output")

	if format == "" {
		cfg, _ := config.Load()
		format = cfg.OutputFormat
	}

	return output.New(format)
}

// -----------------------------------------------------------------------------
// Login Command
// -----------------------------------------------------------------------------

func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Spreaker",
		Long: `Authenticate with your Spreaker account.

You'll need an API token from your Spreaker developer settings.
The token will be saved to your config file for future use.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Enter your Spreaker API token: ")

			var token string
			if _, err := fmt.Scanln(&token); err != nil {
				return fmt.Errorf("failed to read token: %w", err)
			}

			if token == "" {
				return fmt.Errorf("token cannot be empty")
			}

			// Validate token by making a test API call
			client := api.NewClient(token)
			user, err := client.GetMe()
			if err != nil {
				return fmt.Errorf("invalid token: %w", err)
			}

			// Save the token
			if err := config.SaveToken(token); err != nil {
				return fmt.Errorf("failed to save token: %w", err)
			}

			fmt.Printf("✓ Logged in as %s (@%s)\n", user.Fullname, user.Username)
			fmt.Printf("  Token saved to %s\n", config.ConfigFilePath())
			return nil
		},
	}
}

// -----------------------------------------------------------------------------
// Me Command
// -----------------------------------------------------------------------------

func newMeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show current authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			user, err := client.GetMe()
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintUser(user)
			return nil
		},
	}
}

// -----------------------------------------------------------------------------
// Config Commands
// -----------------------------------------------------------------------------

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}

	// config show
	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Display current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			fmt.Printf("Config file: %s\n\n", config.ConfigFilePath())

			// Mask token
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
		},
	})

	// config set
	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value. Available keys:
  - default_show_id: Your default show ID
  - output_format:   Output format (table, json, plain)
  - api_url:         API base URL`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]

			cfg, err := config.Load()
			if err != nil {
				return err
			}

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
		},
	})

	// config path
	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Show config file path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.ConfigFilePath())
		},
	})

	return cmd
}

// -----------------------------------------------------------------------------
// Shows Commands
// -----------------------------------------------------------------------------

func newShowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shows",
		Short: "Manage your podcast shows",
	}

	// shows list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all your shows",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			limit, _ := cmd.Flags().GetInt("limit")
			result, err := client.GetMyShows(api.PaginationParams{Limit: limit})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)

			if len(result.Items) == 0 {
				formatter.PrintMessage("No shows found.")
				return nil
			}

			formatter.PrintShows(result.Items)

			if result.HasMore {
				formatter.PrintMessage("\n(more shows available, use --limit to see more)")
			}

			return nil
		},
	}
	listCmd.Flags().IntP("limit", "l", 20, "Maximum number of shows to list")
	cmd.AddCommand(listCmd)

	// shows get
	getCmd := &cobra.Command{
		Use:   "get <show-id>",
		Short: "Get details of a specific show",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			show, err := client.GetShow(showID)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintShow(show)
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	// shows delete
	deleteCmd := &cobra.Command{
		Use:   "delete <show-id>",
		Short: "Delete a show",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("Are you sure you want to delete show %d? [y/N]: ", showID)
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			if err := client.DeleteShow(showID); err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintSuccess(fmt.Sprintf("Show %d deleted", showID))
			return nil
		},
	}
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	cmd.AddCommand(deleteCmd)

	return cmd
}

// -----------------------------------------------------------------------------
// Episodes Commands
// -----------------------------------------------------------------------------

func newEpisodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "episodes",
		Short: "Manage podcast episodes",
	}

	// episodes list
	listCmd := &cobra.Command{
		Use:   "list [show-id]",
		Short: "List episodes of a show",
		Long: `List episodes of a show.

If no show-id is provided, uses the default_show_id from your config.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			// Determine show ID
			var showID int
			if len(args) > 0 {
				if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
					return fmt.Errorf("invalid show ID: %s", args[0])
				}
			} else {
				cfg, _ := config.Load()
				if cfg.DefaultShowID == 0 {
					return fmt.Errorf("no show ID provided and no default_show_id configured\n" +
						"Either provide a show ID or run: spreaker config set default_show_id <id>")
				}
				showID = cfg.DefaultShowID
			}

			limit, _ := cmd.Flags().GetInt("limit")
			result, err := client.GetShowEpisodes(showID, api.PaginationParams{Limit: limit})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)

			if len(result.Items) == 0 {
				formatter.PrintMessage("No episodes found.")
				return nil
			}

			formatter.PrintEpisodes(result.Items)

			if result.HasMore {
				formatter.PrintMessage("\n(more episodes available, use --limit to see more)")
			}

			return nil
		},
	}
	listCmd.Flags().IntP("limit", "l", 20, "Maximum number of episodes to list")
	cmd.AddCommand(listCmd)

	// episodes get
	getCmd := &cobra.Command{
		Use:   "get <episode-id>",
		Short: "Get details of a specific episode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var episodeID int
			if _, err := fmt.Sscanf(args[0], "%d", &episodeID); err != nil {
				return fmt.Errorf("invalid episode ID: %s", args[0])
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			episode, err := client.GetEpisode(episodeID)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintEpisode(episode)
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	// episodes upload
	uploadCmd := &cobra.Command{
		Use:   "upload <show-id> <audio-file>",
		Short: "Upload a new episode",
		Long: `Upload a new episode to a show.

Example:
  spreaker episodes upload 12345 ./my-episode.mp3 --title "Episode 1"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}
			audioFile := args[1]

			// Check file exists
			if _, err := os.Stat(audioFile); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", audioFile)
			}

			title, _ := cmd.Flags().GetString("title")
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			description, _ := cmd.Flags().GetString("description")
			tags, _ := cmd.Flags().GetStringSlice("tags")
			explicit, _ := cmd.Flags().GetBool("explicit")
			downloadable, _ := cmd.Flags().GetBool("downloadable")

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintMessage(fmt.Sprintf("Uploading %s...", audioFile))

			episode, err := client.UploadEpisode(showID, api.UploadEpisodeParams{
				Title:           title,
				MediaFile:       audioFile,
				Description:     description,
				Tags:            tags,
				Explicit:        explicit,
				DownloadEnabled: downloadable,
			})
			if err != nil {
				return err
			}

			formatter.PrintSuccess("Episode uploaded!")
			formatter.PrintEpisode(episode)
			return nil
		},
	}
	uploadCmd.Flags().StringP("title", "t", "", "Episode title (required)")
	uploadCmd.Flags().StringP("description", "d", "", "Episode description")
	uploadCmd.Flags().StringSlice("tags", nil, "Tags (comma-separated)")
	uploadCmd.Flags().Bool("explicit", false, "Mark as explicit content")
	uploadCmd.Flags().Bool("downloadable", true, "Allow downloads")
	uploadCmd.MarkFlagRequired("title")
	cmd.AddCommand(uploadCmd)

	// episodes delete
	deleteCmd := &cobra.Command{
		Use:   "delete <episode-id>",
		Short: "Delete an episode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var episodeID int
			if _, err := fmt.Sscanf(args[0], "%d", &episodeID); err != nil {
				return fmt.Errorf("invalid episode ID: %s", args[0])
			}

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Printf("Are you sure you want to delete episode %d? [y/N]: ", episodeID)
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			if err := client.DeleteEpisode(episodeID); err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintSuccess(fmt.Sprintf("Episode %d deleted", episodeID))
			return nil
		},
	}
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	cmd.AddCommand(deleteCmd)

	return cmd
}