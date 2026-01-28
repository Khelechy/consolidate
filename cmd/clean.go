package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/khelechy/consolidate/internal/common"
	"github.com/khelechy/consolidate/internal/storage"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean command history",
	Long:  `Remove commands from the history database based on specified criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		fromStr, _ := cmd.Flags().GetString("from")
		toStr, _ := cmd.Flags().GetString("to")
		all, _ := cmd.Flags().GetBool("all")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		_, err := common.InitAndGetDB()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}

		// Validate flags - cannot use --all with --from or --to
		if all && (fromStr != "" || toStr != "") {
			fmt.Printf("Error: cannot use --all with --from or --to flags\n")
			os.Exit(1)
		}

		// Parse datetime ranges
		var fromTime, toTime *time.Time
		if fromStr != "" {
			ft, err := parseDateTime(fromStr, true) // true for from (start of day)
			if err != nil {
				fmt.Printf("Error parsing from datetime: %v\n", err)
				os.Exit(1)
			}
			fromTime = &ft
		}
		if toStr != "" {
			tt, err := parseDateTime(toStr, false) // false for to (end of day)
			if err != nil {
				fmt.Printf("Error parsing to datetime: %v\n", err)
				os.Exit(1)
			}
			toTime = &tt
		}

		// Perform the clean operation
		deleted, err := storage.CleanHistory(fromTime, toTime, all, dryRun)
		if err != nil {
			fmt.Printf("Error cleaning history: %v\n", err)
			os.Exit(1)
		}

		if dryRun {
			fmt.Printf("Dry run: Would delete %d commands\n", deleted)
		} else {
			fmt.Printf("Deleted %d commands from history\n", deleted)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().String("from", "", "Start datetime (RFC3339 or YYYY-MM-DD, e.g., 2023-01-01 or 2023-01-01T00:00:00Z)")
	cleanCmd.Flags().String("to", "", "End datetime (RFC3339 or YYYY-MM-DD, e.g., 2023-12-31 or 2023-12-31T23:59:59Z)")
	cleanCmd.Flags().Bool("all", false, "Delete all commands from history")
	cleanCmd.Flags().Bool("dry-run", false, "Show what would be deleted without actually deleting")
}

// parseDateTime parses a datetime string, accepting both RFC3339 and date-only formats
func parseDateTime(input string, isFrom bool) (time.Time, error) {
	// First try RFC3339 format
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}

	// If that fails, try date-only format (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", input); err == nil {
		if isFrom {
			// For --from, use start of day
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
		} else {
			// For --to, use end of day
			return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid datetime format: %s (use RFC3339 or YYYY-MM-DD)", input)
}
