package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupPermissionRoutes configure les routes des permissions
func SetupPermissionRoutes(router *gin.RouterGroup, permissionHandler *handlers.PermissionHandler) {
	permissions := router.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware())
	{
		permissions.GET("", permissionHandler.GetAll)
		permissions.GET("/code/:code", permissionHandler.GetByCode)
	}
}
