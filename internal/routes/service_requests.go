package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupServiceRequestRoutes configure les routes des demandes de service
func SetupServiceRequestRoutes(router *gin.RouterGroup, serviceRequestHandler *handlers.ServiceRequestHandler) {
	serviceRequests := router.Group("/service-requests")
	serviceRequests.Use(middleware.AuthMiddleware())
	{
		serviceRequests.GET("", serviceRequestHandler.GetAll)
		serviceRequests.GET("/:id", serviceRequestHandler.GetByID)
		serviceRequests.POST("", serviceRequestHandler.Create)
		serviceRequests.DELETE("/:id", serviceRequestHandler.Delete)
		serviceRequests.POST("/:id/validate", serviceRequestHandler.Validate)
	}
}
