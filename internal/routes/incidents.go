package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupIncidentRoutes configure les routes des incidents
func SetupIncidentRoutes(router *gin.RouterGroup, incidentHandler *handlers.IncidentHandler) {
	incidents := router.Group("/incidents")
	incidents.Use(middleware.AuthMiddleware())
	{
		incidents.GET("", incidentHandler.GetAll)
		incidents.GET("/:id", incidentHandler.GetByID)
		incidents.POST("", incidentHandler.Create)
		incidents.PUT("/:id", incidentHandler.Update)
		incidents.DELETE("/:id", incidentHandler.Delete)
		incidents.POST("/:id/qualify", incidentHandler.Qualify)
		incidents.POST("/:id/resolve", incidentHandler.Resolve)
		incidents.GET("/:id/resolution-time", incidentHandler.GetResolutionTime)
		incidents.POST("/:id/link-asset", incidentHandler.LinkAsset)
		incidents.DELETE("/:id/unlink-asset/:assetId", incidentHandler.UnlinkAsset)
		incidents.GET("/:id/linked-assets", incidentHandler.GetLinkedAssets)
		incidents.GET("/by-impact/:impact", incidentHandler.GetByImpact)
		incidents.GET("/by-urgency/:urgency", incidentHandler.GetByUrgency)
	}
}

