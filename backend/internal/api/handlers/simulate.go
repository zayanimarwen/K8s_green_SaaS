package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/analyzer"
)

type simulateRequest struct {
	Deployment    string  `json:"deployment"     binding:"required"`
	Namespace     string  `json:"namespace"      binding:"required"`
	NewCPUReqM    float64 `json:"new_cpu_req_m"`
	NewMemReqMi   float64 `json:"new_mem_req_mi"`
}

func (h *Handler) Simulate(c *gin.Context) {
	clusterID := c.Param("id")

	var req simulateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer les métriques actuelles du pod depuis TimescaleDB
	ctx := c.Request.Context()
	var currentCPUReq, currentMemReq, cpuP95, memP95 float64
	var hasLimits, isHPAManaged bool
	var restartCount int32

	row := h.db.Pool.QueryRow(ctx,
		`SELECT
			AVG(cpu_request_m), AVG(mem_request_mi),
			AVG(cpu_usage_p95_m), AVG(mem_usage_p95_mi),
			BOOL_OR(has_limits), false, 0
		FROM pod_metrics
		WHERE cluster_id = $1 AND namespace = $2
		  AND time > NOW() - INTERVAL '24 hours'`,
		clusterID, req.Namespace,
	)
	row.Scan(&currentCPUReq, &currentMemReq, &cpuP95, &memP95, &hasLimits, &isHPAManaged, &restartCount)

	// Valeurs courantes si pas de données
	if cpuP95 == 0 { cpuP95 = 50 }
	if memP95 == 0 { memP95 = 200 }
	if currentCPUReq == 0 { currentCPUReq = 500 }
	if currentMemReq == 0 { currentMemReq = 512 }

	// Calculer le gaspillage actuel
	detector := analyzer.NewWasteDetector(h.cfg.DefaultCostCPUH, h.cfg.DefaultCostMemH)
	current := detector.Analyze(analyzer.PodInput{
		CPURequestM: currentCPUReq, CPUUsageP95M: cpuP95,
		MemRequestMi: currentMemReq, MemUsageP95Mi: memP95,
		HasLimits: hasLimits, IsHPAManaged: isHPAManaged, RestartCount: restartCount,
	})

	// Calculer avec les nouvelles valeurs
	newCPU := req.NewCPUReqM
	if newCPU == 0 { newCPU = current.CPUOptimalM }
	newMem := req.NewMemReqMi
	if newMem == 0 { newMem = current.MemOptimalMi }

	proposed := detector.Analyze(analyzer.PodInput{
		CPURequestM: newCPU, CPUUsageP95M: cpuP95,
		MemRequestMi: newMem, MemUsageP95Mi: memP95,
		HasLimits: hasLimits, IsHPAManaged: isHPAManaged, RestartCount: restartCount,
	})

	// Delta économies
	savingsDelta := current.AnnualCostWasteEur - proposed.AnnualCostWasteEur
	scoreDelta   := (newCPU / currentCPUReq) * 10.0 // heuristique simple
	co2Delta     := savingsDelta * 0.1               // ~0.1 kgCO2/€ gaspillé

	safeToApply := newCPU >= cpuP95*1.1 && newMem >= memP95*1.1

	c.JSON(http.StatusOK, gin.H{
		"cluster_id":                    clusterID,
		"deployment":                    req.Deployment,
		"namespace":                     req.Namespace,
		"current_cpu_req_m":             currentCPUReq,
		"current_mem_req_mi":            currentMemReq,
		"proposed_cpu_req_m":            newCPU,
		"proposed_mem_req_mi":           newMem,
		"projected_savings_eur_annual":  savingsDelta,
		"projected_score_delta":         scoreDelta,
		"projected_co2_reduction_kg":    co2Delta,
		"safe_to_apply":                 safeToApply,
		"confidence":                    proposed.Confidence,
		"computed_at":                   time.Now(),
	})
}
