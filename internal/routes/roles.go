package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupRoleRoutes configure les routes des rôles
func SetupRoleRoutes(router *gin.RouterGroup, roleHandler *handlers.RoleHandler) {
	roles := router.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	{
		roles.GET("", roleHandler.GetAll)
		roles.GET("/:id", roleHandler.GetByID)
		roles.POST("", roleHandler.Create)
		roles.PUT("/:id", roleHandler.Update)
		roles.DELETE("/:id", roleHandler.Delete)
	}
}

