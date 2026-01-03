package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupAuthRoutes configure les routes d'authentification
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		// Routes publiques (sans authentification)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Routes protégées (avec authentification)
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", authHandler.GetMe)
		}
	}
}
