package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupTicketRoutes configure les routes des tickets
func SetupTicketRoutes(router *gin.RouterGroup, ticketHandler *handlers.TicketHandler) {
	tickets := router.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.GET("", ticketHandler.GetAll)
		tickets.GET("/:id", ticketHandler.GetByID)
		tickets.POST("", ticketHandler.Create)
		tickets.PUT("/:id", ticketHandler.Update)
		tickets.DELETE("/:id", ticketHandler.Delete)
		tickets.POST("/:id/assign", ticketHandler.Assign)
		tickets.PUT("/:id/status", ticketHandler.ChangeStatus)
		tickets.POST("/:id/close", ticketHandler.Close)
		tickets.POST("/:id/comments", ticketHandler.AddComment)
		tickets.GET("/:id/comments", ticketHandler.GetComments)
	}
}

