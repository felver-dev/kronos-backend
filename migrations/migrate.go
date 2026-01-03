package migrations

import (
	"log"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// RunMigrations exécute toutes les migrations pour créer les tables
func RunMigrations() error {
	log.Println("🔄 Démarrage des migrations...")

	// Tables de base (authentification et utilisateurs)
	if err := database.DB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.UserSession{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables d'authentification et utilisateurs créées")

	// Tables de tickets
	if err := database.DB.AutoMigrate(
		&models.Ticket{},
		&models.TicketComment{},
		&models.TicketHistory{},
		&models.TicketAttachment{},
		&models.TicketTag{},
		&models.TicketTagAssignment{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de tickets créées")

	// Tables d'incidents
	if err := database.DB.AutoMigrate(
		&models.Incident{},
		&models.IncidentAsset{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables d'incidents créées")

	// Tables de demandes de service
	if err := database.DB.AutoMigrate(
		&models.ServiceRequestType{},
		&models.ServiceRequest{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de demandes de service créées")

	// Tables de changements
	if err := database.DB.AutoMigrate(
		&models.Change{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de changements créées")

	// Tables de gestion du temps
	if err := database.DB.AutoMigrate(
		&models.TimeEntry{},
		&models.DailyDeclaration{},
		&models.DailyDeclarationTask{},
		&models.WeeklyDeclaration{},
		&models.WeeklyDeclarationTask{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de gestion du temps créées")

	// Tables de retards
	if err := database.DB.AutoMigrate(
		&models.Delay{},
		&models.DelayJustification{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de retards créées")

	// Tables d'actifs IT
	if err := database.DB.AutoMigrate(
		&models.AssetCategory{},
		&models.Asset{},
		&models.TicketAsset{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables d'actifs IT créées")

	// Tables de SLA
	if err := database.DB.AutoMigrate(
		&models.SLA{},
		&models.TicketSLA{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de SLA créées")

	// Tables de notifications
	if err := database.DB.AutoMigrate(
		&models.Notification{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de notifications créées")

	// Tables de base de connaissances
	if err := database.DB.AutoMigrate(
		&models.KnowledgeCategory{},
		&models.KnowledgeArticle{},
		&models.KnowledgeArticleAttachment{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de base de connaissances créées")

	// Tables de projets
	if err := database.DB.AutoMigrate(
		&models.Project{},
		&models.TicketProject{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de projets créées")

	// Tables de paramétrage
	if err := database.DB.AutoMigrate(
		&models.Setting{},
		&models.RequestSource{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables de paramétrage créées")

	// Tables d'audit et sauvegarde
	if err := database.DB.AutoMigrate(
		&models.AuditLog{},
		&models.BackupConfiguration{},
		&models.Backup{},
	); err != nil {
		return err
	}
	log.Println("✅ Tables d'audit et sauvegarde créées")

	log.Println("🎉 Toutes les migrations ont été exécutées avec succès!")
	return nil
}

// SeedData insère les données initiales (rôles, permissions, etc.)
func SeedData() error {
	log.Println("🌱 Démarrage du seeding des données initiales...")

	// Vérifier si les rôles existent déjà
	var roleCount int64
	database.DB.Model(&models.Role{}).Count(&roleCount)
	if roleCount > 0 {
		log.Println("ℹ️  Les données initiales existent déjà, seeding ignoré")
		return nil
	}

	// Créer les rôles système
	roles := []models.Role{
		{Name: "DSI", Description: "Directeur des Systèmes d'Information", IsSystem: true},
		{Name: "RESPONSABLE_IT", Description: "Responsable IT", IsSystem: true},
		{Name: "TECHNICIEN_IT", Description: "Technicien IT", IsSystem: true},
	}

	for _, role := range roles {
		if err := database.DB.Create(&role).Error; err != nil {
			log.Printf("⚠️  Erreur lors de la création du rôle %s: %v", role.Name, err)
		}
	}

	log.Println("✅ Données initiales insérées avec succès!")
	return nil
}
