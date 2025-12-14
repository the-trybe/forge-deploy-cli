package prompts

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/the-trybe/forge-deploy-cli/pkg/models"
)

// PromptBaseConfig prompts for base deployment configuration
func PromptBaseConfig() (*models.DeploymentConfig, error) {
	fmt.Println("\nLaravel Forge Deployment Configuration Generator")
	fmt.Println()
	fmt.Println("Base Configuration")
	fmt.Println(strings.Repeat("-", 50))

	config := &models.DeploymentConfig{}

	questions := []*survey.Question{
		{
			Name:     "organization",
			Prompt:   &survey.Input{Message: "Forge organization name:"},
			Validate: survey.Required,
		},
		{
			Name:     "server",
			Prompt:   &survey.Input{Message: "Forge server name:"},
			Validate: survey.Required,
		},
		{
			Name:     "repository",
			Prompt:   &survey.Input{Message: "GitHub repository (owner/repo):"},
			Validate: survey.Required,
		},
		{
			Name:     "branch",
			Prompt:   &survey.Input{Message: "Default branch:", Default: "main"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Organization string
		Server       string
		Repository   string
		Branch       string
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return nil, err
	}

	config.Organization = answers.Organization
	config.Server = answers.Server
	config.GithubRepository = answers.Repository
	config.GithubBranch = answers.Branch

	return config, nil
}

// PromptSiteBasicInfo prompts for basic site information
func PromptSiteBasicInfo(siteNumber int) (map[string]interface{}, error) {
	fmt.Printf("\nSite %d Configuration\n", siteNumber)
	fmt.Println(strings.Repeat("-", 50))

	var name, domainMode, wwwRedirect string

	questions := []*survey.Question{
		{
			Name: "domainMode",
			Prompt: &survey.Select{
				Message: "Domain mode:",
				Options: []string{"on-forge", "custom"},
				Default: "on-forge",
			},
		},
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Site name:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Name       string
		DomainMode string
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return nil, err
	}

	name = answers.Name
	domainMode = answers.DomainMode

	domainPreview := name
	if domainMode == "on-forge" {
		domainPreview = name + ".on-forge.com"
	}
	fmt.Printf("  -> Domain will be: %s\n", domainPreview)

	wwwQuestion := &survey.Select{
		Message: "WWW redirect type:",
		Options: []string{"none", "from-www", "to-www"},
		Default: "none",
	}
	if err := survey.AskOne(wwwQuestion, &wwwRedirect); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":              name,
		"domain_mode":       domainMode,
		"www_redirect_type": wwwRedirect,
	}, nil
}

// PromptSiteRepositorySettings prompts for repository settings
func PromptSiteRepositorySettings(defaultBranch string) (map[string]interface{}, error) {
	fmt.Println("\nRepository Settings")

	var useCustomBranch bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Use different branch than default (%s)?", defaultBranch),
		Default: false,
	}, &useCustomBranch); err != nil {
		return nil, err
	}

	var githubBranch string
	if useCustomBranch {
		if err := survey.AskOne(&survey.Input{
			Message: "Branch name:",
		}, &githubBranch, survey.WithValidator(survey.Required)); err != nil {
			return nil, err
		}
	}

	var rootDir, webDir string
	var cloneRepo bool

	questions := []*survey.Question{
		{
			Name:   "rootDir",
			Prompt: &survey.Input{Message: "Root directory:", Default: "."},
		},
		{
			Name:   "webDir",
			Prompt: &survey.Input{Message: "Public/web directory:", Default: "public"},
		},
		{
			Name:   "cloneRepo",
			Prompt: &survey.Confirm{Message: "Clone repository during site creation?", Default: true},
		},
	}

	answers := struct {
		RootDir   string
		WebDir    string
		CloneRepo bool
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return nil, err
	}

	rootDir = answers.RootDir
	webDir = answers.WebDir
	cloneRepo = answers.CloneRepo

	return map[string]interface{}{
		"github_branch":    githubBranch,
		"root_dir":         rootDir,
		"web_dir":          webDir,
		"clone_repository": cloneRepo,
	}, nil
}

// PromptSitePHPSettings prompts for PHP settings
func PromptSitePHPSettings() (map[string]interface{}, error) {
	fmt.Println("\nPHP Settings")

	var projectType, phpVersion string
	var installComposer bool

	if err := survey.AskOne(&survey.Select{
		Message: "Project type:",
		Options: []string{"laravel", "other"},
		Default: "laravel",
	}, &projectType); err != nil {
		return nil, err
	}

	var usePHPVersion bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Specify PHP version?",
		Default: false,
	}, &usePHPVersion); err != nil {
		return nil, err
	}

	if usePHPVersion {
		if err := survey.AskOne(&survey.Input{
			Message: "PHP version (e.g., php81, php82, php83, php84):",
			Default: "php84",
		}, &phpVersion); err != nil {
			return nil, err
		}
	}

	if err := survey.AskOne(&survey.Confirm{
		Message: "Install Composer dependencies during site creation?",
		Default: false,
	}, &installComposer); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"project_type":                  projectType,
		"php_version":                   phpVersion,
		"install_composer_dependencies": installComposer,
	}, nil
}

// PromptDeploymentScript prompts for deployment script
func PromptDeploymentScript() (string, error) {
	fmt.Println("\nDeployment Script")

	var addScript bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add custom deployment script?",
		Default: false,
	}, &addScript); err != nil {
		return "", err
	}

	if !addScript {
		return "", nil
	}

	var script string
	if err := survey.AskOne(&survey.Multiline{
		Message: "Enter deployment script:",
	}, &script); err != nil {
		return "", err
	}

	return script, nil
}

// PromptEnvironmentVariables prompts for environment variables
func PromptEnvironmentVariables() (string, string, error) {
	fmt.Println("\nEnvironment Variables")

	var envChoice string
	if err := survey.AskOne(&survey.Select{
		Message: "Environment configuration:",
		Options: []string{"none", "inline", "file"},
		Default: "none",
	}, &envChoice); err != nil {
		return "", "", err
	}

	switch envChoice {
	case "inline":
		var useTemplate bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Copy environment variables from a template file (e.g., .env.example)?",
			Default: false,
		}, &useTemplate); err != nil {
			return "", "", err
		}

		var envVars string
		if useTemplate {
			var templatePath string
			if err := survey.AskOne(&survey.Input{
				Message: "Path to template file (relative to repository root):",
				Default: ".env.example",
			}, &templatePath, survey.WithValidator(survey.Required)); err != nil {
				return "", "", err
			}

			// Read the template file
			templateContent, err := readFileContent(templatePath)
			if err != nil {
				fmt.Printf("Warning: Could not read template file: %v\n", err)
				fmt.Println("You can still enter variables manually below.")
			} else {
				envVars = templateContent
				fmt.Printf("Loaded %d lines from %s\n", len(strings.Split(templateContent, "\n")), templatePath)
			}
		}

		// Allow editing or manual entry
		if err := survey.AskOne(&survey.Multiline{
			Message: "Enter/edit environment variables:",
			Default: envVars,
		}, &envVars); err != nil {
			return "", "", err
		}
		return envVars, "", nil
	case "file":
		var envFile string
		if err := survey.AskOne(&survey.Input{
			Message: "Path to .env file (relative to repository root):",
		}, &envFile, survey.WithValidator(survey.Required)); err != nil {
			return "", "", err
		}
		return "", envFile, nil
	}

	return "", "", nil
}

// readFileContent reads a file and returns its content
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// PromptProcesses prompts for background processes
func PromptProcesses() ([]models.Process, error) {
	fmt.Println("\nBackground Processes")

	var addProcesses bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add background processes?",
		Default: false,
	}, &addProcesses); err != nil {
		return nil, err
	}

	if !addProcesses {
		return nil, nil
	}

	var processes []models.Process

	for {
		var name, command string

		questions := []*survey.Question{
			{
				Name:     "name",
				Prompt:   &survey.Input{Message: "Process name:"},
				Validate: survey.Required,
			},
			{
				Name:     "command",
				Prompt:   &survey.Input{Message: "Process command:"},
				Validate: survey.Required,
			},
		}

		answers := struct {
			Name    string
			Command string
		}{}

		if err := survey.Ask(questions, &answers); err != nil {
			return nil, err
		}

		name = answers.Name
		command = answers.Command

		processes = append(processes, models.Process{
			Name:    name,
			Command: command,
		})

		var addAnother bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add another process?",
			Default: false,
		}, &addAnother); err != nil {
			return nil, err
		}

		if !addAnother {
			break
		}
	}

	return processes, nil
}

// PromptScheduler prompts for Laravel scheduler
func PromptScheduler() (bool, error) {
	fmt.Println("\nLaravel Scheduler")

	var enabled bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Enable Laravel scheduler?",
		Default: false,
	}, &enabled); err != nil {
		return false, err
	}

	return enabled, nil
}

// PromptAliases prompts for domain aliases
func PromptAliases() ([]string, error) {
	fmt.Println("\nDomain Aliases")

	var addAliases bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Add domain aliases?",
		Default: false,
	}, &addAliases); err != nil {
		return nil, err
	}

	if !addAliases {
		return nil, nil
	}

	var aliases []string

	for {
		var alias string
		if err := survey.AskOne(&survey.Input{
			Message: "Alias domain:",
		}, &alias, survey.WithValidator(survey.Required)); err != nil {
			return nil, err
		}

		aliases = append(aliases, alias)

		var addAnother bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add another alias?",
			Default: false,
		}, &addAnother); err != nil {
			return nil, err
		}

		if !addAnother {
			break
		}
	}

	return aliases, nil
}

// PromptNginxConfig prompts for Nginx configuration
func PromptNginxConfig() (map[string]interface{}, error) {
	fmt.Println("\nNginx Configuration")

	var configChoice string
	if err := survey.AskOne(&survey.Select{
		Message: "Nginx configuration:",
		Options: []string{"default", "template", "custom-file"},
		Default: "default",
	}, &configChoice); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	switch configChoice {
	case "template":
		var templateName string
		if err := survey.AskOne(&survey.Input{
			Message: "Template name:",
		}, &templateName, survey.WithValidator(survey.Required)); err != nil {
			return nil, err
		}
		result["nginx_template"] = templateName

		var addVars bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add template variables?",
			Default: false,
		}, &addVars); err != nil {
			return nil, err
		}

		if addVars {
			variables := make(map[string]string)
			for {
				questions := []*survey.Question{
					{
						Name:     "key",
						Prompt:   &survey.Input{Message: "Variable name:"},
						Validate: survey.Required,
					},
					{
						Name:     "value",
						Prompt:   &survey.Input{Message: "Variable value:"},
						Validate: survey.Required,
					},
				}

				answers := struct {
					Key   string
					Value string
				}{}

				if err := survey.Ask(questions, &answers); err != nil {
					return nil, err
				}

				variables[answers.Key] = answers.Value

				var addAnother bool
				if err := survey.AskOne(&survey.Confirm{
					Message: "Add another variable?",
					Default: false,
				}, &addAnother); err != nil {
					return nil, err
				}

				if !addAnother {
					break
				}
			}
			result["nginx_template_variables"] = variables
		}

	case "custom-file":
		var customConfig string
		if err := survey.AskOne(&survey.Input{
			Message: "Path to custom nginx config (relative to repository root):",
		}, &customConfig, survey.WithValidator(survey.Required)); err != nil {
			return nil, err
		}
		result["nginx_custom_config"] = customConfig
	}

	return result, nil
}

// PromptSSLCertificate prompts for SSL certificate
func PromptSSLCertificate() (bool, error) {
	fmt.Println("\nSSL Certificate")

	var enabled bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Create SSL certificate?",
		Default: false,
	}, &enabled); err != nil {
		return false, err
	}

	return enabled, nil
}

// PromptIsolation prompts for site isolation
func PromptIsolation() (map[string]interface{}, error) {
	fmt.Println("\nSite Isolation")

	var isolated bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Run as isolated user?",
		Default: false,
	}, &isolated); err != nil {
		return nil, err
	}

	var isolatedUser string
	if isolated {
		if err := survey.AskOne(&survey.Input{
			Message: "Isolated user name:",
		}, &isolatedUser, survey.WithValidator(survey.Required)); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"isolated":      isolated,
		"isolated_user": isolatedUser,
	}, nil
}

// PromptZeroDowntime prompts for zero-downtime deployment settings
func PromptZeroDowntime() (map[string]interface{}, error) {
	fmt.Println("\nZero-Downtime Deployment")

	var zeroDowntime bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Enable zero-downtime deployments?",
		Default: false,
	}, &zeroDowntime); err != nil {
		return nil, err
	}

	var sharedPaths []models.SharedPath

	if zeroDowntime {
		var addPaths bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Add shared paths?",
			Default: true,
		}, &addPaths); err != nil {
			return nil, err
		}

		if addPaths {
			for {
				var pathType string
				if err := survey.AskOne(&survey.Select{
					Message: "Path type:",
					Options: []string{"simple", "custom"},
					Default: "simple",
				}, &pathType); err != nil {
					return nil, err
				}

				if pathType == "simple" {
					var path string
					if err := survey.AskOne(&survey.Input{
						Message: "Path:",
					}, &path, survey.WithValidator(survey.Required)); err != nil {
						return nil, err
					}
					sharedPaths = append(sharedPaths, models.SharedPath{From: path})
				} else {
					var fromPath, toPath string

					questions := []*survey.Question{
						{
							Name:     "from",
							Prompt:   &survey.Input{Message: "From path:"},
							Validate: survey.Required,
						},
						{
							Name:     "to",
							Prompt:   &survey.Input{Message: "To path:"},
							Validate: survey.Required,
						},
					}

					answers := struct {
						From string
						To   string
					}{}

					if err := survey.Ask(questions, &answers); err != nil {
						return nil, err
					}

					fromPath = answers.From
					toPath = answers.To

					sharedPaths = append(sharedPaths, models.SharedPath{From: fromPath, To: toPath})
				}

				var addAnother bool
				if err := survey.AskOne(&survey.Confirm{
					Message: "Add another shared path?",
					Default: false,
				}, &addAnother); err != nil {
					return nil, err
				}

				if !addAnother {
					break
				}
			}
		}
	}

	return map[string]interface{}{
		"zero_downtime_deployments": zeroDowntime,
		"shared_paths":              sharedPaths,
	}, nil
}

// PromptCompleteSite orchestrates all site prompts
func PromptCompleteSite(defaultBranch string, siteNumber int) (*models.SiteConfig, error) {
	// Basic info
	basicInfo, err := PromptSiteBasicInfo(siteNumber)
	if err != nil {
		return nil, err
	}

	// Repository settings
	repoSettings, err := PromptSiteRepositorySettings(defaultBranch)
	if err != nil {
		return nil, err
	}

	// PHP settings
	phpSettings, err := PromptSitePHPSettings()
	if err != nil {
		return nil, err
	}

	// Deployment script
	deploymentScript, err := PromptDeploymentScript()
	if err != nil {
		return nil, err
	}

	// Environment variables
	environment, envFile, err := PromptEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	// Processes
	processes, err := PromptProcesses()
	if err != nil {
		return nil, err
	}

	// Scheduler
	scheduler, err := PromptScheduler()
	if err != nil {
		return nil, err
	}

	// Aliases
	aliases, err := PromptAliases()
	if err != nil {
		return nil, err
	}

	// Nginx
	nginxConfig, err := PromptNginxConfig()
	if err != nil {
		return nil, err
	}

	// SSL
	certificate, err := PromptSSLCertificate()
	if err != nil {
		return nil, err
	}

	// Isolation
	isolation, err := PromptIsolation()
	if err != nil {
		return nil, err
	}

	// Zero-downtime
	zeroDowntime, err := PromptZeroDowntime()
	if err != nil {
		return nil, err
	}

	// Build site config
	site := &models.SiteConfig{
		Name:                        basicInfo["name"].(string),
		DomainMode:                  basicInfo["domain_mode"].(string),
		WWWRedirectType:             basicInfo["www_redirect_type"].(string),
		RootDir:                     repoSettings["root_dir"].(string),
		WebDir:                      repoSettings["web_dir"].(string),
		CloneRepository:             repoSettings["clone_repository"].(bool),
		ProjectType:                 phpSettings["project_type"].(string),
		InstallComposerDependencies: phpSettings["install_composer_dependencies"].(bool),
		DeploymentScript:            deploymentScript,
		Environment:                 environment,
		EnvFile:                     envFile,
		Processes:                   processes,
		LaravelScheduler:            scheduler,
		Aliases:                     aliases,
		Certificate:                 certificate,
		Isolated:                    isolation["isolated"].(bool),
		ZeroDowntimeDeployments:     zeroDowntime["zero_downtime_deployments"].(bool),
	}

	if repoSettings["github_branch"] != nil && repoSettings["github_branch"].(string) != "" {
		site.GithubBranch = repoSettings["github_branch"].(string)
	}

	if phpSettings["php_version"] != nil && phpSettings["php_version"].(string) != "" {
		site.PHPVersion = phpSettings["php_version"].(string)
	}

	if isolation["isolated_user"] != nil && isolation["isolated_user"].(string) != "" {
		site.IsolatedUser = isolation["isolated_user"].(string)
	}

	if nginxTemplate, ok := nginxConfig["nginx_template"]; ok {
		site.NginxTemplate = nginxTemplate.(string)
	}

	if nginxVars, ok := nginxConfig["nginx_template_variables"]; ok {
		site.NginxTemplateVariables = nginxVars.(map[string]string)
	}

	if nginxCustom, ok := nginxConfig["nginx_custom_config"]; ok {
		site.NginxCustomConfig = nginxCustom.(string)
	}

	if sharedPaths, ok := zeroDowntime["shared_paths"]; ok {
		site.SharedPaths = sharedPaths.([]models.SharedPath)
	}

	site.SetDefaults()

	return site, nil
}
