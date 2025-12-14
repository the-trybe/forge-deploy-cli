package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forge-deploy",
	Short: "Generate Laravel Forge deployment configurations",
	Long: `Interactive CLI tool to generate Laravel Forge deployment configurations.

This tool helps you create GitHub Actions workflow and forge-deploy.yml
files for automated deployment to Laravel Forge.`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
