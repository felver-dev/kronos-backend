package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupDelayRoutes configure les routes des retards
func SetupDelayRoutes(router *gin.RouterGroup, delayHandler *handlers.DelayHandler) {
	delays := router.Group("/delays")
	delays.Use(middleware.AuthMiddleware())
	{
		delays.GET("", delayHandler.GetAll)
		delays.GET("/:id", delayHandler.GetByID)
		delays.POST("/:delay_id/justifications", delayHandler.CreateJustification)
		delays.POST("/justifications/:id/validate", delayHandler.ValidateJustification)
	}
}
