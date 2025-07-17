package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/models"
	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/services"
	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/usdaapi"

	"github.com/gin-gonic/gin"
)

type USDAHandlers struct {
	usdaService *services.USDAService
}

func NewUSDAHandlers() *USDAHandlers {
	return &USDAHandlers{
		usdaService: services.NewUSDAService(),
	}
}

// GetCornProduction obtiene datos de producción de maíz
func (h *USDAHandlers) GetCornProduction(c *gin.Context) {
	year := c.Query("year")
	if year == "" {
		year = "2023" // Default year
	}

	statesParam := c.Query("states")
	if statesParam == "" {
		// Estados del cinturón maicero por defecto
		statesParam = "IA,IL,NE,MN,IN,OH,WI,SD,MO,KS"
	}

	states := []string{"IA", "IL", "NE", "MN", "IN", "OH", "WI", "SD", "MO", "KS"}
	if statesParam != "" {
		// Parse comma-separated states
		states = []string{}
		for _, state := range strings.Split(statesParam, ",") {
			states = append(states, strings.TrimSpace(state))
		}
	}

	data, err := h.usdaService.GetCornProduction(year, states)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo datos de producción: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"year":   year,
		"states": states,
		"data":   data,
	})
}

// SimulateDrought simula el impacto de una sequía
func (h *USDAHandlers) SimulateDrought(c *gin.Context) {
	var params models.DroughtSimulationParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	// Validar parámetros
	if params.Year == "" {
		params.Year = "2023"
	}
	if len(params.States) == 0 {
		params.States = []string{"IA", "IL", "NE", "MN", "IN", "OH", "WI", "SD", "MO", "KS"}
	}
	if params.DroughtSeverity < 0 || params.DroughtSeverity > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La severidad de la sequía debe estar entre 0 y 1",
		})
		return
	}
	if params.AffectedArea < 0 || params.AffectedArea > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El área afectada debe estar entre 0 y 1",
		})
		return
	}

	result, err := h.usdaService.SimulateDrought(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error en la simulación: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulation_params": params,
		"result":            result,
	})
}

// GetDroughtSimulationForm muestra un formulario para la simulación
func (h *USDAHandlers) GetDroughtSimulationForm(c *gin.Context) {
	c.HTML(http.StatusOK, "drought_simulation.html", gin.H{
		"title": "Simulación de Crisis de Oferta - Sequía en el Cinturón Maicero",
	})
}

// GetCornProductionByState obtiene datos de un estado específico
func (h *USDAHandlers) GetCornProductionByState(c *gin.Context) {
	state := c.Param("state")
	year := c.Query("year")
	if year == "" {
		year = "2023"
	}

	data, err := h.usdaService.GetCornProduction(year, []string{state})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo datos: " + err.Error(),
		})
		return
	}

	if len(data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No se encontraron datos para el estado " + state,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"state": state,
		"year":  year,
		"data":  data[0],
	})
}

// GetCap8Simulacion combina los tres ejemplos de simulación para el cap8
func (h *USDAHandlers) GetCap8Simulacion(c *gin.Context) {
	// Obtener parámetros de la query string
	anio := c.Query("anio")
	if anio == "" {
		anio = "2023"
	}
	estado := c.Query("estado")
	if estado == "" {
		estado = "IA"
	}

	// 1. Obtener valor de producción de maíz en el estado y año especificados
	data, err := usdaapi.GetUSDAData("CORN", anio, estado, "PRODUCTION")
	valorProduccion := 0.0
	if err == nil && len(data) > 0 {
		if v, ok := data[0]["Value"].(string); ok {
			v = strings.ReplaceAll(v, ",", "")
			valorProduccion, _ = strconv.ParseFloat(v, 64)
		}
	}

	// 1. Salario ajustado por valor de producción
	salarioBase := 15.0
	salarioAjustado := salarioBase
	if valorProduccion > 0 {
		salarioAjustado += valorProduccion / 1e10 // Ajuste arbitrario
	}

	// 2. Número de ofertas laborales según producción
	numOfertas := 2
	if valorProduccion > 0 {
		numOfertas += int(valorProduccion / 1e10)
	}

	// 3. Demanda laboral ajustada por año
	anios := []string{"2021", "2022", "2023"}
	demandaAnual := []map[string]interface{}{}
	apiRawDemanda := make(map[string]interface{})
	for _, anioConsulta := range anios {
		dataAnio, _ := usdaapi.GetUSDAData("CORN", anioConsulta, estado, "PRODUCTION")
		apiRawDemanda[anioConsulta] = dataAnio
		valorProdAnio := 0.0
		if len(dataAnio) > 0 {
			if v, ok := dataAnio[0]["Value"].(string); ok {
				v = strings.ReplaceAll(v, ",", "")
				valorProdAnio, _ = strconv.ParseFloat(v, 64)
			}
		}
		numDemandas := 2
		if valorProdAnio > 1e10 {
			numDemandas += 1
		}
		demandaAnual = append(demandaAnual, map[string]interface{}{
			"anio":            anioConsulta,
			"valorProduccion": valorProdAnio,
			"numDemandas":     numDemandas,
		})
	}

	// Responder con el JSON esperado por el frontend, incluyendo los datos crudos
	c.JSON(http.StatusOK, gin.H{
		"salarioAjustado": salarioAjustado,
		"valorProduccion": valorProduccion,
		"numOfertas":      numOfertas,
		"demandaAnual":    demandaAnual,
		"apiRaw1":         data,          // datos crudos para sección 1
		"apiRaw2":         data,          // datos crudos para sección 2 (igual que 1)
		"apiRawDemanda":   apiRawDemanda, // datos crudos para cada año en sección 3
	})
}
