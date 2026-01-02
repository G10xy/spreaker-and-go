/*
episodes.go - Episode management commands

This file contains all commands related to podcast episodes:
  - list: List episodes of a show
  - get: Get details of a specific episode
  - upload: Upload a new episode
  - delete: Delete an episode
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/internal/api"
	"github.com/G10xy/spreaker-and-go/internal/config"
)

func newEpisodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "episodes",
		Short: "Manage podcast episodes",
		Long: `Manage episodes for your podcast shows.

Examples:
  spreaker episodes list                    # List episodes (uses default show)
  spreaker episodes list 12345              # List episodes of show 12345
  spreaker episodes get 67890               # Get episode details
  spreaker episodes upload 12345 ./ep.mp3   # Upload a new episode`,
	}

	cmd.AddCommand(
		newEpisodesListCmd(),
		newEpisodesGetCmd(),
		newEpisodesUploadCmd(),
		newEpisodesDeleteCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// episodes list
// -----------------------------------------------------------------------------

func newEpisodesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [show-id]",
		Short: "List episodes of a show",
		Long: `List episodes of a show.

If no show-id is provided, uses the default_show_id from your config.
Set a default with: spreaker config set default_show_id <id>`,
		RunE: runEpisodesList,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of episodes to list")

	return cmd
}

func runEpisodesList(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	// Determine show ID: from argument or default config
	var showID int
	if len(args) > 0 {
		showID, err = parseShowID(args[0])
		if err != nil {
			return err
		}
	} else {
		// Try to use default show ID from config
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
}

// -----------------------------------------------------------------------------
// episodes get
// -----------------------------------------------------------------------------

func newEpisodesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <episode-id>",
		Short: "Get details of a specific episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runEpisodesGet,
	}
}

func runEpisodesGet(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
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
}

// -----------------------------------------------------------------------------
// episodes upload
// -----------------------------------------------------------------------------

func newEpisodesUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <show-id> <audio-file>",
		Short: "Upload a new episode",
		Long: `Upload a new episode to a show.

The audio file should be in a supported format (MP3, WAV, etc.).

Examples:
  spreaker episodes upload 12345 ./episode.mp3 --title "Episode 1"
  
  spreaker episodes upload 12345 ./episode.mp3 \
    --title "Episode 42: The Answer" \
    --description "In this episode we discuss everything." \
    --tags "science,philosophy" \
    --explicit`,
		Args: cobra.ExactArgs(2),
		RunE: runEpisodesUpload,
	}

	// Required flag
	cmd.Flags().StringP("title", "t", "", "Episode title (required)")
	cmd.MarkFlagRequired("title")

	// Optional flags
	cmd.Flags().StringP("description", "d", "", "Episode description")
	cmd.Flags().StringSlice("tags", nil, "Tags (comma-separated)")
	cmd.Flags().Bool("explicit", false, "Mark as explicit content")
	cmd.Flags().Bool("downloadable", true, "Allow downloads")

	return cmd
}

func runEpisodesUpload(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}
	audioFile := args[1]

	// Verify file exists before making API call
	// This gives a better error message than a failed upload
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", audioFile)
	}

	// Get all flag values
	title, _ := cmd.Flags().GetString("title")
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
}

// -----------------------------------------------------------------------------
// episodes delete
// -----------------------------------------------------------------------------

func newEpisodesDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <episode-id>",
		Short: "Delete an episode",
		Long: `Delete an episode permanently.

WARNING: This action cannot be undone.`,
		Args: cobra.ExactArgs(1),
		RunE: runEpisodesDelete,
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runEpisodesDelete(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	force, _ := cmd.Flags().GetBool("force")
	if !force {
		prompt := fmt.Sprintf("Are you sure you want to delete episode %d? [y/N]: ", episodeID)
		if !confirmAction(prompt) {
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
}
