package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupSLARoutes configure les routes des SLA
func SetupSLARoutes(router *gin.RouterGroup, slaHandler *handlers.SLAHandler) {
	sla := router.Group("/sla")
	sla.Use(middleware.AuthMiddleware())
	{
		sla.GET("", slaHandler.GetAll)
		sla.GET("/:id", slaHandler.GetByID)
		sla.POST("", slaHandler.Create)
		sla.PUT("/:id", slaHandler.Update)
		sla.DELETE("/:id", slaHandler.Delete)
		sla.GET("/:id/compliance", slaHandler.GetCompliance)
		sla.GET("/:id/violations", slaHandler.GetViolations)
		sla.GET("/violations", slaHandler.GetAllViolations)
		sla.GET("/compliance-report", slaHandler.GetComplianceReport)
		sla.GET("/tickets/:ticket_id/status", slaHandler.GetTicketSLAStatus)
		sla.POST("/recalculate", slaHandler.RecalculateSLAStatuses)
	}
}
