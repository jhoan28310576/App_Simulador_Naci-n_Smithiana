# 🏗️ Arquitectura del Sistema

## Diagrama de Arquitectura

```mermaid
graph TB
    subgraph "Cliente Web"
        UI[Interfaz Web<br/>drought_simulation.html]
        JS[JavaScript<br/>Frontend Logic]
    end

    subgraph "Servidor Go"
        Main[main.go<br/>Gin Router]
        Handlers[USDA Handlers<br/>HTTP Controllers]
        Services[USDA Service<br/>Business Logic]
        Models[Data Models<br/>Structs]
    end

    subgraph "API Externa"
        USDA[USDA Quick Stats API<br/>https://quickstats.nass.usda.gov]
    end

    subgraph "Datos"
        JSON[Response JSON<br/>Corn Production Data]
    end

    UI -->|HTTP GET/POST| Main
    JS -->|AJAX Calls| Main
    Main -->|Route Requests| Handlers
    Handlers -->|Call Methods| Services
    Services -->|HTTP GET| USDA
    USDA -->|JSON Response| Services
    Services -->|Parse Data| Models
    Models -->|Return Data| Services
    Services -->|Return Results| Handlers
    Handlers -->|JSON Response| Main
    Main -->|HTTP Response| UI
    Main -->|HTTP Response| JS

    style UI fill:#e1f5fe
    style Main fill:#f3e5f5
    style Services fill:#e8f5e8
    style USDA fill:#fff3e0
```

## Flujo de Datos

```mermaid
sequenceDiagram
    participant U as Usuario
    participant UI as Interfaz Web
    participant S as Servidor Go
    participant API as USDA API
    participant DB as Datos

    U->>UI: Configura parámetros
    UI->>S: POST /api/drought-simulation
    S->>API: GET /api_GET/?commodity_desc=CORN
    API->>S: JSON Response
    S->>S: Procesar datos
    S->>S: Aplicar simulación
    S->>S: Calcular impacto económico
    S->>UI: JSON Result
    UI->>U: Mostrar resultados
```

## Estructura de Componentes

```mermaid
graph LR
    subgraph "Frontend Layer"
        A[HTML Templates]
        B[CSS Styles]
        C[JavaScript Logic]
    end

    subgraph "API Layer"
        D[REST Endpoints]
        E[Request Validation]
        F[Response Formatting]
    end

    subgraph "Business Layer"
        G[Simulation Logic]
        H[Economic Calculations]
        I[Data Processing]
    end

    subgraph "External Layer"
        J[USDA API Client]
        K[HTTP Requests]
        L[JSON Parsing]
    end

    A --> D
    B --> D
    C --> D
    D --> E
    E --> F
    F --> G
    G --> H
    H --> I
    I --> J
    J --> K
    K --> L
```

## Modelo de Datos

```mermaid
erDiagram
    USDAQuickStatsResponse ||--o{ USDADataPoint : contains
    USDADataPoint {
        interface year
        string state_name
        string commodity_desc
        string statisticcat_desc
        string Value
        string unit_desc
        string reference_period_desc
        string source_desc
        string agg_level_desc
        string util_practice_desc
        string short_desc
    }

    DroughtSimulationParams ||--o{ CornProductionData : generates
    CornProductionData {
        string year
        string state
        float64 production
        string unit
        float64 area
        string area_unit
    }

    DroughtSimulationResult ||--|| DroughtSimulationParams : based_on
    DroughtSimulationResult {
        float64 original_production
        float64 simulated_production
        float64 production_loss
        float64 price_increase
        float64 economic_impact
    }
```

## Patrones de Diseño Utilizados

### 1. **MVC (Model-View-Controller)**
- **Model**: `internal/models/usda.go`
- **View**: `templates/drought_simulation.html`
- **Controller**: `internal/handlers/usda_handlers.go`

### 2. **Service Layer Pattern**
- **Service**: `internal/services/usda_service.go`
- **Responsabilidad**: Lógica de negocio y comunicación con API externa

### 3. **Repository Pattern** (Simplificado)
- **Data Access**: Manejo de datos de USDA API
- **Abstracción**: Separación de lógica de datos

### 4. **Dependency Injection**
- **Handlers**: Reciben servicios como dependencias
- **Services**: Reciben configuración como dependencias

## Consideraciones de Seguridad

```mermaid
graph TD
    A[Request] --> B{Validación de Entrada}
    B -->|Válido| C[Procesamiento]
    B -->|Inválido| D[Error 400]
    C --> E{API Key Válida}
    E -->|Sí| F[Consulta USDA]
    E -->|No| G[Error 401]
    F --> H[Sanitización de Datos]
    H --> I[Respuesta Segura]
```

## Escalabilidad

### Horizontal
- Múltiples instancias del servidor
- Load balancer
- Caché compartido (Redis)

### Vertical
- Optimización de consultas a API
- Caché de respuestas
- Compresión de datos

## Monitoreo y Logging

```mermaid
graph LR
    A[Request] --> B[Log Entry]
    B --> C[Performance Metrics]
    C --> D[Error Tracking]
    D --> E[Alert System]
```

## Deployment

```mermaid
graph TD
    A[Source Code] --> B[Build]
    B --> C[Test]
    C --> D[Package]
    D --> E[Deploy]
    E --> F[Monitor]
    F --> G[Scale]
``` 