/*
cuepoints.go - Episode cuepoint management commands

Cuepoints are specific points in time within an episode where
audio ads can be injected. This file contains commands to manage them.
*/
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/pkg/models"
)

func newCuepointsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cuepoints",
		Aliases: []string{"cue"},
		Short:   "Manage episode cuepoints for ad injection",
		Long: `Manage cuepoints for episodes. Cuepoints are specific points in time 
within an episode where audio ads can be injected.

Note: Setting cuepoints is not enough to get ads injected. You also need to
enable Ads and Monetization capabilities on your account and show.

Examples:
  spreaker cuepoints list 12345
  spreaker cuepoints set 12345 30000:1 60000:2
  spreaker cuepoints delete 12345`,
	}

	cmd.AddCommand(
		newCuepointsListCmd(),
		newCuepointsSetCmd(),
		newCuepointsDeleteCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// cuepoints list
// -----------------------------------------------------------------------------

func newCuepointsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <episode-id>",
		Short: "List all cuepoints for an episode",
		Long: `List all cuepoints for an episode, sorted chronologically by timecode.

Examples:
  spreaker cuepoints list 12345
  spreaker cue list 12345 --output json`,
		Args: cobra.ExactArgs(1),
		RunE: runCuepointsList,
	}
}

func runCuepointsList(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	cuepoints, err := client.GetEpisodeCuepoints(episodeID)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(cuepoints) == 0 {
		formatter.PrintMessage("No cuepoints found for this episode.")
		return nil
	}

	formatter.PrintCuepoints(cuepoints)
	return nil
}

// -----------------------------------------------------------------------------
// cuepoints set
// -----------------------------------------------------------------------------

func newCuepointsSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <episode-id> [timecode:max_ads]...",
		Short: "Set cuepoints for an episode (replaces all existing)",
		Long: `Set cuepoints for an episode. This replaces all existing cuepoints.
Timecodes are in milliseconds. Format: timecode:max_ads

Examples:
  # Set a single cuepoint at 30 seconds (30000ms) with max 1 ad
  spreaker cuepoints set 12345 30000:1

  # Set multiple cuepoints
  spreaker cuepoints set 12345 30000:1 60000:2 90000:1

  # Clear all cuepoints (set empty list)
  spreaker cuepoints set 12345`,
		Args: cobra.MinimumNArgs(1),
		RunE: runCuepointsSet,
	}
}

func runCuepointsSet(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	// Parse cuepoints from remaining arguments
	cuepoints := make([]models.Cuepoint, 0, len(args)-1)
	for _, arg := range args[1:] {
		var timecode, maxAds int
		if _, err := fmt.Sscanf(arg, "%d:%d", &timecode, &maxAds); err != nil {
			return fmt.Errorf("invalid cuepoint format '%s' (expected timecode:max_ads, e.g., 30000:1)", arg)
		}
		cuepoints = append(cuepoints, models.Cuepoint{
			Timecode:    timecode,
			AdsMaxCount: maxAds,
		})
	}

	if err := client.UpdateEpisodeCuepoints(episodeID, cuepoints); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	if len(cuepoints) == 0 {
		formatter.PrintMessage("All cuepoints deleted successfully.")
	} else {
		formatter.PrintMessage(fmt.Sprintf("Successfully set %d cuepoint(s).", len(cuepoints)))
	}
	return nil
}

// -----------------------------------------------------------------------------
// cuepoints delete
// -----------------------------------------------------------------------------

func newCuepointsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <episode-id>",
		Short: "Delete all cuepoints for an episode",
		Long: `Delete all cuepoints for an episode.

Examples:
  spreaker cuepoints delete 12345
  spreaker cue delete 12345`,
		Args: cobra.ExactArgs(1),
		RunE: runCuepointsDelete,
	}
}

func runCuepointsDelete(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	if err := client.DeleteEpisodeCuepoints(episodeID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage("All cuepoints deleted successfully.")
	return nil
}
