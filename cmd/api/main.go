package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}
	defer database.Close()

	// Initialiser les repositories
	repos := initializeRepositories()

	// Initialiser les services
	svcs := initializeServices(repos)

	// Initialiser les handlers
	handlers := initializeHandlers(svcs)

	// Configurer le routeur Gin
	router := gin.Default()

	// Configurer toutes les routes
	routes.SetupRoutes(router, handlers)

	// Démarrer le serveur HTTP
	port := config.AppConfig.AppPort
	if port == "" {
		port = "8080"
	}

	// Gérer l'arrêt gracieux
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
		}
	}()

	log.Printf("🚀 Serveur démarré sur le port %s", port)
	log.Printf("📚 API disponible sur http://localhost:%s/api/v1", port)

	// Attendre un signal d'arrêt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Arrêt du serveur...")
}

// Repositories contient toutes les instances de repositories
type Repositories struct {
	UserRepository               repositories.UserRepository
	RoleRepository               repositories.RoleRepository
	TicketRepository             repositories.TicketRepository
	TicketCommentRepository      repositories.TicketCommentRepository
	TicketHistoryRepository      repositories.TicketHistoryRepository
	IncidentRepository           repositories.IncidentRepository
	ChangeRepository             repositories.ChangeRepository
	ServiceRequestRepository     repositories.ServiceRequestRepository
	ServiceRequestTypeRepository repositories.ServiceRequestTypeRepository
	TimeEntryRepository          repositories.TimeEntryRepository
	DelayRepository              repositories.DelayRepository
	DelayJustificationRepository repositories.DelayJustificationRepository
	AssetRepository              repositories.AssetRepository
	AssetCategoryRepository      repositories.AssetCategoryRepository
	SLARepository                repositories.SLARepository
	TicketSLARepository          repositories.TicketSLARepository
	NotificationRepository       repositories.NotificationRepository
	KnowledgeArticleRepository   repositories.KnowledgeArticleRepository
	KnowledgeCategoryRepository  repositories.KnowledgeCategoryRepository
	ProjectRepository            repositories.ProjectRepository
	DailyDeclarationRepository   repositories.DailyDeclarationRepository
	WeeklyDeclarationRepository  repositories.WeeklyDeclarationRepository
	UserSessionRepository        repositories.UserSessionRepository
	TicketAssetRepository        repositories.TicketAssetRepository
}

// initializeRepositories initialise tous les repositories
func initializeRepositories() *Repositories {
	return &Repositories{
		UserRepository:               repositories.NewUserRepository(),
		RoleRepository:               repositories.NewRoleRepository(),
		TicketRepository:             repositories.NewTicketRepository(),
		TicketCommentRepository:      repositories.NewTicketCommentRepository(),
		TicketHistoryRepository:      repositories.NewTicketHistoryRepository(),
		IncidentRepository:           repositories.NewIncidentRepository(),
		ChangeRepository:             repositories.NewChangeRepository(),
		ServiceRequestRepository:     repositories.NewServiceRequestRepository(),
		ServiceRequestTypeRepository: repositories.NewServiceRequestTypeRepository(),
		TimeEntryRepository:          repositories.NewTimeEntryRepository(),
		DelayRepository:              repositories.NewDelayRepository(),
		DelayJustificationRepository: repositories.NewDelayJustificationRepository(),
		AssetRepository:              repositories.NewAssetRepository(),
		AssetCategoryRepository:      repositories.NewAssetCategoryRepository(),
		SLARepository:                repositories.NewSLARepository(),
		TicketSLARepository:          repositories.NewTicketSLARepository(),
		NotificationRepository:       repositories.NewNotificationRepository(),
		KnowledgeArticleRepository:   repositories.NewKnowledgeArticleRepository(),
		KnowledgeCategoryRepository:  repositories.NewKnowledgeCategoryRepository(),
		ProjectRepository:            repositories.NewProjectRepository(),
		DailyDeclarationRepository:   repositories.NewDailyDeclarationRepository(),
		WeeklyDeclarationRepository:  repositories.NewWeeklyDeclarationRepository(),
		UserSessionRepository:        repositories.NewUserSessionRepository(),
		TicketAssetRepository:        repositories.NewTicketAssetRepository(),
	}
}

// Services contient toutes les instances de services
type Services struct {
	AuthService              services.AuthService
	UserService              services.UserService
	TicketService            services.TicketService
	IncidentService          services.IncidentService
	ChangeService            services.ChangeService
	ServiceRequestService    services.ServiceRequestService
	TimeEntryService         services.TimeEntryService
	DelayService             services.DelayService
	AssetService             services.AssetService
	SLAService               services.SLAService
	NotificationService      services.NotificationService
	KnowledgeArticleService  services.KnowledgeArticleService
	ProjectService           services.ProjectService
	DailyDeclarationService  services.DailyDeclarationService
	WeeklyDeclarationService services.WeeklyDeclarationService
	PerformanceService       services.PerformanceService
	ReportService            services.ReportService
}

// initializeServices initialise tous les services
func initializeServices(repos *Repositories) *Services {
	return &Services{
		AuthService: services.NewAuthService(repos.UserRepository, repos.UserSessionRepository),
		UserService: services.NewUserService(repos.UserRepository, repos.RoleRepository),
		TicketService: services.NewTicketService(
			repos.TicketRepository,
			repos.UserRepository,
			repos.TicketCommentRepository,
			repos.TicketHistoryRepository,
		),
		IncidentService: services.NewIncidentService(
			repos.IncidentRepository,
			repos.TicketRepository,
			repos.TicketAssetRepository,
			repos.AssetRepository,
		),
		ChangeService: services.NewChangeService(
			repos.ChangeRepository,
			repos.TicketRepository,
			repos.UserRepository,
		),
		ServiceRequestService: services.NewServiceRequestService(
			repos.ServiceRequestRepository,
			repos.ServiceRequestTypeRepository,
			repos.TicketRepository,
			repos.UserRepository,
		),
		TimeEntryService: services.NewTimeEntryService(
			repos.TimeEntryRepository,
			repos.TicketRepository,
			repos.UserRepository,
		),
		DelayService: services.NewDelayService(
			repos.DelayRepository,
			repos.DelayJustificationRepository,
			repos.UserRepository,
		),
		AssetService: services.NewAssetService(
			repos.AssetRepository,
			repos.AssetCategoryRepository,
			repos.UserRepository,
		),
		SLAService: services.NewSLAService(
			repos.SLARepository,
			repos.TicketSLARepository,
			repos.TicketRepository,
		),
		NotificationService: services.NewNotificationService(
			repos.NotificationRepository,
			repos.UserRepository,
		),
		KnowledgeArticleService: services.NewKnowledgeArticleService(
			repos.KnowledgeArticleRepository,
			repos.KnowledgeCategoryRepository,
			repos.UserRepository,
		),
		ProjectService: services.NewProjectService(
			repos.ProjectRepository,
			repos.UserRepository,
		),
		DailyDeclarationService: services.NewDailyDeclarationService(
			repos.DailyDeclarationRepository,
			repos.TimeEntryRepository,
			repos.UserRepository,
		),
		WeeklyDeclarationService: services.NewWeeklyDeclarationService(
			repos.WeeklyDeclarationRepository,
			repos.UserRepository,
		),
		PerformanceService: services.NewPerformanceService(
			repos.TicketRepository,
			repos.TimeEntryRepository,
			repos.DelayRepository,
			repos.UserRepository,
		),
		ReportService: services.NewReportService(
			repos.TicketRepository,
			repos.SLARepository,
			repos.UserRepository,
		),
	}
}

// initializeHandlers initialise tous les handlers
func initializeHandlers(svcs *Services) *routes.Handlers {
	return &routes.Handlers{
		AuthHandler:              handlers.NewAuthHandler(svcs.AuthService),
		UserHandler:              handlers.NewUserHandler(svcs.UserService),
		TicketHandler:            handlers.NewTicketHandler(svcs.TicketService),
		IncidentHandler:          handlers.NewIncidentHandler(svcs.IncidentService),
		ChangeHandler:            handlers.NewChangeHandler(svcs.ChangeService),
		ServiceRequestHandler:    handlers.NewServiceRequestHandler(svcs.ServiceRequestService),
		TimeEntryHandler:         handlers.NewTimeEntryHandler(svcs.TimeEntryService),
		DelayHandler:             handlers.NewDelayHandler(svcs.DelayService),
		AssetHandler:             handlers.NewAssetHandler(svcs.AssetService),
		SLAHandler:               handlers.NewSLAHandler(svcs.SLAService),
		NotificationHandler:      handlers.NewNotificationHandler(svcs.NotificationService),
		KnowledgeArticleHandler:  handlers.NewKnowledgeArticleHandler(svcs.KnowledgeArticleService),
		ProjectHandler:           handlers.NewProjectHandler(svcs.ProjectService),
		DailyDeclarationHandler:  handlers.NewDailyDeclarationHandler(svcs.DailyDeclarationService),
		WeeklyDeclarationHandler: handlers.NewWeeklyDeclarationHandler(svcs.WeeklyDeclarationService),
		PerformanceHandler:       handlers.NewPerformanceHandler(svcs.PerformanceService),
		ReportHandler:            handlers.NewReportHandler(svcs.ReportService),
	}
}
