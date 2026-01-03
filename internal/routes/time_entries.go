package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupTimeEntryRoutes configure les routes des entrées de temps
func SetupTimeEntryRoutes(router *gin.RouterGroup, timeEntryHandler *handlers.TimeEntryHandler) {
	timeEntries := router.Group("/time-entries")
	timeEntries.Use(middleware.AuthMiddleware())
	{
		timeEntries.GET("", timeEntryHandler.GetAll)
		timeEntries.GET("/:id", timeEntryHandler.GetByID)
		timeEntries.POST("", timeEntryHandler.Create)
		timeEntries.DELETE("/:id", timeEntryHandler.Delete)
		timeEntries.POST("/:id/validate", timeEntryHandler.Validate)
	}
}
