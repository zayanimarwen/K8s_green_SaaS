package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
	"github.com/k8s-green/backend/internal/reports"
)

func (h *Handler) GetReports(c *gin.Context) {
	tenantID := middleware.TenantID(c)
	clusterID := c.Query("cluster_id")

	// Lister les rapports générés depuis la DB
	_ = tenantID
	_ = clusterID

	// Pour l'instant retourner un rapport de démo
	c.JSON(http.StatusOK, gin.H{
		"reports": []gin.H{
			{
				"id":           "rpt-2024-w04",
				"type":         "weekly",
				"generated_at": time.Now().AddDate(0, 0, -7),
				"download_url": "/v1/reports/rpt-2024-w04/download",
			},
		},
	})
}

func (h *Handler) DownloadReport(c *gin.Context) {
	reportID  := c.Param("id")
	tenantID  := middleware.TenantID(c)
	clusterID := c.Query("cluster_id")

	// Générer le rapport HTML à la volée
	score, _ := h.db.TimescaleRepo().GetLatestScore(c.Request.Context(), clusterID)

	data := reports.WeeklyReportData{
		TenantName:  tenantID,
		ClusterName: clusterID,
		WeekOf:      time.Now().AddDate(0, 0, -7).Format("02/01/2006"),
		GeneratedAt: time.Now().Format("02/01/2006 15:04"),
	}
	if score != nil {
		data.Score  = score.Score
		data.Grade  = score.Grade
		data.CO2Kg  = score.CO2KgAnnual
		data.WasteEur = score.AnnualWasteEur
		data.Pods   = score.PodsAnalyzed
	}

	html, err := reports.GenerateHTML(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "génération rapport échouée"})
		return
	}

	_ = reportID
	c.Header("Content-Disposition", "attachment; filename=green-report.html")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
