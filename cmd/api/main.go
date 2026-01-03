package main

import (
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

	// Initialiser tous les services
	authService := services.NewAuthService(userRepo, userSessionRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	ticketService := services.NewTicketService(ticketRepo, userRepo, ticketCommentRepo, ticketHistoryRepo)
	incidentService := services.NewIncidentService(incidentRepo, ticketRepo, ticketAssetRepo, assetRepo)
	serviceRequestService := services.NewServiceRequestService(serviceRequestRepo, serviceRequestTypeRepo, ticketRepo, userRepo)
	changeService := services.NewChangeService(changeRepo, ticketRepo, userRepo)
	timeEntryService := services.NewTimeEntryService(timeEntryRepo, ticketRepo, userRepo)
	delayService := services.NewDelayService(delayRepo, delayJustificationRepo, userRepo)
	assetService := services.NewAssetService(assetRepo, assetCategoryRepo, userRepo)
	slaService := services.NewSLAService(slaRepo, ticketSLARepo, ticketRepo)
	notificationService := services.NewNotificationService(notificationRepo, userRepo)
	knowledgeArticleService := services.NewKnowledgeArticleService(knowledgeArticleRepo, knowledgeCategoryRepo, userRepo)
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

	// Initialiser tous les handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	incidentHandler := handlers.NewIncidentHandler(incidentService)
	changeHandler := handlers.NewChangeHandler(changeService)
	serviceRequestHandler := handlers.NewServiceRequestHandler(serviceRequestService)
	timeEntryHandler := handlers.NewTimeEntryHandler(timeEntryService)
	delayHandler := handlers.NewDelayHandler(delayService)
	assetHandler := handlers.NewAssetHandler(assetService)
	slaHandler := handlers.NewSLAHandler(slaService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	knowledgeArticleHandler := handlers.NewKnowledgeArticleHandler(knowledgeArticleService)
	projectHandler := handlers.NewProjectHandler(projectService)
	dailyDeclarationHandler := handlers.NewDailyDeclarationHandler(dailyDeclarationService)
	weeklyDeclarationHandler := handlers.NewWeeklyDeclarationHandler(weeklyDeclarationService)
	performanceHandler := handlers.NewPerformanceHandler(performanceService)
	reportHandler := handlers.NewReportHandler(reportService)

	// Créer la structure Handlers
	appHandlers := &routes.Handlers{
		AuthHandler:              authHandler,
		UserHandler:              userHandler,
		TicketHandler:            ticketHandler,
		IncidentHandler:          incidentHandler,
		ChangeHandler:            changeHandler,
		ServiceRequestHandler:    serviceRequestHandler,
		TimeEntryHandler:         timeEntryHandler,
		DelayHandler:             delayHandler,
		AssetHandler:             assetHandler,
		SLAHandler:               slaHandler,
		NotificationHandler:      notificationHandler,
		KnowledgeArticleHandler:  knowledgeArticleHandler,
		ProjectHandler:           projectHandler,
		DailyDeclarationHandler:  dailyDeclarationHandler,
		WeeklyDeclarationHandler: weeklyDeclarationHandler,
		PerformanceHandler:       performanceHandler,
		ReportHandler:            reportHandler,
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

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("❌ Erreur lors du démarrage du serveur: %v", err)
	}
}
