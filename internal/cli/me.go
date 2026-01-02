/*
me.go - Current user command

A simple command that shows the authenticated user's profile.
This is often the first command users run after login to verify
their authentication is working.
*/
package cli

import (
	"github.com/spf13/cobra"
)

func newMeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show current authenticated user",
		Long: `Display information about the currently authenticated user.

This is useful to verify that your authentication is working correctly
and to see your user ID for other commands.`,
		RunE: runMe,
	}
}

func runMe(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	user, err := client.GetMe()
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintUser(user)
	return nil
}
