package database

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/utils"
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

	// Seed des rôles par défaut
	if err := seedDefaultRoles(); err != nil {
		log.Printf("⚠️  Erreur lors du seeding des rôles: %v", err)
		// Ne pas bloquer les migrations si le seeding échoue
	}

	// Seed de l'utilisateur admin par défaut
	if err := seedDefaultAdmin(); err != nil {
		log.Printf("⚠️  Erreur lors du seeding de l'admin: %v", err)
		// Ne pas bloquer les migrations si le seeding échoue
	}

	log.Println("✅ Migrations automatiques terminées avec succès")
	log.Println("   Toutes les tables ont été créées avec leurs relations")

	return nil
}

// seedDefaultRoles crée les rôles par défaut s'ils n'existent pas
func seedDefaultRoles() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🌱 Seeding des rôles par défaut...")

	defaultRoles := []models.Role{
		{
			Name:        "DSI",
			Description: "DSI / Administrateur - Accès total",
			IsSystem:    true,
		},
		{
			Name:        "RESPONSABLE_IT",
			Description: "Responsable IT - Supervision et validation",
			IsSystem:    true,
		},
		{
			Name:        "TECHNICIEN_IT",
			Description: "Technicien IT - Traitement des tickets",
			IsSystem:    true,
		},
		{
			Name:        "USER",
			Description: "Utilisateur standard - Accès limité",
			IsSystem:    true,
		},
		{
			Name:        "CLIENT",
			Description: "Client - Accès client",
			IsSystem:    true,
		},
	}

	for _, role := range defaultRoles {
		var existingRole models.Role
		result := DB.Where("name = ?", role.Name).First(&existingRole)
		
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Le rôle n'existe pas, le créer
			if err := DB.Create(&role).Error; err != nil {
				log.Printf("⚠️  Erreur lors de la création du rôle %s: %v", role.Name, err)
			} else {
				log.Printf("   ✅ Rôle créé: %s", role.Name)
			}
		} else if result.Error != nil {
			// Autre erreur
			log.Printf("⚠️  Erreur lors de la vérification du rôle %s: %v", role.Name, result.Error)
		} else {
			log.Printf("   ℹ️  Rôle déjà existant: %s", role.Name)
		}
	}

	log.Println("✅ Seeding des rôles terminé")
	return nil
}

// seedDefaultAdmin crée l'utilisateur admin par défaut s'il n'existe pas
func seedDefaultAdmin() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🌱 Seeding de l'utilisateur admin par défaut...")

	// Vérifier si l'admin existe déjà
	var existingAdmin models.User
	result := DB.Where("email = ?", "admin@mcicareci.com").First(&existingAdmin)
	
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if result.Error == nil {
			log.Println("   ℹ️  Utilisateur admin déjà existant")
			return nil
		}
		// Autre erreur
		return fmt.Errorf("erreur lors de la vérification de l'admin: %w", result.Error)
	}

	// Récupérer le rôle DSI
	var dsiRole models.Role
	if err := DB.Where("name = ?", "DSI").First(&dsiRole).Error; err != nil {
		return fmt.Errorf("le rôle DSI n'existe pas, veuillez d'abord exécuter le seeding des rôles: %w", err)
	}

	// Hasher le mot de passe
	passwordHash, err := utils.HashPassword("admin12345")
	if err != nil {
		return fmt.Errorf("erreur lors du hashage du mot de passe: %w", err)
	}

	// Créer l'utilisateur admin
	admin := models.User{
		Username:     "admin",
		Email:        "admin@mcicareci.com",
		PasswordHash: passwordHash,
		FirstName:    "Administrateur",
		LastName:     "Système",
		RoleID:       dsiRole.ID,
		IsActive:     true,
		CreatedByID:  nil, // Pas de créateur pour l'admin système
	}

	if err := DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("erreur lors de la création de l'admin: %w", err)
	}

	log.Println("   ✅ Utilisateur admin créé:")
	log.Printf("      Email: admin@mcicareci.com")
	log.Printf("      Mot de passe: admin12345")
	log.Printf("      Rôle: DSI")
	log.Println("✅ Seeding de l'admin terminé")
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
