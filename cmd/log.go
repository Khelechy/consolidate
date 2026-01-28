package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/khelechy/consolidate/internal/common"
	"github.com/khelechy/consolidate/internal/storage"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [command]",
	Short: "Log a command to the history",
	Long:  `Manually log a command to the consolidate history database.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		sessionID, _ := cmd.Flags().GetString("session")
		cwd, _ := cmd.Flags().GetString("cwd")
		exitCodeStr, _ := cmd.Flags().GetString("exit-code")
		metadata, _ := cmd.Flags().GetString("metadata")
		encoded, _ := cmd.Flags().GetBool("encoded")
		if encoded {
			decoded, err := base64.StdEncoding.DecodeString(command)
			if err != nil {
				fmt.Printf("Error decoding command: %v\n", err)
				os.Exit(1)
			}
			command = string(decoded)
		}

		exitCode := 0
		if exitCodeStr != "" {
			if ec, err := strconv.Atoi(exitCodeStr); err == nil {
				exitCode = ec
			}
		}

		if sessionID == "" {
			sessionID = "default"
		}
		if cwd == "" {
			var err error
			cwd, err = os.Getwd()
			if err != nil {
				cwd = "unknown"
			}
		}

		_, err := common.InitAndGetDB()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}

		if err := storage.SaveCommand(command, sessionID, cwd, exitCode, metadata); err != nil {
			fmt.Printf("Error saving command: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().String("session", "", "Session ID")
	logCmd.Flags().String("cwd", "", "Current working directory")
	logCmd.Flags().String("exit-code", "0", "Exit code")
	logCmd.Flags().String("metadata", "", "Additional metadata")
	logCmd.Flags().Bool("encoded", false, "Command is base64 encoded")
}
