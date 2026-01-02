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
		newShowsDeleteCmd(),
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
