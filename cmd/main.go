package main

import (
	"sonar-to-confluence/internal/confluence"
	"sonar-to-confluence/internal/sonar"

	config "sonar-to-confluence/internal"

	"github.com/joho/godotenv"
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
