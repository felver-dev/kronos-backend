package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupChangeRoutes configure les routes des changements
func SetupChangeRoutes(router *gin.RouterGroup, changeHandler *handlers.ChangeHandler) {
	changes := router.Group("/changes")
	changes.Use(middleware.AuthMiddleware())
	{
		changes.GET("", changeHandler.GetAll)
		changes.GET("/:id", changeHandler.GetByID)
		changes.POST("", changeHandler.Create)
		changes.PUT("/:id", changeHandler.Update)
		changes.DELETE("/:id", changeHandler.Delete)
		changes.POST("/:id/result", changeHandler.RecordResult)
		changes.GET("/:id/result", changeHandler.GetResult)
		changes.PUT("/:id/risk", changeHandler.UpdateRisk)
		changes.POST("/:id/assign-responsible", changeHandler.AssignResponsible)
		changes.GET("/by-risk/:riskLevel", changeHandler.GetByRisk)
		changes.GET("/by-responsible/:userId", changeHandler.GetByResponsible)
	}
}
