package main

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
)

type Page struct {
	Title   string  `json:"title"`
	Type    string  `json:"type"`
	Body    Body    `json:"body"`
	Version Version `json:"version"`
}
type Body struct {
	Storage `json:"storage"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}
type Version struct {
	Number  int8   `json:"number"`
	Message string `json:"message"`
}

func getByPageId() Page {

	apiEndpoint := ConfluenceConfig.Host + "/api/content/" + strconv.Itoa(ConfluenceConfig.PageId) + "?expand=body.storage,version"

	req, _ := http.NewRequest("GET", apiEndpoint, nil)
	key := base64.StdEncoding.EncodeToString([]byte(ConfluenceConfig.ApiKey))
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

func updateByPageId(stats Stats) {
	apiEndpoint := ConfluenceConfig.Host + "/api/content/" + strconv.Itoa(ConfluenceConfig.PageId) + "?expand=body.storage"
	var page Page = getByPageId()

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
	key := base64.StdEncoding.EncodeToString([]byte(ConfluenceConfig.ApiKey))
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
	fmt.Println("Stats updated to confluence page!")
}

func generetaeHTML(stats Stats) string {
	// Mappings for Header Title
	Columns := map[string]string{
		"name":                    "Product",
		"alert_status":            "Quality Gate",
		"code_smells":             "Code Smells",
		"bugs":                    "Bugs",
		"critical_severity_vulns": "Vulnerabilities",
	}

	// Table Header
	Keys := SonarConfig.Metrics
	headers := make([]string, len(Keys)+1)
	headers[0] = "Product" // First column Product name
	for i, k := range Keys {
		if name, ok := Columns[k]; ok {
			headers[i+1] = name
		} else {
			headers[i+1] = k
		}
	}
	// Template Data for HTML parser
	var TemplateData = struct {
		Headers []string
		Stats   Stats
	}{
		Headers: headers,
		Stats:   stats,
	}
	const html = `<table data-table-width="760" data-layout="default" ac:local-id="091ca39e-2b3b-4a0c-8720-7ee499fc6d65">
		<tbody>
				<tr>
					{{ range .Headers }}
					<th><p><strong>{{.}}</strong></p></th>
					{{end}}
				</tr>

				<tr>
				<td> {{ .Stats.Component.Name }} </td>
				{{range .Stats.Component.Measures}}
				<td> {{ .Value }} </td>
				{{end}}
				</tr>
				
		</tbody>
	</table>`

	t, _ := template.New("confluence").Parse(html)

	var buf bytes.Buffer
	t.Execute(&buf, TemplateData)
	return buf.String()
}
