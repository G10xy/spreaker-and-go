/*
helpers.go - Shared utility functions for CLI commands
*/
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/internal/api"
	"github.com/G10xy/spreaker-and-go/internal/config"
	"github.com/G10xy/spreaker-and-go/internal/output"
)

// getClient creates an API client using token from flag, env, or config.
func getClient(cmd *cobra.Command) (*api.Client, error) {
	// Try to get token from --token flag first
	token, _ := cmd.Flags().GetString("token")

	// Fall back to config (which also checks env vars)
	if token == "" {
		var err error
		token, err = config.GetToken()
		if err != nil {
			return nil, err
		}
	}

	// Load config for other settings (base URL, timeout, etc.)
	cfg, _ := config.Load()

	return api.NewClientWithOptions(token, cfg.APIURL, 0), nil
}

// getFormatter creates an output formatter using format from flag or config.
func getFormatter(cmd *cobra.Command) *output.Formatter {
	format, _ := cmd.Flags().GetString("output")

	// Fall back to configured default
	if format == "" {
		cfg, _ := config.Load()
		format = cfg.OutputFormat
	}

	return output.New(format)
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

// confirmAction prompts the user for confirmation.
func confirmAction(prompt string) bool {
	fmt.Print(prompt)
	var confirm string
	fmt.Scanln(&confirm)
	return confirm == "y" || confirm == "Y"
}
