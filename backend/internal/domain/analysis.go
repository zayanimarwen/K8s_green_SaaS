package domain

import "time"

type GreenScore struct {
	ID                 string    `json:"id" db:"id"`
	ClusterID          string    `json:"cluster_id" db:"cluster_id"`
	Time               time.Time `json:"time" db:"time"`
	Score              float64   `json:"score" db:"score"`
	Grade              string    `json:"grade" db:"grade"`
	Label              string    `json:"label" db:"label"`
	CPUEfficiency      float64   `json:"cpu_efficiency" db:"cpu_efficiency"`
	MemEfficiency      float64   `json:"mem_efficiency" db:"mem_efficiency"`
	NodePacking        float64   `json:"node_packing" db:"node_packing"`
	HPACoverage        float64   `json:"hpa_coverage" db:"hpa_coverage"`
	LimitCompliance    float64   `json:"limit_compliance" db:"limit_compliance"`
	PodsAnalyzed       int       `json:"pods_analyzed" db:"pods_analyzed"`
	PodsOverprovisioned int      `json:"pods_overprovisioned" db:"pods_overprovisioned"`
	AnnualWasteEur     float64   `json:"annual_waste_eur" db:"annual_waste_eur"`
	CO2KgAnnual        float64   `json:"co2_kg_annual" db:"co2_kg_annual"`
}
