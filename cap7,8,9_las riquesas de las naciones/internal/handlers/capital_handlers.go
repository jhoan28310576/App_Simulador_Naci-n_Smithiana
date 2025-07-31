package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/capital"
	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/worldbank"
)

const worldBankAPIKey = "1F325726-42E7-3E08-8E7F-C7ED7047890A"

func SimularRetornoHandler(c *gin.Context) {
	var req struct {
		Capital      float64 `json:"capital"`
		Sector       string  `json:"sector"`
		Riesgo       float64 `json:"riesgo"`
		Competencia  float64 `json:"competencia"`
		SalarioMedio float64 `json:"salario_medio"`
		Pais         string  `json:"pais"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de inversión inválidos: " + err.Error()})
		return
	}

	// Convertimos el valor de riesgo a entero (redondeando) porque es un nivel (1-5)
	riesgoInt := int(req.Riesgo + 0.5)

	// Creamos la estructura de inversión con los nuevos campos
	inversion := capital.Inversion{
		Capital:      req.Capital,
		Sector:       req.Sector,
		Riesgo:       riesgoInt,
		Competencia:  req.Competencia, // Ahora es float64
		SalarioMedio: req.SalarioMedio,
	}

	// Si no se proporciona salario medio, consultar World Bank
	if inversion.SalarioMedio == 0 && req.Pais != "" {
		salario, err := worldbank.GetIndicator(req.Pais, "NY.GDP.PCAP.CD", worldBankAPIKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo obtener salario promedio del país: " + err.Error()})
			return
		}
		inversion.SalarioMedio = salario
	}

	añosStr := c.DefaultQuery("años", "5")
	años, err := strconv.Atoi(añosStr)
	if err != nil {
		años = 5 // Valor por defecto si hay error
	}

	// Llamamos a SimularRetorno con la nueva estructura
	resultado := capital.SimularRetorno(inversion, años)
	c.JSON(http.StatusOK, gin.H{"historial": resultado, "salario_medio": inversion.SalarioMedio})
}
