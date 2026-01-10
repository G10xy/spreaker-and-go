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
		// Overall statistics
		newStatsMeCmd(),
		newStatsShowCmd(),
		newStatsEpisodeCmd(),
		// Play statistics
		newStatsPlaysCmd(),
		newStatsPlaysUserCmd(),
		newStatsPlaysEpisodeCmd(),
		newStatsShowsTotalsCmd(),
		newStatsEpisodesTotalsCmd(),
		// Likes statistics
		newStatsLikesCmd(),
		newStatsLikesUserCmd(),
		newStatsLikesEpisodeCmd(),
		// Followers statistics
		newStatsFollowersCmd(),
		// Sources statistics
		newStatsSourcesCmd(),
		newStatsSourcesUserCmd(),
		newStatsSourcesEpisodeCmd(),
		// Devices statistics
		newStatsDevicesCmd(),
		newStatsDevicesUserCmd(),
		newStatsDevicesEpisodeCmd(),
		// OS statistics
		newStatsOSCmd(),
		newStatsOSUserCmd(),
		newStatsOSEpisodeCmd(),
		// Geographic statistics
		newStatsGeoCmd(),
		newStatsGeoUserCmd(),
		// Listeners statistics
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
// stats plays (show)
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
// stats plays-user
// -----------------------------------------------------------------------------

func newStatsPlaysUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plays-user",
		Short: "Show play statistics for authenticated user over time",
		RunE:  runStatsPlaysUser,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsPlaysUser(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	group, _ := cmd.Flags().GetString("group")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	stats, err := client.GetUserPlayStatistics(me.UserID, api.StatisticsParams{
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
// stats plays-episode
// -----------------------------------------------------------------------------

func newStatsPlaysEpisodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plays-episode <episode-id>",
		Short: "Show play statistics for an episode over time",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsPlaysEpisode,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsPlaysEpisode(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
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

	stats, err := client.GetEpisodePlayStatistics(episodeID, api.StatisticsParams{
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
// stats shows-totals
// -----------------------------------------------------------------------------

func newStatsShowsTotalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shows-totals",
		Short: "Show play totals for each of your shows",
		RunE:  runStatsShowsTotals,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().IntP("limit", "l", 20, "Maximum number of shows")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsShowsTotals(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	limit, _ := cmd.Flags().GetInt("limit")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	result, err := client.GetUserShowsPlayTotals(me.UserID, api.StatisticsParams{
		From: from,
		To:   to,
	}, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintShowsPlayTotals(result.Items)
	return nil
}

// -----------------------------------------------------------------------------
// stats episodes-totals
// -----------------------------------------------------------------------------

func newStatsEpisodesTotalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "episodes-totals <show-id>",
		Short: "Show play totals for each episode in a show",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsEpisodesTotals,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().IntP("limit", "l", 20, "Maximum number of episodes")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsEpisodesTotals(cmd *cobra.Command, args []string) error {
	showID, err := parseShowID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	limit, _ := cmd.Flags().GetInt("limit")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	result, err := client.GetShowEpisodesPlayTotals(showID, api.StatisticsParams{
		From: from,
		To:   to,
	}, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintEpisodesPlayTotals(result.Items)
	return nil
}

// -----------------------------------------------------------------------------
// stats likes (show)
// -----------------------------------------------------------------------------

func newStatsLikesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "likes <show-id>",
		Short: "Show likes statistics for a show over time",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsLikes,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsLikes(cmd *cobra.Command, args []string) error {
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

	stats, err := client.GetShowLikesStatistics(showID, api.StatisticsParams{
		From:  from,
		To:    to,
		Group: group,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintLikesStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats likes-user
// -----------------------------------------------------------------------------

func newStatsLikesUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "likes-user",
		Short: "Show likes statistics for authenticated user over time",
		RunE:  runStatsLikesUser,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsLikesUser(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	group, _ := cmd.Flags().GetString("group")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	stats, err := client.GetUserLikesStatistics(me.UserID, api.StatisticsParams{
		From:  from,
		To:    to,
		Group: group,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintLikesStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats likes-episode
// -----------------------------------------------------------------------------

func newStatsLikesEpisodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "likes-episode <episode-id>",
		Short: "Show likes statistics for an episode over time",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsLikesEpisode,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsLikesEpisode(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
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

	stats, err := client.GetEpisodeLikesStatistics(episodeID, api.StatisticsParams{
		From:  from,
		To:    to,
		Group: group,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintLikesStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats followers
// -----------------------------------------------------------------------------

func newStatsFollowersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "followers",
		Short: "Show followers statistics for authenticated user over time",
		RunE:  runStatsFollowers,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsFollowers(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	group, _ := cmd.Flags().GetString("group")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	stats, err := client.GetUserFollowersStatistics(me.UserID, api.StatisticsParams{
		From:  from,
		To:    to,
		Group: group,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintFollowersStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats sources (show)
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
// stats sources-user
// -----------------------------------------------------------------------------

func newStatsSourcesUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sources-user",
		Short: "Show play/download sources for authenticated user",
		RunE:  runStatsSourcesUser,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsSourcesUser(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	group, _ := cmd.Flags().GetString("group")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	stats, err := client.GetUserSourcesStatistics(me.UserID, api.StatisticsParams{
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
// stats sources-episode
// -----------------------------------------------------------------------------

func newStatsSourcesEpisodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sources-episode <episode-id>",
		Short: "Show play/download sources for an episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsSourcesEpisode,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.Flags().String("group", "day", "Group by: day, week, or month")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsSourcesEpisode(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
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

	stats, err := client.GetEpisodeSourcesStatistics(episodeID, api.StatisticsParams{
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
// stats devices (show)
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
// stats devices-user
// -----------------------------------------------------------------------------

func newStatsDevicesUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices-user",
		Short: "Show device breakdown for authenticated user",
		RunE:  runStatsDevicesUser,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsDevicesUser(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	stats, err := client.GetUserDevicesStatistics(me.UserID, api.StatisticsParams{
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
// stats devices-episode
// -----------------------------------------------------------------------------

func newStatsDevicesEpisodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices-episode <episode-id>",
		Short: "Show device breakdown for an episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsDevicesEpisode,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsDevicesEpisode(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetEpisodeDevicesStatistics(episodeID, api.StatisticsParams{
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
// stats os (show)
// -----------------------------------------------------------------------------

func newStatsOSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "os <show-id>",
		Short: "Show operating system breakdown for a show",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsOS,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsOS(cmd *cobra.Command, args []string) error {
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

	stats, err := client.GetShowOSStatistics(showID, api.StatisticsParams{
		From: from,
		To:   to,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintOSStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats os-user
// -----------------------------------------------------------------------------

func newStatsOSUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "os-user",
		Short: "Show operating system breakdown for authenticated user",
		RunE:  runStatsOSUser,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsOSUser(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	stats, err := client.GetUserOSStatistics(me.UserID, api.StatisticsParams{
		From: from,
		To:   to,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintOSStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats os-episode
// -----------------------------------------------------------------------------

func newStatsOSEpisodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "os-episode <episode-id>",
		Short: "Show operating system breakdown for an episode",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatsOSEpisode,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsOSEpisode(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	stats, err := client.GetEpisodeOSStatistics(episodeID, api.StatisticsParams{
		From: from,
		To:   to,
	})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintOSStatistics(stats)
	return nil
}

// -----------------------------------------------------------------------------
// stats geo (show)
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
// stats geo-user
// -----------------------------------------------------------------------------

func newStatsGeoUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "geo-user",
		Short: "Show geographic breakdown for authenticated user",
		RunE:  runStatsGeoUser,
	}

	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, required)")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func runStatsGeoUser(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	me, err := client.GetMe()
	if err != nil {
		return err
	}

	stats, err := client.GetUserGeographicStatistics(me.UserID, api.StatisticsParams{
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
// stats listeners (show only)
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
