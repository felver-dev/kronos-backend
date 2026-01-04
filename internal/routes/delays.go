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
		
		// Routes pour les justifications (sans paramètre de delay) - routes statiques en premier
		delays.GET("/justifications/validated", delayHandler.GetValidatedJustifications)
		delays.GET("/justifications/rejected", delayHandler.GetRejectedJustifications)
		delays.GET("/justifications/history", delayHandler.GetJustificationsHistory)
		delays.POST("/justifications/:id/validate", delayHandler.ValidateJustification)
		
		// Routes spécifiques pour les justifications par delay (doivent être avant la route générique :id)
		// Ces routes ont plus de segments, donc Gin peut les distinguer
		delays.POST("/:id/justifications", delayHandler.CreateJustification)
		delays.GET("/:id/justification", delayHandler.GetJustificationByDelayID)
		delays.PUT("/:id/justification", delayHandler.UpdateJustification)
		delays.DELETE("/:id/justification", delayHandler.DeleteJustification)
		delays.POST("/:id/justification/reject", delayHandler.RejectJustification)
		
		// Route générique pour récupérer un retard par ID (doit être en dernier)
		delays.GET("/:id", delayHandler.GetByID)
	}
}
