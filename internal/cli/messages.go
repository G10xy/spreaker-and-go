/*
messages.go - Episode message management commands

Messages are public textual comments that listeners can leave on episodes
to communicate with the author.
*/
package cli

import (
	"github.com/spf13/cobra"
	
	"github.com/G10xy/spreaker-and-go/internal/api"
)

func newMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "messages",
		Aliases: []string{"message", "msg"},
		Short:   "Manage episode messages",
		Long: `Manage messages for episodes. Messages are public textual comments 
that listeners can leave on episodes to communicate with the author.

Examples:
  spreaker messages list 12345
  spreaker messages create 12345 "Great episode!"
  spreaker messages delete 12345 67890
  spreaker messages report 12345 67890`,
	}

	cmd.AddCommand(
		newMessagesListCmd(),
		newMessagesCreateCmd(),
		newMessagesDeleteCmd(),
		newMessagesReportCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// messages list
// -----------------------------------------------------------------------------

func newMessagesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <episode-id>",
		Short: "List all messages for an episode",
		Long: `List all messages for an episode, sorted from newest to oldest.

Examples:
  spreaker messages list 12345 --limit 50
  spreaker msg list 12345 --output json`,
		Args: cobra.ExactArgs(1),
		RunE: runMessagesList,
	}

	cmd.Flags().IntP("limit", "l", 20, "Maximum number of messages")

	return cmd
}

func runMessagesList(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	result, err := client.GetEpisodeMessages(episodeID, api.PaginationParams{Limit: limit})
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)

	if len(result.Items) == 0 {
		formatter.PrintMessage("No messages found for this episode.")
		return nil
	}

	formatter.PrintMessages(result.Items)

	if result.HasMore {
		formatter.PrintMessage("\n(more messages available, use --limit to see more)")
	}

	return nil
}

// -----------------------------------------------------------------------------
// messages create
// -----------------------------------------------------------------------------

func newMessagesCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <episode-id> <text>",
		Short: "Leave a message on an episode",
		Long: `Leave a public message on an episode.
The message text can be up to 4000 characters.

Examples:
  spreaker messages create 12345 "Great episode!"
  spreaker msg create 12345 "Thanks for the insights, very helpful!"`,
		Args: cobra.ExactArgs(2),
		RunE: runMessagesCreate,
	}
}

func runMessagesCreate(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}
	text := args[1]

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	if err := client.CreateMessage(episodeID, text); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage("Message sent successfully.")
	return nil
}

// -----------------------------------------------------------------------------
// messages delete
// -----------------------------------------------------------------------------

func newMessagesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <episode-id> <message-id>",
		Short: "Delete a message from an episode",
		Long: `Delete a message from an episode.
You can only delete messages that you sent, or messages on your own episodes.

Examples:
  spreaker messages delete 12345 67890
  spreaker msg delete 12345 67890`,
		Args: cobra.ExactArgs(2),
		RunE: runMessagesDelete,
	}
}

func runMessagesDelete(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	messageID, err := parseMessageID(args[1])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	if err := client.DeleteMessage(episodeID, messageID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage("Message deleted successfully.")
	return nil
}

// -----------------------------------------------------------------------------
// messages report
// -----------------------------------------------------------------------------

func newMessagesReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "report <episode-id> <message-id>",
		Short: "Report a message as spam or abuse",
		Long: `Report a message as spam or as violating Spreaker's terms and conditions.
Reported messages will be reviewed by Spreaker's staff within 1 working day.

Examples:
  spreaker messages report 12345 67890
  spreaker msg report 12345 67890`,
		Args: cobra.ExactArgs(2),
		RunE: runMessagesReport,
	}
}

func runMessagesReport(cmd *cobra.Command, args []string) error {
	episodeID, err := parseEpisodeID(args[0])
	if err != nil {
		return err
	}

	messageID, err := parseMessageID(args[1])
	if err != nil {
		return err
	}

	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	if err := client.ReportMessageAbuse(episodeID, messageID); err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintMessage("Message reported successfully. Spreaker staff will review it within 1 working day.")
	return nil
}
