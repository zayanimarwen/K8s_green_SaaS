package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
)

func (h *Handler) ListUsers(c *gin.Context) {
	tenantID := middleware.TenantID(c)
	rows, err := h.db.Pool.Query(c.Request.Context(),
		`SELECT user_id, email, display_name, role, created_at, last_login_at
		 FROM user_settings WHERE tenant_id = $1 ORDER BY email`,
		tenantID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur récupération users"})
		return
	}
	defer rows.Close()

	type userRow struct {
		ID          string  `json:"id"`
		Email       string  `json:"email"`
		DisplayName string  `json:"display_name"`
		Role        string  `json:"role"`
	}
	var users []userRow
	for rows.Next() {
		var u userRow
		var createdAt, lastLogin interface{}
		rows.Scan(&u.ID, &u.Email, &u.DisplayName, &u.Role, &createdAt, &lastLogin)
		users = append(users, u)
	}
	if users == nil { users = []userRow{} }
	c.JSON(http.StatusOK, gin.H{"users": users, "tenant_id": tenantID})
}
