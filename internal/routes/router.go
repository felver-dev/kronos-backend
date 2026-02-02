package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/middleware"
	"github.com/mcicare/itsm-backend/internal/repositories"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configure toutes les routes de l'application
func SetupRoutes(router *gin.Engine, handlers *Handlers, auditLogRepo repositories.AuditLogRepository) {
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

		// Routes publiques pour l'inscription et la création de tickets
		api.GET("/departments/active", handlers.DepartmentHandler.GetActive)
		api.GET("/filiales/active", handlers.FilialeHandler.GetActive)
		api.GET("/software/active", handlers.SoftwareHandler.GetActive)

		// Route WebSocket pour les notifications en temps réel (authentification dans le handler)
		// Note: Cette route doit être avant le middleware AuthMiddleware car elle utilise un protocole différent
		if handlers.WebSocketHandler != nil {
			api.GET("/ws", handlers.WebSocketHandler.HandleWebSocket)
		}

		// Routes protégées (nécessitent authentification)
		api.Use(middleware.AuthMiddleware())
		api.Use(middleware.PerfMiddleware())
		api.Use(middleware.AuditLogMiddleware(auditLogRepo))
		{
			// Diagnostic
			if handlers.DiagnosticHandler != nil {
				api.GET("/diagnostic/it-users", handlers.DiagnosticHandler.GetITUsersInfo)
			}

			// Utilisateurs
			SetupUserRoutes(api, handlers.UserHandler)

			// Rôles
			SetupRoleRoutes(api, handlers.RoleHandler)

			// Permissions
			SetupPermissionRoutes(api, handlers.PermissionHandler)

			// Tickets - Les routes spécifiques doivent être définies avant les routes génériques
			// Donc on définit d'abord les routes de timesheet et delay-justification
			SetupTicketTimesheetRoutes(api, handlers.TimesheetHandler)
			SetupTicketDelayJustificationRoutes(api, handlers.DelayHandler)
			SetupTicketAuditRoutes(api, handlers.AuditHandler)
			// Puis les routes principales des tickets
			SetupTicketRoutes(api, handlers.TicketHandler, handlers.TicketAttachmentHandler, handlers.TicketCategoryHandler, handlers.TicketSolutionHandler)

			// Tickets internes (départements non-IT) — route /panier enregistrée avant le groupe pour éviter que /:id capture "panier"
			if handlers.TicketInternalHandler != nil {
				api.GET("/ticket-internes/panier", handlers.TicketInternalHandler.GetMyPanier)
				SetupTicketInternesRoutes(api, handlers.TicketInternalHandler)
			}

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
			SetupAssetRoutes(api, handlers.AssetHandler, handlers.AssetCategoryHandler, handlers.AssetSoftwareHandler)

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

			// Sièges
			SetupOfficeRoutes(api, handlers.OfficeHandler)

			// Départements
			SetupDepartmentRoutes(api, handlers.DepartmentHandler)

			// Filiales
			SetupFilialeRoutes(api, handlers.FilialeHandler, handlers.FilialeSoftwareHandler)
			SetupFilialeSoftwareRoutes(api, handlers.FilialeSoftwareHandler)

			// Logiciels
			SetupSoftwareRoutes(api, handlers.SoftwareHandler, handlers.FilialeSoftwareHandler)

			// Timesheet
			SetupTimesheetRoutes(api, handlers.TimesheetHandler)
			SetupUserTimesheetRoutes(api, handlers.TimesheetHandler)
			SetupProjectTimesheetRoutes(api, handlers.TimesheetHandler)
		}
	}
}

// Handlers contient toutes les instances de handlers
type Handlers struct {
	AuthHandler               *handlers.AuthHandler
	UserHandler               *handlers.UserHandler
	RoleHandler               *handlers.RoleHandler
	PermissionHandler         *handlers.PermissionHandler
	TicketHandler             *handlers.TicketHandler
	TicketAttachmentHandler   *handlers.TicketAttachmentHandler
	TicketCategoryHandler     *handlers.TicketCategoryHandler
	TicketSolutionHandler     *handlers.TicketSolutionHandler
	TicketInternalHandler     *handlers.TicketInternalHandler
	IncidentHandler           *handlers.IncidentHandler
	ChangeHandler             *handlers.ChangeHandler
	ServiceRequestHandler     *handlers.ServiceRequestHandler
	ServiceRequestTypeHandler *handlers.ServiceRequestTypeHandler
	TimeEntryHandler          *handlers.TimeEntryHandler
	DelayHandler              *handlers.DelayHandler
	AssetHandler              *handlers.AssetHandler
	AssetCategoryHandler      *handlers.AssetCategoryHandler
	AssetSoftwareHandler      *handlers.AssetSoftwareHandler
	SLAHandler                *handlers.SLAHandler
	NotificationHandler       *handlers.NotificationHandler
	KnowledgeArticleHandler   *handlers.KnowledgeArticleHandler
	KnowledgeCategoryHandler  *handlers.KnowledgeCategoryHandler
	ProjectHandler            *handlers.ProjectHandler
	DailyDeclarationHandler   *handlers.DailyDeclarationHandler
	WeeklyDeclarationHandler  *handlers.WeeklyDeclarationHandler
	PerformanceHandler        *handlers.PerformanceHandler
	ReportHandler             *handlers.ReportHandler
	SearchHandler             *handlers.SearchHandler
	StatisticsHandler         *handlers.StatisticsHandler
	AuditHandler              *handlers.AuditHandler
	SettingsHandler           *handlers.SettingsHandler
	RequestSourceHandler      *handlers.RequestSourceHandler
	BackupHandler             *handlers.BackupHandler
	TimesheetHandler          *handlers.TimesheetHandler
	OfficeHandler             *handlers.OfficeHandler
	DepartmentHandler         *handlers.DepartmentHandler
	FilialeHandler            *handlers.FilialeHandler
	SoftwareHandler           *handlers.SoftwareHandler
	FilialeSoftwareHandler    *handlers.FilialeSoftwareHandler
	WebSocketHandler          *handlers.WebSocketHandler
	DiagnosticHandler         *handlers.DiagnosticHandler
}
