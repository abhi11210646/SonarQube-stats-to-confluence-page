package config

import (
	"os"

	"github.com/joho/godotenv"
)

type SonarConfig struct {
	Host     string
	ApiKey   string
	Projects []string
	Metrics  []string
}

type ConfluenceConfig struct {
	Host   string
	ApiKey string
	PageId int
}

type Config struct {
	Sonar      SonarConfig
	Confluence ConfluenceConfig
}

func GetConfig() Config {
	godotenv.Load(".env")
	return Config{
		Confluence: ConfluenceConfig{
			Host:   "https://group-one.atlassian.net/wiki/rest/",
			ApiKey: os.Getenv("CONFLUENCE_API_KEY"),
			PageId: 32589873205,
		},
		Sonar: SonarConfig{
			Host:   "https://sonarqube.one.com",
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
	}
}
