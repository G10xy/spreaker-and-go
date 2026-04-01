/*
root.go - Root command and CLI structure

The root command defines the overall CLI structure
It registers all subcommands and sets up global flags that all commands can use

In Cobra-based CLIs, you typically have a root command that acts as
the parent for all subcommands. Each subcommand is defined in its own
file with a constructor function (e.g., newShowsCmd()).
*/
package cli

import (
	"context"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func Execute(ctx context.Context, version string) error {
	rootCmd = newRootCmd(version)
	return rootCmd.ExecuteContext(ctx)
}

// newRootCmd creates the root command with all subcommands registered.
func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spreaker",
		Short: "A CLI for the Spreaker podcast platform",
		Long: `spreaker-cli is a command line interface for managing your podcasts on Spreaker.

You can manage shows, episodes, view statistics, and more - all from your terminal.

Get started:
  spreaker login          # Authenticate with your API token
  spreaker me             # View your profile
  spreaker shows list     # List your shows
  spreaker episodes list  # List episodes`,
		Version: version,
		// SilenceUsage prevents printing usage on errors.
		// SilenceErrors prevents Cobra from double-printing errors
		// that are already displayed by spinners or formatters.
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags are available to ALL subcommands.
	// PersistentFlags() makes them "inherited" by children.
	cmd.PersistentFlags().StringP("output", "o", "", "Output format: table, json, plain")
	cmd.PersistentFlags().String("token", "", "API token (overrides config) — INSECURE: visible in process listings, prefer SPREAKER_TOKEN env var")
	cmd.PersistentFlags().MarkHidden("token")
	cmd.PersistentFlags().Bool("no-color", false, "Disable colored output")

	cmd.AddCommand(
		newLoginCmd(),
		newMeCmd(),

		newUsersCmd(),
		newShowsCmd(),
		newEpisodesCmd(),

		newStatsCmd(),

		newSearchCmd(),
		newExploreCmd(),
		newTagsCmd(),

		newChaptersCmd(),
		newCuepointsCmd(),
		newMessagesCmd(),

		newMiscCmd(),
		newConfigCmd(),
	)

	return cmd
}
