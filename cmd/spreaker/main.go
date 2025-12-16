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
		newStatsCmd()
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


// -----------------------------------------------------------------------------
// Stats Commands
// -----------------------------------------------------------------------------

func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "View statistics for users, shows, and episodes",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "me",
		Short: "Show your overall statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			stats, err := client.GetMyStatistics()
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintUserStatistics(stats)
			return nil
		},
	})

	// stats show <show-id> - Get show's overall statistics
	showCmd := &cobra.Command{
		Use:   "show <show-id>",
		Short: "Show statistics for a specific show",
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

			stats, err := client.GetShowStatistics(showID)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintShowStatistics(stats)
			return nil
		},
	}
	cmd.AddCommand(showCmd)

	// stats episode <episode-id> - Get episode's overall statistics
	episodeCmd := &cobra.Command{
		Use:   "episode <episode-id>",
		Short: "Show statistics for a specific episode",
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

			stats, err := client.GetEpisodeStatistics(episodeID)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintEpisodeStatistics(stats)
			return nil
		},
	}
	cmd.AddCommand(episodeCmd)

	// stats plays <show-id> - Get play statistics for a show
	playsCmd := &cobra.Command{
		Use:   "plays <show-id>",
		Short: "Show play statistics for a show over time",
		Long: `Show play statistics for a show over a date range.

		Example:
		spreaker stats plays 12345 --from 2024-01-01 --to 2024-01-31 --group day`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}

			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			group, _ := cmd.Flags().GetString("group")

			if from == "" || to == "" {
				return fmt.Errorf("--from and --to are required")
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			stats, err := client.GetShowPlayStatistics(showID, api.StatisticsParams{
				From:  from,
				To:    to,
				Group: group,
			})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintPlayStatistics(stats)
			return nil
		},
	}
	playsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	playsCmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	playsCmd.Flags().String("group", "day", "Group by: day, week, or month")
	playsCmd.MarkFlagRequired("from")
	playsCmd.MarkFlagRequired("to")
	cmd.AddCommand(playsCmd)

	// stats devices <show-id> - Get device statistics for a show
	devicesCmd := &cobra.Command{
		Use:   "devices <show-id>",
		Short: "Show device breakdown for a show",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}

			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")

			if from == "" || to == "" {
				return fmt.Errorf("--from and --to are required")
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			stats, err := client.GetShowDevicesStatistics(showID, api.StatisticsParams{
				From: from,
				To:   to,
			})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintDeviceStatistics(stats)
			return nil
		},
	}
	devicesCmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	devicesCmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	devicesCmd.MarkFlagRequired("from")
	devicesCmd.MarkFlagRequired("to")
	cmd.AddCommand(devicesCmd)

	// stats geo <show-id> - Get geographic statistics for a show
	geoCmd := &cobra.Command{
		Use:   "geo <show-id>",
		Short: "Show geographic breakdown for a show",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}

			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")

			if from == "" || to == "" {
				return fmt.Errorf("--from and --to are required")
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			stats, err := client.GetShowGeographicStatistics(showID, api.StatisticsParams{
				From: from,
				To:   to,
			})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintGeographicStatistics(stats)
			return nil
		},
	}
	geoCmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	geoCmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	geoCmd.MarkFlagRequired("from")
	geoCmd.MarkFlagRequired("to")
	cmd.AddCommand(geoCmd)

	// stats sources <show-id> - Get sources statistics for a show
	sourcesCmd := &cobra.Command{
		Use:   "sources <show-id>",
		Short: "Show play/download sources for a show",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}

			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			group, _ := cmd.Flags().GetString("group")

			if from == "" || to == "" {
				return fmt.Errorf("--from and --to are required")
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			stats, err := client.GetShowSourcesStatistics(showID, api.StatisticsParams{
				From:  from,
				To:    to,
				Group: group,
			})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintSourcesStatistics(stats)
			return nil
		},
	}
	sourcesCmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	sourcesCmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	sourcesCmd.Flags().String("group", "day", "Group by: day, week, or month")
	sourcesCmd.MarkFlagRequired("from")
	sourcesCmd.MarkFlagRequired("to")
	cmd.AddCommand(sourcesCmd)

	// stats listeners <show-id> - Get listeners statistics for a show
	listenersCmd := &cobra.Command{
		Use:   "listeners <show-id>",
		Short: "Show unique listeners for a show over time",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var showID int
			if _, err := fmt.Sscanf(args[0], "%d", &showID); err != nil {
				return fmt.Errorf("invalid show ID: %s", args[0])
			}

			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			group, _ := cmd.Flags().GetString("group")

			if from == "" || to == "" {
				return fmt.Errorf("--from and --to are required")
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			stats, err := client.GetShowListenersStatistics(showID, api.StatisticsParams{
				From:  from,
				To:    to,
				Group: group,
			})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintListenersStatistics(stats)
			return nil
		},
	}
	listenersCmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	listenersCmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	listenersCmd.Flags().String("group", "day", "Group by: day, week, or month")
	listenersCmd.MarkFlagRequired("from")
	listenersCmd.MarkFlagRequired("to")
	cmd.AddCommand(listenersCmd)

	return cmd
}


// -----------------------------------------------------------------------------
// Explore Commands
// -----------------------------------------------------------------------------

func newExploreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explore",
		Short: "Discover podcasts by category",
	}

	categoryCmd := &cobra.Command{
		Use:   "category <category-id>",
		Short: "List shows in a category",
		Long: `List shows in a specific category, ranked by popularity and quality.

Use 'spreaker misc categories' to see available category IDs.
Examples:
  spreaker explore category 14
  spreaker explore category 14 --limit 50`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var categoryID int
			if _, err := fmt.Sscanf(args[0], "%d", &categoryID); err != nil {
				return fmt.Errorf("invalid category ID: %s", args[0])
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			limit, _ := cmd.Flags().GetInt("limit")
			result, err := client.GetCategoryShows(categoryID, api.PaginationParams{Limit: limit})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)

			if len(result.Items) == 0 {
				formatter.PrintMessage("No shows found in this category.")
				return nil
			}

			formatter.PrintExploreShows(result.Items)

			if result.HasMore {
				formatter.PrintMessage("\n(more shows available, use --limit to see more)")
			}

			return nil
		},
	}
	categoryCmd.Flags().IntP("limit", "l", 20, "Maximum number of shows")
	cmd.AddCommand(categoryCmd)

	return cmd
}


// -----------------------------------------------------------------------------
// Miscellaneous Commands 
// -----------------------------------------------------------------------------

func newMiscCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "misc",
		Aliases: []string{"miscellaneous"},
		Short:   "List categories and languages",
	}

	categoriesCmd := &cobra.Command{
		Use:   "categories",
		Short: "List all show categories",
		Long: `List all available show categories.

Examples:
  spreaker misc categories
  spreaker misc categories --locale it_IT`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			locale, _ := cmd.Flags().GetString("locale")
			categories, err := client.GetShowCategories(locale)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintCategories(categories)
			return nil
		},
	}
	categoriesCmd.Flags().String("locale", "", "Locale for category names (e.g., it_IT)")
	cmd.AddCommand(categoriesCmd)

	gpCategoriesCmd := &cobra.Command{
		Use:   "googleplay-categories",
		Short: "List all Google Play podcast categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			categories, err := client.GetGooglePlayCategories()
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintGooglePlayCategories(categories)
			return nil
		},
	}
	cmd.AddCommand(gpCategoriesCmd)

	languagesCmd := &cobra.Command{
		Use:   "languages",
		Short: "List all show languages",
		Long: `List all available languages for shows.

Examples:
  spreaker misc languages
  spreaker misc languages --locale it_IT`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			locale, _ := cmd.Flags().GetString("locale")
			languages, err := client.GetShowLanguagesList(locale)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintLanguages(languages)
			return nil
		},
	}
	languagesCmd.Flags().String("locale", "", "Locale for language names (e.g., it_IT)")
	cmd.AddCommand(languagesCmd)

	return cmd
}

// -----------------------------------------------------------------------------
// Tags Commands 
// -----------------------------------------------------------------------------

func newTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Discover episodes by tag",
	}

	// tags episodes <tag-name> - Get episodes with a specific tag
	episodesCmd := &cobra.Command{
		Use:   "episodes <tag-name>",
		Short: "Get latest episodes with a specific tag",
		Long: `Get the latest episodes that have been tagged with a specific hashtag.
The tag name can contain spaces and special characters.
Examples:
  spreaker tags episodes "breaking news"
  spreaker tags episodes tech
  spreaker tags episodes "machine learning" --limit 50`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tagName := args[0]

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			limit, _ := cmd.Flags().GetInt("limit")
			result, err := client.GetEpisodesByTag(tagName, api.PaginationParams{Limit: limit})
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)

			if len(result.Items) == 0 {
				formatter.PrintMessage(fmt.Sprintf("No episodes found with tag '%s'.", tagName))
				return nil
			}

			formatter.PrintEpisodes(result.Items)

			if result.HasMore {
				formatter.PrintMessage("\n(more episodes available, use --limit to see more)")
			}

			return nil
		},
	}
	episodesCmd.Flags().IntP("limit", "l", 20, "Maximum number of episodes")
	cmd.AddCommand(episodesCmd)

	return cmd
}

// -----------------------------------------------------------------------------
// Episode Cuepoints Commands 
// -----------------------------------------------------------------------------

func newCuepointsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cuepoints",
		Aliases: []string{"cue"},
		Short:   "Manage episode cuepoints for ad injection",
		Long: `Manage cuepoints for episodes. Cuepoints are specific points in time 
within an episode where audio ads can be injected.

Note: Setting cuepoints is not enough to get ads injected. You also need to
enable Ads and Monetization capabilities on your account and show.`,
	}

	listCmd := &cobra.Command{
		Use:   "list <episode-id>",
		Short: "List all cuepoints for an episode",
		Long: `List all cuepoints for an episode, sorted chronologically by timecode.

Examples:
  spreaker cuepoints list 12345
  spreaker cue list 12345 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var episodeID int
			if _, err := fmt.Sscanf(args[0], "%d", &episodeID); err != nil {
				return fmt.Errorf("invalid episode ID: %s", args[0])
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			cuepoints, err := client.GetEpisodeCuepoints(episodeID)
			if err != nil {
				return err
			}

			formatter := getFormatter(cmd)

			if len(cuepoints) == 0 {
				formatter.PrintMessage("No cuepoints found for this episode.")
				return nil
			}

			formatter.PrintCuepoints(cuepoints)
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	setCmd := &cobra.Command{
		Use:   "set <episode-id> <timecode1:max_ads1> [timecode2:max_ads2] ...",
		Short: "Set cuepoints for an episode (replaces all existing)",
		Long: `Set cuepoints for an episode. This replaces all existing cuepoints.
Timecodes are in milliseconds. Format: timecode:max_ads

Examples:
  # Set a single cuepoint at 30 seconds (30000ms) with max 1 ad
  spreaker cuepoints set 12345 30000:1

  # Clear all cuepoints (set empty list)
  spreaker cuepoints set 12345`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var episodeID int
			if _, err := fmt.Sscanf(args[0], "%d", &episodeID); err != nil {
				return fmt.Errorf("invalid episode ID: %s", args[0])
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			// Parse cuepoints from remaining arguments
			cuepoints := make([]models.Cuepoint, 0, len(args)-1)
			for _, arg := range args[1:] {
				var timecode, maxAds int
				if _, err := fmt.Sscanf(arg, "%d:%d", &timecode, &maxAds); err != nil {
					return fmt.Errorf("invalid cuepoint format '%s' (expected timecode:max_ads, e.g., 30000:1)", arg)
				}
				cuepoints = append(cuepoints, models.Cuepoint{
					Timecode:    timecode,
					AdsMaxCount: maxAds,
				})
			}

			if err := client.UpdateEpisodeCuepoints(episodeID, cuepoints); err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			if len(cuepoints) == 0 {
				formatter.PrintMessage("All cuepoints deleted successfully.")
			} else {
				formatter.PrintMessage(fmt.Sprintf("Successfully set %d cuepoint(s).", len(cuepoints)))
			}
			return nil
		},
	}
	cmd.AddCommand(setCmd)

	deleteCmd := &cobra.Command{
		Use:   "delete <episode-id>",
		Short: "Delete all cuepoints for an episode",
		Long: `Delete all cuepoints for an episode.

Examples:
  spreaker cuepoints delete 12345
  spreaker cue delete 12345`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var episodeID int
			if _, err := fmt.Sscanf(args[0], "%d", &episodeID); err != nil {
				return fmt.Errorf("invalid episode ID: %s", args[0])
			}

			client, err := getClient(cmd)
			if err != nil {
				return err
			}

			if err := client.DeleteEpisodeCuepoints(episodeID); err != nil {
				return err
			}

			formatter := getFormatter(cmd)
			formatter.PrintMessage("All cuepoints deleted successfully.")
			return nil
		},
	}
	cmd.AddCommand(deleteCmd)

	return cmd
}