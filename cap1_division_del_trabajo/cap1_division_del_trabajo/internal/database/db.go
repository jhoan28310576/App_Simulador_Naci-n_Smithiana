package database

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Data structure to match data.js
type User struct {
	ID              int     `json:"id"`
	Nombre          string  `json:"nombre"`
	Rol             string  `json:"rol"`
	Especializacion *string `json:"especializacion"`
	Productividad   float64 `json:"productividad"`
	Saldo           float64 `json:"saldo"`
	Inventario      struct {
		Trigo        int     `json:"trigo"`
		Herramientas int     `json:"herramientas"`
		Dinero       float64 `json:"dinero"`
	} `json:"inventario"`
}

type Data struct {
	Usuarios []User `json:"usuarios"`
}

type GeneralStats struct {
	TotalUsuarios         int            `json:"total_usuarios"`
	PromedioProductividad float64        `json:"promedio_productividad"`
	TotalDinero           float64        `json:"total_dinero"`
	TotalTrigo            int            `json:"total_trigo"`
	TotalHerramientas     int            `json:"total_herramientas"`
	Roles                 map[string]int `json:"roles"`
}

var DB *Data

// Estructuras para el sistema de trueque
type Intercambio struct {
	ID              string  `json:"id"`
	UsuarioOrigen   int     `json:"usuario_origen"`
	UsuarioDestino  int     `json:"usuario_destino"`
	ProductoOrigen  string  `json:"producto_origen"`
	ProductoDestino string  `json:"producto_destino"`
	CantidadOrigen  int     `json:"cantidad_origen"`
	CantidadDestino int     `json:"cantidad_destino"`
	ValorOrigen     float64 `json:"valor_origen"`
	ValorDestino    float64 `json:"valor_destino"`
	Estado          string  `json:"estado"` // "pendiente", "aceptado", "rechazado", "completado"
	FechaCreacion   string  `json:"fecha_creacion"`
}

type OfertaTrueque struct {
	ID             string  `json:"id"`
	UsuarioID      int     `json:"usuario_id"`
	ProductoOfrece string  `json:"producto_ofrece"`
	CantidadOfrece int     `json:"cantidad_ofrece"`
	ProductoBusca  string  `json:"producto_busca"`
	CantidadBusca  int     `json:"cantidad_busca"`
	ValorOfrece    float64 `json:"valor_ofrece"`
	ValorBusca     float64 `json:"valor_busca"`
	Activa         bool    `json:"activa"`
}

// Constantes para el cálculo de valores
const (
	HORAS_POR_TRIGO       = 2.0 // 2 horas por unidad de trigo
	HORAS_POR_HERRAMIENTA = 1.0 // 1 hora por herramienta
	HORAS_POR_DINERO      = 0.1 // 0.1 horas por unidad de dinero (valor relativo)
)

// InitDB initializes the data from data.json
func InitDB() {
	// Commented MySQL connection code for future use
	/*
		var err error
		dsn := "root@tcp(localhost:3306)/division_trabajo?parseTime=true"
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		err = DB.Ping()
		if err != nil {
			log.Fatal(err)
		}

		createTables()
	*/

	// Read and parse data.json
	filePath := filepath.Join("internal", "database", "data.json")
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading data.json:", err)
	}

	// Parse JSON directly
	DB = &Data{}
	err = json.Unmarshal(file, DB)
	if err != nil {
		log.Fatal("Error parsing data.json:", err)
	}
}

// GetUserByID returns a user by their ID
func GetUserByID(id string) *User {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}

	for _, user := range DB.Usuarios {
		if user.ID == idInt {
			return &user
		}
	}
	return nil
}

// GetAllUsers returns all users
func GetAllUsers() []User {
	return DB.Usuarios
}

// GetUsersByRole returns all users with a specific role
func GetUsersByRole(rol string) []User {
	var users []User
	for _, user := range DB.Usuarios {
		if user.Rol == rol {
			users = append(users, user)
		}
	}
	return users
}

// GetUsersBySpecialization returns all users with a specific specialization
func GetUsersBySpecialization(especializacion string) []User {
	var users []User
	for _, user := range DB.Usuarios {
		if user.Especializacion != nil && *user.Especializacion == especializacion {
			users = append(users, user)
		}
	}
	return users
}

// GetGeneralStats returns general statistics about all users
func GetGeneralStats() GeneralStats {
	stats := GeneralStats{
		Roles: make(map[string]int),
	}

	for _, user := range DB.Usuarios {
		stats.TotalUsuarios++
		stats.PromedioProductividad += user.Productividad
		stats.TotalDinero += user.Inventario.Dinero
		stats.TotalTrigo += user.Inventario.Trigo
		stats.TotalHerramientas += user.Inventario.Herramientas
		stats.Roles[user.Rol]++
	}

	if stats.TotalUsuarios > 0 {
		stats.PromedioProductividad /= float64(stats.TotalUsuarios)
	}

	return stats
}

// CloseDB is kept for compatibility but doesn't need to do anything for data.js
func CloseDB() {
	// No cleanup needed for in-memory data
}

// Commented MySQL table creation code for future use
/*
func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createRolesTable := `
	CREATE TABLE IF NOT EXISTS roles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		description TEXT,
		productivity_multiplier FLOAT NOT NULL
	);`

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		base_production_time INT NOT NULL,
		role_id INT,
		FOREIGN KEY (role_id) REFERENCES roles(id)
	);`

	createInventoryTable := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		product_id INT NOT NULL,
		quantity INT NOT NULL DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	);`

	tables := []string{createUsersTable, createRolesTable, createProductsTable, createInventoryTable}
	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Fatal(err)
		}
	}
}
*/

/*
import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	// Conexión sin contraseña
	dsn := "root@tcp(localhost:3306)/division_trabajo?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar la conexión
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	createTables()
}

// createTables creates the necessary tables in the database
func createTables() {
	// Tabla de usuarios
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Tabla de roles (especializaciones)
	createRolesTable := `
	CREATE TABLE IF NOT EXISTS roles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		description TEXT,
		productivity_multiplier FLOAT NOT NULL
	);`

	// Tabla de productos
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		base_production_time INT NOT NULL,
		role_id INT,
		FOREIGN KEY (role_id) REFERENCES roles(id)
	);`

	// Tabla de inventario
	createInventoryTable := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		product_id INT NOT NULL,
		quantity INT NOT NULL DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	);`

	// Ejecutar las creaciones de tablas
	tables := []string{createUsersTable, createRolesTable, createProductsTable, createInventoryTable}
	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
*/

// CalcularValorProducto calcula el valor en horas de trabajo de un producto
func CalcularValorProducto(producto string, cantidad int) float64 {
	switch producto {
	case "trigo":
		return float64(cantidad) * HORAS_POR_TRIGO
	case "herramientas":
		return float64(cantidad) * HORAS_POR_HERRAMIENTA
	case "dinero":
		return float64(cantidad) * HORAS_POR_DINERO
	default:
		return 0
	}
}

// ObtenerCantidadProducto obtiene la cantidad disponible de un producto para un usuario
func ObtenerCantidadProducto(userID int, producto string) int {
	user := GetUserByID(strconv.Itoa(userID))
	if user == nil {
		return 0
	}

	switch producto {
	case "trigo":
		return user.Inventario.Trigo
	case "herramientas":
		return user.Inventario.Herramientas
	case "dinero":
		return int(user.Inventario.Dinero)
	default:
		return 0
	}
}

// BuscarIntercambiosViables encuentra intercambios posibles entre usuarios
func BuscarIntercambiosViables(usuarioID int) []Intercambio {
	var intercambios []Intercambio
	usuarioOrigen := GetUserByID(strconv.Itoa(usuarioID))
	if usuarioOrigen == nil {
		return intercambios
	}

	// Obtener todos los usuarios excepto el origen
	for _, usuarioDestino := range DB.Usuarios {
		if usuarioDestino.ID == usuarioID {
			continue
		}

		// Buscar intercambios viables basados en especializaciones
		intercambiosEncontrados := buscarIntercambiosPorEspecializacion(usuarioOrigen, &usuarioDestino)
		intercambios = append(intercambios, intercambiosEncontrados...)
	}

	return intercambios
}

// buscarIntercambiosPorEspecializacion encuentra intercambios basados en especializaciones
func buscarIntercambiosPorEspecializacion(origen, destino *User) []Intercambio {
	var intercambios []Intercambio

	// Si el origen tiene especialización, buscar intercambios con su producto especializado
	if origen.Especializacion != nil {
		productoOrigen := *origen.Especializacion
		cantidadOrigen := ObtenerCantidadProducto(origen.ID, productoOrigen)

		if cantidadOrigen > 0 {
			// Buscar qué puede ofrecer el destino
			intercambios = append(intercambios, generarIntercambios(origen, destino, productoOrigen, cantidadOrigen)...)
		}
	}

	// Si el destino tiene especialización, buscar intercambios con su producto especializado
	if destino.Especializacion != nil {
		productoDestino := *destino.Especializacion
		cantidadDestino := ObtenerCantidadProducto(destino.ID, productoDestino)

		if cantidadDestino > 0 {
			// Buscar qué puede ofrecer el origen
			intercambios = append(intercambios, generarIntercambios(destino, origen, productoDestino, cantidadDestino)...)
		}
	}

	// Buscar intercambios básicos (trigo por herramientas)
	intercambios = append(intercambios, buscarIntercambiosBasicos(origen, destino)...)

	return intercambios
}

// generarIntercambios genera intercambios viables entre dos usuarios
func generarIntercambios(origen, destino *User, productoEspecializado string, cantidadEspecializada int) []Intercambio {
	var intercambios []Intercambio

	// Calcular valor del producto especializado
	valorEspecializado := CalcularValorProducto(productoEspecializado, cantidadEspecializada)

	// Buscar productos que el otro usuario puede ofrecer
	productosDisponibles := []string{"trigo", "herramientas", "dinero"}

	for _, producto := range productosDisponibles {
		if producto == productoEspecializado {
			continue
		}

		cantidadDisponible := ObtenerCantidadProducto(destino.ID, producto)
		if cantidadDisponible > 0 {
			// Calcular cantidad equivalente basada en valor
			valorDisponible := CalcularValorProducto(producto, cantidadDisponible)

			if valorDisponible >= valorEspecializado {
				// Calcular cantidad exacta para intercambio equitativo
				cantidadEquivalente := calcularCantidadEquivalente(producto, valorEspecializado)

				if cantidadEquivalente > 0 && cantidadEquivalente <= cantidadDisponible {
					intercambio := Intercambio{
						ID:              generarIDIntercambio(),
						UsuarioOrigen:   origen.ID,
						UsuarioDestino:  destino.ID,
						ProductoOrigen:  productoEspecializado,
						ProductoDestino: producto,
						CantidadOrigen:  cantidadEspecializada,
						CantidadDestino: cantidadEquivalente,
						ValorOrigen:     valorEspecializado,
						ValorDestino:    CalcularValorProducto(producto, cantidadEquivalente),
						Estado:          "pendiente",
						FechaCreacion:   obtenerFechaActual(),
					}
					intercambios = append(intercambios, intercambio)
				}
			}
		}
	}

	return intercambios
}

// buscarIntercambiosBasicos busca intercambios básicos (trigo por herramientas)
func buscarIntercambiosBasicos(origen, destino *User) []Intercambio {
	var intercambios []Intercambio

	// Intercambio trigo por herramientas
	trigoOrigen := origen.Inventario.Trigo
	herramientasDestino := destino.Inventario.Herramientas

	if trigoOrigen > 0 && herramientasDestino > 0 {
		// Calcular intercambio equitativo: 2 trigo = 1 herramienta
		cantidadTrigo := 10 // Cantidad base para intercambio
		cantidadHerramientas := 5

		valorTrigo := CalcularValorProducto("trigo", cantidadTrigo)
		valorHerramientas := CalcularValorProducto("herramientas", cantidadHerramientas)

		if valorTrigo == valorHerramientas &&
			trigoOrigen >= cantidadTrigo &&
			herramientasDestino >= cantidadHerramientas {

			intercambio := Intercambio{
				ID:              generarIDIntercambio(),
				UsuarioOrigen:   origen.ID,
				UsuarioDestino:  destino.ID,
				ProductoOrigen:  "trigo",
				ProductoDestino: "herramientas",
				CantidadOrigen:  cantidadTrigo,
				CantidadDestino: cantidadHerramientas,
				ValorOrigen:     valorTrigo,
				ValorDestino:    valorHerramientas,
				Estado:          "pendiente",
				FechaCreacion:   obtenerFechaActual(),
			}
			intercambios = append(intercambios, intercambio)
		}
	}

	return intercambios
}

// calcularCantidadEquivalente calcula la cantidad equivalente de un producto basado en valor
func calcularCantidadEquivalente(producto string, valorObjetivo float64) int {
	switch producto {
	case "trigo":
		return int(valorObjetivo / HORAS_POR_TRIGO)
	case "herramientas":
		return int(valorObjetivo / HORAS_POR_HERRAMIENTA)
	case "dinero":
		return int(valorObjetivo / HORAS_POR_DINERO)
	default:
		return 0
	}
}

// generarIDIntercambio genera un ID único para intercambios
func generarIDIntercambio() string {
	return "intercambio_" + strconv.FormatInt(int64(len(DB.Usuarios)), 10) + "_" + strconv.FormatInt(int64(time.Now().Unix()), 10)
}

// obtenerFechaActual obtiene la fecha actual en formato string
func obtenerFechaActual() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// ObtenerOfertasTrueque obtiene todas las ofertas de trueque activas
func ObtenerOfertasTrueque() []OfertaTrueque {
	var ofertas []OfertaTrueque

	for _, usuario := range DB.Usuarios {
		if usuario.Especializacion != nil {
			productoEspecializado := *usuario.Especializacion
			cantidadDisponible := ObtenerCantidadProducto(usuario.ID, productoEspecializado)

			if cantidadDisponible > 0 {
				ofertas = append(ofertas, generarOfertasPorEspecializacion(usuario, productoEspecializado, cantidadDisponible)...)
			}
		}
	}

	return ofertas
}

// generarOfertasPorEspecializacion genera ofertas de trueque basadas en especialización
func generarOfertasPorEspecializacion(usuario User, productoEspecializado string, cantidadDisponible int) []OfertaTrueque {
	var ofertas []OfertaTrueque

	// Generar ofertas para diferentes productos
	productosObjetivo := []string{"trigo", "herramientas", "dinero"}

	for _, productoObjetivo := range productosObjetivo {
		if productoObjetivo == productoEspecializado {
			continue
		}

		// Calcular cantidades para intercambio equitativo
		cantidadOfrece := cantidadDisponible / 2 // Ofrecer la mitad de lo disponible
		if cantidadOfrece > 0 {
			valorOfrece := CalcularValorProducto(productoEspecializado, cantidadOfrece)
			cantidadBusca := calcularCantidadEquivalente(productoObjetivo, valorOfrece)

			if cantidadBusca > 0 {
				oferta := OfertaTrueque{
					ID:             generarIDOferta(),
					UsuarioID:      usuario.ID,
					ProductoOfrece: productoEspecializado,
					CantidadOfrece: cantidadOfrece,
					ProductoBusca:  productoObjetivo,
					CantidadBusca:  cantidadBusca,
					ValorOfrece:    valorOfrece,
					ValorBusca:     CalcularValorProducto(productoObjetivo, cantidadBusca),
					Activa:         true,
				}
				ofertas = append(ofertas, oferta)
			}
		}
	}

	return ofertas
}

// generarIDOferta genera un ID único para ofertas
func generarIDOferta() string {
	return "oferta_" + strconv.FormatInt(int64(len(DB.Usuarios)), 10) + "_" + strconv.FormatInt(int64(time.Now().Unix()), 10)
}
