package cmd

import (
	"fmt"
	"os"

	"github.com/khelechy/consolidate/internal/common"
	"github.com/khelechy/consolidate/internal/storage"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search command history",
	Long:  `Search through the stored command history using full-text search.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		limit, _ := cmd.Flags().GetInt("limit")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		_, err := common.InitAndGetDB()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}

		commands, err := storage.SearchCommands(query, limit)
		if err != nil {
			fmt.Printf("Error searching commands: %v\n", err)
			os.Exit(1)
		}

		if err := common.PrintCommands(commands, jsonOutput); err != nil {
			fmt.Printf("Error printing commands: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Int("limit", 10, "Maximum number of results")
	searchCmd.Flags().Bool("json", false, "Output in JSON format")
}
