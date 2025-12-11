/*
Package output handles formatting and displaying results to the terminal.

It supports multiple output formats:
  - table: Human-readable aligned columns (default)
  - json:  Machine-readable JSON output
  - plain: Simple text, one item per line

Usage:

	formatter := output.New("table")
	formatter.PrintShows(shows)

The format can be set via:
  - CLI flag: --output json
  - Config file: output_format: json
  - Environment: SPREAKER_OUTPUT=json
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

// Format represents the output format type
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatPlain Format = "plain"
)

// Formatter handles output formatting
type Formatter struct {
	format Format
	writer io.Writer
}

// New creates a new Formatter with the specified format.
// Valid formats: "table", "json", "plain"
// Defaults to "table" if invalid format is provided.
func New(format string) *Formatter {
	f := Format(strings.ToLower(format))

	// Validate format, default to table
	switch f {
	case FormatTable, FormatJSON, FormatPlain:
		// valid
	default:
		f = FormatTable
	}

	return &Formatter{
		format: f,
		writer: os.Stdout,
	}
}

// SetWriter sets a custom writer (useful for testing)
func (f *Formatter) SetWriter(w io.Writer) {
	f.writer = w
}

// -----------------------------------------------------------------------------
// User Output
// -----------------------------------------------------------------------------

// PrintUser prints a single user
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

// PrintUsers prints a list of users
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
	tw := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "ID:\t%d\n", user.UserID)
	fmt.Fprintf(tw, "Username:\t%s\n", user.Username)
	fmt.Fprintf(tw, "Name:\t%s\n", user.Fullname)
	fmt.Fprintf(tw, "Kind:\t%s\n", user.Kind)
	fmt.Fprintf(tw, "Plan:\t%s\n", user.Plan)
	fmt.Fprintf(tw, "Followers:\t%d\n", user.FollowersCount)
	fmt.Fprintf(tw, "Following:\t%d\n", user.FollowingsCount)
	fmt.Fprintf(tw, "URL:\t%s\n", user.SiteURL)

	if user.Description != "" {
		// Truncate long descriptions
		desc := user.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Fprintf(tw, "Bio:\t%s\n", desc)
	}

	tw.Flush()
}

func (f *Formatter) printUsersTable(users []models.User) {
	tw := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

	// Header
	fmt.Fprintln(tw, "ID\tUSERNAME\tNAME\tFOLLOWERS")
	fmt.Fprintln(tw, "--\t--------\t----\t---------")

	// Rows
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

// PrintShow prints a single show
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

// PrintShows prints a list of shows
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
	tw := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

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
		fmt.Fprintf(tw, "Last Episode:\t%s\n", show.LastEpisodeAt.Format("2006-01-02"))
	}

	tw.Flush()
}

func (f *Formatter) printShowsTable(shows []models.Show) {
	tw := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

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

// PrintEpisode prints a single episode
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

// PrintEpisodes prints a list of episodes
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
	tw := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

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
		fmt.Fprintf(tw, "Published:\t%s\n", episode.PublishedAt.Format("2006-01-02 15:04"))
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
	tw := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

	// Header
	fmt.Fprintln(tw, "ID\tTITLE\tDURATION\tPLAYS\tSTATUS\tPUBLISHED")
	fmt.Fprintln(tw, "--\t-----\t--------\t-----\t------\t---------")

	// Rows
	for _, e := range episodes {
		published := "-"
		if e.PublishedAt != nil {
			published = e.PublishedAt.Format("2006-01-02")
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
	tw := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Plays:\t%d\n", stats.Plays)
	fmt.Fprintf(tw, "Downloads:\t%d\n", stats.Downloads)
	fmt.Fprintf(tw, "Likes:\t%d\n", stats.Likes)
	fmt.Fprintf(tw, "Messages:\t%d\n", stats.Messages)

	tw.Flush()
}

// -----------------------------------------------------------------------------
// Generic Output
// -----------------------------------------------------------------------------

// PrintMessage prints a simple message
func (f *Formatter) PrintMessage(msg string) {
	fmt.Fprintln(f.writer, msg)
}

// PrintError prints an error message to stderr
func (f *Formatter) PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// PrintSuccess prints a success message
func (f *Formatter) PrintSuccess(msg string) {
	fmt.Fprintf(f.writer, "✓ %s\n", msg)
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

// printJSON outputs any value as formatted JSON
func (f *Formatter) printJSON(v interface{}) {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	encoder.Encode(v)
}

// truncate shortens a string to max length, adding "..." if truncated
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