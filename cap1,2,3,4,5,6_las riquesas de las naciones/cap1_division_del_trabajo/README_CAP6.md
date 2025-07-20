# 📊 Simulador de Componentes del Precio - Capítulo 6 (Adam Smith)

> **"De los componentes del precio de las mercancías"** - La Riqueza de las Naciones

## 🎯 Descripción

Simulador interactivo que demuestra los principios económicos del **Capítulo 6** de Adam Smith, mostrando cómo se descompone el precio de las mercancías en sus tres componentes fundamentales:

- **Salarios** (remuneración del trabajo)
- **Beneficios** (remuneración del capital) 
- **Rentas** (remuneración de la tierra)

## 🚀 Características

### ✨ Funcionalidades Principales
- **Desglose de Precios**: Visualización de los componentes del precio según Adam Smith
- **Simulador de Mercado**: Modificar oferta y demanda para ver cambios en precios
- **Historial de Precios**: Seguimiento de precios natural vs mercado
- **Indicadores Económicos**: Datos por país (Venezuela, Colombia)
- **Análisis Visual**: Gráficos interactivos con Chart.js

### 📈 Componentes del Precio
```json
{
  "precio_mercado": 250.00,
  "componentes": {
    "salarios": 112.50,    // 45% - Trabajo
    "beneficios": 75.00,   // 30% - Capital
    "rentas": 62.50        // 25% - Tierra
  }
}
```

## 🛠️ Tecnologías

- **Backend**: Go + Gin Framework
- **Frontend**: HTML5 + CSS3 + JavaScript
- **Gráficos**: Chart.js
- **UI**: Bootstrap 5
- **Iconos**: Font Awesome

## 📁 Estructura del Proyecto

```
cap6_simulador/
├── main.go                    # Servidor principal
├── internal/
│   ├── handlers/
│   │   └── cap6_handlers.go   # API endpoints
│   └── database/
│       └── data_cap6.json     # Datos de productos
├── templates/
│   └── cap6_simulador.html    # Interfaz web
└── assets/
    └── css/
        └── cap6.css           # Estilos específicos
```

## 🚀 Instalación y Uso

### Prerrequisitos
- Go 1.16+
- Navegador web moderno

### Ejecutar
```bash
# Clonar repositorio
git clone [url-del-repositorio]

# Navegar al directorio
cd cap1_division_del_trabajo/cap1_division_del_trabajo

# Ejecutar servidor
go run main.go
```

### Acceder
- **URL**: `http://localhost:8081/cap6`
- **API**: `http://localhost:8081/api/cap6/*`

## 📡 API Endpoints

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/api/cap6/productos` | GET | Listar todos los productos |
| `/api/cap6/producto/:id` | GET | Obtener producto específico |
| `/api/cap6/simular/:id` | GET | Simular cambios oferta/demanda |
| `/api/cap6/indicadores` | GET | Indicadores económicos |
| `/api/cap6/analizar/:id` | GET | Análisis de componentes |

### Ejemplo de Uso API
```bash
# Obtener productos
curl http://localhost:8081/api/cap6/productos

# Simular mercado para trigo
curl "http://localhost:8081/api/cap6/simular/trigo?oferta=1200&demanda=1000"
```

## 📊 Productos Incluidos

- **Trigo** (Venezuela) - $250/tonelada
- **Maíz** (Venezuela) - $180/tonelada  
- **Herramientas** (Venezuela) - $45/unidad
- **Café** (Colombia) - $320/tonelada
- **Bananas** (Colombia) - $1.20/kg

## 🎨 Características de la UI

- **Nav Blanco**: Diseño limpio con navegación blanca
- **Gradiente de Fondo**: Estilo moderno con gradientes
- **Cards Transparentes**: Efecto glassmorphism
- **Responsive**: Adaptable a móviles y tablets
- **Gráficos Interactivos**: Visualización de datos en tiempo real

## 🔧 Configuración

### Variables de Entorno
```bash
# Puerto del servidor (opcional)
PORT=8081
```

### Personalización de Datos
Editar `internal/database/data_cap6.json` para:
- Agregar nuevos productos
- Modificar precios y componentes
- Cambiar indicadores económicos

## 📚 Conceptos Económicos

### Precio Natural vs Precio de Mercado
- **Precio Natural**: Costo de producción (salarios + beneficios + rentas)
- **Precio de Mercado**: Determinado por oferta y demanda

### Fórmula de Simulación
```
Precio Mercado = Precio Natural × (Demanda / Oferta)
```

## 🤝 Contribuir

1. Fork el proyecto
2. Crear rama para feature (`git checkout -b feature/AmazingFeature`)
3. Commit cambios (`git commit -m 'Add AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 👨‍💻 Autor

**Adam Smith** - *Conceptos económicos originales*
**Desarrollador** - *Implementación técnica*

---

## 🎯 Próximas Mejoras

- [ ] Integración con APIs económicas reales
- [ ] Más países y productos
- [ ] Análisis de tendencias históricas
- [ ] Exportación de datos
- [ ] Modo oscuro/claro

---

⭐ **¡Dale una estrella si te gustó el proyecto!** 