package sonar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	config "sonar-to-confluence/internal"
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

func (s SonarClient) FetchStatsByProjectKey(projectKey string) Stats {
	fmt.Printf("Fetching stats for %s... ", projectKey)
	apiEndpoint := s.config.Host + "/api/measures/component?component=" + projectKey + "&metricKeys=" + strings.Join(s.config.Metrics, ",")

	// Create new request and pass headers
	req, _ := http.NewRequest("GET", apiEndpoint, nil)
	key := base64.StdEncoding.EncodeToString([]byte(s.config.ApiKey + ":"))
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

func (s SonarClient) FetchStats() []Stats {
	var stats []Stats
	// Fetch all stats
	for _, projectKey := range s.config.Projects {
		stats = append(stats, s.FetchStatsByProjectKey(projectKey))
	}
	return stats
}

func (s SonarClient) StatsHTML() string {
	keys := s.config.Metrics
	// Fetch all stats of given projects
	stats := s.FetchStats()
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

	t, _ := template.New("confluence").Parse(html)

	var buf bytes.Buffer
	t.Execute(&buf, TemplateData)
	return buf.String()
}
