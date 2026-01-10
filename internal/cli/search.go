/*
search.go - Search commands

This file contains all commands for searching shows and episodes.
*/
package cli

import (
	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for shows and episodes",
		Long: `Search for shows and episodes on Spreaker.

Examples:
  spreaker search shows "tech podcast"
  spreaker search episodes "artificial intelligence"
  spreaker search user-shows 12345 "interview"
  spreaker search show-episodes 12345 "bonus"`,
	}

	cmd.AddCommand(
		newSearchShowsCmd(),
		newSearchEpisodesCmd(),
		newSearchUserShowsCmd(),
		newSearchUserEpisodesCmd(),
		newSearchShowEpisodesCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// search shows
// -----------------------------------------------------------------------------

func newSearchShowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shows <query>",
		Short: "Search for shows",
		Args:  cobra.ExactArgs(1),
		RunE:  runSearchShows,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of results")
	cmd.Flags().String("filter", "", "Filter: listenable (default) or editable")

	return cmd
}

func runSearchShows(cmd *cobra.Command, args []string) error {
	query := args[0]

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	filter, _ := cmd.Flags().GetString("filter")

	result, err := client.SearchShows(
		api.SearchParams{Query: query, Filter: filter},
		api.PaginationParams{Limit: limit},
	)
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
		formatter.PrintMessage("\n(more results available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// search episodes
// -----------------------------------------------------------------------------

func newSearchEpisodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "episodes <query>",
		Short: "Search for episodes",
		Args:  cobra.ExactArgs(1),
		RunE:  runSearchEpisodes,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of results")
	cmd.Flags().String("filter", "", "Filter: listenable (default) or editable")

	return cmd
}

func runSearchEpisodes(cmd *cobra.Command, args []string) error {
	query := args[0]

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	filter, _ := cmd.Flags().GetString("filter")

	result, err := client.SearchEpisodes(
		api.SearchParams{Query: query, Filter: filter},
		api.PaginationParams{Limit: limit},
	)
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
		formatter.PrintMessage("\n(more results available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// search user-shows
// -----------------------------------------------------------------------------

func newSearchUserShowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-shows <user-id> <query>",
		Short: "Search for shows by a specific user",
		Args:  cobra.ExactArgs(2),
		RunE:  runSearchUserShows,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of results")
	cmd.Flags().String("filter", "", "Filter: listenable (default) or editable")

	return cmd
}

func runSearchUserShows(cmd *cobra.Command, args []string) error {
	userID, err := parseUserID(args[0])
	if err != nil {
		return err
	}
	query := args[1]

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	filter, _ := cmd.Flags().GetString("filter")

	result, err := client.SearchUserShows(
		userID,
		api.SearchParams{Query: query, Filter: filter},
		api.PaginationParams{Limit: limit},
	)
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
		formatter.PrintMessage("\n(more results available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// search user-episodes
// -----------------------------------------------------------------------------

func newSearchUserEpisodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-episodes <user-id> <query>",
		Short: "Search for episodes by a specific user",
		Args:  cobra.ExactArgs(2),
		RunE:  runSearchUserEpisodes,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of results")
	cmd.Flags().String("filter", "", "Filter: listenable (default) or editable")

	return cmd
}

func runSearchUserEpisodes(cmd *cobra.Command, args []string) error {
	userID, err := parseUserID(args[0])
	if err != nil {
		return err
	}
	query := args[1]

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	filter, _ := cmd.Flags().GetString("filter")

	result, err := client.SearchUserEpisodes(
		userID,
		api.SearchParams{Query: query, Filter: filter},
		api.PaginationParams{Limit: limit},
	)
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
		formatter.PrintMessage("\n(more results available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// search show-episodes
// -----------------------------------------------------------------------------

func newSearchShowEpisodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-episodes <show-id> <query>",
		Short: "Search for episodes within a specific show",
		Args:  cobra.ExactArgs(2),
		RunE:  runSearchShowEpisodes,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of results")
	cmd.Flags().String("filter", "", "Filter: listenable (default) or editable")

	return cmd
}

func runSearchShowEpisodes(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}
	query := args[1]

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	filter, _ := cmd.Flags().GetString("filter")

	result, err := client.SearchShowEpisodes(
		showID,
		api.SearchParams{Query: query, Filter: filter},
		api.PaginationParams{Limit: limit},
	)
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
		formatter.PrintMessage("\n(more results available, use --limit to see more)")
	}

	return nil
}
