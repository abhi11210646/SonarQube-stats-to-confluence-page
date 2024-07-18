package sonar

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gitlab.group.one/sonar-to-confluence/internal/config"
)

type SonarClient struct {
	config.Config
}

func NewSonarClient(config config.Config) SonarClient {
	return SonarClient{
		config,
	}
}

func (s SonarClient) FetchStats(projectKey string) Stats {
	fmt.Printf("Fetching stats for %s... ", projectKey)
	apiEndpoint := s.Sonar.Host + "/api/measures/component?component=" + projectKey + "&metricKeys=" + strings.Join(s.Sonar.Metrics, ",")

	// Create new request and pass headers
	req, _ := http.NewRequest("GET", apiEndpoint, nil)
	key := base64.StdEncoding.EncodeToString([]byte(s.Sonar.ApiKey + ":"))
	req.Header.Add("Authorization", "Basic "+key)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error in fetching SonarStats", err)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error in reading response body", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Error response from Sonar API: ", resp.Status, string(body))
	}
	// Unmarshal it to golang struct
	var stats Stats
	if err := json.Unmarshal(body, &stats); err != nil {
		log.Fatal("Error in UnMarshaling: ", err)
	}
	fmt.Println("Done")
	return stats
}
