package main

import (
	"github.com/joho/godotenv"
	config "gitlab.group.one/sonar-to-confluence/internal"
	"gitlab.group.one/sonar-to-confluence/internal/confluence"
	"gitlab.group.one/sonar-to-confluence/internal/sonar"
)

func main() {
	godotenv.Load(".env")
	//Create sonar Client
	sonarConfig := config.GetSonarConfig()
	sonarClient := sonar.NewSonarClient(sonarConfig)
	// Update stats to confluence page
	confluenceConfig := config.GetConfluenceConfig()
	confluenceClient := confluence.NewConfluenceClient(confluenceConfig, sonarClient)
	confluenceClient.UpdatePage()
}
