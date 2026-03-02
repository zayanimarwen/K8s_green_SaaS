package analyzer

import "math"

var carbonIntensityMap = map[string]float64{
	"aws/eu-west-1": 233, "aws/eu-central-1": 338,
	"gcp/europe-west1": 127, "azure/westeurope": 195,
	"azure/francecentral": 56, "on-prem/fr": 65, "on-prem/local": 300,
}
var pueMap = map[string]float64{"aws": 1.18, "gcp": 1.10, "azure": 1.20, "on-prem": 1.58}

const wattCPU = 5.0
const wattRAM = 0.38

type CarbonReport struct {
	Provider, Region     string
	PUE, CarbonIntensity float64
	WattsWasted          float64
	KWhAnnual            float64
	CO2KgAnnual          float64
	EquivalentKmCar      float64
	EquivalentTrees      float64
}

type CarbonCalculator struct{ Provider, Region string }

func NewCarbonCalculator(provider, region string) *CarbonCalculator {
	return &CarbonCalculator{provider, region}
}

func (cc *CarbonCalculator) Compute(reports []WasteReport) CarbonReport {
	var cpuCores, memGB float64
	for _, r := range reports { cpuCores += r.CPUWasteM / 1000; memGB += r.MemWasteMi / 1024 }
	pueFactor := pueMap[cc.Provider]; if pueFactor == 0 { pueFactor = 1.5 }
	ci := carbonIntensityMap[cc.Provider+"/"+cc.Region]; if ci == 0 { ci = 300 }
	watts := cpuCores*wattCPU + memGB*wattRAM
	khw := watts * 8760 / 1000 * pueFactor
	co2 := khw * ci / 1000
	r2 := func(v float64) float64 { return math.Round(v*100) / 100 }
	return CarbonReport{
		Provider: cc.Provider, Region: cc.Region, PUE: pueFactor, CarbonIntensity: ci,
		WattsWasted: r2(watts), KWhAnnual: r2(khw), CO2KgAnnual: r2(co2),
		EquivalentKmCar: r2(co2 * 6.5), EquivalentTrees: r2(co2 / 21.77),
	}
}
