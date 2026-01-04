// @title           ITSM Backend API
// @version         1.0
// @description     API REST pour la gestion des services IT (ITSM) - MCI CARE CI
// @termsOfService  http://swagger.io/terms/

// @contact.name   Support API
// @contact.email  support@mcicare.ci

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" suivi d'un espace puis le token JWT. Exemple: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

package main

import (
	_ "github.com/mcicare/itsm-backend/docs" // Import pour Swagger docs

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/handlers"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/routes"
	"github.com/mcicare/itsm-backend/internal/services"
)

func main() {
	// Charger la configuration
	config.LoadConfig()

	// Se connecter à la base de données
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Erreur de connexion à la base de données: %v", err)
	}
	defer database.Close()

	// Initialiser tous les repositories
	roleRepo := repositories.NewRoleRepository()
	userRepo := repositories.NewUserRepository()
	userSessionRepo := repositories.NewUserSessionRepository()
	ticketRepo := repositories.NewTicketRepository()
	ticketCommentRepo := repositories.NewTicketCommentRepository()
	ticketHistoryRepo := repositories.NewTicketHistoryRepository()
	ticketAttachmentRepo := repositories.NewTicketAttachmentRepository()
	ticketAssetRepo := repositories.NewTicketAssetRepository()
	incidentRepo := repositories.NewIncidentRepository()
	serviceRequestRepo := repositories.NewServiceRequestRepository()
	serviceRequestTypeRepo := repositories.NewServiceRequestTypeRepository()
	changeRepo := repositories.NewChangeRepository()
	timeEntryRepo := repositories.NewTimeEntryRepository()
	delayRepo := repositories.NewDelayRepository()
	delayJustificationRepo := repositories.NewDelayJustificationRepository()
	assetRepo := repositories.NewAssetRepository()
	assetCategoryRepo := repositories.NewAssetCategoryRepository()
	slaRepo := repositories.NewSLARepository()
	ticketSLARepo := repositories.NewTicketSLARepository()
	notificationRepo := repositories.NewNotificationRepository()
	knowledgeArticleRepo := repositories.NewKnowledgeArticleRepository()
	knowledgeCategoryRepo := repositories.NewKnowledgeCategoryRepository()
	projectRepo := repositories.NewProjectRepository()
	dailyDeclarationRepo := repositories.NewDailyDeclarationRepository()
	weeklyDeclarationRepo := repositories.NewWeeklyDeclarationRepository()
	auditLogRepo := repositories.NewAuditLogRepository()
	settingsRepo := repositories.NewSettingsRepository()
	requestSourceRepo := repositories.NewRequestSourceRepository()

	// Initialiser tous les services
	authService := services.NewAuthService(userRepo, userSessionRepo, roleRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	roleService := services.NewRoleService(roleRepo, userRepo)
	ticketService := services.NewTicketService(ticketRepo, userRepo, ticketCommentRepo, ticketHistoryRepo)
	ticketAttachmentService := services.NewTicketAttachmentService(ticketAttachmentRepo, ticketRepo, userRepo)
	incidentService := services.NewIncidentService(incidentRepo, ticketRepo, ticketAssetRepo, assetRepo)
	serviceRequestService := services.NewServiceRequestService(serviceRequestRepo, serviceRequestTypeRepo, ticketRepo, userRepo)
	serviceRequestTypeService := services.NewServiceRequestTypeService(serviceRequestTypeRepo, userRepo)
	changeService := services.NewChangeService(changeRepo, ticketRepo, userRepo)
	timeEntryService := services.NewTimeEntryService(timeEntryRepo, ticketRepo, userRepo)
	delayService := services.NewDelayService(delayRepo, delayJustificationRepo, userRepo)
	assetService := services.NewAssetService(assetRepo, assetCategoryRepo, userRepo, ticketAssetRepo, ticketRepo)
	assetCategoryService := services.NewAssetCategoryService(assetCategoryRepo, userRepo)
	slaService := services.NewSLAService(slaRepo, ticketSLARepo, ticketRepo)
	notificationService := services.NewNotificationService(notificationRepo, userRepo)
	knowledgeArticleService := services.NewKnowledgeArticleService(knowledgeArticleRepo, knowledgeCategoryRepo, userRepo)
	knowledgeCategoryService := services.NewKnowledgeCategoryService(knowledgeCategoryRepo, userRepo)
	projectService := services.NewProjectService(projectRepo, userRepo)
	dailyDeclarationService := services.NewDailyDeclarationService(dailyDeclarationRepo, timeEntryRepo, userRepo)
	weeklyDeclarationService := services.NewWeeklyDeclarationService(weeklyDeclarationRepo, userRepo)
	performanceService := services.NewPerformanceService(
		ticketRepo,
		timeEntryRepo,
		delayRepo,
		userRepo,
	)
	reportService := services.NewReportService(
		ticketRepo,
		slaRepo,
		userRepo,
	)
	searchService := services.NewSearchService(ticketRepo, assetRepo, knowledgeArticleRepo)
	statisticsService := services.NewStatisticsService(ticketRepo, slaRepo, userRepo, timeEntryRepo)
	auditService := services.NewAuditService(auditLogRepo)
	settingsService := services.NewSettingsService(settingsRepo)
	requestSourceService := services.NewRequestSourceService(requestSourceRepo)
	backupService := services.NewBackupService(settingsRepo)
	timesheetService := services.NewTimesheetService(
		timeEntryService,
		dailyDeclarationService,
		weeklyDeclarationService,
		ticketRepo,
		projectRepo,
		delayRepo,
		delayJustificationRepo,
		userRepo,
	)

	// Initialiser tous les handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	ticketAttachmentHandler := handlers.NewTicketAttachmentHandler(ticketAttachmentService)
	incidentHandler := handlers.NewIncidentHandler(incidentService)
	changeHandler := handlers.NewChangeHandler(changeService)
	serviceRequestHandler := handlers.NewServiceRequestHandler(serviceRequestService)
	serviceRequestTypeHandler := handlers.NewServiceRequestTypeHandler(serviceRequestTypeService)
	timeEntryHandler := handlers.NewTimeEntryHandler(timeEntryService)
	delayHandler := handlers.NewDelayHandler(delayService)
	assetHandler := handlers.NewAssetHandler(assetService)
	assetCategoryHandler := handlers.NewAssetCategoryHandler(assetCategoryService)
	slaHandler := handlers.NewSLAHandler(slaService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	knowledgeArticleHandler := handlers.NewKnowledgeArticleHandler(knowledgeArticleService)
	knowledgeCategoryHandler := handlers.NewKnowledgeCategoryHandler(knowledgeCategoryService)
	projectHandler := handlers.NewProjectHandler(projectService)
	dailyDeclarationHandler := handlers.NewDailyDeclarationHandler(dailyDeclarationService)
	weeklyDeclarationHandler := handlers.NewWeeklyDeclarationHandler(weeklyDeclarationService)
	performanceHandler := handlers.NewPerformanceHandler(performanceService)
	reportHandler := handlers.NewReportHandler(reportService)
	searchHandler := handlers.NewSearchHandler(searchService)
	statisticsHandler := handlers.NewStatisticsHandler(statisticsService)
	auditHandler := handlers.NewAuditHandler(auditService)
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	requestSourceHandler := handlers.NewRequestSourceHandler(requestSourceService)
	backupHandler := handlers.NewBackupHandler(backupService)
	timesheetHandler := handlers.NewTimesheetHandler(timesheetService)

	// Créer la structure Handlers
	appHandlers := &routes.Handlers{
		AuthHandler:              authHandler,
		UserHandler:              userHandler,
		RoleHandler:              roleHandler,
		TicketHandler:            ticketHandler,
		TicketAttachmentHandler:  ticketAttachmentHandler,
		IncidentHandler:          incidentHandler,
		ChangeHandler:            changeHandler,
		ServiceRequestHandler:      serviceRequestHandler,
		ServiceRequestTypeHandler:  serviceRequestTypeHandler,
		TimeEntryHandler:           timeEntryHandler,
		DelayHandler:             delayHandler,
		AssetHandler:             assetHandler,
		AssetCategoryHandler:      assetCategoryHandler,
		SLAHandler:               slaHandler,
		NotificationHandler:      notificationHandler,
		KnowledgeArticleHandler:    knowledgeArticleHandler,
		KnowledgeCategoryHandler:   knowledgeCategoryHandler,
		ProjectHandler:             projectHandler,
		DailyDeclarationHandler:  dailyDeclarationHandler,
		WeeklyDeclarationHandler: weeklyDeclarationHandler,
		PerformanceHandler:       performanceHandler,
		ReportHandler:            reportHandler,
		SearchHandler:            searchHandler,
		StatisticsHandler:       statisticsHandler,
		AuditHandler:             auditHandler,
		SettingsHandler:           settingsHandler,
		RequestSourceHandler:      requestSourceHandler,
		BackupHandler:             backupHandler,
		TimesheetHandler:          timesheetHandler,
	}

	// Configurer Gin
	if config.AppConfig.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Créer le routeur
	router := gin.Default()

	// Configurer les routes
	routes.SetupRoutes(router, appHandlers)

	// Démarrer le serveur
	port := ":" + config.AppConfig.AppPort
	log.Printf("🚀 Serveur démarré sur le port %s", config.AppConfig.AppPort)
	log.Printf("📡 API disponible sur http://localhost%s/api/v1", port)
	log.Printf("💚 Health check: http://localhost%s/health", port)
	log.Printf("📚 Swagger UI: http://localhost%s/swagger/index.html", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("❌ Erreur lors du démarrage du serveur: %v", err)
	}
}
