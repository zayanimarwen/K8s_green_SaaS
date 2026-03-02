package analyzer

import "fmt"

type DeployPatch struct {
	Deployment, Namespace string
	SavingsEurAnnual      float64
	ContainerPatches      []ContainerPatch
	CreateHPA             bool
	SuggestedHPA          HPAConfig
}
type ContainerPatch struct {
	Name, CPURequest, CPULimit, MemRequest, MemLimit string
}
type HPAConfig struct { MinReplicas, MaxReplicas int; TargetCPUPercent int }

// GeneratePatches crée les patches K8s pour chaque deployment surdimensionné
func GeneratePatches(reports []WasteReport) []DeployPatch {
	depMap := map[string]*DeployPatch{}
	for _, r := range reports {
		if !r.IsOverprovisioned { continue }
		key := r.Namespace + "/" + guessDeployment(r.PodName)
		if _, ok := depMap[key]; !ok {
			dep := guessDeployment(r.PodName)
			depMap[key] = &DeployPatch{
				Deployment: dep, Namespace: r.Namespace,
				CreateHPA: !r.IsHPAManaged,
				SuggestedHPA: HPAConfig{MinReplicas: 1, MaxReplicas: 3, TargetCPUPercent: 70},
			}
		}
		p := depMap[key]
		p.SavingsEurAnnual += r.AnnualCostWasteEur
		p.ContainerPatches = append(p.ContainerPatches, ContainerPatch{
			Name:       r.ContainerName,
			CPURequest: r.CPURightsizing,
			MemRequest: r.MemRightsizing,
			CPULimit:   fmt.Sprintf("%dm", int(r.CPUOptimalM*2)),
			MemLimit:   fmt.Sprintf("%dMi", int(r.MemOptimalMi*1.5)),
		})
	}
	var result []DeployPatch
	for _, p := range depMap { result = append(result, *p) }
	return result
}
func guessDeployment(podName string) string {
	// Pod name format: <deployment>-<rs-hash>-<pod-hash>
	parts := []byte(podName)
	count := 0
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == '-' { count++; if count == 2 { return string(parts[:i]) } }
	}
	return podName
}
