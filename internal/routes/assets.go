package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupAssetRoutes configure les routes des actifs IT
func SetupAssetRoutes(router *gin.RouterGroup, assetHandler *handlers.AssetHandler) {
	assets := router.Group("/assets")
	assets.Use(middleware.AuthMiddleware())
	{
		assets.GET("", assetHandler.GetAll)
		assets.GET("/:id", assetHandler.GetByID)
		assets.POST("", assetHandler.Create)
		assets.DELETE("/:id", assetHandler.Delete)
		assets.POST("/:id/assign", assetHandler.Assign)
	}
}
