package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	baseURL, _ := url.Parse("https://quickstats.nass.usda.gov/api/api_GET/")

	params := url.Values{}
	params.Add("key", "1F325726-42E7-3E08-8E7F-C7ED7047890A")
	params.Add("commodity_desc", "CORN")
	params.Add("year", "2023")
	params.Add("state_alpha", "IA")
	params.Add("statisticcat_desc", "PRODUCTION")
	params.Add("format", "JSON")

	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response:\n%s\n", string(body))

	// Intentar parsear como JSON genérico para ver la estructura
	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Imprimir la estructura
	prettyJSON, _ := json.MarshalIndent(rawData, "", "  ")
	fmt.Printf("Parsed structure:\n%s\n", string(prettyJSON))
}
