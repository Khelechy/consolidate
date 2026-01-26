package cmd

import (
	"fmt"
	"os"

	"github.com/khelechy/consolidate/internal/common"
	"github.com/khelechy/consolidate/internal/storage"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Display command history",
	Long:  `Display all logged commands in chronological order (recent to oldest).`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		_, err := common.InitAndGetDB()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}

		commands, err := storage.SearchCommands("", limit) // Empty query to get all
		if err != nil {
			fmt.Printf("Error fetching history: %v\n", err)
			os.Exit(1)
		}

		if len(commands) == 0 {
			fmt.Println("No commands in history.")
			return
		}

		if err := common.PrintCommands(commands, jsonOutput); err != nil {
			fmt.Printf("Error printing commands: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().Int("limit", 100, "Maximum number of commands to display")
	historyCmd.Flags().Bool("json", false, "Output in JSON format")
}
