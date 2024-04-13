package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Stats struct {
	Component struct {
		Name     string
		Measures []Measure `json:"measures"`
	} `json:"component"`
}
type Measure struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
}

func SonarStats(projectKey string) Stats {
	sonarConfig := GetSonarConfig()
	apiEndpoint := sonarConfig.Host + "/api/measures/component?component=" + projectKey + "&metricKeys=" + strings.Join(sonarConfig.Metrics, ",")

	// Create new request and pass headers
	req, _ := http.NewRequest("GET", apiEndpoint, nil)
	key := base64.StdEncoding.EncodeToString([]byte(sonarConfig.ApiKey + ":"))
	req.Header.Add("Authorization", "Basic "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in fetching SonarStats", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in reading response body", err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Error response from Sonar API: ", resp.Status, string(body))
		os.Exit(1)
	}
	// Unmarshal it to golang struct
	var stats Stats
	if err := json.Unmarshal(body, &stats); err != nil {
		fmt.Println("Error in UnMarshaling: ", err)
		os.Exit(1)
	}
	return stats
}
