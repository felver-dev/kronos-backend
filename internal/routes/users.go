package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupUserRoutes configure les routes des utilisateurs
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", userHandler.GetAll)
		users.GET("/for-ticket-creation", userHandler.GetForTicketCreation) // Route spécifique avant /:id
		users.GET("/:id", userHandler.GetByID)
		users.POST("", userHandler.Create)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.PUT("/:id/password", userHandler.ChangePassword)
		users.PUT("/:id/reset-password", userHandler.ResetPassword)
		users.GET("/:id/permissions", userHandler.GetPermissions)
		users.PUT("/:id/permissions", userHandler.UpdatePermissions)
		users.POST("/:id/avatar", userHandler.UploadAvatar)
		users.GET("/:id/avatar", userHandler.GetAvatar)
		users.GET("/:id/avatar/thumbnail", userHandler.GetAvatarThumbnail)
		users.DELETE("/:id/avatar", userHandler.DeleteAvatar)
	}
}

// SetupUserDelayJustificationRoutes configure les routes de justification de retard pour les utilisateurs
func SetupUserDelayJustificationRoutes(router *gin.RouterGroup, delayHandler *handlers.DelayHandler) {
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		// Route spécifique avec plus de segments - doit être avant les routes génériques
		users.GET("/:id/delays", delayHandler.GetByUserID)
		users.GET("/:id/delay-justifications", delayHandler.GetJustificationsByUserID)
	}
}
