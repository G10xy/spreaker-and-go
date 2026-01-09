/*
shows.go - Show management commands

This file contains all commands related to podcast shows
*/
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newShowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shows",
		Short: "Manage your podcast shows",
		Long: `Manage your podcast shows on Spreaker.

Examples:
  spreaker shows list              # List all your shows
  spreaker shows get 12345         # Get details of a show
  spreaker shows delete 12345      # Delete a show`,
	}

	cmd.AddCommand(
		newShowsListCmd(),
		newShowsGetCmd(),
		newShowsCreateCmd(),
		newShowsUpdateCmd(),
		newShowsDeleteCmd(),
		newShowsFavoritesCmd(),
		newShowsFavoriteCmd(),
		newShowsUnfavoriteCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// shows list
// -----------------------------------------------------------------------------

func newShowsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all your shows",
		RunE:  runShowsList,
	}

	// Local flags only apply to this specific command, not its children.
	// Use Flags() for local flags, PersistentFlags() for inherited flags.
	cmd.Flags().IntP("limit", "l", 20, "Maximum number of shows to list")

	return cmd
}

func runShowsList(cmd *cobra.Command, args []string) error {
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
}

// -----------------------------------------------------------------------------
// shows get
// -----------------------------------------------------------------------------

func newShowsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <show-id>",
		Short: "Get details of a specific show",
		Args:  cobra.ExactArgs(1),
		RunE:  runShowsGet,
	}
}

func runShowsGet(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
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
}

// -----------------------------------------------------------------------------
// shows delete
// -----------------------------------------------------------------------------

func newShowsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <show-id>",
		Short: "Delete a show",
		Long: `Delete a show permanently.

WARNING: This action cannot be undone. All episodes in the show will also be deleted.`,
		Args: cobra.ExactArgs(1),
		RunE: runShowsDelete,
	}

	// --force flag to skip confirmation
	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runShowsDelete(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	// Check if --force flag was provided
	force, _ := cmd.Flags().GetBool("force")
	if !force {
		prompt := fmt.Sprintf("Are you sure you want to delete show %d? [y/N]: ", showID)
		if !confirmAction(prompt) {
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
}

// -----------------------------------------------------------------------------
// shows create
// -----------------------------------------------------------------------------

func newShowsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new show",
		Long: `Create a new podcast show.

Examples:
  spreaker shows create --title "My Podcast"
  spreaker shows create --title "My Podcast" --language en --category 1`,
		RunE: runShowsCreate,
	}

	cmd.Flags().String("title", "", "Show title (required)")
	cmd.Flags().String("description", "", "Show description")
	cmd.Flags().String("language", "", "Language code (e.g., en, it, es)")
	cmd.Flags().Int("category", 0, "Category ID")
	cmd.Flags().Bool("explicit", false, "Mark as explicit content")

	cmd.MarkFlagRequired("title")

	return cmd
}

func runShowsCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	language, _ := cmd.Flags().GetString("language")
	categoryID, _ := cmd.Flags().GetInt("category")
	explicit, _ := cmd.Flags().GetBool("explicit")

	params := api.CreateShowParams{
		Title:       title,
		Description: description,
		Language:    language,
		CategoryID:  categoryID,
		Explicit:    explicit,
	}

	show, err := client.CreateShow(params)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Show created with ID %d", show.ShowID))
	formatter.PrintShow(show)
	return nil
}

// -----------------------------------------------------------------------------
// shows update
// -----------------------------------------------------------------------------

func newShowsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <show-id>",
		Short: "Update a show",
		Long: `Update an existing show.

Examples:
  spreaker shows update 12345 --title "New Title"
  spreaker shows update 12345 --description "New description"`,
		Args: cobra.ExactArgs(1),
		RunE: runShowsUpdate,
	}

	cmd.Flags().String("title", "", "Show title")
	cmd.Flags().String("description", "", "Show description")
	cmd.Flags().String("language", "", "Language code (e.g., en, it, es)")
	cmd.Flags().Int("category", 0, "Category ID")
	cmd.Flags().Bool("explicit", false, "Mark as explicit content")

	return cmd
}

func runShowsUpdate(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	params := api.UpdateShowParams{}

	if cmd.Flags().Changed("title") {
		val, _ := cmd.Flags().GetString("title")
		params.Title = &val
	}
	if cmd.Flags().Changed("description") {
		val, _ := cmd.Flags().GetString("description")
		params.Description = &val
	}
	if cmd.Flags().Changed("language") {
		val, _ := cmd.Flags().GetString("language")
		params.Language = &val
	}
	if cmd.Flags().Changed("category") {
		val, _ := cmd.Flags().GetInt("category")
		params.CategoryID = &val
	}
	if cmd.Flags().Changed("explicit") {
		val, _ := cmd.Flags().GetBool("explicit")
		params.Explicit = &val
	}

	show, err := client.UpdateShow(showID, params)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess("Show updated")
	formatter.PrintShow(show)
	return nil
}

// -----------------------------------------------------------------------------
// shows favorites
// -----------------------------------------------------------------------------

func newShowsFavoritesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "favorites",
		Short: "List your favorite shows",
		RunE:  runShowsFavorites,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of shows to list")

	return cmd
}

func runShowsFavorites(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetFavoriteShows(me.UserID, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(result.Items) == 0 {
		formatter.PrintMessage("No favorite shows.")
		return nil
	}

	formatter.PrintShows(result.Items)

	if result.HasMore {
		formatter.PrintMessage("\n(more shows available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// shows favorite
// -----------------------------------------------------------------------------

func newShowsFavoriteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "favorite <show-id>",
		Short: "Add a show to your favorites",
		Args:  cobra.ExactArgs(1),
		RunE:  runShowsFavorite,
	}
}

func runShowsFavorite(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
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

	if err := client.AddShowToFavorites(me.UserID, showID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Show %d added to favorites", showID))
	return nil
}

// -----------------------------------------------------------------------------
// shows unfavorite
// -----------------------------------------------------------------------------

func newShowsUnfavoriteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unfavorite <show-id>",
		Short: "Remove a show from your favorites",
		Args:  cobra.ExactArgs(1),
		RunE:  runShowsUnfavorite,
	}
}

func runShowsUnfavorite(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
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

	if err := client.RemoveShowFromFavorites(me.UserID, showID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Show %d removed from favorites", showID))
	return nil
}
