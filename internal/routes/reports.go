package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupReportRoutes configure les routes des rapports
func SetupReportRoutes(router *gin.RouterGroup, reportHandler *handlers.ReportHandler) {
	reports := router.Group("/reports")
	reports.Use(middleware.AuthMiddleware())
	{
		reports.GET("/dashboard", reportHandler.GetDashboard)
		reports.GET("/tickets/count", reportHandler.GetTicketCountReport)
		reports.GET("/tickets/distribution", reportHandler.GetTicketTypeDistribution)
		reports.GET("/tickets/average-resolution-time", reportHandler.GetAverageResolutionTime)
		// reports.GET("/workload/by-agent", reportHandler.GetWorkloadByAgent) // TODO: Implémenter dans le handler
		reports.POST("/custom", reportHandler.GenerateCustomReport)
	}
}
