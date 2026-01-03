package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupWeeklyDeclarationRoutes configure les routes des déclarations hebdomadaires
func SetupWeeklyDeclarationRoutes(router *gin.RouterGroup, weeklyDeclarationHandler *handlers.WeeklyDeclarationHandler) {
	weeklyDeclarations := router.Group("/weekly-declarations")
	weeklyDeclarations.Use(middleware.AuthMiddleware())
	{
		weeklyDeclarations.GET("/:id", weeklyDeclarationHandler.GetByID)
		weeklyDeclarations.GET("/users/:user_id", weeklyDeclarationHandler.GetByUserID)
		weeklyDeclarations.GET("/users/:user_id/by-week", weeklyDeclarationHandler.GetByUserIDAndWeek)
		weeklyDeclarations.POST("/:id/validate", weeklyDeclarationHandler.Validate)
	}
}
