package sonar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	config "gitlab.group.one/sonar-to-confluence/internal"
)

type SonarClient struct {
	config config.SonarConfig
}
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

func NewSonarClient(config config.SonarConfig) SonarClient {
	return SonarClient{
		config,
	}
}

func (s SonarClient) FetchStatsByProjectKey(projectKey string) (Stats, error) {
	fmt.Printf("Fetching stats for %s... ", projectKey)
	apiEndpoint := fmt.Sprintf("%s/api/measures/component?component=%s&metricKeys=%s", s.config.Host, projectKey, strings.Join(s.config.Metrics, ","))

	// Create new request and pass headers
	req, _ := http.NewRequest("GET", apiEndpoint, nil)
	key := base64.StdEncoding.EncodeToString([]byte(s.config.ApiKey + ":"))
	req.Header.Add("Authorization", "Basic "+key)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return Stats{}, fmt.Errorf("error in fetching SonarStats %w", err)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Stats{}, fmt.Errorf("error in reading response body %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Stats{}, fmt.Errorf("error response from Sonar API: %s - %s", resp.Status, string(body))
	}
	// Unmarshal it to golang struct
	var stats Stats
	if err := json.Unmarshal(body, &stats); err != nil {
		return Stats{}, fmt.Errorf("error in UnMarshaling:  %w", err)
	}
	fmt.Println("Done")
	return stats, nil
}

func (s SonarClient) FetchStats() ([]Stats, error) {
	var stats []Stats
	// Fetch all stats
	for _, projectKey := range s.config.Projects {
		projectStats, err := s.FetchStatsByProjectKey(projectKey)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch stats for project %s: %w", projectKey, err)
		}
		stats = append(stats, projectStats)
	}
	return stats, nil
}

func (s SonarClient) StatsHTML() (string, error) {
	keys := s.config.Metrics
	// Fetch all stats of given projects
	stats, err := s.FetchStats()
	if err != nil {
		return "", fmt.Errorf("error in FetchStats:  %w", err)
	}
	// Mappings for Header Title
	Columns := map[string]string{
		"name":                    "Product",
		"alert_status":            "Quality Gate",
		"code_smells":             "Code Smells",
		"bugs":                    "Bugs",
		"critical_severity_vulns": "Vulnerabilities",
	}
	// Table Header
	headers := make([]string, len(keys)+1)
	headers[0] = "Product" // First column Product name
	for i, k := range keys {
		if name, ok := Columns[k]; ok {
			headers[i+1] = name
		} else {
			headers[i+1] = k
		}
	}

	// Table Body
	body := [][]string{}
	for _, stat := range stats {
		metrics := make(map[string]string)
		for _, v := range stat.Component.Measures {
			metrics[v.Metric] = v.Value
		}
		s := make([]string, len(stat.Component.Measures)+1)
		s[0] = stat.Component.Name // First column is Product Name

		for i, k := range keys {
			if metrics[k] == "ERROR" {
				s[i+1] = "Failed"
			} else {
				s[i+1] = metrics[k]
			}
		}
		body = append(body, s)
	}
	// Template Data for HTML parser
	var TemplateData = struct {
		Headers []string
		Body    [][]string
	}{
		Headers: headers,
		Body:    body,
	}
	const html = `<table data-table-width="760" data-layout="default" ac:local-id="091ca39e-2b3b-4a0c-8720-7ee499fc6d65">
		<tbody>
				<tr>
					{{ range .Headers }}
					<th><p><strong>{{.}}</strong></p></th>
					{{end}}
				</tr>

				
				{{range .Body}}
				<tr>
				{{range .}}
				 <td> {{ . }} </td>
				{{end}}
				</tr>
				{{end}}
				
				
		</tbody>
	</table>`

	t, err := template.New("confluence").Parse(html)
	if err != nil {
		return "", fmt.Errorf("error in template.Parse:  %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, TemplateData); err != nil {
		return "", fmt.Errorf("error in Execute:  %w", err)
	}
	return buf.String(), nil
}
