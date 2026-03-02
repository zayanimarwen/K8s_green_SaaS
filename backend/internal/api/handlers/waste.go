package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
)

func (h *Handler) GetWaste(c *gin.Context) {
	clusterID := c.Param("id")
	tenantID  := middleware.TenantID(c)
	cacheKey  := fmt.Sprintf("tenant:%s:waste:%s", tenantID, clusterID)
	ctx       := c.Request.Context()

	if cached, err := h.rdb.Client().Get(ctx, cacheKey).Bytes(); err == nil {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", cached)
		return
	}

	// Requête TimescaleDB : agréger le gaspillage des 24 dernières heures
	rows, err := h.db.Pool.Query(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE cpu_waste_m > cpu_request_m * 0.4)  AS pods_overprovisioned,
			COALESCE(SUM(cpu_waste_m) / 1000, 0)                       AS total_cpu_waste_cores,
			COALESCE(SUM(mem_waste_mi) / 1024, 0)                      AS total_mem_waste_gb,
			COALESCE(SUM(cost_waste_eur), 0)                           AS annual_waste_eur
		FROM pod_metrics
		WHERE cluster_id = $1 AND time > NOW() - INTERVAL '24 hours'`,
		clusterID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête waste"})
		return
	}
	defer rows.Close()

	var podsOver int
	var cpuCores, memGB, wasteEur float64
	if rows.Next() {
		rows.Scan(&podsOver, &cpuCores, &memGB, &wasteEur)
	}

	resp := gin.H{
		"cluster_id":            clusterID,
		"pods_overprovisioned":  podsOver,
		"total_cpu_waste_cores": cpuCores,
		"total_mem_waste_gb":    memGB,
		"annual_waste_eur":      wasteEur,
		"computed_at":           time.Now(),
	}

	if b, err := json.Marshal(resp); err == nil {
		h.rdb.Client().Set(ctx, cacheKey, b, 5*time.Minute)
		c.Header("X-Cache", "MISS")
		c.Data(http.StatusOK, "application/json", b)
		return
	}
	c.JSON(http.StatusOK, resp)
}
