package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func updateByPageId() {
	apiEndpoint := ConfluenceConfig.Host + "/api/content/" + strconv.Itoa(ConfluenceConfig.PageId) + "?expand=body.storage"
	var page Page = getByPageId()

	newPage := Page{
		Title: page.Title,
		Type:  page.Type,
		Version: Version{
			Number:  page.Version.Number + 1,
			Message: page.Version.Message,
		},
		Body: Body{
			Storage: Storage{
				Representation: "storage",
				Value:          "hhh4444444h",
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
}
