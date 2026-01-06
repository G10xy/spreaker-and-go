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
	"path/filepath"
	"net/http"
	"io"
	"strings"

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
  spreaker episodes upload 12345 ./ep.mp3   # Upload a new episode,
  spreaker episodes download 67890          # Download an episode`,
	}

	cmd.AddCommand(
		newEpisodesListCmd(),
		newEpisodesGetCmd(),
		newEpisodesUploadCmd(),
		newEpisodesDeleteCmd(),
		newEpisodesDownloadCmd(),
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


// -----------------------------------------------------------------------------
// episodes download
// -----------------------------------------------------------------------------

func newEpisodesDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <episode-id>",
		Short: "Download an episode's audio file",
		Long: `Download an episode's audio file to your local machine.

By default, the file is saved with the episode title as filename.
Use --output to specify a custom filename or path.
Use --url-only to just print the download URL without downloading.

Examples:
  spreaker episodes download 67890

  spreaker episodes download 67890 --output ~/podcasts/episode.mp3

  # Just get the download URL
  spreaker episodes download 67890 --url-only`,
		Args: cobra.ExactArgs(1),
		RunE: runEpisodesDownload,
	}

	cmd.Flags().StringP("output", "O", "", "Output file path (default: episode title)")
	cmd.Flags().BoolP("url-only", "u", false, "Only print the download URL, don't download")

	return cmd
}

func runEpisodesDownload(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	downloadURL, err := client.GetEpisodeDownloadURL(episodeID)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	// If --url-only flag is set, just print the URL and exit
	urlOnly, _ := cmd.Flags().GetBool("url-only")
	if urlOnly {
		fmt.Println(downloadURL)
		return nil
	}

	// Determine output filename
	outputPath, _ := cmd.Flags().GetString("output")
	if outputPath == "" {
		episode, err := client.GetEpisode(episodeID)
		if err != nil {
			outputPath = fmt.Sprintf("episode_%d.mp3", episodeID)
		} else {
			outputPath = sanitizeFilename(episode.Title) + ".mp3"
		}
	}

	// Ensure directory exists if path contains directories
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	formatter.PrintMessage(fmt.Sprintf("Downloading episode %d to %s...", episodeID, outputPath))

	if err := downloadFile(downloadURL, outputPath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	formatter.PrintSuccess(fmt.Sprintf("Downloaded to %s", outputPath))
	return nil
}

// downloadFile downloads a file from the given URL to the specified path.
func downloadFile(url, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()


	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("  Size: %.2f MB\n", float64(written)/(1024*1024))

	return nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
		"\n", " ",
		"\r", "",
		"\t", " ",
	)

	sanitized := replacer.Replace(name)

	sanitized = strings.TrimSpace(sanitized)
	sanitized = strings.Trim(sanitized, ".")

	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}

	if sanitized == "" {
		sanitized = "episode"
	}

	return sanitized
}
