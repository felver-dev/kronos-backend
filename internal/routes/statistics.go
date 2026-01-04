package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupStatisticsRoutes configure les routes des statistiques
func SetupStatisticsRoutes(router *gin.RouterGroup, statisticsHandler *handlers.StatisticsHandler) {
	stats := router.Group("/stats")
	stats.Use(middleware.AuthMiddleware())
	{
		stats.GET("/overview", statisticsHandler.GetOverview)
		stats.GET("/workload", statisticsHandler.GetWorkload)
		stats.GET("/performance", statisticsHandler.GetPerformance)
		stats.GET("/trends", statisticsHandler.GetTrends)
		stats.GET("/kpi", statisticsHandler.GetKPI)
	}
}

