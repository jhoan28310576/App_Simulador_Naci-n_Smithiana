package models

// USDAQuickStatsResponse representa la respuesta de la API de USDA Quick Stats
type USDAQuickStatsResponse struct {
	Data []USDADataPoint `json:"data"`
}

// USDADataPoint representa un punto de datos individual de la API
type USDADataPoint struct {
	Year                interface{} `json:"year"` // Puede ser string o number
	StateName           string      `json:"state_name"`
	CommodityDesc       string      `json:"commodity_desc"`
	StatisticcatDesc    string      `json:"statisticcat_desc"`
	Value               string      `json:"Value"`
	UnitDesc            string      `json:"unit_desc"`
	ReferencePeriodDesc string      `json:"reference_period_desc"`
	SourceDesc          string      `json:"source_desc"`
	AggLevelDesc        string      `json:"agg_level_desc"`
	UtilPracticeDesc    string      `json:"util_practice_desc"`
	ShortDesc           string      `json:"short_desc"`
}

// CornProductionData representa datos específicos de producción de maíz
type CornProductionData struct {
	Year       string  `json:"year"`
	State      string  `json:"state"`
	Production float64 `json:"production"`
	Unit       string  `json:"unit"`
	Area       float64 `json:"area,omitempty"`
	AreaUnit   string  `json:"area_unit,omitempty"`
}

// DroughtSimulationParams parámetros para la simulación de sequía
type DroughtSimulationParams struct {
	Year            string   `json:"year"`
	States          []string `json:"states"`
	DroughtSeverity float64  `json:"drought_severity"` // 0.0 a 1.0 (0% a 100% de reducción)
	AffectedArea    float64  `json:"affected_area"`    // Porcentaje del área afectada
}

// DroughtSimulationResult resultado de la simulación
type DroughtSimulationResult struct {
	OriginalProduction  float64              `json:"original_production"`
	SimulatedProduction float64              `json:"simulated_production"`
	ProductionLoss      float64              `json:"production_loss"`
	PriceIncrease       float64              `json:"price_increase"`
	EconomicImpact      float64              `json:"economic_impact"`
	States              []CornProductionData `json:"states"`
}
