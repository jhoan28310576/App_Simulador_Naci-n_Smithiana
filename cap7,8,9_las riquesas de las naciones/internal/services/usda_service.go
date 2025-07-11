package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/models"
)

const (
	USDA_API_BASE_URL = "https://quickstats.nass.usda.gov/api/api_GET/"
	USDA_API_KEY      = "1F325726-42E7-3E08-8E7F-C7ED7047890A"
)

type USDAService struct {
	apiKey string
	client *http.Client
}

func NewUSDAService() *USDAService {
	return &USDAService{
		apiKey: USDA_API_KEY,
		client: &http.Client{},
	}
}

// GetCornProduction obtiene datos de producción de maíz para estados específicos
func (s *USDAService) GetCornProduction(year string, states []string) ([]models.CornProductionData, error) {
	var allData []models.CornProductionData

	for _, state := range states {
		// Consultar producción
		productionData, err := s.queryUSDA("CORN", year, state, "PRODUCTION")
		if err != nil {
			return nil, fmt.Errorf("error querying production for %s: %v", state, err)
		}

		// Consultar área cosechada
		areaData, err := s.queryUSDA("CORN", year, state, "AREA HARVESTED")
		if err != nil {
			return nil, fmt.Errorf("error querying area for %s: %v", state, err)
		}

		// Procesar datos de producción
		for _, data := range productionData.Data {
			// Filtrar solo datos de nivel estatal y de grano (no silage)
			if data.AggLevelDesc != "STATE" || data.UtilPracticeDesc != "GRAIN" {
				continue
			}

			production, err := parseNumericValue(data.Value)
			if err != nil {
				continue // Skip invalid data
			}

			// Convertir year a string
			yearStr := ""
			switch v := data.Year.(type) {
			case string:
				yearStr = v
			case float64:
				yearStr = fmt.Sprintf("%.0f", v)
			case int:
				yearStr = fmt.Sprintf("%d", v)
			default:
				yearStr = fmt.Sprintf("%v", v)
			}

			cornData := models.CornProductionData{
				Year:       yearStr,
				State:      data.StateName,
				Production: production,
				Unit:       data.UnitDesc,
			}

			// Buscar área correspondiente
			for _, area := range areaData.Data {
				// Filtrar solo datos de nivel estatal y de grano
				if area.AggLevelDesc != "STATE" || area.UtilPracticeDesc != "GRAIN" {
					continue
				}

				// Convertir year del área a string para comparar
				areaYearStr := ""
				switch v := area.Year.(type) {
				case string:
					areaYearStr = v
				case float64:
					areaYearStr = fmt.Sprintf("%.0f", v)
				case int:
					areaYearStr = fmt.Sprintf("%d", v)
				default:
					areaYearStr = fmt.Sprintf("%v", v)
				}

				if areaYearStr == yearStr && area.StateName == data.StateName {
					if areaValue, err := parseNumericValue(area.Value); err == nil {
						cornData.Area = areaValue
						cornData.AreaUnit = area.UnitDesc
					}
					break
				}
			}

			allData = append(allData, cornData)
		}
	}

	return allData, nil
}

// queryUSDA hace una consulta a la API de USDA Quick Stats
func (s *USDAService) queryUSDA(commodity, year, state, statistic string) (*models.USDAQuickStatsResponse, error) {
	baseURL, err := url.Parse(USDA_API_BASE_URL)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("key", s.apiKey)
	params.Add("commodity_desc", commodity)
	params.Add("year", year)
	params.Add("state_alpha", state)
	params.Add("statisticcat_desc", statistic)
	params.Add("format", "JSON")

	baseURL.RawQuery = params.Encode()

	resp, err := s.client.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var response models.USDAQuickStatsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &response, nil
}

// parseNumericValue convierte el valor string de la API a float64
func parseNumericValue(valueStr string) (float64, error) {
	// Remover comas y espacios
	cleanValue := strings.ReplaceAll(valueStr, ",", "")
	cleanValue = strings.TrimSpace(cleanValue)

	return strconv.ParseFloat(cleanValue, 64)
}

// SimulateDrought simula el impacto de una sequía en la producción de maíz
func (s *USDAService) SimulateDrought(params models.DroughtSimulationParams) (*models.DroughtSimulationResult, error) {
	// Obtener datos reales
	realData, err := s.GetCornProduction(params.Year, params.States)
	if err != nil {
		return nil, err
	}

	var originalProduction, simulatedProduction float64
	var simulatedStates []models.CornProductionData

	for _, data := range realData {
		originalProduction += data.Production

		// Calcular producción simulada
		affectedProduction := data.Production * params.AffectedArea * params.DroughtSeverity
		unaffectedProduction := data.Production * (1 - params.AffectedArea)
		newProduction := unaffectedProduction + (affectedProduction * (1 - params.DroughtSeverity))

		simulatedData := data
		simulatedData.Production = newProduction
		simulatedStates = append(simulatedStates, simulatedData)
		simulatedProduction += newProduction
	}

	productionLoss := originalProduction - simulatedProduction
	priceIncrease := (productionLoss / originalProduction) * 0.5 // Elasticidad de precio simplificada
	economicImpact := productionLoss * 4.5                       // Valor aproximado por bushel

	return &models.DroughtSimulationResult{
		OriginalProduction:  originalProduction,
		SimulatedProduction: simulatedProduction,
		ProductionLoss:      productionLoss,
		PriceIncrease:       priceIncrease,
		EconomicImpact:      economicImpact,
		States:              simulatedStates,
	}, nil
}
