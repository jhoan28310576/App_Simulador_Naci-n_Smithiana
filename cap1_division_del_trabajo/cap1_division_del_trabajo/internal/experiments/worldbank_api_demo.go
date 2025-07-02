package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Estructura para parsear la respuesta del World Bank
// La respuesta es un array: [metadata, [data...]]
type WBResponse []interface{}

// Estructura para parsear cada registro de datos
type CPIRecord struct {
	Country struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"country"`
	Date  string      `json:"date"`
	Value interface{} `json:"value"` // Usamos interface{} para manejar diferentes tipos
}

// Función para extraer el valor numérico de manera segura
func (c *CPIRecord) GetNumericValue() (float64, bool) {
	if c.Value == nil {
		return 0, false
	}

	switch v := c.Value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		if v == "" {
			return 0, false
		}
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func main() {
	// Venezuela: country code 'VE', Colombia: 'CO'
	countryCodes := []string{"VE", "CO"}
	indicator := "FP.CPI.TOTL" // CPI (Consumer Price Index)

	for _, country := range countryCodes {
		fmt.Printf("\n==== CPI para %s ===\n", country)
		url := fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/%s?format=json&per_page=10", country, indicator)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error en la petición para %s: %v\n", country, err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error leyendo respuesta para %s: %v\n", country, err)
			continue
		}

		var wbResp WBResponse
		if err := json.Unmarshal(body, &wbResp); err != nil {
			fmt.Printf("Error parseando JSON para %s: %v\n", country, err)
			continue
		}

		if len(wbResp) < 2 {
			fmt.Printf("Respuesta inesperada para %s: %s\n", country, string(body))
			continue
		}

		// El segundo elemento es el array de datos
		recordsRaw, ok := wbResp[1].([]interface{})
		if !ok {
			fmt.Printf("Formato de datos inesperado para %s\n", country)
			continue
		}

		fmt.Printf("Registros encontrados: %d\n", len(recordsRaw))

		for i, rec := range recordsRaw {
			recBytes, _ := json.Marshal(rec)
			var cpi CPIRecord

			if err := json.Unmarshal(recBytes, &cpi); err != nil {
				fmt.Printf("Error parseando registro %d: %v\n", i, err)
				continue
			}

			if value, ok := cpi.GetNumericValue(); ok {
				fmt.Printf("Año: %s | CPI: %.2f\n", cpi.Date, value)
			} else {
				fmt.Printf("Año: %s | CPI: Sin datos\n", cpi.Date)
			}
		}
	}
}
