package handlers

import (
	"net/http"
	"strings"

	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/models"
	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/services"

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
