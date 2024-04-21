package confluence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"gitlab.group.one/sonar-to-confluence/config"
	"gitlab.group.one/sonar-to-confluence/sonar"
)

var confluenceConfig = config.GetConfluenceConfig()

func getByPageId() Page {

	apiEndpoint := confluenceConfig.Host + "/api/content/" + strconv.Itoa(confluenceConfig.PageId) + "?expand=body.storage,version"

	req, _ := http.NewRequest("GET", apiEndpoint, nil)
	key := base64.StdEncoding.EncodeToString([]byte(confluenceConfig.ApiKey))
	req.Header.Add("Authorization", "Basic "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[getByPageId]Error in fetching Confluence API", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[getByPageId]Error in reading response body", err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("[getByPageId]Error response from Confluence API: ", resp.Status, string(body))
		os.Exit(1)
	}
	// Unmarshal it to golang struct
	var page Page
	if err := json.Unmarshal(body, &page); err != nil {
		fmt.Println("[getByPageId]Error in UnMarshaling: ", err)
		os.Exit(1)
	}
	return page
}

func UpdateStats(stats []sonar.Stats) {
	apiEndpoint := confluenceConfig.Host + "/api/content/" + strconv.Itoa(confluenceConfig.PageId) + "?expand=body.storage"
	page := getByPageId()
	newPage := Page{
		Title: page.Title,
		Type:  page.Type,
		Version: Version{
			Number:  page.Version.Number + 1,
			Message: "Updated by CronJob",
		},
		Body: Body{
			Storage: Storage{
				Representation: "storage",
				Value:          generetaeHTML(stats),
			},
		},
	}

	page_bytes, _ := json.Marshal(newPage)

	req, _ := http.NewRequest("PUT", apiEndpoint, bytes.NewReader(page_bytes))
	key := base64.StdEncoding.EncodeToString([]byte(confluenceConfig.ApiKey))
	req.Header.Add("Authorization", "Basic "+key)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[updateByPageId]Error in fetching Confluence API", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[updateByPageId]Error in reading response body", err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("[updateByPageId]Error response from Confluence API: ", resp.Status, string(body))
		os.Exit(1)
	}
	fmt.Println("Stats updated to confluence page! Success.")
}

func generetaeHTML(stats []sonar.Stats) string {
	// Mappings for Header Title
	Columns := map[string]string{
		"name":                    "Product",
		"alert_status":            "Quality Gate",
		"code_smells":             "Code Smells",
		"bugs":                    "Bugs",
		"critical_severity_vulns": "Vulnerabilities",
	}

	// Keys := config.SonarConfig.Metrics
	Keys := []string{}
	// Table Header
	headers := make([]string, len(Keys)+1)
	headers[0] = "Product" // First column Product name
	for i, k := range Keys {
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

		for i, k := range Keys {

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
