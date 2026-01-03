package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupNotificationRoutes configure les routes des notifications
func SetupNotificationRoutes(router *gin.RouterGroup, notificationHandler *handlers.NotificationHandler) {
	notifications := router.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware())
	{
		notifications.GET("", notificationHandler.GetByUserID)
		notifications.GET("/unread/count", notificationHandler.GetUnreadCount)
		notifications.POST("/:id/read", notificationHandler.MarkAsRead)
		notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
	}
}
