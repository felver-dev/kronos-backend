package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupTicketRoutes configure les routes des tickets
func SetupTicketRoutes(router *gin.RouterGroup, ticketHandler *handlers.TicketHandler, ticketAttachmentHandler *handlers.TicketAttachmentHandler, ticketCategoryHandler *handlers.TicketCategoryHandler, ticketSolutionHandler *handlers.TicketSolutionHandler) {
	tickets := router.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.GET("", ticketHandler.GetAll)
		tickets.POST("", ticketHandler.Create)

		// Routes statiques (sans paramètres) en premier
		tickets.GET("/my-tickets", ticketHandler.GetMyTickets)
		tickets.GET("/panier", ticketHandler.GetMyPanier)
		tickets.GET("/by-source/:source", ticketHandler.GetBySource)
		tickets.GET("/by-category/:category", ticketHandler.GetByCategory)
		tickets.GET("/by-status/:status", ticketHandler.GetByStatus)
		tickets.GET("/by-department/:departmentId", ticketHandler.GetByDepartment)
		tickets.GET("/by-assignee/:userId", ticketHandler.GetByAssignee)

		// Routes pour les catégories (doivent être avant les routes avec :id)
		tickets.GET("/categories", ticketCategoryHandler.GetAll)
		tickets.POST("/categories", ticketCategoryHandler.Create)
		tickets.GET("/categories/slug/:slug", ticketCategoryHandler.GetBySlug)
		tickets.GET("/categories/:id", ticketCategoryHandler.GetByID)
		tickets.PUT("/categories/:id", ticketCategoryHandler.Update)
		tickets.DELETE("/categories/:id", ticketCategoryHandler.Delete)

		// Routes spécifiques avec plus de segments - doivent être avant la route générique :id
		// Routes pour les pièces jointes
		tickets.POST("/:id/attachments", ticketAttachmentHandler.UploadAttachment)
		tickets.GET("/:id/attachments", ticketAttachmentHandler.GetAttachments)
		tickets.GET("/:id/attachments/images", ticketAttachmentHandler.GetImages)
		tickets.GET("/:id/attachments/:attachmentId", ticketAttachmentHandler.GetByID)
		tickets.GET("/:id/attachments/:attachmentId/download", ticketAttachmentHandler.Download)
		tickets.GET("/:id/attachments/:attachmentId/thumbnail", ticketAttachmentHandler.GetThumbnail)
		tickets.PUT("/:id/attachments/:attachmentId", ticketAttachmentHandler.Update)
		tickets.PUT("/:id/attachments/:attachmentId/set-primary", ticketAttachmentHandler.SetPrimary)
		tickets.DELETE("/:id/attachments/:attachmentId", ticketAttachmentHandler.Delete)
		tickets.PUT("/:id/attachments/reorder", ticketAttachmentHandler.Reorder)

		// Routes pour les solutions (doivent être avant les routes génériques)
		tickets.GET("/:id/solutions", ticketSolutionHandler.GetByTicketID)
		tickets.POST("/:id/solutions", ticketSolutionHandler.Create)
		tickets.GET("/solutions/:id", ticketSolutionHandler.GetByID)
		tickets.PUT("/solutions/:id", ticketSolutionHandler.Update)
		tickets.DELETE("/solutions/:id", ticketSolutionHandler.Delete)
		tickets.POST("/solutions/:id/publish-to-kb", ticketSolutionHandler.PublishToKB)

		// Autres routes spécifiques
		tickets.POST("/:id/assign", ticketHandler.Assign)
		tickets.PUT("/:id/status", ticketHandler.ChangeStatus)
		tickets.POST("/:id/validate", ticketHandler.ValidateTicket) // Valider un ticket résolu
		tickets.POST("/:id/close", ticketHandler.Close)
		tickets.POST("/:id/comments", ticketHandler.AddComment)
		tickets.GET("/:id/comments", ticketHandler.GetComments)
		tickets.PUT("/:id/comments/:commentId", ticketHandler.UpdateComment)
		tickets.DELETE("/:id/comments/:commentId", ticketHandler.DeleteComment)
		tickets.POST("/:id/reassign", ticketHandler.Reassign)
		tickets.GET("/:id/history", ticketHandler.GetHistory)

		// Routes génériques (doivent être en dernier)
		tickets.GET("/:id", ticketHandler.GetByID)
		tickets.PUT("/:id", ticketHandler.Update)
		tickets.DELETE("/:id", ticketHandler.Delete)
	}
}

// SetupTicketDelayJustificationRoutes configure les routes de justification de retard pour les tickets
func SetupTicketDelayJustificationRoutes(router *gin.RouterGroup, delayHandler *handlers.DelayHandler) {
	tickets := router.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		// Route spécifique avec plus de segments - doit être avant les routes génériques
		tickets.GET("/:id/delay-justification", delayHandler.GetJustificationByTicketID)
	}
}

// SetupTicketAuditRoutes configure les routes d'audit pour les tickets
func SetupTicketAuditRoutes(router *gin.RouterGroup, auditHandler *handlers.AuditHandler) {
	tickets := router.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.GET("/:id/audit-trail", auditHandler.GetTicketAuditTrail)
	}
}
