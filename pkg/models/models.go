package models

import (
	"fmt"
	"strings"
)

// SharedPath represents a shared path for zero-downtime deployments
type SharedPath struct {
	From string `yaml:"-"`
	To   string `yaml:"-"`
}

// MarshalYAML implements custom YAML marshaling
func (sp SharedPath) MarshalYAML() (interface{}, error) {
	if sp.To == "" || sp.From == sp.To {
		return sp.From, nil
	}
	return map[string]string{
		"from": sp.From,
		"to":   sp.To,
	}, nil
}

// Process represents a background process
type Process struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

// SiteConfig represents configuration for a single site
type SiteConfig struct {
	Name                        string            `yaml:"name"`
	DomainMode                  string            `yaml:"domain_mode,omitempty"`
	WWWRedirectType             string            `yaml:"www_redirect_type,omitempty"`
	GithubBranch                string            `yaml:"github_branch,omitempty"`
	RootDir                     string            `yaml:"root_dir,omitempty"`
	WebDir                      string            `yaml:"web_dir,omitempty"`
	ProjectType                 string            `yaml:"project_type,omitempty"`
	PHPVersion                  string            `yaml:"php_version,omitempty"`
	InstallComposerDependencies bool              `yaml:"install_composer_dependencies,omitempty"`
	DeploymentScript            string            `yaml:"deployment_script,omitempty"`
	Processes                   []Process         `yaml:"processes,omitempty"`
	LaravelScheduler            bool              `yaml:"laravel_scheduler,omitempty"`
	Environment                 string            `yaml:"environment,omitempty"`
	EnvFile                     string            `yaml:"env_file,omitempty"`
	Aliases                     []string          `yaml:"aliases,omitempty"`
	NginxTemplate               string            `yaml:"nginx_template,omitempty"`
	NginxTemplateVariables      map[string]string `yaml:"nginx_template_variables,omitempty"`
	NginxCustomConfig           string            `yaml:"nginx_custom_config,omitempty"`
	Certificate                 bool              `yaml:"certificate,omitempty"`
	Isolated                    bool              `yaml:"isolated,omitempty"`
	IsolatedUser                string            `yaml:"isolated_user,omitempty"`
	ZeroDowntimeDeployments     bool              `yaml:"zero_downtime_deployments,omitempty"`
	SharedPaths                 []SharedPath      `yaml:"shared_paths,omitempty"`
	CloneRepository             bool              `yaml:"clone_repository,omitempty"`
}

// Validate validates the site configuration
func (s *SiteConfig) Validate() []string {
	var errors []string

	if s.Name == "" {
		errors = append(errors, "Site name is required")
	}

	if s.DomainMode != "" && s.DomainMode != "on-forge" && s.DomainMode != "custom" {
		errors = append(errors, "domain_mode must be 'on-forge' or 'custom'")
	}

	if s.WWWRedirectType != "" && s.WWWRedirectType != "none" && s.WWWRedirectType != "from-www" && s.WWWRedirectType != "to-www" {
		errors = append(errors, "www_redirect_type must be 'none', 'from-www', or 'to-www'")
	}

	if s.ProjectType != "" && s.ProjectType != "laravel" && s.ProjectType != "other" {
		errors = append(errors, "project_type must be 'laravel' or 'other'")
	}

	if s.Isolated && s.IsolatedUser == "" {
		errors = append(errors, "isolated_user is required when isolated is true")
	}

	if s.PHPVersion != "" && !strings.HasPrefix(s.PHPVersion, "php") {
		errors = append(errors, "php_version must start with 'php' (e.g., 'php81', 'php84')")
	}

	return errors
}

// DeploymentConfig represents the complete deployment configuration
type DeploymentConfig struct {
	Organization     string       `yaml:"organization"`
	Server           string       `yaml:"server"`
	GithubRepository string       `yaml:"github_repository"`
	GithubBranch     string       `yaml:"github_branch"`
	Sites            []SiteConfig `yaml:"sites"`
}

// Validate validates the deployment configuration
func (d *DeploymentConfig) Validate() []string {
	var errors []string

	if d.Organization == "" {
		errors = append(errors, "Organization is required")
	}

	if d.Server == "" {
		errors = append(errors, "Server is required")
	}

	if d.GithubRepository == "" {
		errors = append(errors, "GitHub repository is required")
	} else if !strings.Contains(d.GithubRepository, "/") {
		errors = append(errors, "GitHub repository must be in format 'owner/repo'")
	}

	if len(d.Sites) == 0 {
		errors = append(errors, "At least one site must be configured")
	}

	for i, site := range d.Sites {
		siteErrors := site.Validate()
		for _, err := range siteErrors {
			errors = append(errors, fmt.Sprintf("Site %d (%s): %s", i+1, site.Name, err))
		}
	}

	return errors
}

// SetDefaults sets default values for optional fields
func (s *SiteConfig) SetDefaults() {
	if s.DomainMode == "" {
		s.DomainMode = "on-forge"
	}
	if s.WWWRedirectType == "" {
		s.WWWRedirectType = "none"
	}
	if s.RootDir == "" {
		s.RootDir = "."
	}
	if s.WebDir == "" {
		s.WebDir = "public"
	}
	if s.ProjectType == "" {
		s.ProjectType = "laravel"
	}
	if !s.CloneRepository {
		s.CloneRepository = true
	}
}
