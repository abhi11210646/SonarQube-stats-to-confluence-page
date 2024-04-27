package main

import (
	"gitlab.group.one/sonar-to-confluence/config"
	"gitlab.group.one/sonar-to-confluence/confluence"
	"gitlab.group.one/sonar-to-confluence/sonar"
)

func main() {
	config := config.GetConfig()
	sonarClient := sonar.NewSonarClient(config)
	// Fetch all stats
	var stats []sonar.Stats
	for _, projectKey := range config.Sonar.Projects {
		stats = append(stats, sonarClient.FetchStats(projectKey))
	}
	// Update stats to confluence page
	confluence.NewConfluenceClient(config).UpdateStats(stats)
}
