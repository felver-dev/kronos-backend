package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
)

// SetupRoutes configure toutes les routes de l'application
func SetupRoutes(router *gin.Engine, handlers *Handlers) {
	// Middleware global
	router.Use(middleware.CORSMiddleware())

	// Route de santé
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Route Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Groupe API v1
	api := router.Group("/api/v1")
	{
		// Routes d'authentification (publiques)
		SetupAuthRoutes(api, handlers.AuthHandler)

		// Routes protégées (nécessitent authentification)
		api.Use(middleware.AuthMiddleware())
		{
			// Utilisateurs
			SetupUserRoutes(api, handlers.UserHandler)

			// Tickets
			SetupTicketRoutes(api, handlers.TicketHandler)

			// Incidents
			SetupIncidentRoutes(api, handlers.IncidentHandler)

			// Changements
			SetupChangeRoutes(api, handlers.ChangeHandler)

			// Demandes de service
			SetupServiceRequestRoutes(api, handlers.ServiceRequestHandler)

			// Entrées de temps
			SetupTimeEntryRoutes(api, handlers.TimeEntryHandler)

			// Retards
			SetupDelayRoutes(api, handlers.DelayHandler)

			// Actifs IT
			SetupAssetRoutes(api, handlers.AssetHandler)

			// SLA
			SetupSLARoutes(api, handlers.SLAHandler)

			// Notifications
			SetupNotificationRoutes(api, handlers.NotificationHandler)

			// Base de connaissances
			SetupKnowledgeBaseRoutes(api, handlers.KnowledgeArticleHandler)

			// Projets
			SetupProjectRoutes(api, handlers.ProjectHandler)

			// Déclarations journalières
			SetupDailyDeclarationRoutes(api, handlers.DailyDeclarationHandler)

			// Déclarations hebdomadaires
			SetupWeeklyDeclarationRoutes(api, handlers.WeeklyDeclarationHandler)

			// Performances
			SetupPerformanceRoutes(api, handlers.PerformanceHandler)

			// Rapports
			SetupReportRoutes(api, handlers.ReportHandler)
		}
	}
}

// Handlers contient toutes les instances de handlers
type Handlers struct {
	AuthHandler                *handlers.AuthHandler
	UserHandler                *handlers.UserHandler
	TicketHandler              *handlers.TicketHandler
	IncidentHandler            *handlers.IncidentHandler
	ChangeHandler              *handlers.ChangeHandler
	ServiceRequestHandler      *handlers.ServiceRequestHandler
	TimeEntryHandler           *handlers.TimeEntryHandler
	DelayHandler               *handlers.DelayHandler
	AssetHandler               *handlers.AssetHandler
	SLAHandler                 *handlers.SLAHandler
	NotificationHandler        *handlers.NotificationHandler
	KnowledgeArticleHandler    *handlers.KnowledgeArticleHandler
	ProjectHandler             *handlers.ProjectHandler
	DailyDeclarationHandler    *handlers.DailyDeclarationHandler
	WeeklyDeclarationHandler   *handlers.WeeklyDeclarationHandler
	PerformanceHandler         *handlers.PerformanceHandler
	ReportHandler              *handlers.ReportHandler
}

