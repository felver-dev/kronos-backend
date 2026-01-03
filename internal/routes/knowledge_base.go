package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupKnowledgeBaseRoutes configure les routes de la base de connaissances
func SetupKnowledgeBaseRoutes(router *gin.RouterGroup, knowledgeArticleHandler *handlers.KnowledgeArticleHandler) {
	kb := router.Group("/knowledge-base")
	{
		// Routes publiques (articles publiés)
		kb.GET("/articles/published", knowledgeArticleHandler.GetPublished)
		kb.GET("/articles/search", knowledgeArticleHandler.Search)

		// Routes protégées (gestion des articles)
		kb.Use(middleware.AuthMiddleware())
		{
			kb.GET("/articles/:id", knowledgeArticleHandler.GetByID)
			kb.POST("/articles", knowledgeArticleHandler.Create)
			kb.DELETE("/articles/:id", knowledgeArticleHandler.Delete)
		}
	}
}
