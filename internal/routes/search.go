package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupSearchRoutes configure les routes de recherche
func SetupSearchRoutes(router *gin.RouterGroup, searchHandler *handlers.SearchHandler) {
	search := router.Group("/search")
	search.Use(middleware.AuthMiddleware())
	{
		search.GET("", searchHandler.GlobalSearch)
		search.GET("/tickets", searchHandler.SearchTickets)
		search.GET("/assets", searchHandler.SearchAssets)
		search.GET("/knowledge-base", searchHandler.SearchKnowledgeBase)
		search.GET("/users", searchHandler.SearchUsers)
		search.GET("/time-entries", searchHandler.SearchTimeEntries)
	}
}

