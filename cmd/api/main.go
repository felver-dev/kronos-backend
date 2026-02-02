// @title           ITSM Backend API
// @version         1.0
// @description     API REST pour la gestion des services IT (ITSM)
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
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
	"github.com/mcicare/itsm-backend/internal/routes"
	"github.com/mcicare/itsm-backend/internal/scope"
	"github.com/mcicare/itsm-backend/internal/services"
	"github.com/mcicare/itsm-backend/internal/websocket"
)

func main() {
	// Charger la configuration
	config.LoadConfig()

	// Se connecter à la base de données
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Erreur de connexion à la base de données: %v", err)
	}
	defer database.Close()

	// Initialiser le checker de table assignees pour le package scope
	// Cela évite les cycles d'importation
	scope.SetAssigneesTableChecker(func() bool {
		return database.DB.Migrator().HasTable(&models.TicketAssignee{})
	})

	// S'assurer que la table audit_logs existe
	if err := database.DB.AutoMigrate(&models.AuditLog{}); err != nil {
		log.Printf("⚠️  Avertissement: Impossible de migrer audit_logs: %v", err)
	}

	// S'assurer que les tables pour assignations et sous-tickets existent
	if err := database.DB.AutoMigrate(&models.Ticket{}, &models.TicketAssignee{}); err != nil {
		log.Printf("⚠️  Avertissement: Impossible de migrer ticket_assignees: %v", err)
	}

	// S'assurer que les tables de retards existent
	if err := database.DB.AutoMigrate(&models.Delay{}, &models.DelayJustification{}); err != nil {
		log.Printf("⚠️  Avertissement: Impossible de migrer delays: %v", err)
	}

	// Initialiser tous les repositories
	roleRepo := repositories.NewRoleRepository()
	permissionRepo := repositories.NewPermissionRepository()

	// Initialiser le getter de permissions pour le package scope
	// Cela évite les cycles d'importation
	// IMPORTANT: Doit être fait après la création de roleRepo
	scope.SetPermissionsGetter(func(roleName string) []string {
		role, err := roleRepo.FindByName(roleName)
		if err != nil {
			return []string{"tickets.view_own"}
		}
		permissions, err := roleRepo.GetPermissionsByRoleID(role.ID)
		if err != nil {
			return []string{"tickets.view_own"}
		}
		if len(permissions) == 0 {
			return []string{"tickets.view_own"}
		}
		return permissions
	})

	userRepo := repositories.NewUserRepository()
	userSessionRepo := repositories.NewUserSessionRepository()
	ticketRepo := repositories.NewTicketRepository()
	ticketCommentRepo := repositories.NewTicketCommentRepository()
	ticketHistoryRepo := repositories.NewTicketHistoryRepository()
	ticketAttachmentRepo := repositories.NewTicketAttachmentRepository()
	ticketCategoryRepo := repositories.NewTicketCategoryRepository()
	ticketSolutionRepo := repositories.NewTicketSolutionRepository()
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
	assetSoftwareRepo := repositories.NewAssetSoftwareRepository()
	slaRepo := repositories.NewSLARepository()
	ticketSLARepo := repositories.NewTicketSLARepository()
	notificationRepo := repositories.NewNotificationRepository()
	knowledgeArticleRepo := repositories.NewKnowledgeArticleRepository()
	knowledgeCategoryRepo := repositories.NewKnowledgeCategoryRepository()
	projectRepo := repositories.NewProjectRepository()
	projectBudgetExtRepo := repositories.NewProjectBudgetExtensionRepository()
	projectPhaseRepo := repositories.NewProjectPhaseRepository()
	projectFunctionRepo := repositories.NewProjectFunctionRepository()
	projectMemberRepo := repositories.NewProjectMemberRepository()
	projectPhaseMemberRepo := repositories.NewProjectPhaseMemberRepository()
	projectTaskRepo := repositories.NewProjectTaskRepository()
	dailyDeclarationRepo := repositories.NewDailyDeclarationRepository()
	weeklyDeclarationRepo := repositories.NewWeeklyDeclarationRepository()
	auditLogRepo := repositories.NewAuditLogRepository()
	settingsRepo := repositories.NewSettingsRepository()
	requestSourceRepo := repositories.NewRequestSourceRepository()
	officeRepo := repositories.NewOfficeRepository()
	departmentRepo := repositories.NewDepartmentRepository()
	filialeRepo := repositories.NewFilialeRepository()
	ticketInternalRepo := repositories.NewTicketInternalRepository()

	// Initialiser tous les services
	authService := services.NewAuthService(userRepo, userSessionRepo, roleRepo)
	userService := services.NewUserService(userRepo, roleRepo, departmentRepo, ticketRepo)
	roleService := services.NewRoleService(roleRepo, userRepo, permissionRepo, filialeRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Créer et démarrer le hub WebSocket pour les notifications en temps réel
	wsHub := websocket.NewHub()
	go wsHub.Run()
	log.Println("✅ Hub WebSocket démarré pour les notifications en temps réel")

	// Créer le service de notifications AVANT le ticketService (car ticketService en a besoin)
	notificationService := services.NewNotificationService(notificationRepo, userRepo, wsHub)

	ticketService := services.NewTicketService(ticketRepo, userRepo, ticketCommentRepo, ticketHistoryRepo, slaRepo, ticketSLARepo, ticketCategoryRepo, notificationRepo, notificationService, departmentRepo, filialeRepo, timeEntryRepo)
	ticketAttachmentService := services.NewTicketAttachmentService(ticketAttachmentRepo, ticketRepo, userRepo)
	ticketCategoryService := services.NewTicketCategoryService(ticketCategoryRepo)
	ticketSolutionService := services.NewTicketSolutionService(ticketSolutionRepo, ticketRepo, userRepo, roleRepo, knowledgeArticleRepo, knowledgeCategoryRepo)
	ticketInternalService := services.NewTicketInternalService(ticketInternalRepo, userRepo, departmentRepo, notificationService)
	incidentService := services.NewIncidentService(incidentRepo, ticketRepo, ticketAssetRepo, assetRepo)
	serviceRequestService := services.NewServiceRequestService(serviceRequestRepo, serviceRequestTypeRepo, ticketRepo, userRepo)
	serviceRequestTypeService := services.NewServiceRequestTypeService(serviceRequestTypeRepo, userRepo)
	changeService := services.NewChangeService(changeRepo, ticketRepo, userRepo)
	timeEntryService := services.NewTimeEntryService(timeEntryRepo, ticketRepo, userRepo, delayRepo)
	delayService := services.NewDelayService(delayRepo, delayJustificationRepo, userRepo, ticketRepo)
	assetService := services.NewAssetService(assetRepo, assetCategoryRepo, userRepo, ticketAssetRepo, ticketRepo)
	assetCategoryService := services.NewAssetCategoryService(assetCategoryRepo, assetRepo, userRepo)
	assetSoftwareService := services.NewAssetSoftwareService(assetSoftwareRepo, assetRepo)
	slaService := services.NewSLAService(slaRepo, ticketSLARepo, ticketRepo, ticketCategoryRepo)
	knowledgeArticleService := services.NewKnowledgeArticleService(knowledgeArticleRepo, knowledgeCategoryRepo, userRepo)
	knowledgeCategoryService := services.NewKnowledgeCategoryService(knowledgeCategoryRepo, userRepo)
	projectService := services.NewProjectService(projectRepo, userRepo, projectBudgetExtRepo, projectPhaseRepo, projectFunctionRepo, projectMemberRepo, projectPhaseMemberRepo, projectTaskRepo, notificationService)
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
		ticketInternalRepo,
		slaRepo,
		userRepo,
	)
	searchService := services.NewSearchService(ticketRepo, assetRepo, knowledgeArticleRepo, userRepo, timeEntryRepo)
	statisticsService := services.NewStatisticsService(ticketRepo, slaRepo, userRepo, timeEntryRepo)
	auditService := services.NewAuditService(auditLogRepo)
	settingsService := services.NewSettingsService(settingsRepo)
	requestSourceService := services.NewRequestSourceService(requestSourceRepo)
	backupService := services.NewBackupService(settingsRepo)
	officeService := services.NewOfficeService(officeRepo, filialeRepo)
	departmentService := services.NewDepartmentService(departmentRepo, officeRepo, filialeRepo)
	softwareRepo := repositories.NewSoftwareRepository()
	filialeSoftwareRepo := repositories.NewFilialeSoftwareRepository()
	filialeService := services.NewFilialeService(filialeRepo)
	softwareService := services.NewSoftwareService(softwareRepo)
	filialeSoftwareService := services.NewFilialeSoftwareService(filialeSoftwareRepo, filialeRepo, softwareRepo)
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
	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	ticketAttachmentHandler := handlers.NewTicketAttachmentHandler(ticketAttachmentService)
	ticketCategoryHandler := handlers.NewTicketCategoryHandler(ticketCategoryService)
	ticketSolutionHandler := handlers.NewTicketSolutionHandler(ticketSolutionService)
	ticketInternalHandler := handlers.NewTicketInternalHandler(ticketInternalService)
	incidentHandler := handlers.NewIncidentHandler(incidentService)
	changeHandler := handlers.NewChangeHandler(changeService)
	serviceRequestHandler := handlers.NewServiceRequestHandler(serviceRequestService)
	serviceRequestTypeHandler := handlers.NewServiceRequestTypeHandler(serviceRequestTypeService)
	timeEntryHandler := handlers.NewTimeEntryHandler(timeEntryService)
	delayHandler := handlers.NewDelayHandler(delayService)
	assetHandler := handlers.NewAssetHandler(assetService)
	assetCategoryHandler := handlers.NewAssetCategoryHandler(assetCategoryService)
	assetSoftwareHandler := handlers.NewAssetSoftwareHandler(assetSoftwareService)
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
	officeHandler := handlers.NewOfficeHandler(officeService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)
	filialeHandler := handlers.NewFilialeHandler(filialeService)
	softwareHandler := handlers.NewSoftwareHandler(softwareService)
	filialeSoftwareHandler := handlers.NewFilialeSoftwareHandler(filialeSoftwareService)
	wsHandler := handlers.NewWebSocketHandler(wsHub)
	diagnosticHandler := handlers.NewDiagnosticHandler(filialeRepo)

	// Créer la structure Handlers
	appHandlers := &routes.Handlers{
		AuthHandler:               authHandler,
		UserHandler:               userHandler,
		RoleHandler:               roleHandler,
		PermissionHandler:         permissionHandler,
		TicketHandler:             ticketHandler,
		TicketAttachmentHandler:   ticketAttachmentHandler,
		TicketCategoryHandler:     ticketCategoryHandler,
		TicketSolutionHandler:     ticketSolutionHandler,
		TicketInternalHandler:     ticketInternalHandler,
		IncidentHandler:           incidentHandler,
		ChangeHandler:             changeHandler,
		ServiceRequestHandler:     serviceRequestHandler,
		ServiceRequestTypeHandler: serviceRequestTypeHandler,
		TimeEntryHandler:          timeEntryHandler,
		DelayHandler:              delayHandler,
		AssetHandler:              assetHandler,
		AssetCategoryHandler:      assetCategoryHandler,
		AssetSoftwareHandler:      assetSoftwareHandler,
		SLAHandler:                slaHandler,
		NotificationHandler:       notificationHandler,
		KnowledgeArticleHandler:   knowledgeArticleHandler,
		KnowledgeCategoryHandler:  knowledgeCategoryHandler,
		ProjectHandler:            projectHandler,
		DailyDeclarationHandler:   dailyDeclarationHandler,
		WeeklyDeclarationHandler:  weeklyDeclarationHandler,
		PerformanceHandler:        performanceHandler,
		ReportHandler:             reportHandler,
		SearchHandler:             searchHandler,
		StatisticsHandler:         statisticsHandler,
		AuditHandler:              auditHandler,
		SettingsHandler:           settingsHandler,
		RequestSourceHandler:      requestSourceHandler,
		BackupHandler:             backupHandler,
		TimesheetHandler:          timesheetHandler,
		OfficeHandler:             officeHandler,
		DepartmentHandler:         departmentHandler,
		FilialeHandler:            filialeHandler,
		SoftwareHandler:           softwareHandler,
		FilialeSoftwareHandler:    filialeSoftwareHandler,
		WebSocketHandler:          wsHandler,
		DiagnosticHandler:         diagnosticHandler,
	}

	// Configurer Gin
	if config.AppConfig.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Créer le routeur
	router := gin.Default()

	// Configurer les routes
	routes.SetupRoutes(router, appHandlers, auditLogRepo)

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
