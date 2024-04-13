package main

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

func GetSonarConfig() SonarConfig {
	return SonarConfig{
		Host:   "https://sonarqube.one.com",
		ApiKey: os.Getenv("SONAR_API_KEY"),
		Projects: []string{
			"app.webmail",
			"CompanionApp",
			"Webshop",
			"one.com-wp-addons-assets",
		},
		Metrics: []string{
			"code_smells",
			"critical_severity_vulns",
			"bugs",
			"alert_status",
			"quality_gate_details",
		},
	}
}

type ConfluenceConfig struct {
	Host   string
	ApiKey string
	PageId int
}

func GetConfluenceConfig() ConfluenceConfig {
	return ConfluenceConfig{
		Host:   "https://group-one.atlassian.net/wiki/rest/",
		ApiKey: os.Getenv("CONFLUENCE_API_KEY"),
		PageId: 32589873205,
	}
}

func init() {
	godotenv.Load()
}
