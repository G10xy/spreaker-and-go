/*
stats.go - Statistics commands

This file contains all commands for viewing statistics
*/
package cli

import (
	"github.com/spf13/cobra"
	
	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "View statistics for users, shows, and episodes",
		Long: `View various statistics about your podcasts.

Overall statistics:
  spreaker stats me                    # Personal stats
  spreaker stats show 12345            # Show's overall stats
  spreaker stats episode 67890         # Episode's overall stats

Time-series statistics (require --from and --to):
  spreaker stats plays 12345 --from 2024-01-01 --to 2024-01-31
  spreaker stats devices 12345 --from 2024-01-01 --to 2024-01-31
  spreaker stats listeners 12345 --from 2024-01-01 --to 2024-01-31`,
	}

	cmd.AddCommand(
		newStatsMeCmd(),
		newStatsShowCmd(),
		newStatsEpisodeCmd(),
		newStatsPlaysCmd(),
		newStatsDevicesCmd(),
		newStatsGeoCmd(),
		newStatsSourcesCmd(),
		newStatsListenersCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// stats me
// -----------------------------------------------------------------------------

func newStatsMeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show your overall statistics",
		RunE:  runStatsMe,
	}
}

func runStatsMe(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetMyStatistics()
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintUserStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats show
// -----------------------------------------------------------------------------

func newStatsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <show-id>",
		Short: "Show statistics for a specific show",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsShow,
	}
}

func runStatsShow(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetShowStatistics(showID)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintShowStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats episode
// -----------------------------------------------------------------------------

func newStatsEpisodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "episode <episode-id>",
		Short: "Show statistics for a specific episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsEpisode,
	}
}

func runStatsEpisode(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetEpisodeStatistics(episodeID)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintEpisodeStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats plays
// -----------------------------------------------------------------------------

func newStatsPlaysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plays <show-id>",
		Short: "Show play statistics for a show over time",
		Long: `Show play statistics for a show over a date range.

Example:
  spreaker stats plays 12345 --from 2024-01-01 --to 2024-01-31 --group day`,
		Args: cobra.ExactArgs(1),
		RunE: runStatsPlays,
	}

	// Date range flags - marked as required
	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsPlays(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	group, _ := cmd.Flags().GetString("group")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetShowPlayStatistics(showID, api.StatisticsParams{
		From:  from,
		To:    to,
		Group: group,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintPlayStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats devices
// -----------------------------------------------------------------------------

func newStatsDevicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices <show-id>",
		Short: "Show device breakdown for a show",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsDevices,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsDevices(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetShowDevicesStatistics(showID, api.StatisticsParams{
		From: from,
		To:   to,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintDeviceStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats geo
// -----------------------------------------------------------------------------

func newStatsGeoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "geo <show-id>",
		Short: "Show geographic breakdown for a show",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsGeo,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsGeo(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetShowGeographicStatistics(showID, api.StatisticsParams{
		From: from,
		To:   to,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintGeographicStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats sources
// -----------------------------------------------------------------------------

func newStatsSourcesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sources <show-id>",
		Short: "Show play/download sources for a show",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsSources,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsSources(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	group, _ := cmd.Flags().GetString("group")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetShowSourcesStatistics(showID, api.StatisticsParams{
		From:  from,
		To:    to,
		Group: group,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintSourcesStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats listeners
// -----------------------------------------------------------------------------

func newStatsListenersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listeners <show-id>",
		Short: "Show unique listeners for a show over time",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsListeners,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsListeners(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	group, _ := cmd.Flags().GetString("group")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetShowListenersStatistics(showID, api.StatisticsParams{
		From:  from,
		To:    to,
		Group: group,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintListenersStatistics(stats)
	return nil
}
