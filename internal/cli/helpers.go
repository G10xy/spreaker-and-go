/*
helpers.go - Shared utility functions for CLI commands
*/
package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/G10xy/spreaker-and-go/internal/api"
	"github.com/G10xy/spreaker-and-go/internal/config"
	"github.com/G10xy/spreaker-and-go/internal/output"
)

// getClient creates an API client using token from flag, env, or config.
func getClient(cmd *cobra.Command) (*api.Client, error) {
	// Try to get token from --token flag first
	token, _ := cmd.Flags().GetString("token")

	if token != "" {
		fmt.Fprintln(os.Stderr, "WARNING: Passing tokens via --token exposes them in process listings. Use SPREAKER_TOKEN env var or 'spreaker login' instead.")
	}

	// Fall back to config (which also checks env vars)
	if token == "" {
		var err error
		token, err = config.GetToken()
		if err != nil {
			return nil, err
		}
	}

	// Load config for other settings (base URL, timeout, etc.)
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return api.NewClientWithOptions(token, cfg.APIURL, 0), nil
}

// getFormatter creates an output formatter using format from flag or config.
func getFormatter(cmd *cobra.Command) *output.Formatter {
	format, _ := cmd.Flags().GetString("output")

	// Fall back to configured default
	if format == "" {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
		}
		format = cfg.OutputFormat
	}

	color := resolveColor(cmd, format)
	if !color {
		pterm.DisableColor()
	} else {
		pterm.EnableColor()
	}

	return output.New(format, color)
}

// resolveColor determines whether color output should be enabled.
func resolveColor(cmd *cobra.Command, format string) bool {
	// Only table format gets color
	if format != "" && format != "table" {
		return false
	}

	// Respect --no-color flag
	noColor, _ := cmd.Flags().GetBool("no-color")
	if noColor {
		return false
	}

	// Respect NO_COLOR env var (https://no-color.org/)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}

	// Only color when stdout is a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	return true
}


// getMyUserID returns the authenticated user's ID from cached config,
// avoiding an extra API round-trip to /v2/users/self.
func getMyUserID() (int, error) {
	return config.GetUserID()
}

func parseShowID(arg string) (int, error) {
    return parseIntArg(arg, "show ID")
}

func parseEpisodeID(arg string) (int, error) {
    return parseIntArg(arg, "episode ID")
}

func parseUserID(arg string) (int, error) {
    return parseIntArg(arg, "user ID")
}

func parseChapterID(arg string) (int, error) {
    return parseIntArg(arg, "chapter ID")
}

func parseMessageID(arg string) (int, error) {
    return parseIntArg(arg, "message ID")
}

func parseIntArg(arg string, fieldName string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(arg))
    if err != nil {
        return 0, fmt.Errorf("invalid %s: %s", fieldName, arg)
    }
    return n, nil
}

// validateFilter checks that the filter flag value is one of the allowed values.
func validateFilter(filter string) error {
	if filter == "" {
		return nil
	}
	switch filter {
	case "listenable", "editable":
		return nil
	default:
		return fmt.Errorf("invalid filter %q: must be 'listenable' or 'editable'", filter)
	}
}

// confirmAction prompts the user for confirmation.
func confirmAction(prompt string) bool {
	pterm.FgYellow.Print(prompt)
	var confirm string
	if _, err := fmt.Scanln(&confirm); err != nil {
		fmt.Fprintln(os.Stderr, "No input received, aborting.")
		return false
	}
	return confirm == "y" || confirm == "Y"
}
