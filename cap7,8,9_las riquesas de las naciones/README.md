# 🌾 Simulación de Crisis de Oferta - Las Riquezas de las Naciones

> Aplicación web para simular crisis de oferta en el mercado del maíz, basada en el Capítulo 7 de "Las Riquezas de las Naciones" de Adam Smith.

## 🎯 Descripción

Esta aplicación simula el impacto de sequías en el cinturón maicero de EE.UU. utilizando datos reales de la USDA Quick Stats API. Permite ajustar parámetros como severidad de la sequía y área afectada para ver el impacto económico en tiempo real.

## 🚀 Características

- **Datos Reales**: Integración con USDA Quick Stats API
- **Simulación Interactiva**: Parámetros ajustables en tiempo real
- **Cálculos Económicos**: Pérdida de producción, incremento de precios, impacto económico
- **Interfaz Moderna**: Diseño responsivo y intuitivo
- **API REST**: Endpoints para integración con otros sistemas

## 🛠️ Tecnologías

- **Backend**: Go 1.23.2 + Gin Framework
- **Frontend**: HTML5, CSS3, JavaScript ES6+
- **API Externa**: USDA Quick Stats API
- **Formato**: JSON

## 📦 Instalación

### Prerrequisitos
- Go 1.23.2 o superior
- API Key de USDA Quick Stats

### Pasos
1. Clonar el repositorio
```bash
git clone <repository-url>
cd cap7,8,9_las_riquesas_de_las_naciones
```

2. Instalar dependencias
```bash
go mod tidy
```

3. Configurar API Key (ya incluida en el código)
```go
// En internal/services/usda_service.go
const USDA_API_KEY = "1F325726-42E7-3E08-8E7F-C7ED7047890A"
```

4. Ejecutar la aplicación
```bash
go run cmd/webserver/main.go
```

5. Acceder a la aplicación
```
http://localhost:8080/simulation
```

## 📡 API Endpoints

### GET `/api/corn-production`
Obtiene datos de producción de maíz.

**Parámetros:**
- `year` (opcional): Año de referencia (default: 2023)
- `states` (opcional): Estados separados por coma (default: cinturón maicero)

**Ejemplo:**
```bash
curl "http://localhost:8080/api/corn-production?year=2023&states=IA,IL,NE"
```

### GET `/api/corn-production/:state`
Obtiene datos de producción para un estado específico.

**Ejemplo:**
```bash
curl "http://localhost:8080/api/corn-production/IA?year=2023"
```

### POST `/api/drought-simulation`
Ejecuta simulación de sequía.

**Body:**
```json
{
  "year": "2023",
  "states": ["IA", "IL", "NE"],
  "drought_severity": 0.3,
  "affected_area": 0.5
}
```

## 🎮 Uso de la Interfaz Web

1. **Acceder**: `http://localhost:8080/simulation`
2. **Configurar Parámetros**:
   - Año de referencia
   - Estados del cinturón maicero
   - Severidad de la sequía (0-100%)
   - Área afectada (0-100%)
3. **Ejecutar Simulación**: Hacer clic en "Ejecutar Simulación"
4. **Ver Resultados**: Métricas económicas y datos por estado

## 📊 Modelo Económico

### Cálculo de Producción Simulada
```
Producción Afectada = Producción Original × Área Afectada × Severidad
Producción No Afectada = Producción Original × (1 - Área Afectada)
Producción Final = Producción No Afectada + (Producción Afectada × (1 - Severidad))
```

### Cálculo de Impacto Económico
```
Pérdida de Producción = Producción Original - Producción Simulada
Incremento de Precio = (Pérdida / Producción Original) × 0.5
Impacto Económico = Pérdida × $4.50/bushel
```

## 🏗️ Estructura del Proyecto

```
├── cmd/webserver/
│   └── main.go                 # Punto de entrada
├── internal/
│   ├── models/
│   │   └── usda.go            # Modelos de datos
│   ├── services/
│   │   └── usda_service.go    # Lógica de negocio
│   └── handlers/
│       └── usda_handlers.go   # Controladores HTTP
├── templates/
│   ├── index.html             # Página principal
│   └── drought_simulation.html # Interfaz de simulación
├── assets/                    # Archivos estáticos
├── go.mod                     # Dependencias
└── README.md                  # Este archivo
```

## 🔧 Configuración

### Variables de Entorno (Opcional)
```bash
export GIN_MODE=release  # Para producción
export PORT=8080         # Puerto del servidor
```

### Personalización
- **Estados por Defecto**: Modificar en `internal/handlers/usda_handlers.go`
- **Parámetros Económicos**: Ajustar en `internal/services/usda_service.go`
- **Interfaz**: Personalizar `templates/drought_simulation.html`

## 🧪 Pruebas

### Prueba Rápida
```bash
# Verificar que el servidor responde
curl http://localhost:8080/

# Probar endpoint de producción
curl http://localhost:8080/api/corn-production?year=2023&states=IA

# Probar simulación
curl -X POST http://localhost:8080/api/drought-simulation \
  -H "Content-Type: application/json" \
  -d '{"year":"2023","states":["IA"],"drought_severity":0.3,"affected_area":0.5}'
```

## 📈 Ejemplo de Resultados

```json
{
  "simulation_params": {
    "year": "2023",
    "states": ["IA", "IL", "NE"],
    "drought_severity": 0.3,
    "affected_area": 0.5
  },
  "result": {
    "original_production": 12500000000,
    "simulated_production": 10625000000,
    "production_loss": 1875000000,
    "price_increase": 0.075,
    "economic_impact": 8437500000,
    "states": [...]
  }
}
```

## 🐛 Solución de Problemas

### Error de API
- Verificar conectividad a internet
- Confirmar que la API Key es válida
- Revisar logs del servidor

### Datos No Encontrados
- Verificar que el año solicitado tiene datos disponibles
- Confirmar que los códigos de estado son correctos
- Revisar filtros en el servicio

### Errores de Cálculo
- Verificar que los parámetros están en rangos válidos (0-1)
- Revisar la lógica de simulación en el servicio

## 🤝 Contribuciones

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 📚 Referencias

- [Adam Smith - Las Riquezas de las Naciones](https://es.wikipedia.org/wiki/La_riqueza_de_las_naciones)
- [USDA Quick Stats API](https://quickstats.nass.usda.gov/api)
- [Gin Framework](https://gin-gonic.com/)
- [Go Documentation](https://golang.org/doc/)

## 📞 Contacto

Para preguntas o soporte, crear un issue en el repositorio.

---
muestratra cap7


https://github.com/user-attachments/assets/daf7dd9d-2212-4848-a9c0-ec0e9e24bea8


**Desarrollado con ❤️ para el estudio de la economía clásica** 

--- 

siguientes capitulos en desarrollo 8,9
