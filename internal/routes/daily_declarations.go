package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupDailyDeclarationRoutes configure les routes des déclarations journalières
func SetupDailyDeclarationRoutes(router *gin.RouterGroup, dailyDeclarationHandler *handlers.DailyDeclarationHandler) {
	dailyDeclarations := router.Group("/daily-declarations")
	dailyDeclarations.Use(middleware.AuthMiddleware())
	{
		dailyDeclarations.GET("/:id", dailyDeclarationHandler.GetByID)
		dailyDeclarations.GET("/users/:user_id", dailyDeclarationHandler.GetByUserID)
		dailyDeclarations.GET("/users/:user_id/by-date", dailyDeclarationHandler.GetByUserIDAndDate)
		dailyDeclarations.POST("/:id/validate", dailyDeclarationHandler.Validate)
	}
}

