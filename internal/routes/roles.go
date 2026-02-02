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
		// Routes spécifiques avant les routes génériques pour éviter les conflits
		roles.GET("/assignable-permissions", roleHandler.GetAssignablePermissions)
		roles.GET("/my-delegations", roleHandler.GetMyDelegations)
		roles.GET("/for-delegation", roleHandler.GetForDelegationPage)
		roles.GET("/:id", roleHandler.GetByID)
		roles.POST("", roleHandler.Create)
		roles.PUT("/:id", roleHandler.Update)
		roles.DELETE("/:id", roleHandler.Delete)

		// Routes pour la gestion des permissions
		roles.GET("/:id/permissions", roleHandler.GetRolePermissions)
		roles.PUT("/:id/permissions", roleHandler.UpdateRolePermissions)
	}
}
