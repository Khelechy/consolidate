package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/khelechy/consolidate/internal/common"
	"github.com/spf13/cobra"
)

// hookCmd represents the hook command
var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Install shell hooks for automatic command logging",
	Long:  `Automatically install the appropriate hook script for your shell to enable automatic command logging.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the executable path
		execPath, err := os.Executable()
		if err != nil {
			fmt.Printf("Error getting executable path: %v\n", err)
			os.Exit(1)
		}
		execDir := filepath.Dir(execPath)

		// Detect shell
		shell := common.DetectShell()
		var profilePath string
		var hookLine string

		switch shell {
		case "bash", "zsh":
			homeDir, _ := os.UserHomeDir()
			if shell == "zsh" {
				profilePath = filepath.Join(homeDir, ".zshrc")
			} else {
				profilePath = filepath.Join(homeDir, ".bashrc")
			}
			hookLine = fmt.Sprintf("export CONSOLIDATE_BIN='%s'; source %s/scripts/hook.sh", execPath, execDir)
		case "powershell":
			// For PowerShell, get the profile path
			cmd := exec.Command("powershell", "-Command", "$PROFILE")
			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("Error getting PowerShell profile: %v\n", err)
				os.Exit(1)
			}
			profilePath = strings.TrimSpace(string(output))
			hookLine = fmt.Sprintf("$env:CONSOLIDATE_BIN='%s'; . %s/scripts/hook.ps1", execPath, execDir)
		default:
			fmt.Printf("Unsupported shell: %s\n", shell)
			os.Exit(1)
		}

		// Check if profile exists, create if not
		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			file, err := os.Create(profilePath)
			if err != nil {
				fmt.Printf("Error creating profile file: %v\n", err)
				os.Exit(1)
			}
			file.Close()
		}

		// Read existing profile
		content, err := os.ReadFile(profilePath)
		if err != nil {
			fmt.Printf("Error reading profile file: %v\n", err)
			os.Exit(1)
		}

		// Check if hook is already installed
		if strings.Contains(string(content), hookLine) {
			fmt.Println("Hook already installed.")
			return
		}

		// Append the hook line
		newContent := string(content) + "\n" + hookLine + "\n"
		err = os.WriteFile(profilePath, []byte(newContent), 0644)
		if err != nil {
			fmt.Printf("Error writing to profile file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Hook installed successfully. Restart your shell or run 'source %s' to activate.\n", profilePath)
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
