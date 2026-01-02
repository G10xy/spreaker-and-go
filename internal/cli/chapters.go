/*
chapters.go - Episode chapter management commands

Chapters are bookmarks within an episode that help listeners
fast forward to specific points, especially useful for long audio files.
*/
package cli

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newChaptersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chapters",
		Aliases: []string{"chapter"},
		Short:   "Manage episode chapters",
		Long: `Manage chapters for episodes. Chapters are bookmarks within your episode 
that help listeners fast forward to specific points, especially useful 
for long audio files.

Examples:
  spreaker chapters list 12345
  spreaker chapters add 12345 --starts-at 30000 --title "Introduction"
  spreaker chapters update 12345 67890 --title "New Title"
  spreaker chapters delete 12345 67890
  spreaker chapters delete-all 12345`,
	}

	cmd.AddCommand(
		newChaptersListCmd(),
		newChaptersAddCmd(),
		newChaptersUpdateCmd(),
		newChaptersDeleteCmd(),
		newChaptersDeleteAllCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// chapters list
// -----------------------------------------------------------------------------

func newChaptersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <episode-id>",
		Short: "List all chapters for an episode",
		Long: `List all chapters for an episode, sorted chronologically by start time.

Examples:
  spreaker chapters list 12345
  spreaker chapters list 12345 --limit 50
  spreaker chapter list 12345 --output json`,
		Args: cobra.ExactArgs(1),
		RunE: runChaptersList,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of chapters")

	return cmd
}

func runChaptersList(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetEpisodeChapters(episodeID, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(result.Items) == 0 {
		formatter.PrintMessage("No chapters found for this episode.")
		return nil
	}

	formatter.PrintChapters(result.Items)

	if result.HasMore {
		formatter.PrintMessage("\n(more chapters available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// chapters add
// -----------------------------------------------------------------------------

func newChaptersAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <episode-id>",
		Short: "Add a new chapter to an episode",
		Long: `Add a new chapter to an episode.

Required flags:
  --starts-at: Position in milliseconds where chapter begins
  --title: Chapter title (max 120 characters)

Optional flags:
  --url: External URL for extra information
  --image: Path to image file (400x400+, max 5MB, JPG/PNG)
  --crop: Crop coordinates "x1,y1,x2,y2"

Examples:
  # Add chapter at 30 seconds
  spreaker chapters add 12345 --starts-at 30000 --title "Introduction"

  # Add chapter with URL and image
  spreaker chapters add 12345 --starts-at 120000 --title "Main Topic" \
    --url "https://example.com" --image chapter.jpg`,
		Args: cobra.ExactArgs(1),
		RunE: runChaptersAdd,
	}

	cmd.Flags().Int("starts-at", -1, "Position in milliseconds (required)")
	cmd.Flags().String("title", "", "Chapter title (required)")
	cmd.Flags().String("url", "", "External URL")
	cmd.Flags().String("image", "", "Image file path")
	cmd.Flags().String("crop", "", "Crop coordinates: x1,y1,x2,y2")

	return cmd
}

func runChaptersAdd(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	startsAt, _ := cmd.Flags().GetInt("starts-at")
	title, _ := cmd.Flags().GetString("title")
	url, _ := cmd.Flags().GetString("url")
	image, _ := cmd.Flags().GetString("image")
	crop, _ := cmd.Flags().GetString("crop")

	// Validate required fields
	if startsAt < 0 {
		return fmt.Errorf("--starts-at is required")
	}
	if title == "" {
		return fmt.Errorf("--title is required")
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	params := api.ChapterParams{
		StartsAt:    &startsAt,
		Title:       title,
		ExternalURL: url,
		ImageFile:   image,
		ImageCrop:   crop,
	}

	chapter, err := client.AddChapter(episodeID, params)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage(fmt.Sprintf("Chapter added successfully (ID: %d).", chapter.ChapterID))
	return nil
}

// -----------------------------------------------------------------------------
// chapters update
// -----------------------------------------------------------------------------

func newChaptersUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <episode-id> <chapter-id>",
		Short: "Update an existing chapter",
		Long: `Update an existing chapter. Specify only the fields you want to update.

Examples:
  # Update title only
  spreaker chapters update 12345 67890 --title "New Title"

  # Remove image
  spreaker chapters update 12345 67890 --image remove`,
		Args: cobra.ExactArgs(2),
		RunE: runChaptersUpdate,
	}

	cmd.Flags().Int("starts-at", 0, "Position in milliseconds")
	cmd.Flags().String("title", "", "Chapter title")
	cmd.Flags().String("url", "", "External URL")
	cmd.Flags().String("image", "", "Image file path (or 'remove' to delete)")
	cmd.Flags().String("crop", "", "Crop coordinates: x1,y1,x2,y2")

	return cmd
}

func runChaptersUpdate(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	chapterID, err := parseChapterID(args[1])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	// Build params only with flags that were explicitly set
	// cmd.Flags().Changed() tells if the user provided the flag
	params := api.ChapterParams{}

	if cmd.Flags().Changed("starts-at") {
		startsAt, _ := cmd.Flags().GetInt("starts-at")
		params.StartsAt = &startsAt
	}
	if cmd.Flags().Changed("title") {
		params.Title, _ = cmd.Flags().GetString("title")
	}
	if cmd.Flags().Changed("url") {
		params.ExternalURL, _ = cmd.Flags().GetString("url")
	}
	if cmd.Flags().Changed("image") {
		params.ImageFile, _ = cmd.Flags().GetString("image")
	}
	if cmd.Flags().Changed("crop") {
		params.ImageCrop, _ = cmd.Flags().GetString("crop")
	}

	_, err = client.UpdateChapter(episodeID, chapterID, params)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage("Chapter updated successfully.")
	return nil
}

// -----------------------------------------------------------------------------
// chapters delete
// -----------------------------------------------------------------------------

func newChaptersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <episode-id> <chapter-id>",
		Short: "Delete a single chapter",
		Long: `Delete a single chapter from an episode.

Examples:
  spreaker chapters delete 12345 67890`,
		Args: cobra.ExactArgs(2),
		RunE: runChaptersDelete,
	}
}

func runChaptersDelete(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	chapterID, err := parseChapterID(args[1])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	if err := client.DeleteChapter(episodeID, chapterID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage("Chapter deleted successfully.")
	return nil
}

// -----------------------------------------------------------------------------
// chapters delete-all
// -----------------------------------------------------------------------------

func newChaptersDeleteAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-all <episode-id>",
		Short: "Delete all chapters from an episode",
		Long: `Delete all chapters from an episode.

Examples:
  spreaker chapters delete-all 12345`,
		Args: cobra.ExactArgs(1),
		RunE: runChaptersDeleteAll,
	}
}

func runChaptersDeleteAll(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	if err := client.DeleteAllChapters(episodeID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage("All chapters deleted successfully.")
	return nil
}
