package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupProjectRoutes configure les routes des projets
func SetupProjectRoutes(router *gin.RouterGroup, projectHandler *handlers.ProjectHandler) {
	projects := router.Group("/projects")
	projects.Use(middleware.AuthMiddleware())
	{
		projects.GET("", projectHandler.GetAll)
		projects.GET("/:id", projectHandler.GetByID)
		projects.POST("", projectHandler.Create)
	}
}

