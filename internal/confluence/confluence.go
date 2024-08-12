package confluence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	config "sonar-to-confluence/internal"
)

type SonarClient interface {
	StatsHTML() string
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
	Number  int8   `json:"number"`
	Message string `json:"message"`
}

func NewConfluenceClient(config config.ConfluenceConfig, sonar SonarClient) Confluence {
	return Confluence{
		config,
		sonar,
	}
}

func (c Confluence) FetchPage() (*Page, error) {
	apiEndpoint := c.config.Host + "/api/content/" + strconv.Itoa(c.config.PageId) + "?expand=body.storage,version"

	req, _ := http.NewRequest("GET", apiEndpoint, nil)
	key := base64.StdEncoding.EncodeToString([]byte(c.config.ApiKey))
	req.Header.Add("Authorization", "Basic "+key)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[getByPageId]Error in fetching Confluence API %w", err)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("[getByPageId]Error in reading response body %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[getByPageId]Error response from Confluence API: %v, %v", resp.Status, string(body))
	}
	// Unmarshal it to golang struct
	page := &Page{}
	if err := json.Unmarshal(body, page); err != nil {
		return nil, fmt.Errorf("[getByPageId]Error in UnMarshaling: %w", err)
	}
	return page, nil
}

func (c Confluence) UpdatePage() {
	apiEndpoint := c.config.Host + "/api/content/" + strconv.Itoa(c.config.PageId) + "?expand=body.storage"
	page, err := c.FetchPage()
	if err != nil {
		log.Fatal("[UpdatePage]Error in FetchPage: ", err)
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
				Value:          c.sonar.StatsHTML(),
			},
		},
	}

	page_bytes, _ := json.Marshal(newPage)

	req, _ := http.NewRequest("PUT", apiEndpoint, bytes.NewReader(page_bytes))
	key := base64.StdEncoding.EncodeToString([]byte(c.config.ApiKey))
	req.Header.Add("Authorization", "Basic "+key)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("[updateByPageId]Error in fetching Confluence API", err)
	}
	defer resp.Body.Close()
	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("[updateByPageId]Error in reading response body", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("[updateByPageId]Error response from Confluence API: ", resp.Status, string(body))
	}
	fmt.Println("Stats updated to confluence page! Success.")
}
