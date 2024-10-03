package confluence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	config "gitlab.group.one/sonar-to-confluence/internal"
)

type SonarClient interface {
	StatsHTML() (string, error)
}

type Confluence struct {
	config config.ConfluenceConfig
	sonar  SonarClient
}

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
	Number  int    `json:"number"`
	Message string `json:"message"`
}

func NewConfluenceClient(config config.ConfluenceConfig, sonar SonarClient) *Confluence {
	return &Confluence{
		config,
		sonar,
	}
}

func (c *Confluence) addAuthHeader(req *http.Request) {
	key := base64.StdEncoding.EncodeToString([]byte(c.config.ApiKey))
	req.Header.Add("Authorization", "Basic "+key)
}

func (c *Confluence) makeRequest(method string, apiEndpoint string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, apiEndpoint, body)

	// Add auth
	c.addAuthHeader(req)

	if method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return client.Do(req)
}

func (c *Confluence) FetchPage() (Page, error) {
	apiEndpoint := fmt.Sprintf("%s/api/content/%s?expand=body.storage,version", c.config.Host, c.config.PageID)

	resp, err := c.makeRequest("GET", apiEndpoint, nil)
	if err != nil {
		return Page{}, fmt.Errorf("failed to fetch Confluence page: %w", err)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Page{}, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("Confluence API error: %s - %s", resp.Status, string(body))
	}
	// Unmarshal it to golang struct
	var page Page
	if err := json.Unmarshal(body, &page); err != nil {
		return Page{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return page, nil
}

func (c *Confluence) UpdatePage() error {
	apiEndpoint := fmt.Sprintf("%s/api/content/%s?expand=body.storage", c.config.Host, c.config.PageID)

	page, err := c.FetchPage()
	if err != nil {
		return fmt.Errorf("failed to fetch existing page: %w", err)
	}
	html, err := c.sonar.StatsHTML()
	if err != nil {
		return err
	}
	newPage := Page{
		Title: page.Title,
		Type:  page.Type,
		Version: Version{
			Number:  page.Version.Number + 1,
			Message: "Updated by CronJob...",
		},
		Body: Body{
			Storage: Storage{
				Representation: "storage",
				Value:          html,
			},
		},
	}

	page_bytes, err := json.Marshal(newPage)
	if err != nil {
		return fmt.Errorf("failed to marshal new page data: %w", err)
	}

	resp, err := c.makeRequest("PUT", apiEndpoint, bytes.NewReader(page_bytes))
	if err != nil {
		return fmt.Errorf("failed to update Confluence page: %w", err)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Confluence API error: %s - %s", resp.Status, string(body))
	}
	log.Println("Confluence page updated successfully.")
	return nil
}
