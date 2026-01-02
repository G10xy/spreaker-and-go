/*
explore.go - Podcast discovery commands

Commands for discovering podcasts by category.
*/
package cli

import (
	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newExploreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explore",
		Short: "Discover podcasts by category",
		Long: `Discover podcasts by browsing categories.

Use 'spreaker misc categories' to see available category IDs.

Examples:
  spreaker explore category 14
  spreaker explore category 14 --limit 50`,
	}

	cmd.AddCommand(newExploreCategoryCmd())

	return cmd
}

func newExploreCategoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category <category-id>",
		Short: "List shows in a category",
		Long: `List shows in a specific category, ranked by popularity and quality.

Use 'spreaker misc categories' to see available category IDs.`,
		Args: cobra.ExactArgs(1),
		RunE: runExploreCategory,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of shows")

	return cmd
}

func runExploreCategory(cmd *cobra.Command, args []string) error {
	categoryID, err := parseShowID(args[0]) // Reusing parseShowID for category ID
	if err != nil {
		return err
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
}
