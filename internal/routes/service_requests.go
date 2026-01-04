package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupServiceRequestRoutes configure les routes des demandes de service
func SetupServiceRequestRoutes(router *gin.RouterGroup, serviceRequestHandler *handlers.ServiceRequestHandler, serviceRequestTypeHandler *handlers.ServiceRequestTypeHandler) {
	serviceRequests := router.Group("/service-requests")
	serviceRequests.Use(middleware.AuthMiddleware())
	{
		serviceRequests.GET("", serviceRequestHandler.GetAll)
		serviceRequests.GET("/:id", serviceRequestHandler.GetByID)
		serviceRequests.POST("", serviceRequestHandler.Create)
		serviceRequests.PUT("/:id", serviceRequestHandler.Update)
		serviceRequests.DELETE("/:id", serviceRequestHandler.Delete)
		serviceRequests.POST("/:id/validate", serviceRequestHandler.Validate)
		serviceRequests.GET("/:id/deadline", serviceRequestHandler.GetDeadline)
		serviceRequests.GET("/:id/validation-status", serviceRequestHandler.GetValidationStatus)

		// Types de demandes de service
		serviceRequests.GET("/types", serviceRequestTypeHandler.GetAll)
		serviceRequests.GET("/types/:id", serviceRequestTypeHandler.GetByID)
		serviceRequests.POST("/types", serviceRequestTypeHandler.Create)
		serviceRequests.PUT("/types/:id", serviceRequestTypeHandler.Update)
		serviceRequests.DELETE("/types/:id", serviceRequestTypeHandler.Delete)
	}
}
