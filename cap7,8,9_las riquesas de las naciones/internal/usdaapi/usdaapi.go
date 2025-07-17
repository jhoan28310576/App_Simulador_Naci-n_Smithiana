package usdaapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// USDAResponse representa la estructura de la respuesta de la API de USDA
// El campo Data contiene los registros devueltos por la consulta
// Cada registro es un mapa de clave-valor
type USDAResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// GetUSDAData consulta la API de USDA con los parámetros dados y retorna los datos
func GetUSDAData(commodity, year, state, statisticcat string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf(
		"https://quickstats.nass.usda.gov/api/api_GET/?key=1F325726-42E7-3E08-8E7F-C7ED7047890A&commodity_desc=%s&year=%s&state_alpha=%s&statisticcat_desc=%s&format=JSON",
		commodity, year, state, statisticcat,
	)

	// Hacer la petición HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al hacer la petición: %w", err)
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo el cuerpo de la respuesta: %w", err)
	}

	// Parsear el JSON
	var parsed USDAResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("error parseando JSON: %w", err)
	}

	// Validar que haya datos
	if len(parsed.Data) == 0 {
		return nil, fmt.Errorf("la respuesta no contiene datos")
	}

	return parsed.Data, nil
}
