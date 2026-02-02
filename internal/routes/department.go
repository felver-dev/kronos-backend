package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupDepartmentRoutes configure les routes des départements
func SetupDepartmentRoutes(router *gin.RouterGroup, departmentHandler *handlers.DepartmentHandler) {
	departments := router.Group("/departments")
	departments.Use(middleware.AuthMiddleware())
	{
		departments.POST("", departmentHandler.Create)
		departments.GET("", departmentHandler.GetAll)
		departments.GET("/:id", departmentHandler.GetByID)
		departments.PUT("/:id", departmentHandler.Update)
		departments.DELETE("/:id", departmentHandler.Delete)
		departments.GET("/office/:office_id", departmentHandler.GetByOfficeID)
		departments.GET("/filiale/:filiale_id", departmentHandler.GetByFilialeID)
	}
}
