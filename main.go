package main

func main() {
	var stats []Stats
	for _, projectKey := range SonarConfig.Projects {
		stats = append(stats, SonarStats(projectKey))
	}
	updateByPageId(stats)
}
