/*
login.go - Authentication command

WHY A SEPARATE FILE FOR LOGIN:
Even though login is a small command, it's a distinct feature.
Keeping it separate means:
  - Easy to find when debugging auth issues
  - Can be modified without touching other commands
  - Clear ownership in a team setting
*/
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/G10xy/spreaker-and-go/internal/api"
	"github.com/G10xy/spreaker-and-go/internal/config"
)

// newLoginCmd creates the login command.
func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Spreaker",
		Long: `Authenticate with your Spreaker account.

You'll need an API token from your Spreaker developer settings.`,
		RunE: runLogin,
	}
}


func runLogin(cmd *cobra.Command, args []string) error {
	// Use plain fmt to avoid ANSI codes before color mode is resolved.
	fmt.Print("Enter your Spreaker API token: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		return fmt.Errorf("no input received")
	}
	token := strings.TrimSpace(scanner.Text())

	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Validate token by making a test API call.
	client := api.NewClient(token)
	user, err := client.GetMe()
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	if err := config.SaveToken(token, user.UserID); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	formatter := getFormatter(cmd)
	formatter.PrintSuccess(fmt.Sprintf("Logged in as %s (@%s)", user.Fullname, user.Username))
	formatter.PrintMessage(fmt.Sprintf("Token saved to %s", config.ConfigFilePath()))
	return nil
}
