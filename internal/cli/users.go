/*
users.go - User management commands

This file contains all commands related to users
*/
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users",
		Long: `Manage users on Spreaker.

Examples:
  spreaker users get 12345              # Get a user's profile
  spreaker users shows 12345            # List a user's shows
  spreaker users followers 12345        # List a user's followers
  spreaker users follow 12345           # Follow a user
  spreaker users block 12345            # Block a user`,
	}

	cmd.AddCommand(
		newUsersGetCmd(),
		newUsersUpdateCmd(),
		newUsersShowsCmd(),
		newUsersFollowersCmd(),
		newUsersFollowingsCmd(),
		newUsersFollowCmd(),
		newUsersUnfollowCmd(),
		newUsersBlocksCmd(),
		newUsersBlockCmd(),
		newUsersUnblockCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// users get
// -----------------------------------------------------------------------------

func newUsersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <user-id>",
		Short: "Get a user's public profile",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersGet,
	}
}

func runUsersGet(cmd *cobra.Command, args []string) error {
	userID, err := parseUserID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	user, err := client.GetUser(userID)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintUser(user)
	return nil
}

// -----------------------------------------------------------------------------
// users update
// -----------------------------------------------------------------------------

func newUsersUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update your profile",
		Long: `Update the authenticated user's profile.

Examples:
  spreaker users update --fullname "John Doe"
  spreaker users update --description "Podcast enthusiast"
  spreaker users update --username johndoe`,
		RunE: runUsersUpdate,
	}

	cmd.Flags().String("fullname", "", "Display name")
	cmd.Flags().String("description", "", "Bio/description")
	cmd.Flags().String("username", "", "Username")
	cmd.Flags().String("gender", "", "Gender (male, female, other)")
	cmd.Flags().String("birthday", "", "Birthday (YYYY-MM-DD)")
	cmd.Flags().String("location", "", "Location")
	cmd.Flags().String("contact-email", "", "Contact email")

	return cmd
}

func runUsersUpdate(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	// Get current user ID
	me, err := client.GetMe()
	if err != nil {
		return err
	}

	params := api.UpdateUserParams{}

	if val, _ := cmd.Flags().GetString("fullname"); val != "" {
		params.Fullname = &val
	}
	if val, _ := cmd.Flags().GetString("description"); val != "" {
		params.Description = &val
	}
	if val, _ := cmd.Flags().GetString("username"); val != "" {
		params.Username = &val
	}
	if val, _ := cmd.Flags().GetString("gender"); val != "" {
		params.Gender = &val
	}
	if val, _ := cmd.Flags().GetString("birthday"); val != "" {
		params.Birthday = &val
	}
	if val, _ := cmd.Flags().GetString("location"); val != "" {
		params.Location = &val
	}
	if val, _ := cmd.Flags().GetString("contact-email"); val != "" {
		params.ContactEmail = &val
	}

	user, err := client.UpdateUser(me.UserID, params)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess("Profile updated")
	formatter.PrintUser(user)
	return nil
}

// -----------------------------------------------------------------------------
// users shows
// -----------------------------------------------------------------------------

func newUsersShowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shows <user-id>",
		Short: "List a user's shows",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersShows,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of shows to list")

	return cmd
}

func runUsersShows(cmd *cobra.Command, args []string) error {
	userID, err := parseUserID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetUserShows(userID, api.PaginationParams{Limit: limit})
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
// users followers
// -----------------------------------------------------------------------------

func newUsersFollowersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "followers <user-id>",
		Short: "List a user's followers",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersFollowers,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of followers to list")

	return cmd
}

func runUsersFollowers(cmd *cobra.Command, args []string) error {
	userID, err := parseUserID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetUserFollowers(userID, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(result.Items) == 0 {
		formatter.PrintMessage("No followers found.")
		return nil
	}

	formatter.PrintUsers(result.Items)

	if result.HasMore {
		formatter.PrintMessage("\n(more followers available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// users followings
// -----------------------------------------------------------------------------

func newUsersFollowingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "followings <user-id>",
		Short: "List who a user follows",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersFollowings,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of users to list")

	return cmd
}

func runUsersFollowings(cmd *cobra.Command, args []string) error {
	userID, err := parseUserID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetUserFollowings(userID, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(result.Items) == 0 {
		formatter.PrintMessage("No followings found.")
		return nil
	}

	formatter.PrintUsers(result.Items)

	if result.HasMore {
		formatter.PrintMessage("\n(more users available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// users follow
// -----------------------------------------------------------------------------

func newUsersFollowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "follow <user-id>",
		Short: "Follow a user",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersFollow,
	}
}

func runUsersFollow(cmd *cobra.Command, args []string) error {
	followingID, err := parseUserID(args[0])
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

	if err := client.FollowUser(me.UserID, followingID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Now following user %d", followingID))
	return nil
}

// -----------------------------------------------------------------------------
// users unfollow
// -----------------------------------------------------------------------------

func newUsersUnfollowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unfollow <user-id>",
		Short: "Unfollow a user",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersUnfollow,
	}
}

func runUsersUnfollow(cmd *cobra.Command, args []string) error {
	followingID, err := parseUserID(args[0])
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

	if err := client.UnfollowUser(me.UserID, followingID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Unfollowed user %d", followingID))
	return nil
}

// -----------------------------------------------------------------------------
// users blocks
// -----------------------------------------------------------------------------

func newUsersBlocksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocks",
		Short: "List your blocked users",
		RunE:  runUsersBlocks,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of users to list")

	return cmd
}

func runUsersBlocks(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetUserBlocks(me.UserID, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(result.Items) == 0 {
		formatter.PrintMessage("No blocked users.")
		return nil
	}

	formatter.PrintUsers(result.Items)

	if result.HasMore {
		formatter.PrintMessage("\n(more users available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// users block
// -----------------------------------------------------------------------------

func newUsersBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block <user-id>",
		Short: "Block a user",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersBlock,
	}
}

func runUsersBlock(cmd *cobra.Command, args []string) error {
	blockedID, err := parseUserID(args[0])
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

	if err := client.BlockUser(me.UserID, blockedID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Blocked user %d", blockedID))
	return nil
}

// -----------------------------------------------------------------------------
// users unblock
// -----------------------------------------------------------------------------

func newUsersUnblockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unblock <user-id>",
		Short: "Unblock a user",
		Args:  cobra.ExactArgs(1),
		RunE:  runUsersUnblock,
	}
}

func runUsersUnblock(cmd *cobra.Command, args []string) error {
	blockedID, err := parseUserID(args[0])
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

	if err := client.UnblockUser(me.UserID, blockedID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Unblocked user %d", blockedID))
	return nil
}
