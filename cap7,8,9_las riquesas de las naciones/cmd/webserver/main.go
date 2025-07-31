package main

import (
	"net/http"

	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Crear una nueva instancia de Gin
	r := gin.Default()

	// Inicializar handlers
	usdaHandlers := handlers.NewUSDAHandlers()

	// Servir archivos estáticos desde la carpeta assets
	r.Static("/assets", "./assets")

	// Servir archivos HTML desde la carpeta templates
	r.LoadHTMLGlob("templates/*")
	r.GET("/html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Las Riquezas de las Naciones",
		})
	})

	// API endpoints agrupados
	api := r.Group("/api")
	{
		// Endpoints para USDA (Capítulo 8)
		api.GET("/corn-production", usdaHandlers.GetCornProduction)
		api.GET("/corn-production/:state", usdaHandlers.GetCornProductionByState)
		api.POST("/drought-simulation", usdaHandlers.SimulateDrought)
		api.GET("/cap8/simulacion", usdaHandlers.GetCap8Simulacion)

		// Endpoint para simulación de capital (Capítulo 9)
		api.POST("/simular-retorno", handlers.SimularRetornoHandler)
	}

	// Página de simulación Capítulo 8
	r.GET("/simulation", usdaHandlers.GetDroughtSimulationForm)
	// Nueva ruta para la simulación del cap8
	r.GET("/cap8/simulacion", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cap8_simulacion.html", gin.H{})
	})
	// Nueva ruta para la simulación del cap9
	r.GET("/capital/simulacion", func(c *gin.Context) {
		c.HTML(http.StatusOK, "capital_simulacion.html", gin.H{})
	})

	// Iniciar el servidor en el puerto 8080
	r.Run(":8080")
}
