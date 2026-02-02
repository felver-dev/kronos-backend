package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupFilialeRoutes configure les routes des filiales
func SetupFilialeRoutes(router *gin.RouterGroup, filialeHandler *handlers.FilialeHandler, filialeSoftwareHandler *handlers.FilialeSoftwareHandler) {
	filiales := router.Group("/filiales")
	filiales.Use(middleware.AuthMiddleware())
	{
		filiales.POST("", filialeHandler.Create)
		filiales.GET("", filialeHandler.GetAll)
		// Note: /active est défini comme route publique dans router.go
		filiales.GET("/software-provider", filialeHandler.GetSoftwareProvider)
		filiales.GET("/code/:code", filialeHandler.GetByCode)

		// Routes pour les déploiements de logiciels par filiale (DOIVENT être avant /:filiale_id)
		filiales.GET("/:filiale_id/software", filialeSoftwareHandler.GetByFilialeID)
		filiales.POST("/:filiale_id/software", filialeSoftwareHandler.Create)

		// Routes génériques (utilisent :filiale_id pour éviter le conflit avec les routes ci-dessus)
		filiales.GET("/:filiale_id", filialeHandler.GetByID)
		filiales.PUT("/:filiale_id", filialeHandler.Update)
		filiales.DELETE("/:filiale_id", filialeHandler.Delete)
	}
}

// SetupFilialeSoftwareRoutes configure les routes des déploiements de logiciels
func SetupFilialeSoftwareRoutes(router *gin.RouterGroup, filialeSoftwareHandler *handlers.FilialeSoftwareHandler) {
	deployments := router.Group("/filiales-software")
	deployments.Use(middleware.AuthMiddleware())
	{
		deployments.GET("", filialeSoftwareHandler.GetAll)
		deployments.GET("/active", filialeSoftwareHandler.GetActive)
		deployments.GET("/:id", filialeSoftwareHandler.GetByID)
		deployments.PUT("/:id", filialeSoftwareHandler.Update)
		deployments.DELETE("/:id", filialeSoftwareHandler.Delete)
	}
}
