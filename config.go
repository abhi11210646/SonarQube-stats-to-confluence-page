package main

import (
	"github.com/joho/godotenv"
)

var env, _ = godotenv.Read(".env")

var SonarConfig = struct {
	Host     string
	ApiKey   string
	Projects []string
	Metrics  []string
}{
	Host:   "https://sonarqube.one.com",
	ApiKey: env["SONAR_API_KEY"],
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
	},
}

var ConfluenceConfig = struct {
	Host   string
	ApiKey string
	PageId int
}{
	Host:   "https://group-one.atlassian.net/wiki/rest/",
	ApiKey: env["CONFLUENCE_API_KEY"],
	PageId: 32954647094,
}
