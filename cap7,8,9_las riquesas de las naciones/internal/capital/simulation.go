package capital

import (
	"math/rand"
	"time"
)

type Inversion struct {
	Capital      float64
	Sector       string
	Riesgo       int
	Competencia  float64
	SalarioMedio float64
}

type ResultadoRetorno struct {
	Año     int
	Retorno float64
	Capital float64
}

func SimularRetorno(inv Inversion, años int) []ResultadoRetorno {
	rand.Seed(time.Now().UnixNano())
	resultados := make([]ResultadoRetorno, 0, años)
	capitalAcumulado := inv.Capital
	beneficio := calcularBeneficio(inv)

	for año := 1; año <= años; año++ {
		variacion := -0.02 + rand.Float64()*(0.03+0.02)
		retornoAnual := capitalAcumulado * (beneficio + variacion)
		capitalAcumulado += retornoAnual

		resultados = append(resultados, ResultadoRetorno{
			Año:     año,
			Retorno: retornoAnual,
			Capital: capitalAcumulado,
		})
	}

	return resultados
}

func calcularBeneficio(inv Inversion) float64 {
	base := 0.08

	ajusteRiesgo := []float64{0.02, 0.05, 0.08, 0.12, 0.15}
	if inv.Riesgo < 1 || inv.Riesgo > 5 {
		panic("Riesgo debe estar entre 1 y 5")
	}

	ajustesSector := map[string]float64{
		"agricultura": -0.01,
		"manufactura": 0.02,
		"comercio":    0.03,
	}

	ajusteSector := ajustesSector[inv.Sector] // Default 0 if not found

	// Efecto de la competencia (mayor competencia = menor beneficio)
	factorCompetencia := 1.0 - inv.Competencia*0.3 // Reducción hasta 30%

	// Relación inversa con salarios (Adam Smith)
	factorSalario := 1.0
	if inv.SalarioMedio > 0 {
		factorSalario = 1.0 - (inv.SalarioMedio / 100000 * 0.2) // Ajuste del 20%
	}

	return (base + ajusteRiesgo[inv.Riesgo-1] + ajusteSector) * factorCompetencia * factorSalario
}
