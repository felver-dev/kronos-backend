package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupSettingsRoutes configure les routes des paramètres
func SetupSettingsRoutes(router *gin.RouterGroup, settingsHandler *handlers.SettingsHandler, requestSourceHandler *handlers.RequestSourceHandler, backupHandler *handlers.BackupHandler) {
	settings := router.Group("/settings")
	settings.Use(middleware.AuthMiddleware())
	{
		// Paramètres généraux
		settings.GET("", settingsHandler.GetAll)
		settings.PUT("", settingsHandler.Update)

		// Sources de demande
		sources := settings.Group("/sources")
		{
			sources.GET("", requestSourceHandler.GetAll)
			sources.GET("/:id", requestSourceHandler.GetByID)
			sources.POST("", requestSourceHandler.Create)
			sources.PUT("/:id", requestSourceHandler.Update)
			sources.DELETE("/:id", requestSourceHandler.Delete)
		}

		// Sauvegarde
		backup := settings.Group("/backup")
		{
			backup.GET("", backupHandler.GetConfiguration)
			backup.PUT("", backupHandler.UpdateConfiguration)
			backup.POST("/execute", backupHandler.ExecuteBackup)
		}
	}
}

