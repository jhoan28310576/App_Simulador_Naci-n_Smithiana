package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Estructuras para el Capítulo 6
type Producto struct {
	ID               string            `json:"id"`
	Nombre           string            `json:"nombre"`
	Categoria        string            `json:"categoria"`
	PrecioMercado    float64           `json:"precio_mercado"`
	PrecioNatural    float64           `json:"precio_natural"`
	Componentes      ComponentesPrecio `json:"componentes"`
	HistorialPrecios HistorialPrecios  `json:"historial_precios"`
	Oferta           int               `json:"oferta"`
	Demanda          int               `json:"demanda"`
	Pais             string            `json:"pais"`
	Unidad           string            `json:"unidad"`
}

type ComponentesPrecio struct {
	Salarios   float64 `json:"salarios"`
	Beneficios float64 `json:"beneficios"`
	Rentas     float64 `json:"rentas"`
}

type HistorialPrecios struct {
	Natural []float64 `json:"natural"`
	Mercado []float64 `json:"mercado"`
}

type Mercado struct {
	ID            string               `json:"id"`
	Nombre        string               `json:"nombre"`
	Pais          string               `json:"pais"`
	Productos     []string             `json:"productos"`
	Configuracion ConfiguracionMercado `json:"configuracion"`
}

type ConfiguracionMercado struct {
	FactorOfertaDemanda float64 `json:"factor_oferta_demanda"`
	VolatilidadPrecio   float64 `json:"volatilidad_precio"`
	InflacionAnual      float64 `json:"inflacion_anual"`
}

type IndicadoresEconomicos struct {
	Venezuela IndicadorPais `json:"venezuela"`
	Colombia  IndicadorPais `json:"colombia"`
}

type IndicadorPais struct {
	InflacionAnual float64 `json:"inflacion_anual"`
	PibPerCapita   float64 `json:"pib_per_capita"`
	TasaDesempleo  float64 `json:"tasa_desempleo"`
	SalarioMinimo  float64 `json:"salario_minimo"`
}

type DataCap6 struct {
	Usuarios              []interface{}         `json:"usuarios"`
	Productos             []Producto            `json:"productos"`
	Mercados              []Mercado             `json:"mercados"`
	IndicadoresEconomicos IndicadoresEconomicos `json:"indicadores_economicos"`
}

var dataCap6 DataCap6

// Función para cargar datos del Capítulo 6
func CargarDatosCap6() error {
	// En un entorno real, cargarías desde un archivo
	// Por ahora, usaremos datos hardcodeados para demostración
	dataCap6 = DataCap6{
		Productos: []Producto{
			{
				ID:            "trigo",
				Nombre:        "Trigo",
				Categoria:     "agricultura",
				PrecioMercado: 250.00,
				PrecioNatural: 220.00,
				Componentes: ComponentesPrecio{
					Salarios:   112.50,
					Beneficios: 75.00,
					Rentas:     62.50,
				},
				HistorialPrecios: HistorialPrecios{
					Natural: []float64{200.00, 210.00, 215.00, 218.00, 220.00},
					Mercado: []float64{210.00, 225.00, 235.00, 245.00, 250.00},
				},
				Oferta:  1000,
				Demanda: 1200,
				Pais:    "VE",
				Unidad:  "tonelada",
			},
			{
				ID:            "maiz",
				Nombre:        "Maíz",
				Categoria:     "agricultura",
				PrecioMercado: 180.00,
				PrecioNatural: 160.00,
				Componentes: ComponentesPrecio{
					Salarios:   81.00,
					Beneficios: 54.00,
					Rentas:     45.00,
				},
				HistorialPrecios: HistorialPrecios{
					Natural: []float64{150.00, 155.00, 158.00, 159.00, 160.00},
					Mercado: []float64{160.00, 165.00, 170.00, 175.00, 180.00},
				},
				Oferta:  800,
				Demanda: 900,
				Pais:    "VE",
				Unidad:  "tonelada",
			},
			{
				ID:            "herramientas",
				Nombre:        "Herramientas",
				Categoria:     "manufactura",
				PrecioMercado: 45.00,
				PrecioNatural: 40.00,
				Componentes: ComponentesPrecio{
					Salarios:   20.25,
					Beneficios: 13.50,
					Rentas:     11.25,
				},
				HistorialPrecios: HistorialPrecios{
					Natural: []float64{38.00, 39.00, 39.50, 39.80, 40.00},
					Mercado: []float64{40.00, 42.00, 43.00, 44.00, 45.00},
				},
				Oferta:  500,
				Demanda: 600,
				Pais:    "VE",
				Unidad:  "unidad",
			},
		},
		Mercados: []Mercado{
			{
				ID:        "mercado_ve",
				Nombre:    "Mercado de Venezuela",
				Pais:      "VE",
				Productos: []string{"trigo", "maiz", "herramientas"},
				Configuracion: ConfiguracionMercado{
					FactorOfertaDemanda: 1.2,
					VolatilidadPrecio:   0.15,
					InflacionAnual:      0.25,
				},
			},
		},
		IndicadoresEconomicos: IndicadoresEconomicos{
			Venezuela: IndicadorPais{
				InflacionAnual: 25.0,
				PibPerCapita:   3500,
				TasaDesempleo:  15.0,
				SalarioMinimo:  150,
			},
			Colombia: IndicadorPais{
				InflacionAnual: 8.0,
				PibPerCapita:   6500,
				TasaDesempleo:  10.0,
				SalarioMinimo:  300,
			},
		},
	}
	return nil
}

// Handler para obtener todos los productos
func GetProductos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"productos": dataCap6.Productos,
	})
}

// Handler para obtener un producto específico
func GetProducto(c *gin.Context) {
	productoID := c.Param("id")

	for _, producto := range dataCap6.Productos {
		if producto.ID == productoID {
			c.JSON(http.StatusOK, gin.H{
				"success":  true,
				"producto": producto,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   "Producto no encontrado",
	})
}

// Handler para obtener productos por país
func GetProductosPorPais(c *gin.Context) {
	pais := c.Param("pais")

	var productosPais []Producto
	for _, producto := range dataCap6.Productos {
		if producto.Pais == pais {
			productosPais = append(productosPais, producto)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"productos": productosPais,
		"pais":      pais,
	})
}

// Handler para simular cambio en oferta/demanda
func SimularMercado(c *gin.Context) {
	productoID := c.Param("id")

	// Buscar el producto
	var producto *Producto
	for i := range dataCap6.Productos {
		if dataCap6.Productos[i].ID == productoID {
			producto = &dataCap6.Productos[i]
			break
		}
	}

	if producto == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Producto no encontrado",
		})
		return
	}

	// Obtener parámetros de la query
	ofertaStr := c.Query("oferta")
	demandaStr := c.Query("demanda")

	var nuevaOferta, nuevaDemanda int
	var err error

	if ofertaStr != "" {
		nuevaOferta, err = strconv.Atoi(ofertaStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Oferta debe ser un número válido",
			})
			return
		}
	} else {
		nuevaOferta = producto.Oferta
	}

	if demandaStr != "" {
		nuevaDemanda, err = strconv.Atoi(demandaStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Demanda debe ser un número válido",
			})
			return
		}
	} else {
		nuevaDemanda = producto.Demanda
	}

	// Calcular nuevo precio de mercado según Adam Smith
	precioAnterior := producto.PrecioMercado
	producto.Oferta = nuevaOferta
	producto.Demanda = nuevaDemanda

	// Fórmula simple: precio = precio_natural * (demanda / oferta)
	if nuevaOferta > 0 {
		factorMercado := float64(nuevaDemanda) / float64(nuevaOferta)
		producto.PrecioMercado = producto.PrecioNatural * factorMercado
	}

	// Actualizar componentes proporcionalmente
	producto.Componentes.Salarios = producto.PrecioMercado * 0.45
	producto.Componentes.Beneficios = producto.PrecioMercado * 0.30
	producto.Componentes.Rentas = producto.PrecioMercado * 0.25

	// Agregar al historial
	producto.HistorialPrecios.Natural = append(producto.HistorialPrecios.Natural, producto.PrecioNatural)
	producto.HistorialPrecios.Mercado = append(producto.HistorialPrecios.Mercado, producto.PrecioMercado)

	// Limitar historial a 10 elementos
	if len(producto.HistorialPrecios.Natural) > 10 {
		producto.HistorialPrecios.Natural = producto.HistorialPrecios.Natural[1:]
		producto.HistorialPrecios.Mercado = producto.HistorialPrecios.Mercado[1:]
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"producto": producto,
		"cambios": gin.H{
			"precio_anterior":  precioAnterior,
			"precio_nuevo":     producto.PrecioMercado,
			"oferta_anterior":  producto.Oferta,
			"oferta_nueva":     nuevaOferta,
			"demanda_anterior": producto.Demanda,
			"demanda_nueva":    nuevaDemanda,
		},
	})
}

// Handler para obtener indicadores económicos
func GetIndicadoresEconomicos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"indicadores": dataCap6.IndicadoresEconomicos,
	})
}

// Handler para obtener mercados
func GetMercados(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"mercados": dataCap6.Mercados,
	})
}

// Handler para análisis de componentes de precio
func AnalizarComponentes(c *gin.Context) {
	productoID := c.Param("id")

	var producto *Producto
	for i := range dataCap6.Productos {
		if dataCap6.Productos[i].ID == productoID {
			producto = &dataCap6.Productos[i]
			break
		}
	}

	if producto == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Producto no encontrado",
		})
		return
	}

	// Calcular porcentajes
	total := producto.PrecioMercado
	porcentajeSalarios := (producto.Componentes.Salarios / total) * 100
	porcentajeBeneficios := (producto.Componentes.Beneficios / total) * 100
	porcentajeRentas := (producto.Componentes.Rentas / total) * 100

	analisis := gin.H{
		"producto": producto,
		"analisis": gin.H{
			"porcentajes": gin.H{
				"salarios":   porcentajeSalarios,
				"beneficios": porcentajeBeneficios,
				"rentas":     porcentajeRentas,
			},
			"explicacion": gin.H{
				"salarios":   "Remuneración del trabajo (según Adam Smith)",
				"beneficios": "Remuneración del capital invertido",
				"rentas":     "Remuneración de la tierra y recursos naturales",
			},
			"precio_natural": producto.PrecioNatural,
			"precio_mercado": producto.PrecioMercado,
			"diferencia":     producto.PrecioMercado - producto.PrecioNatural,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analisis,
	})
}
