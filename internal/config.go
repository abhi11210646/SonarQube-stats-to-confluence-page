package config

import (
	"fmt"
	"os"
)

const (
	SonarHost      = "https://sonarqube.one.com"
	ConfluenceHost = "https://group-one.atlassian.net/wiki/rest/"
)

type Config struct {
	Sonar      SonarConfig
	Confluence ConfluenceConfig
}

func (c *Config) Validate() error {
	if c.Sonar.ApiKey == "" {
		return fmt.Errorf("SONAR_API_KEY is required")
	}
	if c.Confluence.ApiKey == "" {
		return fmt.Errorf("CONFLUENCE_API_KEY is required")
	}
	return nil
}

type SonarConfig struct {
	Host     string
	ApiKey   string
	Projects []string
	Metrics  []string
}

type ConfluenceConfig struct {
	Host   string
	ApiKey string
	PageID string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Sonar: SonarConfig{
			Host:   SonarHost,
			ApiKey: os.Getenv("SONAR_API_KEY"),
			Projects: []string{
				"app.webmail",
				"CompanionApp",
				"Webshop",
				"one.com-wp-addons-assets",
			},
			Metrics: []string{
				"alert_status",
				"code_smells",
				"critical_severity_vulns",
				"bugs",
			},
		},
		Confluence: ConfluenceConfig{
			Host:   ConfluenceHost,
			ApiKey: os.Getenv("CONFLUENCE_API_KEY"),
			PageID: "32589873205",
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
