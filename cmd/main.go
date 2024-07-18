package main

import (
	"gitlab.group.one/sonar-to-confluence/internal/config"
	"gitlab.group.one/sonar-to-confluence/internal/confluence"
	"gitlab.group.one/sonar-to-confluence/internal/sonar"
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
