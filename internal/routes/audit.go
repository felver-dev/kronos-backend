package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupAuditRoutes configure les routes des logs d'audit
func SetupAuditRoutes(router *gin.RouterGroup, auditHandler *handlers.AuditHandler) {
	audit := router.Group("/audit-logs")
	audit.Use(middleware.AuthMiddleware())
	{
		audit.GET("", auditHandler.GetAll)
		audit.GET("/:id", auditHandler.GetByID)
		audit.GET("/by-user/:userId", auditHandler.GetByUserID)
		audit.GET("/by-action/:action", auditHandler.GetByAction)
		audit.GET("/by-entity/:entityType/:entityId", auditHandler.GetByEntity)
	}
}

