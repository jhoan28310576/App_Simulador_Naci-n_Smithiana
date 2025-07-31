package worldbank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type wbResponse struct {
	Page    int    `json:"page"`
	Pages   int    `json:"pages"`
	PerPage string `json:"per_page"`
	Total   int    `json:"total"`
	Source  string `json:"source"`
}

type wbIndicator struct {
	Country struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"country"`
	Date  string   `json:"date"`
	Value *float64 `json:"value"`
}

// GetIndicator consulta el valor más reciente disponible para un país e indicador
// Si se pasa apiKey, la agrega como parámetro a la URL
func GetIndicator(country, indicator string, apiKey ...string) (float64, error) {
	url := fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/%s?format=json&per_page=5", country, indicator)
	if len(apiKey) > 0 && apiKey[0] != "" {
		url += "&key=" + apiKey[0]
	}
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result []json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	if len(result) < 2 {
		return 0, fmt.Errorf("no data")
	}
	var indicators []wbIndicator
	if err := json.Unmarshal(result[1], &indicators); err != nil {
		return 0, err
	}
	for _, ind := range indicators {
		if ind.Value != nil {
			return *ind.Value, nil
		}
	}
	return 0, fmt.Errorf("no value found for %s/%s", country, indicator)
}
