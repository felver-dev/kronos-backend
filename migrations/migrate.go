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
	return runMigrationsWithRetry(0)
}

// runMigrationsWithRetry exécute les migrations avec un mécanisme de retry limité
func runMigrationsWithRetry(retryCount int) error {
	if retryCount > 1 {
		return fmt.Errorf("trop de tentatives de recréation de la base de données (max: 1)")
	}

	log.Println("🔄 Démarrage des migrations...")

	// Vérifier que la connexion est valide
	sqlDB, err := database.DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance SQL: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("la connexion à la base de données n'est pas valide: %w", err)
	}

	// Supprimer toutes les tables existantes pour repartir sur une base propre
	// (utile en développement, à désactiver en production)
	log.Println("🧹 Nettoyage des tables existantes...")
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Liste des tables à supprimer (dans l'ordre inverse des dépendances)
	tables := []string{
		"backups", "backup_configurations", "audit_logs",
		"request_sources", "settings",
		"ticket_projects", "projects",
		"knowledge_article_attachments", "knowledge_articles", "knowledge_categories",
		"notifications",
		"ticket_slas", "slas",
		"ticket_assets", "assets", "asset_categories",
		"delay_justifications", "delays",
		"weekly_declaration_tasks", "weekly_declarations",
		"daily_declaration_tasks", "daily_declarations",
		"time_entries",
		"changes",
		"service_requests", "service_request_types",
		"incident_assets", "incidents",
		"ticket_tag_assignments", "ticket_tags",
		"ticket_attachments", "ticket_comments", "ticket_histories",
		"tickets",
		"user_sessions", "users", "role_permissions", "permissions", "roles",
	}

	for _, table := range tables {
		_, _ = sqlDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	}

	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	log.Println("✅ Tables existantes supprimées")

	// Tables de base (authentification et utilisateurs)
	if err := database.DB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.User{},
		&models.UserSession{},
	); err != nil {
		// Si l'erreur est liée aux tablespaces, essayer de nettoyer d'abord
		errMsg := strings.ToLower(err.Error())
		if (strings.Contains(errMsg, "doesn't exist in engine") ||
			strings.Contains(errMsg, "tablespace") ||
			strings.Contains(errMsg, "discard the tablespace")) && retryCount == 0 {
			log.Println("⚠️  Erreur de tablespace détectée")
			log.Println("💡 Tentative de nettoyage des tablespaces orphelins...")

			// Essayer de nettoyer les tablespaces orphelins directement
			if err := cleanupOrphanedTablespaces(); err == nil {
				log.Println("✅ Nettoyage réussi, nouvelle tentative...")
				return runMigrationsWithRetry(retryCount + 1)
			}

			// Si le nettoyage échoue, recréer la base
			log.Println("🔄 Le nettoyage automatique a échoué, recréation de la base de données...")
			database.Close()
			if err := recreateDatabase(); err != nil {
				return fmt.Errorf("erreur lors de la recréation: %w", err)
			}
			if err := database.Connect(); err != nil {
				return fmt.Errorf("erreur lors de la reconnexion: %w", err)
			}
			// Réessayer une seule fois
			log.Println("🔄 Nouvelle tentative de migration...")
			return runMigrationsWithRetry(retryCount + 1)
		}

		// Si c'est toujours une erreur de tablespace après retry, donner des instructions
		if strings.Contains(errMsg, "tablespace") && retryCount > 0 {
			log.Println("")
			log.Println("❌ ERREUR: Les fichiers de tablespace persistent dans le répertoire de données MySQL/MariaDB.")
			log.Println("")
			log.Println("📋 SOLUTION MANUELLE:")
			log.Println("   1. Arrêtez MySQL/MariaDB (via XAMPP Control Panel)")
			log.Println("   2. Supprimez le répertoire de la base de données:")
			log.Printf("      - XAMPP: C:\\xampp\\mysql\\data\\%s\\", config.AppConfig.DBName)
			log.Println("      - Ou le répertoire de données MySQL configuré")
			log.Println("   3. Redémarrez MySQL/MariaDB")
			log.Println("   4. Relancez les migrations")
			log.Println("")
			return fmt.Errorf("impossible de résoudre le problème de tablespace automatiquement")
		}

		return err
	}
	log.Println("✅ Tables d'authentification et utilisateurs créées")

	// Tables de tickets - créer Ticket seul d'abord
	log.Println("🔄 Création de la table tickets...")
	// Désactiver temporairement les contraintes de clé étrangère
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err := database.DB.AutoMigrate(&models.Ticket{}); err != nil {
		_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
		return fmt.Errorf("erreur lors de la création de la table tickets: %w", err)
	}
	// Réactiver les contraintes
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	log.Println("✅ Table tickets créée")
	
	// Ensuite créer TicketAttachment qui dépend de Ticket
	log.Println("🔄 Création de la table ticket_attachments...")
	if err := database.DB.AutoMigrate(&models.TicketAttachment{}); err != nil {
		return fmt.Errorf("erreur lors de la création de la table ticket_attachments: %w", err)
	}
	log.Println("✅ Table ticket_attachments créée")
	
	// Ensuite créer les autres tables de tickets qui dépendent de Ticket
	log.Println("🔄 Création des autres tables de tickets...")
	if err := database.DB.AutoMigrate(
		&models.TicketComment{},
		&models.TicketHistory{},
		&models.TicketTag{},
		&models.TicketTagAssignment{},
	); err != nil {
		return fmt.Errorf("erreur lors de la création des autres tables de tickets: %w", err)
	}
	log.Println("✅ Tables de tickets créées")

	// Tables d'incidents
	log.Println("🔄 Création des tables d'incidents...")
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err := database.DB.AutoMigrate(
		&models.Incident{},
		&models.IncidentAsset{},
	); err != nil {
		_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
		return fmt.Errorf("erreur lors de la création des tables d'incidents: %w", err)
	}
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
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

// cleanupOrphanedTablespaces nettoie les tablespaces orphelins
func cleanupOrphanedTablespaces() error {
	sqlDB, err := database.DB.DB()
	if err != nil {
		return err
	}

	// Liste des tables à vérifier (les premières tables créées)
	tablesToCheck := []string{"roles", "permissions", "role_permissions", "users", "user_sessions"}

	for _, tableName := range tablesToCheck {
		// Vérifier si la table existe
		var exists int
		query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s'",
			config.AppConfig.DBName, tableName)
		if err := sqlDB.QueryRow(query).Scan(&exists); err == nil && exists == 0 {
			// La table n'existe pas, mais le tablespace peut exister
			// Essayer de créer une table temporaire avec le même nom pour forcer MySQL à nettoyer
			// Puis la supprimer immédiatement
			tempQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s_temp_cleanup` (id INT) ENGINE=InnoDB", tableName)
			sqlDB.Exec(tempQuery)
			sqlDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s_temp_cleanup`", tableName))
		}
	}

	return nil
}

// recreateDatabase supprime et recrée la base de données (méthode rapide)
// ATTENTION: Cette fonction supprime TOUTES les données de la base de données
func recreateDatabase() error {
	log.Printf("🗑️  Nettoyage de la base de données '%s' (toutes les données seront perdues)...", config.AppConfig.DBName)

	// DSN avec base de données
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
		config.AppConfig.DBCharset,
		config.AppConfig.DBParseTime,
		config.AppConfig.DBLoc,
	)

	// D'abord, essayer de se connecter à la base pour nettoyer les tables
	dbWithDB, err := sql.Open("mysql", dsnWithDB)
	if err == nil {
		// Tester la connexion
		if err := dbWithDB.Ping(); err == nil {
			log.Println("🧹 Nettoyage des tables et tablespaces...")
			// Désactiver les contraintes de clés étrangères
			dbWithDB.Exec("SET FOREIGN_KEY_CHECKS = 0")

			// Lister toutes les tables
			rows, err := dbWithDB.Query("SHOW TABLES")
			if err == nil {
				var tables []string
				for rows.Next() {
					var tableName string
					if err := rows.Scan(&tableName); err == nil {
						tables = append(tables, tableName)
					}
				}
				rows.Close()

				// Supprimer les tablespaces et les tables
				for _, table := range tables {
					// Essayer de supprimer le tablespace d'abord (ignore les erreurs)
					dbWithDB.Exec(fmt.Sprintf("ALTER TABLE `%s` DISCARD TABLESPACE", table))
					// Supprimer la table
					dbWithDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table))
				}
				if len(tables) > 0 {
					log.Printf("✅ %d table(s) supprimée(s)", len(tables))
				}
				// Forcer MySQL à libérer les fichiers
				dbWithDB.Exec("FLUSH TABLES")
			}
			dbWithDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
		}
		dbWithDB.Close()
	}

	// Maintenant, supprimer et recréer la base de données
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=%t&loc=%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBCharset,
		config.AppConfig.DBParseTime,
		config.AppConfig.DBLoc,
	)

	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("erreur de connexion: %w", err)
	}
	defer db.Close()

	// Essayer DROP DATABASE avec FORCE (MySQL 8.0.17+)
	dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.AppConfig.DBName)
	if _, err := db.Exec(dropQuery); err != nil {
		// Si ça échoue, essayer avec FORCE (si supporté)
		dropQueryForce := fmt.Sprintf("DROP DATABASE IF EXISTS %s FORCE", config.AppConfig.DBName)
		if _, err := db.Exec(dropQueryForce); err != nil {
			log.Printf("⚠️  Impossible de supprimer la base (les fichiers peuvent rester): %v", err)
		}
	}

	// Recréer la base de données
	log.Printf("🔄 Création de la base de données '%s'...", config.AppConfig.DBName)
	createQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.AppConfig.DBName)
	if _, err := db.Exec(createQuery); err != nil {
		return fmt.Errorf("erreur lors de la création de la base: %w", err)
	}

	log.Printf("✅ Base de données '%s' recréée avec succès", config.AppConfig.DBName)
	return nil
}
