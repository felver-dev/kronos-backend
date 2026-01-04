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

			// Rôles
			SetupRoleRoutes(api, handlers.RoleHandler)

			// Tickets - Les routes spécifiques doivent être définies avant les routes génériques
			// Donc on définit d'abord les routes de timesheet et delay-justification
			SetupTicketTimesheetRoutes(api, handlers.TimesheetHandler)
			SetupTicketDelayJustificationRoutes(api, handlers.DelayHandler)
			SetupTicketAuditRoutes(api, handlers.AuditHandler)
			// Puis les routes principales des tickets
			SetupTicketRoutes(api, handlers.TicketHandler, handlers.TicketAttachmentHandler)

			// Incidents
			SetupIncidentRoutes(api, handlers.IncidentHandler)

			// Changements
			SetupChangeRoutes(api, handlers.ChangeHandler)

			// Demandes de service
			SetupServiceRequestRoutes(api, handlers.ServiceRequestHandler, handlers.ServiceRequestTypeHandler)

			// Entrées de temps
			SetupTimeEntryRoutes(api, handlers.TimeEntryHandler)

			// Retards
			SetupDelayRoutes(api, handlers.DelayHandler)
			SetupUserDelayJustificationRoutes(api, handlers.DelayHandler)

			// Actifs IT
			SetupAssetRoutes(api, handlers.AssetHandler, handlers.AssetCategoryHandler)

			// SLA
			SetupSLARoutes(api, handlers.SLAHandler)

			// Notifications
			SetupNotificationRoutes(api, handlers.NotificationHandler)

			// Base de connaissances
			SetupKnowledgeBaseRoutes(api, handlers.KnowledgeArticleHandler, handlers.KnowledgeCategoryHandler)

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

			// Recherche globale
			SetupSearchRoutes(api, handlers.SearchHandler)

			// Statistiques
			SetupStatisticsRoutes(api, handlers.StatisticsHandler)

			// Logs d'audit
			SetupAuditRoutes(api, handlers.AuditHandler)

			// Paramétrage
			SetupSettingsRoutes(api, handlers.SettingsHandler, handlers.RequestSourceHandler, handlers.BackupHandler)

			// Timesheet
			SetupTimesheetRoutes(api, handlers.TimesheetHandler)
			SetupUserTimesheetRoutes(api, handlers.TimesheetHandler)
			SetupProjectTimesheetRoutes(api, handlers.TimesheetHandler)
		}
	}
}

// Handlers contient toutes les instances de handlers
type Handlers struct {
	AuthHandler                *handlers.AuthHandler
	UserHandler                *handlers.UserHandler
	RoleHandler                *handlers.RoleHandler
	TicketHandler              *handlers.TicketHandler
	TicketAttachmentHandler     *handlers.TicketAttachmentHandler
	IncidentHandler            *handlers.IncidentHandler
	ChangeHandler              *handlers.ChangeHandler
	ServiceRequestHandler      *handlers.ServiceRequestHandler
	ServiceRequestTypeHandler   *handlers.ServiceRequestTypeHandler
	TimeEntryHandler           *handlers.TimeEntryHandler
	DelayHandler               *handlers.DelayHandler
	AssetHandler               *handlers.AssetHandler
	AssetCategoryHandler        *handlers.AssetCategoryHandler
	SLAHandler                 *handlers.SLAHandler
	NotificationHandler        *handlers.NotificationHandler
	KnowledgeArticleHandler    *handlers.KnowledgeArticleHandler
	KnowledgeCategoryHandler    *handlers.KnowledgeCategoryHandler
	ProjectHandler             *handlers.ProjectHandler
	DailyDeclarationHandler    *handlers.DailyDeclarationHandler
	WeeklyDeclarationHandler   *handlers.WeeklyDeclarationHandler
		PerformanceHandler         *handlers.PerformanceHandler
		ReportHandler              *handlers.ReportHandler
		SearchHandler              *handlers.SearchHandler
		StatisticsHandler          *handlers.StatisticsHandler
		AuditHandler               *handlers.AuditHandler
		SettingsHandler            *handlers.SettingsHandler
		RequestSourceHandler       *handlers.RequestSourceHandler
		BackupHandler              *handlers.BackupHandler
		TimesheetHandler           *handlers.TimesheetHandler
	}

