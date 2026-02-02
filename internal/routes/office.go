package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupOfficeRoutes configure les routes des sièges
func SetupOfficeRoutes(router *gin.RouterGroup, officeHandler *handlers.OfficeHandler) {
	offices := router.Group("/offices")
	offices.Use(middleware.AuthMiddleware())
	{
		offices.POST("", officeHandler.Create)
		offices.GET("", officeHandler.GetAll)
		offices.GET("/:id", officeHandler.GetByID)
		offices.PUT("/:id", officeHandler.Update)
		offices.DELETE("/:id", officeHandler.Delete)
		offices.GET("/country/:country", officeHandler.GetByCountry)
		offices.GET("/city/:city", officeHandler.GetByCity)
	}
}
