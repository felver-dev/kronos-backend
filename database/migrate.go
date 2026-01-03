package database

import (
	"fmt"
	"log"

	"github.com/mcicare/itsm-backend/internal/models"
)

// AutoMigrate exécute les migrations automatiques pour créer les tables
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🔄 Démarrage des migrations automatiques...")

	// Désactiver temporairement les contraintes de clé étrangère
	sqlDB, _ := DB.DB()
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Créer toutes les tables dans le bon ordre
	err := DB.AutoMigrate(
		// Tables de base (authentification et utilisateurs) - en premier
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.UserSession{},

		// Tables de tickets
		&models.Ticket{},
		&models.TicketAttachment{},
		&models.TicketComment{},
		&models.TicketHistory{},
		&models.TicketTag{},
		&models.TicketTagAssignment{},

		// Tables d'incidents
		&models.Incident{},
		&models.IncidentAsset{},

		// Tables de demandes de service
		&models.ServiceRequestType{},
		&models.ServiceRequest{},

		// Tables de changements
		&models.Change{},

		// Tables de gestion du temps
		&models.TimeEntry{},
		&models.DailyDeclaration{},
		&models.DailyDeclarationTask{},
		&models.WeeklyDeclaration{},
		&models.WeeklyDeclarationTask{},

		// Tables de retards
		&models.Delay{},
		&models.DelayJustification{},

		// Tables d'actifs IT
		&models.AssetCategory{},
		&models.Asset{},
		&models.TicketAsset{},

		// Tables de SLA
		&models.SLA{},
		&models.TicketSLA{},

		// Tables de notifications
		&models.Notification{},

		// Tables de base de connaissances
		&models.KnowledgeCategory{},
		&models.KnowledgeArticle{},
		&models.KnowledgeArticleAttachment{},

		// Tables de projets
		&models.Project{},
		&models.TicketProject{},

		// Tables de paramétrage
		&models.Setting{},
		&models.RequestSource{},

		// Tables d'audit et sauvegarde - en dernier car elles peuvent référencer d'autres tables
		&models.AuditLog{},
		&models.BackupConfiguration{},
		&models.Backup{},
	)

	// Réactiver les contraintes
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	if err != nil {
		return fmt.Errorf("échec des migrations: %w", err)
	}

	log.Println("✅ Migrations automatiques terminées avec succès")
	log.Println("   Toutes les tables ont été créées avec leurs relations")

	return nil
}

// DropAllTables supprime toutes les tables (ATTENTION: à utiliser uniquement en développement!)
func DropAllTables() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("⚠️  Suppression de toutes les tables...")

	// Désactiver les contraintes de clé étrangère
	sqlDB, _ := DB.DB()
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Supprimer les tables dans l'ordre inverse des dépendances
	err := DB.Migrator().DropTable(
		&models.Backup{},
		&models.BackupConfiguration{},
		&models.AuditLog{},
		&models.RequestSource{},
		&models.Setting{},
		&models.TicketProject{},
		&models.Project{},
		&models.KnowledgeArticleAttachment{},
		&models.KnowledgeArticle{},
		&models.KnowledgeCategory{},
		&models.Notification{},
		&models.TicketSLA{},
		&models.SLA{},
		&models.TicketAsset{},
		&models.Asset{},
		&models.AssetCategory{},
		&models.DelayJustification{},
		&models.Delay{},
		&models.WeeklyDeclarationTask{},
		&models.WeeklyDeclaration{},
		&models.DailyDeclarationTask{},
		&models.DailyDeclaration{},
		&models.TimeEntry{},
		&models.Change{},
		&models.ServiceRequest{},
		&models.ServiceRequestType{},
		&models.IncidentAsset{},
		&models.Incident{},
		&models.TicketTagAssignment{},
		&models.TicketTag{},
		&models.TicketHistory{},
		&models.TicketComment{},
		&models.TicketAttachment{},
		&models.Ticket{},
		&models.UserSession{},
		&models.User{},
		&models.RolePermission{},
		&models.Permission{},
		&models.Role{},
	)

	// Réactiver les contraintes
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	if err != nil {
		return fmt.Errorf("échec de la suppression des tables: %w", err)
	}

	log.Println("✅ Toutes les tables ont été supprimées")
	return nil
}

// ResetDatabase supprime et recrée toutes les tables (ATTENTION: à utiliser uniquement en développement!)
func ResetDatabase() error {
	if err := DropAllTables(); err != nil {
		return err
	}
	return AutoMigrate()
}
