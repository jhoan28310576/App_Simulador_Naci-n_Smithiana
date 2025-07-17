package main

import (
	"fmt"
	"log"

	"github.com/jhoan28310576/cap7-8-9_las_riquesas_de_las_naciones/internal/usdaapi"
)

func main() {
	// Parámetros de ejemplo: maíz, año 2023, Iowa, producción
	commodity := "CORN"
	year := "2023"
	state := "IA"
	statisticcat := "PRODUCTION"

	// Llamar a la función que consulta la API de USDA
	data, err := usdaapi.GetUSDAData(commodity, year, state, statisticcat)
	if err != nil {
		log.Fatalf("Error obteniendo datos de USDA: %v", err)
	}

	// Mostrar la cantidad de registros y un ejemplo
	fmt.Printf("Se obtuvieron %d registros de la API de USDA.\n", len(data))
	if len(data) > 0 {
		fmt.Printf("Ejemplo de registro:\n")
		for k, v := range data[0] {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
}
