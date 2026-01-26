package cmd

import (
	"fmt"
	"os"

	"github.com/khelechy/consolidate/internal/common"
	"github.com/khelechy/consolidate/internal/storage"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize consolidate by setting up storage and hooks",
	Long:  `This command sets up the necessary directories, database, and shell hooks for consolidate to work.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := common.EnsureConfigDir(); err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		dbPath, err := common.GetDBPath()
		if err != nil {
			fmt.Printf("Error getting database path: %v\n", err)
			os.Exit(1)
		}

		if err := storage.InitDB(dbPath); err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Consolidate initialized successfully. Database created at %s\n", dbPath)
		fmt.Println("To activate hooks automatically, run: consolidate hook")
		fmt.Println("This will install hooks with the correct path, so consolidate doesn't need to be in PATH.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
