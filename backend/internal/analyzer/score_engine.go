package analyzer

import "fmt"

// ScoreBreakdown détaille les 5 dimensions du score
type ScoreBreakdown struct {
	CPUEfficiency  float64 `json:"cpu_efficiency"`
	MemEfficiency  float64 `json:"mem_efficiency"`
	NodePacking    float64 `json:"node_packing"`
	HPACoverage    float64 `json:"hpa_coverage"`
	LimitCompliance float64 `json:"limit_compliance"`
}

// GreenScore est le score global d'un cluster
type GreenScore struct {
	Value     float64        `json:"score"`
	Grade     string         `json:"grade"`
	Label     string         `json:"label"`
	Breakdown ScoreBreakdown `json:"breakdown"`
}

var weights = map[string]float64{
	"cpu":   0.35,
	"mem":   0.25,
	"node":  0.20,
	"hpa":   0.10,
	"limit": 0.10,
}

// ComputeScore calcule le Green Score à partir de la liste de pods
func ComputeScore(pods []PodInput, nodeCount, totalPods int) GreenScore {
	if len(pods) == 0 {
		return GreenScore{Value: 0, Grade: "N/A", Label: "Aucune donnée"}
	}

	// CPU efficiency — ratio usage/request, cible 70%
	var cpuRatioSum float64
	for _, p := range pods {
		if p.CPURequestM > 0 {
			cpuRatioSum += p.CPUUsageP95M / p.CPURequestM
		}
	}
	cpuEff := clamp(cpuRatioSum/float64(len(pods))/0.70*100, 0, 100)

	// Mem efficiency — cible 75%
	var memRatioSum float64
	for _, p := range pods {
		if p.MemRequestMi > 0 {
			memRatioSum += p.MemUsageP95Mi / p.MemRequestMi
		}
	}
	memEff := clamp(memRatioSum/float64(len(pods))/0.75*100, 0, 100)

	// Node packing — stub (nécessite info nodes)
	nodePacking := 70.0
	if nodeCount > 0 && totalPods > 0 {
		avgPodsPerNode := float64(totalPods) / float64(nodeCount)
		nodePacking = clamp(avgPodsPerNode/15.0*100, 0, 100)
	}

	// HPA coverage
	hpaCount := 0
	deployments := map[string]bool{}
	for _, p := range pods {
		if p.IsHPAManaged {
			hpaCount++
		}
		deployments[guessDeployment(p.PodName)] = true
	}
	hpaCov := 0.0
	if len(deployments) > 0 {
		hpaCov = clamp(float64(hpaCount)/float64(len(pods))*100, 0, 100)
	}

	// Limit compliance
	limitCount := 0
	for _, p := range pods {
		if p.HasLimits {
			limitCount++
		}
	}
	limitComp := clamp(float64(limitCount)/float64(len(pods))*100, 0, 100)

	total := cpuEff*weights["cpu"] +
		memEff*weights["mem"] +
		nodePacking*weights["node"] +
		hpaCov*weights["hpa"] +
		limitComp*weights["limit"]

	grade, label := gradeFrom(total)

	return GreenScore{
		Value: total,
		Grade: grade,
		Label: label,
		Breakdown: ScoreBreakdown{
			CPUEfficiency:  cpuEff,
			MemEfficiency:  memEff,
			NodePacking:    nodePacking,
			HPACoverage:    hpaCov,
			LimitCompliance: limitComp,
		},
	}
}

func gradeFrom(score float64) (string, string) {
	switch {
	case score >= 90: return "A+", "Excellent — cluster exemplaire"
	case score >= 80: return "A",  "Très bien — quelques optimisations possibles"
	case score >= 70: return "B+", "Bien — améliorations identifiées"
	case score >= 60: return "B",  "Correct — surprovisionnement modéré"
	case score >= 45: return "C",  "À améliorer — gaspillage significatif"
	case score >= 30: return "D",  "Mauvais — action recommandée"
	default:          return "F",  "Critique — intervention urgente"
	}
}

func clamp(v, min, max float64) float64 {
	if v < min { return min }
	if v > max { return max }
	return v
}

func guessDeploymentLocal(podName string) string {
	// Heuristique: pod = "deploy-rs-xyz" → enlever les 2 derniers segments
	parts := []rune(podName)
	count := 0
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == '-' {
			count++
			if count == 2 {
				return fmt.Sprintf("%s", string(parts[:i]))
			}
		}
	}
	return podName
}
