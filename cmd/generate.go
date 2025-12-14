package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/the-trybe/forge-deploy-cli/pkg/generators"
	"github.com/the-trybe/forge-deploy-cli/pkg/prompts"
)

var (
	outputDir        string
	workflowFilename string
	forgeConfigFile  string
	triggerBranch    string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate deployment configuration files interactively",
	Long:  `Generate deployment configuration files interactively by prompting for all configuration options.`,
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Output directory for generated files")
	generateCmd.Flags().StringVarP(&workflowFilename, "workflow-file", "w", "deploy.yml", "GitHub Actions workflow filename")
	generateCmd.Flags().StringVarP(&forgeConfigFile, "forge-config", "f", "forge-deploy.yml", "Forge deployment config filename")
	generateCmd.Flags().StringVarP(&triggerBranch, "trigger-branch", "b", "main", "Branch that triggers deployment")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Get base configuration
	config, err := prompts.PromptBaseConfig()
	if err != nil {
		return fmt.Errorf("failed to get base config: %w", err)
	}

	// Add sites
	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))

	var siteCount int
	if err := survey.AskOne(&survey.Input{
		Message: "How many sites do you want to configure?",
		Default: "1",
	}, &siteCount); err != nil {
		return fmt.Errorf("failed to get site count: %w", err)
	}

	for i := 0; i < siteCount; i++ {
		site, err := prompts.PromptCompleteSite(config.GithubBranch, i+1)
		if err != nil {
			return fmt.Errorf("failed to configure site %d: %w", i+1, err)
		}

		config.Sites = append(config.Sites, *site)

		if i < siteCount-1 {
			fmt.Println("\nSite configured successfully!")
			fmt.Println()
		}
	}

	// Validate configuration
	fmt.Println("\nValidating configuration...")
	errors := config.Validate()
	if len(errors) > 0 {
		fmt.Println("\nConfiguration validation failed:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		return fmt.Errorf("configuration validation failed")
	}

	fmt.Println("Configuration valid!")

	// Generate files
	fmt.Println("\nGenerating files...")

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate forge config file
	forgeConfigPath := filepath.Join(outputDir, forgeConfigFile)
	forgeConfig, err := generators.GenerateForgeDeployYAML(config)
	if err != nil {
		return fmt.Errorf("failed to generate forge config: %w", err)
	}

	if err := os.WriteFile(forgeConfigPath, []byte(forgeConfig), 0644); err != nil {
		return fmt.Errorf("failed to write forge config: %w", err)
	}
	fmt.Printf("  Created %s\n", forgeConfigPath)

	// Generate GitHub workflow
	workflowDir := filepath.Join(outputDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	workflowPath := filepath.Join(workflowDir, workflowFilename)
	workflow := generators.GenerateGitHubWorkflow(config, "Deploy to Forge", triggerBranch, forgeConfigFile)

	if err := os.WriteFile(workflowPath, []byte(workflow), 0644); err != nil {
		return fmt.Errorf("failed to write workflow: %w", err)
	}
	fmt.Printf("  Created %s\n", workflowPath)

	// Success message
	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Configuration generated successfully!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Review the generated files")
	fmt.Println("  2. Add required secrets to your GitHub repository:")
	fmt.Println("     Settings > Secrets and variables > Actions")
	fmt.Println("  3. Commit and push the files to your repository")
	fmt.Printf("  4. Push to '%s' branch to trigger deployment\n", triggerBranch)
	fmt.Println()
	fmt.Println("Happy deploying!")

	return nil
}
