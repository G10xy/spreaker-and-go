/*
misc.go - Miscellaneous utility commands

Commands for listing categories and languages.
*/
package cli

import (
	"github.com/spf13/cobra"
)

func newMiscCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "misc",
		Aliases: []string{"miscellaneous"},
		Short:   "List categories and languages",
		Long: `List available categories and languages for shows.

These are useful reference data when creating or searching for shows.`,
	}

	cmd.AddCommand(
		newMiscCategoriesCmd(),
		newMiscGooglePlayCategoriesCmd(),
		newMiscLanguagesCmd(),
	)

	return cmd
}

// -----------------------------------------------------------------------------
// misc categories
// -----------------------------------------------------------------------------

func newMiscCategoriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "categories",
		Short: "List all show categories",
		Long: `List all available show categories.

Examples:
  spreaker misc categories
  spreaker misc categories --locale it_IT`,
		RunE: runMiscCategories,
	}

	cmd.Flags().String("locale", "", "Locale for category names (e.g., it_IT)")

	return cmd
}

func runMiscCategories(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	locale, _ := cmd.Flags().GetString("locale")
	categories, err := client.GetShowCategories(locale)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintCategories(categories)
	return nil
}

// -----------------------------------------------------------------------------
// misc googleplay-categories
// -----------------------------------------------------------------------------

func newMiscGooglePlayCategoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "googleplay-categories",
		Short: "List all Google Play podcast categories",
		RunE:  runMiscGooglePlayCategories,
	}
}

func runMiscGooglePlayCategories(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	categories, err := client.GetGooglePlayCategories()
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintGooglePlayCategories(categories)
	return nil
}

// -----------------------------------------------------------------------------
// misc languages
// -----------------------------------------------------------------------------

func newMiscLanguagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "languages",
		Short: "List all show languages",
		Long: `List all available languages for shows.

Examples:
  spreaker misc languages
  spreaker misc languages --locale it_IT`,
		RunE: runMiscLanguages,
	}

	cmd.Flags().String("locale", "", "Locale for language names (e.g., it_IT)")

	return cmd
}

func runMiscLanguages(cmd *cobra.Command, args []string) error {
	client, err := getClient(cmd)
	if err != nil {
		return err
	}

	locale, _ := cmd.Flags().GetString("locale")
	languages, err := client.GetShowLanguagesList(locale)
	if err != nil {
		return err
	}

	formatter := getFormatter(cmd)
	formatter.PrintLanguages(languages)
	return nil
}
