package analyzer

import (
	"fmt"
	"math"
)

const (
	CPUWasteThreshold = 0.40
	MemWasteThreshold = 0.50
	SafetyBuffer      = 1.20
)

type PodInput struct {
	PodName, ContainerName, Namespace string
	CPURequestM, CPUUsageP95M         float64
	MemRequestMi, MemUsageP95Mi       float64
	HasLimits, IsHPAManaged           bool
	RestartCount                      int32
}

type WasteReport struct {
	PodName, ContainerName, Namespace string
	IsOverprovisioned, HasLimits, IsHPAManaged bool
	CPURequestM, CPUUsageP95M, CPUOptimalM, CPUWasteM, CPUWastePct float64
	MemRequestMi, MemUsageP95Mi, MemOptimalMi, MemWasteMi, MemWastePct float64
	AnnualCostWasteEur float64
	CPURightsizing, MemRightsizing string
	Priority string; Confidence float64
}

type WasteDetector struct{ CostPerCPUH, CostPerGBH float64 }

func NewWasteDetector(cpu, mem float64) *WasteDetector { return &WasteDetector{cpu, mem} }

func (wd *WasteDetector) Analyze(p PodInput) WasteReport {
	cpuOpt := roundCPU(math.Max(p.CPUUsageP95M*SafetyBuffer, 10))
	memOpt := roundMem(math.Max(p.MemUsageP95Mi*SafetyBuffer, 32))
	cpuW := math.Max(0, p.CPURequestM-cpuOpt)
	memW := math.Max(0, p.MemRequestMi-memOpt)
	cpuWPct, memWPct := 0.0, 0.0
	if p.CPURequestM > 0 { cpuWPct = cpuW / p.CPURequestM * 100 }
	if p.MemRequestMi > 0 { memWPct = memW / p.MemRequestMi * 100 }
	annual := math.Round((cpuW/1000*wd.CostPerCPUH+memW/1024*wd.CostPerGBH)*8760*100) / 100
	isOver := (p.CPURequestM > 0 && p.CPUUsageP95M/p.CPURequestM < CPUWasteThreshold) ||
		(p.MemRequestMi > 0 && p.MemUsageP95Mi/p.MemRequestMi < MemWasteThreshold)
	return WasteReport{
		PodName: p.PodName, ContainerName: p.ContainerName, Namespace: p.Namespace,
		IsOverprovisioned: isOver, HasLimits: p.HasLimits, IsHPAManaged: p.IsHPAManaged,
		CPURequestM: p.CPURequestM, CPUUsageP95M: p.CPUUsageP95M,
		CPUOptimalM: cpuOpt, CPUWasteM: cpuW, CPUWastePct: cpuWPct,
		CPURightsizing: fmt.Sprintf("%dm", int(cpuOpt)),
		MemRequestMi: p.MemRequestMi, MemUsageP95Mi: p.MemUsageP95Mi,
		MemOptimalMi: memOpt, MemWasteMi: memW, MemWastePct: memWPct,
		MemRightsizing: fmt.Sprintf("%dMi", int(memOpt)),
		AnnualCostWasteEur: annual,
		Priority:   priority(cpuWPct, annual, p.HasLimits),
		Confidence: confidence(p),
	}
}
func roundCPU(m float64) float64 {
	if m <= 100 { return math.Ceil(m/10)*10 }
	return math.Ceil(m/50)*50
}
func roundMem(mi float64) float64 {
	for _, p := range []float64{32,64,128,256,512,1024,2048,4096} { if mi <= p { return p } }
	return math.Ceil(mi/512)*512
}
func priority(cpuWastePct, annual float64, hasLimits bool) string {
	if annual > 500 || cpuWastePct > 70 || !hasLimits { return "HIGH" }
	if annual > 100 || cpuWastePct > 50 { return "MEDIUM" }
	return "LOW"
}
func confidence(p PodInput) float64 {
	if p.RestartCount > 5 { return 0.6 }
	if p.IsHPAManaged { return 0.75 }
	return 0.90
}
