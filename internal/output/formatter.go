/*
Package output handles formatting and displaying results to the terminal.

It supports multiple output formats:
  - table: Human-readable aligned columns (default)
  - json:  Machine-readable JSON output
  - plain: Simple text, one item per line
*/
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/G10xy/spreaker-and-go/pkg/models"
	"github.com/pterm/pterm"
)


type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatPlain Format = "plain"
)

type Formatter struct {
	format Format
	writer io.Writer
	color  bool
}

// New creates a new Formatter with the specified format and color support.
func New(format string, color bool) *Formatter {
	f := Format(strings.ToLower(strings.TrimSpace(format)))

	switch f {
	case FormatTable, FormatJSON, FormatPlain:
	default:
		f = FormatTable
	}

	// Only enable color for table format
	if f != FormatTable {
		color = false
	}

	return &Formatter{
		format: f,
		writer: os.Stdout,
		color:  color,
	}
}

func (f *Formatter) tabw() *tabwriter.Writer {
    return tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
}

// -----------------------------------------------------------------------------
// User Output
// -----------------------------------------------------------------------------


func (f *Formatter) PrintUser(user *models.User) {
	switch f.format {
	case FormatJSON:
		f.printJSON(user)
	case FormatPlain:
		fmt.Fprintf(f.writer, "%d\t%s\n", user.UserID, user.Fullname)
	default:
		f.printUserTable(user)
	}
}

func (f *Formatter) PrintUsers(users []models.User) {
	switch f.format {
	case FormatJSON:
		f.printJSON(users)
	case FormatPlain:
		for _, u := range users {
			fmt.Fprintf(f.writer, "%d\t%s\n", u.UserID, u.Fullname)
		}
	default:
		f.printUsersTable(users)
	}
}

func (f *Formatter) printUserTable(user *models.User) {
	pairs := [][2]string{
		{"ID:", fmt.Sprintf("%d", user.UserID)},
		{"Username:", user.Username},
		{"Name:", user.Fullname},
		{"Kind:", user.Kind},
		{"Plan:", user.Plan},
		{"Followers:", fmt.Sprintf("%d", user.FollowersCount)},
		{"Following:", fmt.Sprintf("%d", user.FollowingsCount)},
		{"URL:", user.SiteURL},
	}

	if user.Description != "" {
		desc := user.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		pairs = append(pairs, [2]string{"Bio:", desc})
	}

	f.PrintKeyValue(pairs)
}

func (f *Formatter) printUsersTable(users []models.User) {
	header := []string{"ID", "USERNAME", "NAME", "FOLLOWERS"}
	rows := make([][]string, len(users))
	for i, u := range users {
		rows[i] = []string{
			fmt.Sprintf("%d", u.UserID),
			u.Username,
			truncate(u.Fullname, 30),
			fmt.Sprintf("%d", u.FollowersCount),
		}
	}
	f.renderTable(header, rows)
}

// -----------------------------------------------------------------------------
// Show Output
// -----------------------------------------------------------------------------


func (f *Formatter) PrintShow(show *models.Show) {
	switch f.format {
	case FormatJSON:
		f.printJSON(show)
	case FormatPlain:
		fmt.Fprintf(f.writer, "%d\t%s\n", show.ShowID, show.Title)
	default:
		f.printShowTable(show)
	}
}

func (f *Formatter) PrintShows(shows []models.Show) {
	switch f.format {
	case FormatJSON:
		f.printJSON(shows)
	case FormatPlain:
		for _, s := range shows {
			fmt.Fprintf(f.writer, "%d\t%s\n", s.ShowID, s.Title)
		}
	default:
		f.printShowsTable(shows)
	}
}

func (f *Formatter) printShowTable(show *models.Show) {
	pairs := [][2]string{
		{"ID:", fmt.Sprintf("%d", show.ShowID)},
		{"Title:", show.Title},
		{"Language:", show.Language},
		{"Episodes:", fmt.Sprintf("%d", show.EpisodesCount)},
		{"Followers:", fmt.Sprintf("%d", show.FollowersCount)},
		{"Plays:", fmt.Sprintf("%d", show.PlayCount)},
		{"Explicit:", fmt.Sprintf("%v", show.Explicit)},
		{"URL:", show.SiteURL},
	}

	if show.Description != "" {
		desc := show.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		pairs = append(pairs, [2]string{"Description:", desc})
	}

	if show.LastEpisodeAt != nil {
		pairs = append(pairs, [2]string{"Last Episode:", show.LastEpisodeAt.Format(time.DateTime)})
	}

	f.PrintKeyValue(pairs)
}

func (f *Formatter) printShowsTable(shows []models.Show) {
	header := []string{"ID", "TITLE", "EPISODES", "FOLLOWERS", "PLAYS"}
	rows := make([][]string, len(shows))
	for i, s := range shows {
		rows[i] = []string{
			fmt.Sprintf("%d", s.ShowID),
			truncate(s.Title, 40),
			fmt.Sprintf("%d", s.EpisodesCount),
			fmt.Sprintf("%d", s.FollowersCount),
			fmt.Sprintf("%d", s.PlayCount),
		}
	}
	f.renderTable(header, rows)
}

// -----------------------------------------------------------------------------
// Episode Output
// -----------------------------------------------------------------------------

func (f *Formatter) PrintEpisode(episode *models.Episode) {
	switch f.format {
	case FormatJSON:
		f.printJSON(episode)
	case FormatPlain:
		fmt.Fprintf(f.writer, "%d\t%s\n", episode.EpisodeID, episode.Title)
	default:
		f.printEpisodeTable(episode)
	}
}

func (f *Formatter) PrintEpisodes(episodes []models.Episode) {
	switch f.format {
	case FormatJSON:
		f.printJSON(episodes)
	case FormatPlain:
		for _, e := range episodes {
			fmt.Fprintf(f.writer, "%d\t%s\n", e.EpisodeID, e.Title)
		}
	default:
		f.printEpisodesTable(episodes)
	}
}

func (f *Formatter) printEpisodeTable(episode *models.Episode) {
	pairs := [][2]string{
		{"ID:", fmt.Sprintf("%d", episode.EpisodeID)},
		{"Title:", episode.Title},
		{"Show ID:", fmt.Sprintf("%d", episode.ShowID)},
		{"Duration:", formatDuration(episode.Duration)},
		{"Plays:", fmt.Sprintf("%d", episode.PlayCount)},
		{"Likes:", fmt.Sprintf("%d", episode.LikesCount)},
		{"Status:", episode.EncodingStatus},
		{"Explicit:", fmt.Sprintf("%v", episode.Explicit)},
		{"Downloads:", fmt.Sprintf("%v", episode.DownloadEnabled)},
		{"URL:", episode.SiteURL},
	}

	if episode.PublishedAt != nil {
		pairs = append(pairs, [2]string{"Published:", episode.PublishedAt.Format(time.DateTime)})
	}

	if len(episode.Tags) > 0 {
		pairs = append(pairs, [2]string{"Tags:", strings.Join(episode.Tags, ", ")})
	}

	if episode.Description != "" {
		desc := episode.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		pairs = append(pairs, [2]string{"Description:", desc})
	}

	f.PrintKeyValue(pairs)
}

func (f *Formatter) printEpisodesTable(episodes []models.Episode) {
	header := []string{"ID", "TITLE", "DURATION", "PLAYS", "STATUS", "PUBLISHED"}
	rows := make([][]string, len(episodes))
	for i, e := range episodes {
		published := "-"
		if e.PublishedAt != nil {
			published = e.PublishedAt.Format(time.DateTime)
		}
		rows[i] = []string{
			fmt.Sprintf("%d", e.EpisodeID),
			truncate(e.Title, 35),
			formatDuration(e.Duration),
			fmt.Sprintf("%d", e.PlayCount),
			e.EncodingStatus,
			published,
		}
	}
	f.renderTable(header, rows)
}

// -----------------------------------------------------------------------------
// Statistics Output
// -----------------------------------------------------------------------------

// PrintStatistics prints overall statistics
func (f *Formatter) PrintStatistics(stats *models.Statistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		fmt.Fprintf(f.writer, "plays=%d downloads=%d likes=%d messages=%d\n",
			stats.Plays, stats.Downloads, stats.Likes, stats.Messages)
	default:
		f.printStatisticsTable(stats)
	}
}

func (f *Formatter) printStatisticsTable(stats *models.Statistics) {
	f.PrintKeyValue([][2]string{
		{"Plays:", fmt.Sprintf("%d", stats.Plays)},
		{"Downloads:", fmt.Sprintf("%d", stats.Downloads)},
		{"Likes:", fmt.Sprintf("%d", stats.Likes)},
		{"Messages:", fmt.Sprintf("%d", stats.Messages)},
	})
}

// -----------------------------------------------------------------------------
// Generic Output
// -----------------------------------------------------------------------------

func (f *Formatter) PrintMessage(msg string) {
	if f.color {
		pterm.Info.Println(msg)
	} else {
		fmt.Fprintln(f.writer, msg)
	}
}

func (f *Formatter) PrintError(err error) {
	if f.color {
		pterm.Error.Println(err.Error())
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func (f *Formatter) PrintSuccess(msg string) {
	if f.color {
		pterm.Success.Println(msg)
	} else {
		fmt.Fprintf(f.writer, "✓ %s\n", msg)
	}
}

func (f *Formatter) PrintWarning(msg string) {
	if f.color {
		pterm.Warning.Println(msg)
	} else {
		fmt.Fprintf(os.Stderr, "WARNING: %s\n", msg)
	}
}

// -----------------------------------------------------------------------------
// Styled rendering helpers
// -----------------------------------------------------------------------------

// renderTable renders a list table with a header row.
func (f *Formatter) renderTable(header []string, rows [][]string) {
	if f.color {
		data := pterm.TableData{header}
		data = append(data, rows...)
		pterm.DefaultTable.WithHasHeader(true).WithData(data).Render()
		return
	}
	tw := f.tabw()
	fmt.Fprintln(tw, strings.Join(header, "\t"))
	seps := make([]string, len(header))
	for i, h := range header {
		seps[i] = strings.Repeat("-", max(len(h), 2))
	}
	fmt.Fprintln(tw, strings.Join(seps, "\t"))
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	tw.Flush()
}

// PrintKeyValue renders a detail view with key-value pairs.
func (f *Formatter) PrintKeyValue(pairs [][2]string) {
	if f.color {
		data := pterm.TableData{}
		for _, p := range pairs {
			data = append(data, []string{pterm.FgCyan.Sprint(p[0]), p[1]})
		}
		pterm.DefaultTable.WithData(data).Render()
		return
	}
	tw := f.tabw()
	for _, p := range pairs {
		fmt.Fprintf(tw, "%s\t%s\n", p[0], p[1])
	}
	tw.Flush()
}

// renderSection renders a section header.
func (f *Formatter) renderSection(title string) {
	if f.color {
		pterm.DefaultSection.Println(title)
	} else {
		fmt.Fprintf(f.writer, "=== %s ===\n", title)
	}
}

// StartSpinner starts a spinner with the given message. Returns nil if color is disabled.
func (f *Formatter) StartSpinner(msg string) *pterm.SpinnerPrinter {
	if !f.color {
		fmt.Fprintln(f.writer, msg)
		return nil
	}
	spinner, _ := pterm.DefaultSpinner.Start(msg)
	return spinner
}

// StopSpinner stops a spinner with a success or failure message.
func (f *Formatter) StopSpinner(spinner *pterm.SpinnerPrinter, success bool, msg string) {
	if spinner == nil {
		if success {
			f.PrintSuccess(msg)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
		}
		return
	}
	if success {
		spinner.Success(msg)
	} else {
		spinner.Fail(msg)
	}
}

// StartProgressBar starts a progress bar. Returns nil if color is disabled.
func (f *Formatter) StartProgressBar(total int, title string) *pterm.ProgressbarPrinter {
	if !f.color {
		return nil
	}
	bar, _ := pterm.DefaultProgressbar.WithTotal(total).WithTitle(title).Start()
	return bar
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func (f *Formatter) printJSON(v interface{}) {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	encoder.Encode(v)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// formatDuration converts milliseconds to human-readable duration
func formatDuration(ms int) string {
	d := time.Duration(ms) * time.Millisecond

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}


// -----------------------------------------------------------------------------
// Statistics Output (add to internal/output/formatter.go)
// -----------------------------------------------------------------------------

// PrintUserStatistics prints user overall statistics.
func (f *Formatter) PrintUserStatistics(stats *models.UserOverallStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		fmt.Fprintf(f.writer, "plays=%d downloads=%d likes=%d followers=%d shows=%d episodes=%d\n",
			stats.PlaysCount, stats.DownloadsCount, stats.LikesCount,
			stats.FollowersCount, stats.ShowsCount, stats.EpisodesCount)
	default:
		f.printUserStatisticsTable(stats)
	}
}

func (f *Formatter) printUserStatisticsTable(stats *models.UserOverallStatistics) {
	f.renderSection("Overall Statistics")
	f.PrintKeyValue([][2]string{
		{"Total Plays:", fmt.Sprintf("%d", stats.PlaysCount)},
		{"  On Demand:", fmt.Sprintf("%d", stats.PlaysOndemandCount)},
		{"  Live:", fmt.Sprintf("%d", stats.PlaysLiveCount)},
		{"Downloads:", fmt.Sprintf("%d", stats.DownloadsCount)},
		{"Likes:", fmt.Sprintf("%d", stats.LikesCount)},
		{"Followers:", fmt.Sprintf("%d", stats.FollowersCount)},
		{"Shows:", fmt.Sprintf("%d", stats.ShowsCount)},
		{"Episodes:", fmt.Sprintf("%d", stats.EpisodesCount)},
	})
}

// PrintShowStatistics prints show overall statistics.
func (f *Formatter) PrintShowStatistics(stats *models.ShowOverallStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		fmt.Fprintf(f.writer, "plays=%d downloads=%d likes=%d episodes=%d\n",
			stats.PlaysCount, stats.DownloadsCount, stats.LikesCount, stats.EpisodesCount)
	default:
		f.printShowStatisticsTable(stats)
	}
}

func (f *Formatter) printShowStatisticsTable(stats *models.ShowOverallStatistics) {
	if stats.Title != "" {
		f.PrintKeyValue([][2]string{{"Show:", stats.Title}})
	}
	f.renderSection("Overall Statistics")
	f.PrintKeyValue([][2]string{
		{"Total Plays:", fmt.Sprintf("%d", stats.PlaysCount)},
		{"  On Demand:", fmt.Sprintf("%d", stats.PlaysOndemandCount)},
		{"  Live:", fmt.Sprintf("%d", stats.PlaysLiveCount)},
		{"Downloads:", fmt.Sprintf("%d", stats.DownloadsCount)},
		{"Likes:", fmt.Sprintf("%d", stats.LikesCount)},
		{"Episodes:", fmt.Sprintf("%d", stats.EpisodesCount)},
	})
}

// PrintEpisodeStatistics prints episode overall statistics.
func (f *Formatter) PrintEpisodeStatistics(stats *models.EpisodeOverallStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		fmt.Fprintf(f.writer, "plays=%d downloads=%d likes=%d messages=%d\n",
			stats.PlaysCount, stats.DownloadsCount, stats.LikesCount, stats.MessagesCount)
	default:
		f.printEpisodeStatisticsTable(stats)
	}
}

func (f *Formatter) printEpisodeStatisticsTable(stats *models.EpisodeOverallStatistics) {
	f.renderSection("Overall Statistics")
	f.PrintKeyValue([][2]string{
		{"Total Plays:", fmt.Sprintf("%d", stats.PlaysCount)},
		{"  On Demand:", fmt.Sprintf("%d", stats.PlaysOndemandCount)},
		{"  Live:", fmt.Sprintf("%d", stats.PlaysLiveCount)},
		{"Downloads:", fmt.Sprintf("%d", stats.DownloadsCount)},
		{"Likes:", fmt.Sprintf("%d", stats.LikesCount)},
		{"Messages:", fmt.Sprintf("%d", stats.MessagesCount)},
		{"Chapters:", fmt.Sprintf("%d", stats.ChaptersCount)},
	})
}

// PrintPlayStatistics prints time-series play statistics.
func (f *Formatter) PrintPlayStatistics(stats []models.PlayStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats {
			fmt.Fprintf(f.writer, "%s\t%d\t%d\n", s.Date, s.PlaysCount, s.DownloadsCount)
		}
	default:
		f.printPlayStatisticsTable(stats)
	}
}

func (f *Formatter) printPlayStatisticsTable(stats []models.PlayStatistics) {
	header := []string{"DATE", "PLAYS", "ON DEMAND", "LIVE", "DOWNLOADS"}
	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{
			s.Date,
			fmt.Sprintf("%d", s.PlaysCount),
			fmt.Sprintf("%d", s.PlaysOndemandCount),
			fmt.Sprintf("%d", s.PlaysLiveCount),
			fmt.Sprintf("%d", s.DownloadsCount),
		}
	}
	f.renderTable(header, rows)
}

// PrintDeviceStatistics prints device breakdown statistics.
func (f *Formatter) PrintDeviceStatistics(stats []models.DeviceStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats {
			fmt.Fprintf(f.writer, "%s\t%.1f%%\n", s.Name, s.Percentage)
		}
	default:
		f.printDeviceStatisticsTable(stats)
	}
}

func (f *Formatter) printDeviceStatisticsTable(stats []models.DeviceStatistics) {
	header := []string{"DEVICE", "PERCENTAGE"}
	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{s.Name, fmt.Sprintf("%.1f%%", s.Percentage)}
	}
	f.renderTable(header, rows)
}

// PrintGeographicStatistics prints geographic breakdown statistics.
func (f *Formatter) PrintGeographicStatistics(stats *models.GeographicStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, c := range stats.Country {
			fmt.Fprintf(f.writer, "country\t%s\t%.1f%%\n", c.Name, c.Percentage)
		}
		for _, c := range stats.City {
			fmt.Fprintf(f.writer, "city\t%s\t%.1f%%\n", c.Name, c.Percentage)
		}
	default:
		f.printGeographicStatisticsTable(stats)
	}
}

func (f *Formatter) printGeographicStatisticsTable(stats *models.GeographicStatistics) {
	f.renderSection("By Country")
	countryRows := make([][]string, len(stats.Country))
	for i, c := range stats.Country {
		countryRows[i] = []string{c.Name, fmt.Sprintf("%.1f%%", c.Percentage)}
	}
	f.renderTable([]string{"COUNTRY", "PERCENTAGE"}, countryRows)

	fmt.Fprintln(f.writer)

	f.renderSection("By City")
	cityRows := make([][]string, len(stats.City))
	for i, c := range stats.City {
		cityRows[i] = []string{c.Name, fmt.Sprintf("%.1f%%", c.Percentage)}
	}
	f.renderTable([]string{"CITY", "PERCENTAGE"}, cityRows)
}

// PrintSourcesStatistics prints sources breakdown statistics.
func (f *Formatter) PrintSourcesStatistics(stats *models.SourcesStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats.Overall {
			fmt.Fprintf(f.writer, "%s\t%d\t%d%%\n", s.Name, s.PlaysCount, s.Percentage)
		}
	default:
		f.printSourcesStatisticsTable(stats)
	}
}

func (f *Formatter) printSourcesStatisticsTable(stats *models.SourcesStatistics) {
	header := []string{"SOURCE", "PLAYS", "PERCENTAGE"}
	rows := make([][]string, len(stats.Overall))
	for i, s := range stats.Overall {
		rows[i] = []string{s.Name, fmt.Sprintf("%d", s.PlaysCount), fmt.Sprintf("%d%%", s.Percentage)}
	}
	f.renderTable(header, rows)
}

// PrintListenersStatistics prints time-series listeners statistics.
func (f *Formatter) PrintListenersStatistics(stats []models.ListenersStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats {
			fmt.Fprintf(f.writer, "%s\t%d\n", s.Date, s.ListenersCount)
		}
	default:
		f.printListenersStatisticsTable(stats)
	}
}

func (f *Formatter) printListenersStatisticsTable(stats []models.ListenersStatistics) {
	header := []string{"DATE", "LISTENERS"}
	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{s.Date, fmt.Sprintf("%d", s.ListenersCount)}
	}
	f.renderTable(header, rows)
}

// PrintShowsPlayTotals prints play totals for each show.
func (f *Formatter) PrintShowsPlayTotals(stats []models.ShowPlayTotals) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats {
			fmt.Fprintf(f.writer, "%d\t%s\t%d\t%d\n", s.ShowID, s.Title, s.PlaysCount, s.DownloadsCount)
		}
	default:
		f.printShowsPlayTotalsTable(stats)
	}
}

func (f *Formatter) printShowsPlayTotalsTable(stats []models.ShowPlayTotals) {
	header := []string{"SHOW ID", "TITLE", "PLAYS", "ON DEMAND", "LIVE", "DOWNLOADS"}
	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{
			fmt.Sprintf("%d", s.ShowID),
			truncate(s.Title, 30),
			fmt.Sprintf("%d", s.PlaysCount),
			fmt.Sprintf("%d", s.PlaysOndemandCount),
			fmt.Sprintf("%d", s.PlaysLiveCount),
			fmt.Sprintf("%d", s.DownloadsCount),
		}
	}
	f.renderTable(header, rows)
}

// PrintEpisodesPlayTotals prints play totals for each episode.
func (f *Formatter) PrintEpisodesPlayTotals(stats []models.EpisodePlayTotals) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats {
			fmt.Fprintf(f.writer, "%d\t%s\t%d\t%d\n", s.EpisodeID, s.Title, s.PlaysCount, s.DownloadsCount)
		}
	default:
		f.printEpisodesPlayTotalsTable(stats)
	}
}

func (f *Formatter) printEpisodesPlayTotalsTable(stats []models.EpisodePlayTotals) {
	header := []string{"EPISODE ID", "TITLE", "PLAYS", "ON DEMAND", "LIVE", "DOWNLOADS"}
	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{
			fmt.Sprintf("%d", s.EpisodeID),
			truncate(s.Title, 30),
			fmt.Sprintf("%d", s.PlaysCount),
			fmt.Sprintf("%d", s.PlaysOndemandCount),
			fmt.Sprintf("%d", s.PlaysLiveCount),
			fmt.Sprintf("%d", s.DownloadsCount),
		}
	}
	f.renderTable(header, rows)
}

// PrintLikesStatistics prints time-series likes statistics.
func (f *Formatter) PrintLikesStatistics(stats []models.LikesStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats {
			fmt.Fprintf(f.writer, "%s\t%d\n", s.Date, s.LikesCount)
		}
	default:
		f.printLikesStatisticsTable(stats)
	}
}

func (f *Formatter) printLikesStatisticsTable(stats []models.LikesStatistics) {
	header := []string{"DATE", "LIKES"}
	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{s.Date, fmt.Sprintf("%d", s.LikesCount)}
	}
	f.renderTable(header, rows)
}

// PrintFollowersStatistics prints time-series followers statistics.
func (f *Formatter) PrintFollowersStatistics(stats []models.FollowersStatistics) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats {
			fmt.Fprintf(f.writer, "%s\t%d\n", s.Date, s.FollowersCount)
		}
	default:
		f.printFollowersStatisticsTable(stats)
	}
}

func (f *Formatter) printFollowersStatisticsTable(stats []models.FollowersStatistics) {
	header := []string{"DATE", "FOLLOWERS"}
	rows := make([][]string, len(stats))
	for i, s := range stats {
		rows[i] = []string{s.Date, fmt.Sprintf("%d", s.FollowersCount)}
	}
	f.renderTable(header, rows)
}

// PrintOSStatistics prints operating system breakdown statistics.
func (f *Formatter) PrintOSStatistics(stats *models.OSStatisticsBreakdown) {
	switch f.format {
	case FormatJSON:
		f.printJSON(stats)
	case FormatPlain:
		for _, s := range stats.Desktop {
			fmt.Fprintf(f.writer, "desktop\t%s\t%.1f%%\n", s.Name, s.Percentage)
		}
		for _, s := range stats.Mobile {
			fmt.Fprintf(f.writer, "mobile\t%s\t%.1f%%\n", s.Name, s.Percentage)
		}
	default:
		f.printOSStatisticsTable(stats)
	}
}

func (f *Formatter) printOSStatisticsTable(stats *models.OSStatisticsBreakdown) {
	f.renderSection("Desktop")
	desktopRows := make([][]string, len(stats.Desktop))
	for i, s := range stats.Desktop {
		desktopRows[i] = []string{s.Name, fmt.Sprintf("%.1f%%", s.Percentage)}
	}
	f.renderTable([]string{"OS", "PERCENTAGE"}, desktopRows)

	fmt.Fprintln(f.writer)

	f.renderSection("Mobile")
	mobileRows := make([][]string, len(stats.Mobile))
	for i, s := range stats.Mobile {
		mobileRows[i] = []string{s.Name, fmt.Sprintf("%.1f%%", s.Percentage)}
	}
	f.renderTable([]string{"OS", "PERCENTAGE"}, mobileRows)
}

// PrintExploreShows prints a list of shows from explore endpoints.
func (f *Formatter) PrintExploreShows(shows []models.ExploreShow) {
	switch f.format {
	case FormatJSON:
		f.printJSON(shows)
	case FormatPlain:
		for _, s := range shows {
			fmt.Fprintf(f.writer, "%d\t%s\n", s.ShowID, s.Title)
		}
	default:
		f.printExploreShowsTable(shows)
	}
}

func (f *Formatter) printExploreShowsTable(shows []models.ExploreShow) {
	header := []string{"ID", "TITLE", "URL"}
	rows := make([][]string, len(shows))
	for i, s := range shows {
		rows[i] = []string{
			fmt.Sprintf("%d", s.ShowID),
			truncate(s.Title, 40),
			s.SiteURL,
		}
	}
	f.renderTable(header, rows)
}


// -----------------------------------------------------------------------------
// Miscellaneous Output
// -----------------------------------------------------------------------------

func (f *Formatter) PrintCategories(categories []models.Category) {
	switch f.format {
	case FormatJSON:
		f.printJSON(categories)
	case FormatPlain:
		for _, c := range categories {
			fmt.Fprintf(f.writer, "%d\t%s\t%d\n", c.CategoryID, c.Name, c.Level)
		}
	default:
		f.printCategoriesTable(categories)
	}
}

func (f *Formatter) printCategoriesTable(categories []models.Category) {
	header := []string{"ID", "NAME", "LEVEL"}
	rows := make([][]string, len(categories))
	for i, c := range categories {
		name := c.Name
		if c.Level == 2 {
			name = "  └─ " + name
		}
		rows[i] = []string{
			fmt.Sprintf("%d", c.CategoryID),
			name,
			fmt.Sprintf("%d", c.Level),
		}
	}
	f.renderTable(header, rows)
}

func (f *Formatter) PrintGooglePlayCategories(categories []models.GooglePlayCategory) {
	switch f.format {
	case FormatJSON:
		f.printJSON(categories)
	case FormatPlain:
		for _, c := range categories {
			fmt.Fprintf(f.writer, "%d\t%s\n", c.CategoryID, c.Name)
		}
	default:
		f.printGooglePlayCategoriesTable(categories)
	}
}

func (f *Formatter) printGooglePlayCategoriesTable(categories []models.GooglePlayCategory) {
	header := []string{"ID", "NAME"}
	rows := make([][]string, len(categories))
	for i, c := range categories {
		rows[i] = []string{fmt.Sprintf("%d", c.CategoryID), c.Name}
	}
	f.renderTable(header, rows)
}

func (f *Formatter) PrintLanguages(languages []models.Language) {
	switch f.format {
	case FormatJSON:
		f.printJSON(languages)
	case FormatPlain:
		for _, l := range languages {
			fmt.Fprintf(f.writer, "%s\t%s\n", l.Code, l.Name)
		}
	default:
		f.printLanguagesTable(languages)
	}
}

func (f *Formatter) printLanguagesTable(languages []models.Language) {
	header := []string{"CODE", "LANGUAGE"}
	rows := make([][]string, len(languages))
	for i, l := range languages {
		rows[i] = []string{l.Code, l.Name}
	}
	f.renderTable(header, rows)
}


// -----------------------------------------------------------------------------
// Episode Cuepoints Output 
// -----------------------------------------------------------------------------

func (f *Formatter) PrintCuepoints(cuepoints []models.Cuepoint) {
	switch f.format {
	case FormatJSON:
		f.printJSON(cuepoints)
	case FormatPlain:
		for _, c := range cuepoints {
			fmt.Fprintf(f.writer, "%d\t%d\n", c.Timecode, c.AdsMaxCount)
		}
	default:
		f.printCuepointsTable(cuepoints)
	}
}

func (f *Formatter) printCuepointsTable(cuepoints []models.Cuepoint) {
	header := []string{"TIMECODE (ms)", "TIME", "MAX ADS"}
	rows := make([][]string, len(cuepoints))
	for i, c := range cuepoints {
		totalSeconds := c.Timecode / 1000
		minutes := totalSeconds / 60
		seconds := totalSeconds % 60
		rows[i] = []string{
			fmt.Sprintf("%d", c.Timecode),
			fmt.Sprintf("%d:%02d", minutes, seconds),
			fmt.Sprintf("%d", c.AdsMaxCount),
		}
	}
	f.renderTable(header, rows)
}

// -----------------------------------------------------------------------------
// Episode Chapters Output
// -----------------------------------------------------------------------------

func (f *Formatter) PrintChapters(chapters []models.Chapter) {
	switch f.format {
	case FormatJSON:
		f.printJSON(chapters)
	case FormatPlain:
		for _, c := range chapters {
			fmt.Fprintf(f.writer, "%d\t%d\t%s\n", c.ChapterID, c.StartsAt, c.Title)
		}
	default:
		f.printChaptersTable(chapters)
	}
}

func (f *Formatter) printChaptersTable(chapters []models.Chapter) {
	header := []string{"ID", "STARTS AT (ms)", "TIME", "TITLE", "URL"}
	rows := make([][]string, len(chapters))
	for i, c := range chapters {
		totalSeconds := c.StartsAt / 1000
		minutes := totalSeconds / 60
		seconds := totalSeconds % 60

		url := c.ExternalURL
		if url == "" {
			url = "-"
		}

		rows[i] = []string{
			fmt.Sprintf("%d", c.ChapterID),
			fmt.Sprintf("%d", c.StartsAt),
			fmt.Sprintf("%d:%02d", minutes, seconds),
			truncate(c.Title, 40),
			truncate(url, 40),
		}
	}
	f.renderTable(header, rows)
}

// -----------------------------------------------------------------------------
// Episode Messages Output
// -----------------------------------------------------------------------------

func (f *Formatter) PrintMessages(messages []models.Message) {
	switch f.format {
	case FormatJSON:
		f.printJSON(messages)
	case FormatPlain:
		for _, m := range messages {
			fmt.Fprintf(f.writer, "%d\t%s\t%s\t%s\n",
				m.MessageID,
				m.AuthorUsername,
				m.CreatedAt,
				m.Text,
			)
		}
	default:
		f.printMessagesTable(messages)
	}
}

func (f *Formatter) printMessagesTable(messages []models.Message) {
	header := []string{"ID", "AUTHOR", "DATE", "MESSAGE"}
	rows := make([][]string, len(messages))
	for i, m := range messages {
		date := m.CreatedAt
		if len(date) > 10 {
			date = date[:10]
		}

		author := m.AuthorFullname
		if m.AuthorIsOwner {
			author = author + " ★"
		}

		rows[i] = []string{
			fmt.Sprintf("%d", m.MessageID),
			truncate(author, 20),
			date,
			truncate(m.Text, 50),
		}
	}
	f.renderTable(header, rows)
}