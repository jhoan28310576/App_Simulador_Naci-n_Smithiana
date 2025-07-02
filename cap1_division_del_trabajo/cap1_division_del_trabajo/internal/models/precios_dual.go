package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Producto representa un bien con valor real (horas de trabajo) y nominal (precio en moneda)
type Producto struct {
	ID               string            `json:"id"`
	Nombre           string            `json:"nombre"`
	HorasTrabajo     float64           `json:"horas_trabajo"`  // Valor real
	PrecioNominal    float64           `json:"precio_nominal"` // Valor nominal actual
	Moneda           string            `json:"moneda"`
	Pais             string            `json:"pais"`
	FechaCreacion    time.Time         `json:"fecha_creacion"`
	HistorialPrecios []PrecioHistorico `json:"historial_precios"`
}

// PrecioHistorico registra cambios en precios nominales
type PrecioHistorico struct {
	Fecha           time.Time `json:"fecha"`
	PrecioNominal   float64   `json:"precio_nominal"`
	FactorInflacion float64   `json:"factor_inflacion"`
	PrecioReal      float64   `json:"precio_real"`
}

// DatosInflacion representa datos de inflación del World Bank
type DatosInflacion struct {
	Pais      string          `json:"pais"`
	Indicador string          `json:"indicador"`
	Datos     []DatoInflacion `json:"datos"`
}

type DatoInflacion struct {
	Ano   string  `json:"ano"`
	Valor float64 `json:"valor"`
}

// Productos globales para el simulador
var Productos = map[string]*Producto{
	"trigo": {
		ID:            "trigo",
		Nombre:        "Trigo",
		HorasTrabajo:  2.0,
		PrecioNominal: 10.0,
		Moneda:        "USD",
		Pais:          "Venezuela",
		FechaCreacion: time.Now(),
	},
	"herramientas": {
		ID:            "herramientas",
		Nombre:        "Herramientas",
		HorasTrabajo:  8.0,
		PrecioNominal: 40.0,
		Moneda:        "USD",
		Pais:          "Venezuela",
		FechaCreacion: time.Now(),
	},
	"ropa": {
		ID:            "ropa",
		Nombre:        "Ropa",
		HorasTrabajo:  4.0,
		PrecioNominal: 20.0,
		Moneda:        "USD",
		Pais:          "Venezuela",
		FechaCreacion: time.Now(),
	},
	"vivienda": {
		ID:            "vivienda",
		Nombre:        "Vivienda",
		HorasTrabajo:  100.0,
		PrecioNominal: 500.0,
		Moneda:        "USD",
		Pais:          "Venezuela",
		FechaCreacion: time.Now(),
	},
}

// ObtenerDatosInflacionWorldBank obtiene datos de inflación del World Bank API
func ObtenerDatosInflacionWorldBank(pais, indicador string) (*DatosInflacion, error) {
	// URL del World Bank API para CPI (Consumer Price Index)
	url := fmt.Sprintf("http://api.worldbank.org/v2/country/%s/indicator/%s?format=json&per_page=20", pais, indicador)

	fmt.Printf("Consultando API para %s: %s\n", pais, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al obtener datos de inflación: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	// Parsear respuesta del World Bank API
	var response []interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error al parsear JSON: %v", err)
	}

	if len(response) < 2 {
		return nil, fmt.Errorf("respuesta inesperada del API")
	}

	// Extraer datos del segundo elemento del array
	datosRaw, ok := response[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de datos inesperado")
	}

	fmt.Printf("Encontrados %d registros para %s\n", len(datosRaw), pais)

	var datos []DatoInflacion
	for i, datoRaw := range datosRaw {
		datoMap, ok := datoRaw.(map[string]interface{})
		if !ok {
			fmt.Printf("Registro %d no es un mapa válido\n", i)
			continue
		}

		// Extraer año y valor
		ano, ok1 := datoMap["date"].(string)
		if !ok1 {
			fmt.Printf("Registro %d: año no encontrado\n", i)
			continue
		}

		// El valor puede venir como string, número o nulo
		var valor float64
		switch v := datoMap["value"].(type) {
		case string:
			if v == "" {
				fmt.Printf("Registro %d: valor vacío\n", i)
				continue // Ignorar valores vacíos
			}
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				fmt.Printf("Registro %d: error parseando valor '%s': %v\n", i, v, err)
				continue // Ignorar valores que no se pueden parsear
			}
			valor = parsed
		case float64:
			valor = v
		case int:
			valor = float64(v)
		case nil:
			fmt.Printf("Registro %d: valor nulo\n", i)
			continue // Ignorar valores nulos
		default:
			fmt.Printf("Registro %d: tipo de valor inesperado: %T\n", i, v)
			continue // Ignorar otros tipos
		}

		fmt.Printf("Registro %d: %s = %.2f\n", i, ano, valor)
		datos = append(datos, DatoInflacion{
			Ano:   ano,
			Valor: valor,
		})
	}

	fmt.Printf("Total de datos válidos para %s: %d\n", pais, len(datos))

	return &DatosInflacion{
		Pais:      pais,
		Indicador: indicador,
		Datos:     datos,
	}, nil
}

// CalcularFactorInflacion calcula el factor de inflación acumulada
func CalcularFactorInflacion(datos *DatosInflacion, anoBase string) (float64, error) {
	if len(datos.Datos) == 0 {
		return 1.0, nil
	}

	// Encontrar el valor base
	var valorBase float64
	encontrado := false
	for _, dato := range datos.Datos {
		if dato.Ano == anoBase {
			valorBase = dato.Valor
			encontrado = true
			break
		}
	}

	if !encontrado || valorBase == 0 {
		// Usar el valor más antiguo disponible como base
		valorBase = datos.Datos[len(datos.Datos)-1].Valor
	}

	// Calcular factor de inflación (valor actual / valor base)
	valorActual := datos.Datos[0].Valor
	factor := valorActual / valorBase

	return factor, nil
}

// ActualizarPrecioNominal actualiza el precio nominal basado en inflación
func (p *Producto) ActualizarPrecioNominal(factorInflacion float64) {
	// El precio real se mantiene constante (horas de trabajo)
	// El precio nominal se ajusta por inflación
	precioReal := p.HorasTrabajo * 5.0 // 5 USD por hora de trabajo como base

	// Nuevo precio nominal = precio real * factor de inflación
	nuevoPrecio := precioReal * factorInflacion

	// Registrar en historial
	historial := PrecioHistorico{
		Fecha:           time.Now(),
		PrecioNominal:   nuevoPrecio,
		FactorInflacion: factorInflacion,
		PrecioReal:      precioReal,
	}

	p.HistorialPrecios = append(p.HistorialPrecios, historial)
	p.PrecioNominal = nuevoPrecio
}

// ObtenerPrecioReal devuelve el precio real en horas de trabajo
func (p *Producto) ObtenerPrecioReal() float64 {
	return p.HorasTrabajo
}

// ObtenerPrecioNominal devuelve el precio nominal actual
func (p *Producto) ObtenerPrecioNominal() float64 {
	return p.PrecioNominal
}

// CalcularPoderAdquisitivo calcula cuántas unidades se pueden comprar con cierta cantidad de dinero
func (p *Producto) CalcularPoderAdquisitivo(cantidadDinero float64) float64 {
	if p.PrecioNominal <= 0 {
		return 0
	}
	return cantidadDinero / p.PrecioNominal
}

// CompararPreciosRealVsNominal compara precios reales y nominales
func CompararPreciosRealVsNominal() map[string]interface{} {
	resultado := make(map[string]interface{})

	for id, producto := range Productos {
		precioReal := producto.ObtenerPrecioReal()
		precioNominal := producto.ObtenerPrecioNominal()

		resultado[id] = map[string]interface{}{
			"nombre":         producto.Nombre,
			"precio_real":    precioReal,
			"precio_nominal": precioNominal,
			"unidad_real":    "horas de trabajo",
			"unidad_nominal": producto.Moneda,
			"relacion":       precioNominal / precioReal,
		}
	}

	return resultado
}

// ObtenerEstadisticasInflacion obtiene estadísticas de inflación
func ObtenerEstadisticasInflacion() map[string]interface{} {
	// Obtener datos de inflación para Venezuela
	datosVenezuela, err := ObtenerDatosInflacionWorldBank("VE", "FP.CPI.TOTL")
	if err != nil {
		return map[string]interface{}{
			"error": "No se pudieron obtener datos de inflación de Venezuela: " + err.Error(),
		}
	}

	// Obtener datos de inflación para Colombia
	datosColombia, err2 := ObtenerDatosInflacionWorldBank("CO", "FP.CPI.TOTL")
	if err2 != nil {
		return map[string]interface{}{
			"error": "No se pudieron obtener datos de inflación de Colombia: " + err2.Error(),
		}
	}

	// Verificar que hay datos disponibles
	if len(datosVenezuela.Datos) == 0 || len(datosColombia.Datos) == 0 {
		return map[string]interface{}{
			"error": "No hay datos de inflación disponibles",
		}
	}

	// Calcular factores de inflación usando años disponibles
	anoBaseVenezuela := "2015" // Usar 2015 como base para Venezuela
	if len(datosVenezuela.Datos) > 0 {
		anoBaseVenezuela = datosVenezuela.Datos[len(datosVenezuela.Datos)-1].Ano // Usar el año más antiguo disponible
	}

	anoBaseColombia := "2015" // Usar 2015 como base para Colombia
	if len(datosColombia.Datos) > 0 {
		anoBaseColombia = datosColombia.Datos[len(datosColombia.Datos)-1].Ano // Usar el año más antiguo disponible
	}

	factorVenezuela, _ := CalcularFactorInflacion(datosVenezuela, anoBaseVenezuela)
	factorColombia, _ := CalcularFactorInflacion(datosColombia, anoBaseColombia)

	// Preparar datos recientes (máximo 5 años)
	datosRecientesVenezuela := datosVenezuela.Datos
	if len(datosRecientesVenezuela) > 5 {
		datosRecientesVenezuela = datosRecientesVenezuela[:5]
	}

	datosRecientesColombia := datosColombia.Datos
	if len(datosRecientesColombia) > 5 {
		datosRecientesColombia = datosRecientesColombia[:5]
	}

	return map[string]interface{}{
		"venezuela": map[string]interface{}{
			"factor_inflacion": factorVenezuela,
			"datos_recientes":  datosRecientesVenezuela,
			"ano_base":         anoBaseVenezuela,
		},
		"colombia": map[string]interface{}{
			"factor_inflacion": factorColombia,
			"datos_recientes":  datosRecientesColombia,
			"ano_base":         anoBaseColombia,
		},
		"comparacion": map[string]interface{}{
			"diferencia_factores": factorVenezuela - factorColombia,
			"ratio_factores":      factorVenezuela / factorColombia,
		},
	}
}
