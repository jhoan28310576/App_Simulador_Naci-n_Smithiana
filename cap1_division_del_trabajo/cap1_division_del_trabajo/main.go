package main

import (
	"net/http"
	"strconv"

	"cap1_division_del_trabajo/internal/database"
	"cap1_division_del_trabajo/internal/handlers"
	"cap1_division_del_trabajo/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	router := gin.Default()

	// Serve static files from the "assets" directory
	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Simulador de División del Trabajo",
		})
	})

	// Ruta para la página de trueque
	router.GET("/trueque", func(c *gin.Context) {
		c.HTML(http.StatusOK, "trueque.html", gin.H{
			"title": "Sistema de Trueque Inteligente",
		})
	})

	// Ruta para la página de expansión de mercado
	router.GET("/expansion", func(c *gin.Context) {
		c.HTML(http.StatusOK, "expansion.html", gin.H{
			"title": "Módulo de Expansión de Mercado",
		})
	})

	// Ruta para la página de dinero
	router.GET("/dinero", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dinero.html", gin.H{
			"title": "Simulador del Dinero - Capítulo 4",
		})
	})

	// Ruta para la página de precios reales vs nominales
	router.GET("/precios", func(c *gin.Context) {
		c.HTML(http.StatusOK, "precios.html", gin.H{
			"title": "Precios Reales vs Nominales - Capítulo 5",
		})
	})

	// Endpoint para obtener todos los usuarios
	router.GET("/api/users", func(c *gin.Context) {
		users := database.GetAllUsers()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"users":   users,
		})
	})

	// Endpoint para probar la conexión con data.json
	router.GET("/api/test-connection", func(c *gin.Context) {
		user := database.GetUserByID("1")
		if user != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"user":    user,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "No se pudo obtener el usuario de prueba",
			})
		}
	})

	// Endpoint para buscar usuario por ID
	router.GET("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		user := database.GetUserByID(id)
		if user != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"user":    user,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Usuario no encontrado",
			})
		}
	})

	// Endpoint para buscar usuarios por rol
	router.GET("/api/users/role/:rol", func(c *gin.Context) {
		rol := c.Param("rol")
		users := database.GetUsersByRole(rol)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"users":   users,
			"count":   len(users),
		})
	})

	// Endpoint para buscar usuarios por especialización
	router.GET("/api/users/specialization/:especializacion", func(c *gin.Context) {
		especializacion := c.Param("especializacion")
		users := database.GetUsersBySpecialization(especializacion)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"users":   users,
			"count":   len(users),
		})
	})

	// Endpoint para obtener estadísticas generales
	router.GET("/api/stats", func(c *gin.Context) {
		stats := database.GetGeneralStats()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"stats":   stats,
		})
	})

	// ===== ENDPOINTS DEL SISTEMA DE TRUEQUE INTELIGENTE =====

	// Endpoint para buscar intercambios viables para un usuario
	router.GET("/api/trueque/intercambios/:usuarioID", func(c *gin.Context) {
		usuarioID := c.Param("usuarioID")
		id, err := strconv.Atoi(usuarioID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "ID de usuario inválido",
			})
			return
		}

		intercambios := database.BuscarIntercambiosViables(id)
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"intercambios": intercambios,
			"count":        len(intercambios),
		})
	})

	// Endpoint para obtener todas las ofertas de trueque activas
	router.GET("/api/trueque/ofertas", func(c *gin.Context) {
		ofertas := database.ObtenerOfertasTrueque()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"ofertas": ofertas,
			"count":   len(ofertas),
		})
	})

	// Endpoint para obtener ofertas de trueque de un usuario específico
	router.GET("/api/trueque/ofertas/:usuarioID", func(c *gin.Context) {
		usuarioID := c.Param("usuarioID")
		id, err := strconv.Atoi(usuarioID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "ID de usuario inválido",
			})
			return
		}

		// Filtrar ofertas por usuario
		todasOfertas := database.ObtenerOfertasTrueque()
		var ofertasUsuario []database.OfertaTrueque
		for _, oferta := range todasOfertas {
			if oferta.UsuarioID == id {
				ofertasUsuario = append(ofertasUsuario, oferta)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"ofertas": ofertasUsuario,
			"count":   len(ofertasUsuario),
		})
	})

	// Endpoint para calcular el valor de un producto
	router.GET("/api/trueque/valor/:producto/:cantidad", func(c *gin.Context) {
		producto := c.Param("producto")
		cantidadStr := c.Param("cantidad")

		cantidad, err := strconv.Atoi(cantidadStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Cantidad inválida",
			})
			return
		}

		valor := database.CalcularValorProducto(producto, cantidad)
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"producto": producto,
			"cantidad": cantidad,
			"valor":    valor,
			"unidad":   "horas de trabajo",
		})
	})

	// Endpoint para obtener información de trueque de un usuario
	router.GET("/api/trueque/usuario/:usuarioID", func(c *gin.Context) {
		usuarioID := c.Param("usuarioID")
		usuario := database.GetUserByID(usuarioID)
		if usuario == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Usuario no encontrado",
			})
			return
		}

		// Calcular valores de inventario
		valorTrigo := database.CalcularValorProducto("trigo", usuario.Inventario.Trigo)
		valorHerramientas := database.CalcularValorProducto("herramientas", usuario.Inventario.Herramientas)
		valorDinero := database.CalcularValorProducto("dinero", int(usuario.Inventario.Dinero))

		infoTrueque := gin.H{
			"usuario": gin.H{
				"id":              usuario.ID,
				"nombre":          usuario.Nombre,
				"rol":             usuario.Rol,
				"especializacion": usuario.Especializacion,
				"productividad":   usuario.Productividad,
			},
			"inventario": gin.H{
				"trigo": gin.H{
					"cantidad": usuario.Inventario.Trigo,
					"valor":    valorTrigo,
				},
				"herramientas": gin.H{
					"cantidad": usuario.Inventario.Herramientas,
					"valor":    valorHerramientas,
				},
				"dinero": gin.H{
					"cantidad": usuario.Inventario.Dinero,
					"valor":    valorDinero,
				},
			},
			"valor_total_inventario": valorTrigo + valorHerramientas + valorDinero,
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    infoTrueque,
		})
	})

	// Endpoint para obtener estadísticas del mercado de trueque
	router.GET("/api/trueque/estadisticas", func(c *gin.Context) {
		ofertas := database.ObtenerOfertasTrueque()

		// Calcular estadísticas
		totalOfertas := len(ofertas)
		productosOfrecidos := make(map[string]int)
		productosBuscados := make(map[string]int)
		valorTotalOfrecido := 0.0
		valorTotalBuscado := 0.0

		for _, oferta := range ofertas {
			productosOfrecidos[oferta.ProductoOfrece]++
			productosBuscados[oferta.ProductoBusca]++
			valorTotalOfrecido += oferta.ValorOfrece
			valorTotalBuscado += oferta.ValorBusca
		}

		estadisticas := gin.H{
			"total_ofertas":        totalOfertas,
			"productos_ofrecidos":  productosOfrecidos,
			"productos_buscados":   productosBuscados,
			"valor_total_ofrecido": valorTotalOfrecido,
			"valor_total_buscado":  valorTotalBuscado,
			"balance_mercado":      valorTotalOfrecido - valorTotalBuscado,
		}

		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"estadisticas": estadisticas,
		})
	})

	// ===== ENDPOINTS DEL SISTEMA DE MONEDAS (CAPÍTULO 4) =====

	// Endpoint para consultar los valores actuales de las monedas
	router.GET("/api/monedas/valores", func(c *gin.Context) {
		valores := make(map[string]gin.H)
		for k, v := range database.Monedas {
			valores[k] = gin.H{
				"nombre":             v.Nombre,
				"valor_relativo_oro": v.ValorRelativoOro,
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"valores": valores,
		})
	})

	// Endpoint para consultar el historial de valores de una moneda
	router.GET("/api/monedas/historial/:moneda", func(c *gin.Context) {
		moneda := c.Param("moneda")
		m, ok := database.Monedas[moneda]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Moneda no encontrada",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"historial": m.HistorialValores,
		})
	})

	// Endpoint para convertir entre monedas
	router.GET("/api/monedas/convertir/:cantidad/:origen/:destino", func(c *gin.Context) {
		cantidadStr := c.Param("cantidad")
		origen := c.Param("origen")
		destino := c.Param("destino")
		cantidad, err := strconv.ParseFloat(cantidadStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Cantidad inválida",
			})
			return
		}
		resultado := database.ConvertirMoneda(cantidad, origen, destino)
		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"cantidad_origen": cantidad,
			"moneda_origen":   origen,
			"moneda_destino":  destino,
			"resultado":       resultado,
		})
	})

	// ===== ENDPOINTS DEL SISTEMA DE PRECIOS DUAL (CAPÍTULO 5) =====

	// Endpoint para obtener todos los productos con precios reales y nominales
	router.GET("/api/precios/productos", func(c *gin.Context) {
		productos := models.CompararPreciosRealVsNominal()
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"productos": productos,
		})
	})

	// Endpoint para obtener datos de inflación del World Bank
	router.GET("/api/precios/inflacion", func(c *gin.Context) {
		estadisticas := models.ObtenerEstadisticasInflacion()
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"estadisticas": estadisticas,
		})
	})

	// Endpoint para actualizar precios nominales basado en inflación
	router.GET("/api/precios/actualizar/:pais", func(c *gin.Context) {
		pais := c.Param("pais")

		// Obtener datos de inflación
		datosInflacion, err := models.ObtenerDatosInflacionWorldBank(pais, "FP.CPI.TOTL")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Error al obtener datos de inflación: " + err.Error(),
			})
			return
		}

		// Calcular factor de inflación
		factorInflacion, err := models.CalcularFactorInflacion(datosInflacion, "2010")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Error al calcular factor de inflación: " + err.Error(),
			})
			return
		}

		// Actualizar precios de todos los productos
		for _, producto := range models.Productos {
			producto.ActualizarPrecioNominal(factorInflacion)
		}

		c.JSON(http.StatusOK, gin.H{
			"success":                true,
			"pais":                   pais,
			"factor_inflacion":       factorInflacion,
			"productos_actualizados": len(models.Productos),
		})
	})

	// Endpoint para calcular poder adquisitivo
	router.GET("/api/precios/poder-adquisitivo/:producto/:cantidad", func(c *gin.Context) {
		productoID := c.Param("producto")
		cantidadStr := c.Param("cantidad")

		cantidad, err := strconv.ParseFloat(cantidadStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Cantidad inválida",
			})
			return
		}

		producto, existe := models.Productos[productoID]
		if !existe {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Producto no encontrado",
			})
			return
		}

		poderAdquisitivo := producto.CalcularPoderAdquisitivo(cantidad)

		c.JSON(http.StatusOK, gin.H{
			"success":           true,
			"producto":          producto.Nombre,
			"cantidad_dinero":   cantidad,
			"moneda":            producto.Moneda,
			"poder_adquisitivo": poderAdquisitivo,
			"unidades":          poderAdquisitivo,
		})
	})

	// Endpoint para obtener historial de precios de un producto
	router.GET("/api/precios/historial/:producto", func(c *gin.Context) {
		productoID := c.Param("producto")

		producto, existe := models.Productos[productoID]
		if !existe {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Producto no encontrado",
			})
			return
		}

		historial := producto.HistorialPrecios
		if historial == nil {
			historial = []models.PrecioHistorico{}
		}

		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"producto":  producto.Nombre,
			"historial": historial,
		})
	})

	// ===== ENDPOINTS DEL CAPÍTULO 6 - COMPONENTES DEL PRECIO =====

	// Cargar datos del Capítulo 6
	handlers.CargarDatosCap6()

	// Ruta para la página del simulador del Capítulo 6
	router.GET("/cap6", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cap6_simulador.html", gin.H{
			"title": "Simulador de Componentes del Precio - Capítulo 6",
		})
	})

	// Endpoint para obtener todos los productos del Capítulo 6
	router.GET("/api/cap6/productos", handlers.GetProductos)

	// Endpoint para obtener un producto específico del Capítulo 6
	router.GET("/api/cap6/producto/:id", handlers.GetProducto)

	// Endpoint para obtener productos por país del Capítulo 6
	router.GET("/api/cap6/productos-pais/:pais", handlers.GetProductosPorPais)

	// Endpoint para simular cambios en oferta/demanda del Capítulo 6
	router.GET("/api/cap6/simular/:id", handlers.SimularMercado)

	// Endpoint para obtener indicadores económicos del Capítulo 6
	router.GET("/api/cap6/indicadores", handlers.GetIndicadoresEconomicos)

	// Endpoint para obtener mercados del Capítulo 6
	router.GET("/api/cap6/mercados", handlers.GetMercados)

	// Endpoint para análisis de componentes de precio del Capítulo 6
	router.GET("/api/cap6/analizar/:id", handlers.AnalizarComponentes)

	router.Run(":0808")
}
