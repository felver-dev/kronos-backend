package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupKnowledgeBaseRoutes configure les routes de la base de connaissances
func SetupKnowledgeBaseRoutes(router *gin.RouterGroup, knowledgeArticleHandler *handlers.KnowledgeArticleHandler, knowledgeCategoryHandler *handlers.KnowledgeCategoryHandler) {
	kb := router.Group("/knowledge-base")
	{
		// Routes publiques (articles publiés)
		kb.GET("/articles/published", knowledgeArticleHandler.GetPublished)
		kb.GET("/articles/search", knowledgeArticleHandler.Search)

		// Routes protégées (gestion des articles)
		kb.Use(middleware.AuthMiddleware())
		{
			kb.GET("/articles", knowledgeArticleHandler.GetAll)
			kb.GET("/articles/:id", knowledgeArticleHandler.GetByID)
			kb.POST("/articles", knowledgeArticleHandler.Create)
			kb.PUT("/articles/:id", knowledgeArticleHandler.Update)
			kb.DELETE("/articles/:id", knowledgeArticleHandler.Delete)
			kb.POST("/articles/:id/publish", knowledgeArticleHandler.Publish)
			kb.POST("/articles/:id/view", knowledgeArticleHandler.IncrementViewCount)
			kb.GET("/articles/by-category/:categoryId", knowledgeArticleHandler.GetByCategory)
			kb.GET("/articles/by-author/:authorId", knowledgeArticleHandler.GetByAuthor)

			// Catégories
			kb.GET("/categories", knowledgeCategoryHandler.GetAll)
			kb.GET("/categories/:id", knowledgeCategoryHandler.GetByID)
			kb.POST("/categories", knowledgeCategoryHandler.Create)
			kb.PUT("/categories/:id", knowledgeCategoryHandler.Update)
			kb.DELETE("/categories/:id", knowledgeCategoryHandler.Delete)
		}
	}
}
