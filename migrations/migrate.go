package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL
	"github.com/mcicare/itsm-backend/config"
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// RunMigrations exécute toutes les migrations pour créer les tables
func RunMigrations() error {
	log.Println("🔄 Démarrage des migrations...")

	// Vérifier que la connexion est valide
	sqlDB, err := database.DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance SQL: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("la connexion à la base de données n'est pas valide: %w", err)
	}

	// Tables de base (authentification et utilisateurs)
	if err := database.DB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.UserSession{},
	); err != nil {
		// Si l'erreur est "table doesn't exist in engine" ou "Tablespace exists", la base est corrompue
		// Il faut supprimer et recréer la base de données
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "doesn't exist in engine") ||
			strings.Contains(errMsg, "tablespace") ||
			strings.Contains(errMsg, "discard the tablespace") {
			log.Println("⚠️  Détection d'une incohérence dans la base de données")
			log.Println("🔄 Suppression et recréation de la base de données...")

			// Fermer la connexion actuelle
			database.Close()

			// Supprimer et recréer la base de données
			if err := recreateDatabase(); err != nil {
				return fmt.Errorf("erreur lors de la recréation de la base de données: %w", err)
			}

			// Se reconnecter
			if err := database.Connect(); err != nil {
				return fmt.Errorf("erreur lors de la reconnexion: %w", err)
			}

			// Réessayer les migrations
			log.Println("🔄 Nouvelle tentative de migration...")
			return RunMigrations()
		}
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

// recreateDatabase supprime toutes les tables et recrée la base de données
func recreateDatabase() error {
	// D'abord, essayer de supprimer toutes les tables
	log.Println("🗑️  Suppression de toutes les tables...")

	// Récupérer la liste de toutes les tables
	var tables []string
	rows, err := database.DB.Raw("SHOW TABLES").Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err == nil {
				tables = append(tables, tableName)
			}
		}
	}

	// Supprimer toutes les tables une par une avec gestion des tablespaces
	for _, table := range tables {
		// Essayer de supprimer le tablespace d'abord (pour les tables InnoDB)
		discardQuery := fmt.Sprintf("ALTER TABLE `%s` DISCARD TABLESPACE", table)
		database.DB.Exec(discardQuery) // Ignorer l'erreur si la table n'existe pas ou n'a pas de tablespace

		// Supprimer la table
		dropTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
		if err := database.DB.Exec(dropTableQuery).Error; err != nil {
			log.Printf("⚠️  Erreur lors de la suppression de la table %s: %v", table, err)
			// Essayer de forcer la suppression
			forceDropQuery := fmt.Sprintf("DROP TABLE `%s`", table)
			database.DB.Exec(forceDropQuery) // Ignorer l'erreur
		}
	}

	// Maintenant, se connecter sans base de données spécifiée pour supprimer la base
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=%t&loc=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBCharset,
		config.AppConfig.DBParseTime,
		config.AppConfig.DBLoc,
	)

	// Utiliser database/sql pour supprimer et recréer
	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("erreur de connexion: %w", err)
	}
	defer db.Close()

	// Supprimer la base de données si elle existe (maintenant qu'elle est vide)
	dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.AppConfig.DBName)
	if _, err := db.Exec(dropQuery); err != nil {
		// Si la suppression échoue, ce n'est pas grave, on continue
		log.Printf("⚠️  Impossible de supprimer la base de données (peut être déjà vide): %v", err)
	} else {
		log.Printf("🗑️  Base de données '%s' supprimée", config.AppConfig.DBName)
	}

	// Recréer la base de données
	createQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.AppConfig.DBName)
	if _, err := db.Exec(createQuery); err != nil {
		return fmt.Errorf("erreur lors de la création de la base: %w", err)
	}
	log.Printf("✅ Base de données '%s' recréée", config.AppConfig.DBName)

	return nil
}
