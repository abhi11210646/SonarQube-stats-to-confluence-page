package main

import (
	"log"

	"github.com/joho/godotenv"
	config "gitlab.group.one/sonar-to-confluence/internal"
	"gitlab.group.one/sonar-to-confluence/internal/confluence"
	"gitlab.group.one/sonar-to-confluence/internal/sonar"
)

func main() {
	godotenv.Load(".env")
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or error loading it. Proceeding with environment variables.")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	//Create sonar Client
	sonarClient := sonar.NewSonarClient(cfg.Sonar)

	// Update stats to confluence page
	confluenceClient := confluence.NewConfluenceClient(cfg.Confluence, sonarClient)
	if err := confluenceClient.UpdatePage(); err != nil {
		log.Fatalf("Failed to update Confluence page: %v", err)
	}
}
