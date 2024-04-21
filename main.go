package main

import (
	"gitlab.group.one/sonar-to-confluence/config"
	"gitlab.group.one/sonar-to-confluence/confluence"
	"gitlab.group.one/sonar-to-confluence/sonar"
)

func main() {
	sonarConfig := config.GetSonarConfig()
	var stats []sonar.Stats
	for _, projectKey := range sonarConfig.Projects {
		stats = append(stats, sonar.FetchStats(projectKey))
	}
	confluence.UpdateStats(stats)
}
