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
}

// New creates a new Formatter with the specified format.
func New(format string) *Formatter {
	f := Format(strings.ToLower(strings.TrimSpace(format)))

	switch f {
	case FormatTable, FormatJSON, FormatPlain:
	default:
		f = FormatTable
	}

	return &Formatter{
		format: f,
		writer: os.Stdout,
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
	tw := f.tabw()

	fmt.Fprintf(tw, "ID:\t%d\n", user.UserID)
	fmt.Fprintf(tw, "Username:\t%s\n", user.Username)
	fmt.Fprintf(tw, "Name:\t%s\n", user.Fullname)
	fmt.Fprintf(tw, "Kind:\t%s\n", user.Kind)
	fmt.Fprintf(tw, "Plan:\t%s\n", user.Plan)
	fmt.Fprintf(tw, "Followers:\t%d\n", user.FollowersCount)
	fmt.Fprintf(tw, "Following:\t%d\n", user.FollowingsCount)
	fmt.Fprintf(tw, "URL:\t%s\n", user.SiteURL)

	if user.Description != "" {
		desc := user.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Fprintf(tw, "Bio:\t%s\n", desc)
	}

	tw.Flush()
}

func (f *Formatter) printUsersTable(users []models.User) {
	tw := f.tabw()

	fmt.Fprintln(tw, "ID\tUSERNAME\tNAME\tFOLLOWERS")
	fmt.Fprintln(tw, "--\t--------\t----\t---------")

	for _, u := range users {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%d\n",
			u.UserID,
			u.Username,
			truncate(u.Fullname, 30),
			u.FollowersCount,
		)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintf(tw, "ID:\t%d\n", show.ShowID)
	fmt.Fprintf(tw, "Title:\t%s\n", show.Title)
	fmt.Fprintf(tw, "Language:\t%s\n", show.Language)
	fmt.Fprintf(tw, "Episodes:\t%d\n", show.EpisodesCount)
	fmt.Fprintf(tw, "Followers:\t%d\n", show.FollowersCount)
	fmt.Fprintf(tw, "Plays:\t%d\n", show.PlayCount)
	fmt.Fprintf(tw, "Explicit:\t%v\n", show.Explicit)
	fmt.Fprintf(tw, "URL:\t%s\n", show.SiteURL)

	if show.Description != "" {
		desc := show.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Fprintf(tw, "Description:\t%s\n", desc)
	}

	if show.LastEpisodeAt != nil {
		fmt.Fprintf(tw, "Last Episode:\t%s\n", show.LastEpisodeAt.Format(time.DateTime))
	}

	tw.Flush()
}

func (f *Formatter) printShowsTable(shows []models.Show) {
	tw := f.tabw()

	// Header
	fmt.Fprintln(tw, "ID\tTITLE\tEPISODES\tFOLLOWERS\tPLAYS")
	fmt.Fprintln(tw, "--\t-----\t--------\t---------\t-----")

	// Rows
	for _, s := range shows {
		fmt.Fprintf(tw, "%d\t%s\t%d\t%d\t%d\n",
			s.ShowID,
			truncate(s.Title, 40),
			s.EpisodesCount,
			s.FollowersCount,
			s.PlayCount,
		)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintf(tw, "ID:\t%d\n", episode.EpisodeID)
	fmt.Fprintf(tw, "Title:\t%s\n", episode.Title)
	fmt.Fprintf(tw, "Show ID:\t%d\n", episode.ShowID)
	fmt.Fprintf(tw, "Duration:\t%s\n", formatDuration(episode.Duration))
	fmt.Fprintf(tw, "Plays:\t%d\n", episode.PlayCount)
	fmt.Fprintf(tw, "Likes:\t%d\n", episode.LikesCount)
	fmt.Fprintf(tw, "Status:\t%s\n", episode.EncodingStatus)
	fmt.Fprintf(tw, "Explicit:\t%v\n", episode.Explicit)
	fmt.Fprintf(tw, "Downloads:\t%v\n", episode.DownloadEnabled)
	fmt.Fprintf(tw, "URL:\t%s\n", episode.SiteURL)

	if episode.PublishedAt != nil {
		fmt.Fprintf(tw, "Published:\t%s\n", episode.PublishedAt.Format(time.DateTime))
	}

	if len(episode.Tags) > 0 {
		fmt.Fprintf(tw, "Tags:\t%s\n", strings.Join(episode.Tags, ", "))
	}

	if episode.Description != "" {
		desc := episode.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Fprintf(tw, "Description:\t%s\n", desc)
	}

	tw.Flush()
}

func (f *Formatter) printEpisodesTable(episodes []models.Episode) {
	tw := f.tabw()

	// Header
	fmt.Fprintln(tw, "ID\tTITLE\tDURATION\tPLAYS\tSTATUS\tPUBLISHED")
	fmt.Fprintln(tw, "--\t-----\t--------\t-----\t------\t---------")

	// Rows
	for _, e := range episodes {
		published := "-"
		if e.PublishedAt != nil {
			published = e.PublishedAt.Format(time.DateTime)
		}

		fmt.Fprintf(tw, "%d\t%s\t%s\t%d\t%s\t%s\n",
			e.EpisodeID,
			truncate(e.Title, 35),
			formatDuration(e.Duration),
			e.PlayCount,
			e.EncodingStatus,
			published,
		)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintf(tw, "Plays:\t%d\n", stats.Plays)
	fmt.Fprintf(tw, "Downloads:\t%d\n", stats.Downloads)
	fmt.Fprintf(tw, "Likes:\t%d\n", stats.Likes)
	fmt.Fprintf(tw, "Messages:\t%d\n", stats.Messages)

	tw.Flush()
}

// -----------------------------------------------------------------------------
// Generic Output
// -----------------------------------------------------------------------------

func (f *Formatter) PrintMessage(msg string) {
	fmt.Fprintln(f.writer, msg)
}

func (f *Formatter) PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

func (f *Formatter) PrintSuccess(msg string) {
	fmt.Fprintf(f.writer, "✓ %s\n", msg)
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
	tw := f.tabw()

	fmt.Fprintln(tw, "=== Overall Statistics ===")
	fmt.Fprintf(tw, "Total Plays:\t%d\n", stats.PlaysCount)
	fmt.Fprintf(tw, "  On Demand:\t%d\n", stats.PlaysOndemandCount)
	fmt.Fprintf(tw, "  Live:\t%d\n", stats.PlaysLiveCount)
	fmt.Fprintf(tw, "Downloads:\t%d\n", stats.DownloadsCount)
	fmt.Fprintf(tw, "Likes:\t%d\n", stats.LikesCount)
	fmt.Fprintf(tw, "Followers:\t%d\n", stats.FollowersCount)
	fmt.Fprintf(tw, "Shows:\t%d\n", stats.ShowsCount)
	fmt.Fprintf(tw, "Episodes:\t%d\n", stats.EpisodesCount)

	tw.Flush()
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
	tw := f.tabw()

	if stats.Title != "" {
		fmt.Fprintf(tw, "Show:\t%s\n", stats.Title)
	}
	fmt.Fprintln(tw, "=== Overall Statistics ===")
	fmt.Fprintf(tw, "Total Plays:\t%d\n", stats.PlaysCount)
	fmt.Fprintf(tw, "  On Demand:\t%d\n", stats.PlaysOndemandCount)
	fmt.Fprintf(tw, "  Live:\t%d\n", stats.PlaysLiveCount)
	fmt.Fprintf(tw, "Downloads:\t%d\n", stats.DownloadsCount)
	fmt.Fprintf(tw, "Likes:\t%d\n", stats.LikesCount)
	fmt.Fprintf(tw, "Episodes:\t%d\n", stats.EpisodesCount)

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "=== Overall Statistics ===")
	fmt.Fprintf(tw, "Total Plays:\t%d\n", stats.PlaysCount)
	fmt.Fprintf(tw, "  On Demand:\t%d\n", stats.PlaysOndemandCount)
	fmt.Fprintf(tw, "  Live:\t%d\n", stats.PlaysLiveCount)
	fmt.Fprintf(tw, "Downloads:\t%d\n", stats.DownloadsCount)
	fmt.Fprintf(tw, "Likes:\t%d\n", stats.LikesCount)
	fmt.Fprintf(tw, "Messages:\t%d\n", stats.MessagesCount)
	fmt.Fprintf(tw, "Chapters:\t%d\n", stats.ChaptersCount)

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "DATE\tPLAYS\tON DEMAND\tLIVE\tDOWNLOADS")
	fmt.Fprintln(tw, "----\t-----\t---------\t----\t---------")

	for _, s := range stats {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\n",
			s.Date, s.PlaysCount, s.PlaysOndemandCount, s.PlaysLiveCount, s.DownloadsCount)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "DEVICE\tPERCENTAGE")
	fmt.Fprintln(tw, "------\t----------")

	for _, s := range stats {
		fmt.Fprintf(tw, "%s\t%.1f%%\n", s.Name, s.Percentage)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "=== By Country ===")
	fmt.Fprintln(tw, "COUNTRY\tPERCENTAGE")
	fmt.Fprintln(tw, "-------\t----------")
	for _, c := range stats.Country {
		fmt.Fprintf(tw, "%s\t%.1f%%\n", c.Name, c.Percentage)
	}

	fmt.Fprintln(tw, "")
	fmt.Fprintln(tw, "=== By City ===")
	fmt.Fprintln(tw, "CITY\tPERCENTAGE")
	fmt.Fprintln(tw, "----\t----------")
	for _, c := range stats.City {
		fmt.Fprintf(tw, "%s\t%.1f%%\n", c.Name, c.Percentage)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "SOURCE\tPLAYS\tPERCENTAGE")
	fmt.Fprintln(tw, "------\t-----\t----------")

	for _, s := range stats.Overall {
		fmt.Fprintf(tw, "%s\t%d\t%d%%\n", s.Name, s.PlaysCount, s.Percentage)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "DATE\tLISTENERS")
	fmt.Fprintln(tw, "----\t---------")

	for _, s := range stats {
		fmt.Fprintf(tw, "%s\t%d\n", s.Date, s.ListenersCount)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "ID\tTITLE\tURL")
	fmt.Fprintln(tw, "--\t-----\t---")

	for _, s := range shows {
		fmt.Fprintf(tw, "%d\t%s\t%s\n",
			s.ShowID,
			truncate(s.Title, 40),
			s.SiteURL,
		)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "ID\tNAME\tLEVEL")
	fmt.Fprintln(tw, "--\t----\t-----")

	for _, c := range categories {
		levelIndicator := ""
		if c.Level == 2 {
			levelIndicator = "  └─ " // Indent subcategories
		}
		fmt.Fprintf(tw, "%d\t%s%s\t%d\n", c.CategoryID, levelIndicator, c.Name, c.Level)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "ID\tNAME")
	fmt.Fprintln(tw, "--\t----")

	for _, c := range categories {
		fmt.Fprintf(tw, "%d\t%s\n", c.CategoryID, c.Name)
	}

	tw.Flush()
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
	tw := f.tabw()

	fmt.Fprintln(tw, "CODE\tLANGUAGE")
	fmt.Fprintln(tw, "----\t--------")

	for _, l := range languages {
		fmt.Fprintf(tw, "%s\t%s\n", l.Code, l.Name)
	}

	tw.Flush()
}