package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupTicketInternesRoutes configure les routes des tickets internes
// La route GET /panier est enregistrée dans router.go (api.GET) pour éviter que /:id ne capture "panier".
func SetupTicketInternesRoutes(router *gin.RouterGroup, handler *handlers.TicketInternalHandler) {
	ti := router.Group("/ticket-internes")
	ti.Use(middleware.AuthMiddleware())
	ti.GET("", handler.GetAll)
	ti.POST("", handler.Create)
	ti.GET("/performance/mine", handler.GetMyPerformance)
	ti.GET("/:id", handler.GetByID)
	ti.PUT("/:id", handler.Update)
	ti.POST("/:id/assign", handler.Assign)
	ti.PUT("/:id/status", handler.ChangeStatus)
	ti.POST("/:id/validate", handler.Validate)
	ti.POST("/:id/close", handler.Close)
	ti.DELETE("/:id", handler.Delete)
}
