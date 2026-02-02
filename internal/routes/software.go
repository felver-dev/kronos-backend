package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupSoftwareRoutes configure les routes des logiciels
func SetupSoftwareRoutes(router *gin.RouterGroup, softwareHandler *handlers.SoftwareHandler, filialeSoftwareHandler *handlers.FilialeSoftwareHandler) {
	software := router.Group("/software")
	software.Use(middleware.AuthMiddleware())
	{
		software.POST("", softwareHandler.Create)
		software.GET("", softwareHandler.GetAll)
		// Note: /active est défini comme route publique dans router.go
		software.GET("/code/:code", softwareHandler.GetByCode)

		// Routes pour les déploiements d'un logiciel (DOIVENT être avant /:software_id)
		software.GET("/:software_id/deployments", filialeSoftwareHandler.GetBySoftwareID)

		// Routes génériques (utilisent :software_id pour éviter le conflit avec les routes ci-dessus)
		software.GET("/:software_id", softwareHandler.GetByID)
		software.PUT("/:software_id", softwareHandler.Update)
		software.DELETE("/:software_id", softwareHandler.Delete)
	}
}
