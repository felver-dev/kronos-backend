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
		users.GET("/:id", userHandler.GetByID)
		users.POST("", userHandler.Create)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.PUT("/:id/password", userHandler.ChangePassword)
	}
}

