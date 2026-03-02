package analyzer

import "math"

type CostReport struct {
	AnnualSavingsEur  float64
	MonthlySavingsEur float64
	Breakdown         map[string]float64
}

// ComputeSavings calcule les économies potentielles depuis les waste reports
func ComputeSavings(reports []WasteReport) CostReport {
	var cpuSavings, memSavings float64
	for _, r := range reports {
		cpuSavings += r.CPUWasteM / 1000 * 0.05 * 8760
		memSavings += r.MemWasteMi / 1024 * 0.006 * 8760
	}
	total := math.Round((cpuSavings+memSavings)*100) / 100
	return CostReport{
		AnnualSavingsEur:  total,
		MonthlySavingsEur: math.Round(total/12*100) / 100,
		Breakdown: map[string]float64{
			"cpu_rightsizing_eur": math.Round(cpuSavings*100) / 100,
			"mem_rightsizing_eur": math.Round(memSavings*100) / 100,
		},
	}
}
