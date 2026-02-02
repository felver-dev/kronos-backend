package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupAssetRoutes configure les routes des actifs IT
func SetupAssetRoutes(router *gin.RouterGroup, assetHandler *handlers.AssetHandler, assetCategoryHandler *handlers.AssetCategoryHandler, assetSoftwareHandler *handlers.AssetSoftwareHandler) {
	assets := router.Group("/assets")
	assets.Use(middleware.AuthMiddleware())
	{
		// Routes statiques et spécifiques en premier
		assets.GET("", assetHandler.GetAll)
		assets.POST("", assetHandler.Create)
		assets.GET("/inventory", assetHandler.GetInventory)
		assets.GET("/by-category/:categoryId", assetHandler.GetByCategory)
		assets.GET("/by-user/:userId", assetHandler.GetByUser)

		// Catégories d'actifs (doivent être avant les routes avec :id)
		assets.GET("/categories", assetCategoryHandler.GetAll)
		assets.POST("/categories", assetCategoryHandler.Create)
		assets.GET("/categories/:id", assetCategoryHandler.GetByID)
		assets.PUT("/categories/:id", assetCategoryHandler.Update)
		assets.DELETE("/categories/:id", assetCategoryHandler.Delete)

		// Routes pour les logiciels installés
		assets.GET("/software/statistics", assetSoftwareHandler.GetStatistics)
		assets.GET("/software", assetSoftwareHandler.GetAll)
		assets.GET("/software/by-name/:softwareName", assetSoftwareHandler.GetBySoftwareName)
		assets.GET("/software/by-name/:softwareName/version/:version", assetSoftwareHandler.GetBySoftwareNameAndVersion)
		assets.POST("/software", assetSoftwareHandler.Create)
		assets.GET("/software/:id", assetSoftwareHandler.GetByID)
		assets.PUT("/software/:id", assetSoftwareHandler.Update)
		assets.DELETE("/software/:id", assetSoftwareHandler.Delete)
		assets.GET("/:id/software", assetSoftwareHandler.GetByAssetID)

		// Routes génériques avec :id en dernier
		assets.GET("/:id", assetHandler.GetByID)
		assets.PUT("/:id", assetHandler.Update)
		assets.DELETE("/:id", assetHandler.Delete)
		assets.POST("/:id/assign", assetHandler.Assign)
		assets.DELETE("/:id/unassign-user", assetHandler.Unassign)
		assets.GET("/:id/assigned-user", assetHandler.GetAssignedUser)
		assets.GET("/:id/tickets", assetHandler.GetLinkedTickets)
		assets.POST("/:id/link-ticket/:ticketId", assetHandler.LinkTicket)
		assets.DELETE("/:id/unlink-ticket/:ticketId", assetHandler.UnlinkTicket)
	}
}
