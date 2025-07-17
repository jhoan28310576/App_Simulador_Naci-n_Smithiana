package test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type usdaResponse struct {
	Data []interface{} `json:"data"`
}

func TestUSDAApiReturnsData(t *testing.T) {
	url := "https://quickstats.nass.usda.gov/api/api_GET/?key=1F325726-42E7-3E08-8E7F-C7ED7047890A&commodity_desc=CORN&year=2023&state_alpha=IA&statisticcat_desc=PRODUCTION&format=JSON"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Error al hacer la petición: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("Código de estado inesperado: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error leyendo el cuerpo de la respuesta: %v", err)
	}
	var parsed usdaResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("Error parseando JSON: %v", err)
	}
	if len(parsed.Data) == 0 {
		t.Fatalf("La respuesta no contiene datos en el campo 'data'")
	}
	t.Logf("La respuesta contiene %d registros de datos", len(parsed.Data))
}
