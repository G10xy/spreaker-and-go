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
		newEpisodesUpdateCmd(),
		newEpisodesDraftCmd(),
		newEpisodesDeleteCmd(),
		newEpisodesDownloadCmd(),
		newEpisodesDownloadAllCmd(),
		newEpisodesLikesCmd(),
		newEpisodesLikeCmd(),
		newEpisodesUnlikeCmd(),
		newEpisodesBookmarkCmd(),
		newEpisodesUnbookmarkCmd(),
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

// -----------------------------------------------------------------------------
// episodes download-all
// -----------------------------------------------------------------------------

func newEpisodesDownloadAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-all <show-id>",
		Short: "Download all episodes of a show",
		Long: `Download all episodes of a show to your local machine.

By default, episodes are saved to a directory named after the show title.
Files that already exist are skipped (resume capability).

Examples:
  spreaker episodes download-all 12345

  spreaker episodes download-all 12345 --output-dir ~/podcasts/myshow

  spreaker episodes download-all 12345 --limit 10

  # Force re-download of existing files
  spreaker episodes download-all 12345 --no-skip-existing`,
		Args: cobra.ExactArgs(1),
		RunE: runEpisodesDownloadAll,
	}

	cmd.Flags().StringP("output-dir", "O", "", "Output directory (default: ./<show-title>/)")
	cmd.Flags().Bool("skip-existing", true, "Skip episodes that already exist locally")
	cmd.Flags().IntP("limit", "l", 0, "Maximum number of episodes to download (0 = all)")

	return cmd
}

func runEpisodesDownloadAll(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	// Get show details for directory name
	show, err := client.GetShow(showID)
	if err != nil {
		return fmt.Errorf("failed to get show details: %w", err)
	}

	// Determine output directory
	outputDir, _ := cmd.Flags().GetString("output-dir")
	if outputDir == "" {
		outputDir = sanitizeFilename(show.Title)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	skipExisting, _ := cmd.Flags().GetBool("skip-existing")
	limit, _ := cmd.Flags().GetInt("limit")

	formatter.PrintMessage(fmt.Sprintf("Fetching episodes for show: %s", show.Title))

	// Fetch all episodes using pagination
	var allEpisodes []struct {
		ID    int
		Title string
	}

	pageLimit := 100
	if limit > 0 && limit < pageLimit {
		pageLimit = limit
	}

	result, err := client.GetShowEpisodes(showID, api.PaginationParams{Limit: pageLimit})
	if err != nil {
		return fmt.Errorf("failed to fetch episodes: %w", err)
	}

	for _, ep := range result.Items {
		allEpisodes = append(allEpisodes, struct {
			ID    int
			Title string
		}{ID: ep.EpisodeID, Title: ep.Title})
		if limit > 0 && len(allEpisodes) >= limit {
			break
		}
	}

	// Continue fetching if there are more episodes and we haven't hit the limit
	for result.HasMore && (limit == 0 || len(allEpisodes) < limit) {
		nextLimit := pageLimit
		if limit > 0 && limit-len(allEpisodes) < nextLimit {
			nextLimit = limit - len(allEpisodes)
		}

		result, err = client.GetShowEpisodes(showID, api.PaginationParams{
			Limit:  nextLimit,
			LastID: result.Items[len(result.Items)-1].EpisodeID,
		})
		if err != nil {
			return fmt.Errorf("failed to fetch episodes: %w", err)
		}

		for _, ep := range result.Items {
			allEpisodes = append(allEpisodes, struct {
				ID    int
				Title string
			}{ID: ep.EpisodeID, Title: ep.Title})
			if limit > 0 && len(allEpisodes) >= limit {
				break
			}
		}
	}

	if len(allEpisodes) == 0 {
		formatter.PrintMessage("No episodes found.")
		return nil
	}

	formatter.PrintMessage(fmt.Sprintf("Found %d episodes to download", len(allEpisodes)))

	// Download statistics
	var downloaded, skipped, failed int

	for i, ep := range allEpisodes {
		filename := sanitizeFilename(ep.Title) + ".mp3"
		filePath := filepath.Join(outputDir, filename)

		
		if skipExisting {
			if _, err := os.Stat(filePath); err == nil {
				formatter.PrintMessage(fmt.Sprintf("[%d/%d] Skipping (exists): %s", i+1, len(allEpisodes), filename))
				skipped++
				continue
			}
		}

		formatter.PrintMessage(fmt.Sprintf("[%d/%d] Downloading: %s", i+1, len(allEpisodes), filename))

		
		downloadURL, err := client.GetEpisodeDownloadURL(ep.ID)
		if err != nil {
			formatter.PrintError(fmt.Sprintf("  Failed to get download URL: %v", err))
			failed++
			continue
		}

		
		if err := downloadFile(downloadURL, filePath); err != nil {
			formatter.PrintError(fmt.Sprintf("  Download failed: %v", err))
			failed++
			continue
		}

		downloaded++
	}

	
	formatter.PrintMessage("")
	formatter.PrintMessage("Download complete!")
	formatter.PrintMessage(fmt.Sprintf("  Downloaded: %d", downloaded))
	if skipped > 0 {
		formatter.PrintMessage(fmt.Sprintf("  Skipped:    %d", skipped))
	}
	if failed > 0 {
		formatter.PrintMessage(fmt.Sprintf("  Failed:     %d", failed))
	}
	formatter.PrintMessage(fmt.Sprintf("  Location:   %s", outputDir))

	return nil
}

// -----------------------------------------------------------------------------
// episodes update
// -----------------------------------------------------------------------------

func newEpisodesUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <episode-id>",
		Short: "Update an episode",
		Long: `Update an existing episode.

Examples:
  spreaker episodes update 67890 --title "New Title"
  spreaker episodes update 67890 --description "New description"
  spreaker episodes update 67890 --hidden`,
		Args: cobra.ExactArgs(1),
		RunE: runEpisodesUpdate,
	}

	cmd.Flags().String("title", "", "Episode title")
	cmd.Flags().String("description", "", "Episode description")
	cmd.Flags().StringSlice("tags", nil, "Tags (comma-separated)")
	cmd.Flags().Bool("explicit", false, "Mark as explicit content")
	cmd.Flags().Bool("downloadable", false, "Allow downloads")
	cmd.Flags().Bool("hidden", false, "Hide the episode")

	return cmd
}

func runEpisodesUpdate(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	params := api.UpdateEpisodeParams{}

	if cmd.Flags().Changed("title") {
		val, _ := cmd.Flags().GetString("title")
		params.Title = &val
	}
	if cmd.Flags().Changed("description") {
		val, _ := cmd.Flags().GetString("description")
		params.Description = &val
	}
	if cmd.Flags().Changed("tags") {
		val, _ := cmd.Flags().GetStringSlice("tags")
		params.Tags = &val
	}
	if cmd.Flags().Changed("explicit") {
		val, _ := cmd.Flags().GetBool("explicit")
		params.Explicit = &val
	}
	if cmd.Flags().Changed("downloadable") {
		val, _ := cmd.Flags().GetBool("downloadable")
		params.DownloadEnabled = &val
	}
	if cmd.Flags().Changed("hidden") {
		val, _ := cmd.Flags().GetBool("hidden")
		params.Hidden = &val
	}

	episode, err := client.UpdateEpisode(episodeID, params)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess("Episode updated")
	formatter.PrintEpisode(episode)
	return nil
}

// -----------------------------------------------------------------------------
// episodes draft
// -----------------------------------------------------------------------------

func newEpisodesDraftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draft <show-id>",
		Short: "Create a draft episode",
		Long: `Create a draft episode without an audio file.

The audio file can be uploaded later.

Examples:
  spreaker episodes draft 12345 --title "Upcoming Episode"
  spreaker episodes draft 12345 --title "Draft" --description "Work in progress"`,
		Args: cobra.ExactArgs(1),
		RunE: runEpisodesDraft,
	}

	cmd.Flags().String("title", "", "Episode title (required)")
	cmd.Flags().String("description", "", "Episode description")
	cmd.Flags().StringSlice("tags", nil, "Tags (comma-separated)")
	cmd.Flags().Bool("explicit", false, "Mark as explicit content")
	cmd.Flags().Bool("downloadable", true, "Allow downloads")
	cmd.Flags().Bool("hidden", false, "Hide the episode")

	cmd.MarkFlagRequired("title")

	return cmd
}

func runEpisodesDraft(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	explicit, _ := cmd.Flags().GetBool("explicit")
	downloadable, _ := cmd.Flags().GetBool("downloadable")
	hidden, _ := cmd.Flags().GetBool("hidden")

	params := api.CreateDraftEpisodeParams{
		Title:           title,
		ShowID:          showID,
		Description:     description,
		Tags:            tags,
		Explicit:        explicit,
		DownloadEnabled: downloadable,
		Hidden:          hidden,
	}

	episode, err := client.CreateDraftEpisode(params)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Draft episode created with ID %d", episode.EpisodeID))
	formatter.PrintEpisode(episode)
	return nil
}

// -----------------------------------------------------------------------------
// episodes likes
// -----------------------------------------------------------------------------

func newEpisodesLikesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "likes",
		Short: "List your liked episodes",
		RunE:  runEpisodesLikes,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of episodes to list")

	return cmd
}

func runEpisodesLikes(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetLikedEpisodes(me.UserID, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(result.Items) == 0 {
		formatter.PrintMessage("No liked episodes.")
		return nil
	}

	formatter.PrintEpisodes(result.Items)

	if result.HasMore {
		formatter.PrintMessage("\n(more episodes available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// episodes like
// -----------------------------------------------------------------------------

func newEpisodesLikeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "like <episode-id>",
		Short: "Like an episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runEpisodesLike,
	}
}

func runEpisodesLike(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	if err := client.LikeEpisode(me.UserID, episodeID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Liked episode %d", episodeID))
	return nil
}

// -----------------------------------------------------------------------------
// episodes unlike
// -----------------------------------------------------------------------------

func newEpisodesUnlikeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlike <episode-id>",
		Short: "Unlike an episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runEpisodesUnlike,
	}
}

func runEpisodesUnlike(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	if err := client.UnlikeEpisode(me.UserID, episodeID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Unliked episode %d", episodeID))
	return nil
}

// -----------------------------------------------------------------------------
// episodes bookmark
// -----------------------------------------------------------------------------

func newEpisodesBookmarkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bookmark <episode-id>",
		Short: "Bookmark an episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runEpisodesBookmark,
	}
}

func runEpisodesBookmark(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	if err := client.BookmarkEpisode(me.UserID, episodeID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Bookmarked episode %d", episodeID))
	return nil
}

// -----------------------------------------------------------------------------
// episodes unbookmark
// -----------------------------------------------------------------------------

func newEpisodesUnbookmarkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unbookmark <episode-id>",
		Short: "Remove an episode from bookmarks",
		Args:  cobra.ExactArgs(1),
		RunE:  runEpisodesUnbookmark,
	}
}

func runEpisodesUnbookmark(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	if err := client.UnbookmarkEpisode(me.UserID, episodeID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Removed episode %d from bookmarks", episodeID))
	return nil
}
