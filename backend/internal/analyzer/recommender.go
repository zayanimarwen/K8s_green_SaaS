package analyzer

import "sort"

type Recommendation struct {
	Priority  string  `json:"priority"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	Description string `json:"description"`
	SavingsEurAnnual float64 `json:"savings_eur_annual"`
	Target    string  `json:"target"`
	Confidence float64 `json:"confidence"`
}

// GenerateRecommendations produit des recommandations priorisées depuis les waste reports
func GenerateRecommendations(reports []WasteReport) []Recommendation {
	var recs []Recommendation
	for _, r := range reports {
		if !r.IsOverprovisioned { continue }
		recs = append(recs, Recommendation{
			Priority:  r.Priority,
			Type:      "rightsizing",
			Title:     "Rightsizing " + r.PodName + "/" + r.ContainerName,
			Description: "CPU: " + r.CPURightsizing + " | RAM: " + r.MemRightsizing,
			SavingsEurAnnual: r.AnnualCostWasteEur,
			Target:    r.Namespace + "/" + r.PodName,
			Confidence: r.Confidence,
		})
		if !r.HasLimits {
			recs = append(recs, Recommendation{
				Priority: "HIGH", Type: "missing_limits",
				Title: "Ajouter limits sur " + r.ContainerName,
				Description: "Container sans CPU/memory limits — risque OOM et noisy neighbor",
				Target: r.Namespace + "/" + r.PodName,
			})
		}
	}
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].SavingsEurAnnual > recs[j].SavingsEurAnnual
	})
	if len(recs) > 10 { recs = recs[:10] }
	return recs
}
