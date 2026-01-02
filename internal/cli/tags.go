/*
tags.go - Tag-based discovery commands

Commands for discovering episodes by tag/hashtag.
*/
package cli

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Discover episodes by tag",
		Long: `Discover episodes by searching for specific tags/hashtags.

Examples:
  spreaker tags episodes "breaking news"
  spreaker tags episodes tech
  spreaker tags episodes "machine learning" --limit 50`,
	}

	cmd.AddCommand(newTagsEpisodesCmd())

	return cmd
}

func newTagsEpisodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "episodes <tag-name>",
		Short: "Get latest episodes with a specific tag",
		Long: `Get the latest episodes that have been tagged with a specific hashtag.

The tag name can contain spaces and special characters.

Examples:
  spreaker tags episodes "breaking news"
  spreaker tags episodes tech
  spreaker tags episodes "machine learning" --limit 50`,
		Args: cobra.ExactArgs(1),
		RunE: runTagsEpisodes,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of episodes")

	return cmd
}

func runTagsEpisodes(cmd *cobra.Command, args []string) error {
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
}
