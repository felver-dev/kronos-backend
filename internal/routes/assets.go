package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupAssetRoutes configure les routes des actifs IT
func SetupAssetRoutes(router *gin.RouterGroup, assetHandler *handlers.AssetHandler, assetCategoryHandler *handlers.AssetCategoryHandler) {
	assets := router.Group("/assets")
	assets.Use(middleware.AuthMiddleware())
	{
		assets.GET("", assetHandler.GetAll)
		assets.GET("/:id", assetHandler.GetByID)
		assets.POST("", assetHandler.Create)
		assets.PUT("/:id", assetHandler.Update)
		assets.DELETE("/:id", assetHandler.Delete)
		assets.POST("/:id/assign", assetHandler.Assign)
		assets.DELETE("/:id/unassign-user", assetHandler.Unassign)
		assets.GET("/:id/assigned-user", assetHandler.GetAssignedUser)
		assets.GET("/:id/tickets", assetHandler.GetLinkedTickets)
		assets.POST("/:id/link-ticket/:ticketId", assetHandler.LinkTicket)
		assets.DELETE("/:id/unlink-ticket/:ticketId", assetHandler.UnlinkTicket)
		assets.GET("/by-category/:categoryId", assetHandler.GetByCategory)
		assets.GET("/by-user/:userId", assetHandler.GetByUser)
		assets.GET("/inventory", assetHandler.GetInventory)

		// Catégories d'actifs
		assets.GET("/categories", assetCategoryHandler.GetAll)
		assets.GET("/categories/:id", assetCategoryHandler.GetByID)
		assets.POST("/categories", assetCategoryHandler.Create)
		assets.PUT("/categories/:id", assetCategoryHandler.Update)
		assets.DELETE("/categories/:id", assetCategoryHandler.Delete)
	}
}
