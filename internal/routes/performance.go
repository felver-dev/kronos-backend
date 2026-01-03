package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupPerformanceRoutes configure les routes des performances
func SetupPerformanceRoutes(router *gin.RouterGroup, performanceHandler *handlers.PerformanceHandler) {
	performance := router.Group("/performance")
	performance.Use(middleware.AuthMiddleware())
	{
		performance.GET("/users/:user_id", performanceHandler.GetPerformanceByUserID)
		performance.GET("/users/:user_id/efficiency", performanceHandler.GetEfficiencyByUserID)
		performance.GET("/users/:user_id/productivity", performanceHandler.GetProductivityByUserID)
		performance.GET("/ranking", performanceHandler.GetPerformanceRanking)
	}
}

