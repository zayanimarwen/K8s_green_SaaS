package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/handlers"
	"github.com/k8s-green/backend/internal/api/middleware"
	"github.com/k8s-green/backend/internal/config"
	"github.com/k8s-green/backend/internal/repository"
)

func NewRouter(cfg *config.Config, db *repository.Postgres, rdb *repository.Redis) *gin.Engine {
	if cfg.Env == "production" { gin.SetMode(gin.ReleaseMode) }

	r := gin.New()
	h := handlers.NewHandler(db, rdb, cfg)

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg))

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/ready", func(c *gin.Context) {
		if err := db.Pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	v1 := r.Group("/v1",
		middleware.Auth(cfg),
		middleware.Tenant(),
		middleware.RateLimit(rdb),
		middleware.Audit(db),
	)

	// Clusters
	clusters := v1.Group("/clusters")
	{
		clusters.GET("",  h.ListClusters)
		clusters.POST("", h.CreateCluster)
		clusters.GET("/:id/score",           h.GetScore)
		clusters.GET("/:id/waste",           h.GetWaste)
		clusters.GET("/:id/carbon",          h.GetCarbon)
		clusters.GET("/:id/savings",         h.GetSavings)
		clusters.GET("/:id/recommendations", h.GetRecommendations)
		clusters.POST("/:id/simulate",       h.Simulate)
		clusters.GET("/:id/history",         h.GetScoreHistory)
		clusters.GET("/:id/diagnostics",     h.GetDiagnostics)
		clusters.POST("/:id/analyze-issue",  h.AnalyzeSingleIssue)
	}

	// Gestion modeles Ollama (admin seulement)
	models := v1.Group("/ai/models", middleware.RequireRole("admin"))
	{
		models.GET("",             h.ListModels)
		models.POST("/pull",       h.PullModel)
		models.DELETE("/:name",    h.DeleteModel)
		models.POST("/switch",     h.SwitchModel)
	}

	v1.GET("/reports", h.GetReports)

	// Endpoint interne pour l agent — auth par header X-Internal-Key
	internal := r.Group("/internal")
	internal.POST("/metrics", middleware.InternalKey(cfg), h.IngestMetrics)

	ws := r.Group("/v1/ws", middleware.Auth(cfg), middleware.Tenant())
	ws.GET("/clusters/:id/live", h.LiveWebSocket)

	admin := v1.Group("/admin", middleware.RequireRole("superadmin"))
	{
		admin.GET("/tenants",     h.ListTenants)
		admin.POST("/tenants",    h.CreateTenant)
		admin.GET("/tenants/:id", h.GetTenant)
		admin.GET("/users",       h.ListUsers)
	}

	return r
}
