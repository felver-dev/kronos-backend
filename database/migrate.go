package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/utils"
	"gorm.io/gorm"
)

// ResetDatabase supprime toutes les tables et recrée la base de données
func ResetDatabase() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🗑️  Suppression de toutes les tables...")

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}

	// Désactiver les contraintes de clés étrangères
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer func() {
		_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}()

	// Récupérer toutes les tables
	rows, err := sqlDB.Query(`
		SELECT TABLE_NAME 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_TYPE = 'BASE TABLE'
	`)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}
		tables = append(tables, tableName)
	}

	// Supprimer toutes les tables
	for _, table := range tables {
		log.Printf("   🗑️  Suppression de la table: %s", table)
		_, _ = sqlDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table))
	}

	log.Println("✅ Toutes les tables supprimées")

	// Recréer toutes les tables
	return AutoMigrate()
}

// AutoMigrate exécute les migrations automatiques pour créer les tables
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🔄 Démarrage des migrations automatiques...")

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}

	// Désactiver complètement les contraintes de clés étrangères pendant toute la migration
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer func() {
		_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}()

	// Étape 1: Créer toutes les tables SANS contraintes de clés étrangères
	log.Println("📋 Étape 1: Création de toutes les tables (sans contraintes FK)...")

	// Toutes les tables dans l'ordre logique
	err = DB.AutoMigrate(
		// Tables de base (sans dépendances)
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.Filiale{},         // Nouvelle table : filiales
		&models.Software{},        // Nouvelle table : software
		&models.FilialeSoftware{}, // Nouvelle table : filiale_software
		&models.Office{},
		&models.Department{},

		// Table users (sans contraintes auto-référentielles)
		&models.User{},

		// Tables de tickets
		&models.TicketCategory{},
		&models.Ticket{},
		&models.TicketAttachment{},
		&models.TicketComment{},
		&models.TicketHistory{},
		&models.TicketTag{},
		&models.TicketTagAssignment{},
		&models.TicketAssignee{},
		&models.TicketSolution{},
		&models.TicketInternal{},

		// Tables de sessions
		&models.UserSession{},

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
		&models.AssetSoftware{},
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
		&models.ProjectPhase{},
		&models.ProjectFunction{},
		&models.ProjectMember{},
		&models.ProjectMemberFunction{},
		&models.ProjectPhaseMember{},
		&models.ProjectTask{},
		&models.ProjectTaskAssignee{},
		&models.ProjectTaskComment{},
		&models.ProjectTaskAttachment{},
		&models.ProjectTaskHistory{},
		&models.ProjectBudgetExtension{},

		// Tables de paramétrage
		&models.Setting{},
		&models.RequestSource{},

		// Tables d'audit et sauvegarde
		&models.AuditLog{},
		&models.BackupConfiguration{},
		&models.Backup{},
	)

	if err != nil {
		return fmt.Errorf("échec de la création des tables: %w", err)
	}
	log.Println("✅ Toutes les tables créées")

	// Étape 2: Supprimer toutes les contraintes incorrectes créées par GORM
	log.Println("🔧 Étape 2: Nettoyage des contraintes incorrectes...")
	if err := removeAllIncorrectForeignKeys(); err != nil {
		log.Printf("⚠️  Erreur lors du nettoyage: %v", err)
	}

	// Étape 3: Ajouter toutes les contraintes correctes manuellement
	log.Println("🔧 Étape 3: Ajout des contraintes de clés étrangères...")
	if err := addAllForeignKeys(); err != nil {
		log.Printf("⚠️  Erreur lors de l'ajout des contraintes: %v", err)
		// Ne pas bloquer, continuer
	}

	// Étape 4: Seeding des données par défaut
	log.Println("🌱 Étape 4: Seeding des données par défaut...")
	if err := seedDefaultPermissions(); err != nil {
		log.Printf("⚠️  Erreur lors du seeding des permissions: %v", err)
	}
	if err := seedDefaultUserRole(); err != nil {
		log.Printf("⚠️  Erreur lors du seeding du rôle USER: %v", err)
	}
	if err := seedUserRoleProjectPermissions(); err != nil {
		log.Printf("⚠️  Erreur lors de l'attribution des permissions projets au rôle USER: %v", err)
	}
	if err := seedDefaultAdmin(); err != nil {
		log.Printf("⚠️  Erreur lors du seeding de l'admin: %v", err)
	}
	if err := seedDefaultTicketCategories(); err != nil {
		log.Printf("⚠️  Erreur lors du seeding des catégories: %v", err)
	}

	// Générer les codes pour les tickets existants
	if err := generateTicketCodes(); err != nil {
		log.Printf("⚠️  Erreur lors de la génération des codes: %v", err)
	}

	// Migrer les requester_id
	if err := migrateRequesterIDs(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration des requester_id: %v", err)
	}

	// Modifier asset_software.asset_id pour le rendre nullable
	if err := makeAssetSoftwareAssetIDNullable(); err != nil {
		log.Printf("⚠️  Erreur lors de la modification de asset_software.asset_id: %v", err)
	}

	// project_functions.type et project_member_functions (rétrocompat)
	if err := migrateProjectFunctionTypesAndMemberFunctions(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration project_functions / project_member_functions: %v", err)
	}

	// Préremplir Chef de projet et Lead pour les projets existants
	if err := migrateEnsureDefaultDirectionFunctions(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration des fonctions direction par défaut: %v", err)
	}

	// project_tasks: contrainte unique (code) -> (project_id, code) pour permettre TAP-YYYY-NNNN par projet
	if err := migrateProjectTasksCodeUniquePerProject(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration project_tasks code unique: %v", err)
	}

	// projects: colonnes start_date et end_date si absentes (période prévue)
	if err := migrateProjectsStartEndDates(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration projects start_date/end_date: %v", err)
	}

	// project_budget_extensions: colonnes start_date et end_date (période de chaque extension)
	if err := migrateProjectBudgetExtensionsStartEndDates(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration project_budget_extensions start_date/end_date: %v", err)
	}

	// Migrations multi-filiales : ajouter les colonnes filiale_id, software_id, etc.
	if err := migrateMultiFiliales(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration multi-filiales: %v", err)
	}

	// software: contrainte unique (code) -> (code, version) pour permettre plusieurs versions du même logiciel
	if err := migrateSoftwareCodeVersionUnique(); err != nil {
		log.Printf("⚠️  Erreur lors de la migration software code+version unique: %v", err)
	}

	log.Println("✅ Migrations terminées avec succès")
	return nil
}

// removeAllIncorrectForeignKeys supprime toutes les contraintes incorrectes créées par GORM
func removeAllIncorrectForeignKeys() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}

	// Liste des contraintes incorrectes connues sur users
	badConstraints := []string{
		"fk_ticket_solutions_created_by",
		"fk_assets_created_by",
		"fk_backups_created_by",
		"fk_projects_created_by",
		"fk_service_request_types_created_by",
		"fk_sla_created_by",
		"fk_tickets_created_by",
		"fk_backup_configurations_updated_by",
		"fk_settings_updated_by",
	}

	removedCount := 0
	for _, constraintName := range badConstraints {
		var exists int
		err = sqlDB.QueryRow(`
			SELECT COUNT(*) 
			FROM information_schema.TABLE_CONSTRAINTS 
			WHERE CONSTRAINT_SCHEMA = DATABASE() 
			AND TABLE_NAME = 'users' 
			AND CONSTRAINT_NAME = ?
		`, constraintName).Scan(&exists)

		if err == nil && exists > 0 {
			log.Printf("   🗑️  Suppression de la contrainte incorrecte: %s", constraintName)
			_, _ = sqlDB.Exec(fmt.Sprintf("ALTER TABLE `users` DROP FOREIGN KEY `%s`", constraintName))
			removedCount++
		}
	}

	// Supprimer toutes les contraintes sur users.created_by_id et users.updated_by_id qui ne référencent pas users.id
	for _, columnName := range []string{"created_by_id", "updated_by_id"} {
		rows, err := sqlDB.Query(`
			SELECT CONSTRAINT_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
			FROM information_schema.KEY_COLUMN_USAGE
			WHERE CONSTRAINT_SCHEMA = DATABASE()
			AND TABLE_NAME = 'users'
			AND COLUMN_NAME = ?
			AND REFERENCED_TABLE_NAME IS NOT NULL
		`, columnName)
		if err != nil {
			continue
		}

		for rows.Next() {
			var constraintName, referencedTable, referencedColumn string
			if err := rows.Scan(&constraintName, &referencedTable, &referencedColumn); err != nil {
				continue
			}

			if referencedTable != "users" || referencedColumn != "id" {
				log.Printf("   🗑️  Suppression de la contrainte incorrecte: %s (référence %s.%s)", constraintName, referencedTable, referencedColumn)
				_, _ = sqlDB.Exec(fmt.Sprintf("ALTER TABLE `users` DROP FOREIGN KEY `%s`", constraintName))
				removedCount++
			}
		}
		rows.Close()
	}

	if removedCount > 0 {
		log.Printf("   ✅ %d contrainte(s) incorrecte(s) supprimée(s)", removedCount)
	}

	return nil
}

// addAllForeignKeys ajoute toutes les contraintes de clés étrangères correctes
func addAllForeignKeys() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}

	// Fonction helper pour ajouter une contrainte si elle n'existe pas
	addFK := func(table, constraint, column, refTable, refColumn string) error {
		var exists int
		err := sqlDB.QueryRow(`
			SELECT COUNT(*) 
			FROM information_schema.TABLE_CONSTRAINTS 
			WHERE CONSTRAINT_SCHEMA = DATABASE() 
			AND TABLE_NAME = ? 
			AND CONSTRAINT_NAME = ?
		`, table, constraint).Scan(&exists)

		if err != nil {
			return err
		}

		if exists > 0 {
			return nil // Déjà existe
		}

		_, err = sqlDB.Exec(fmt.Sprintf(`
			ALTER TABLE %s 
			ADD CONSTRAINT %s 
			FOREIGN KEY (%s) REFERENCES %s(%s)
			ON DELETE RESTRICT ON UPDATE CASCADE
		`, table, constraint, column, refTable, refColumn))
		return err
	}

	// Contraintes users
	_ = addFK("users", "fk_users_role", "role_id", "roles", "id")
	_ = addFK("users", "fk_users_department", "department_id", "departments", "id")
	_ = addFK("users", "fk_users_created_by", "created_by_id", "users", "id")
	_ = addFK("users", "fk_users_updated_by", "updated_by_id", "users", "id")

	// Contraintes tickets
	_ = addFK("tickets", "fk_tickets_created_by", "created_by_id", "users", "id")
	_ = addFK("tickets", "fk_tickets_assigned_to", "assigned_to_id", "users", "id")
	_ = addFK("tickets", "fk_tickets_requester", "requester_id", "users", "id")
	_ = addFK("tickets", "fk_tickets_category", "category_id", "ticket_categories", "id")
	_ = addFK("tickets", "fk_tickets_primary_image", "primary_image_id", "ticket_attachments", "id")
	_ = addFK("tickets", "fk_tickets_parent", "parent_id", "tickets", "id")

	// Contraintes ticket_attachments
	_ = addFK("ticket_attachments", "fk_ticket_attachments_ticket", "ticket_id", "tickets", "id")
	_ = addFK("ticket_attachments", "fk_ticket_attachments_user", "user_id", "users", "id")

	// Contraintes ticket_solutions
	_ = addFK("ticket_solutions", "fk_ticket_solutions_ticket", "ticket_id", "tickets", "id")
	_ = addFK("ticket_solutions", "fk_ticket_solutions_created_by", "created_by_id", "users", "id")

	// Contraintes ticket_comments
	_ = addFK("ticket_comments", "fk_ticket_comments_ticket", "ticket_id", "tickets", "id")
	_ = addFK("ticket_comments", "fk_ticket_comments_user", "user_id", "users", "id")

	// Contraintes ticket_history
	_ = addFK("ticket_history", "fk_ticket_history_ticket", "ticket_id", "tickets", "id")
	_ = addFK("ticket_history", "fk_ticket_history_user", "user_id", "users", "id")

	// Contraintes ticket_assignees
	_ = addFK("ticket_assignees", "fk_ticket_assignees_ticket", "ticket_id", "tickets", "id")
	_ = addFK("ticket_assignees", "fk_ticket_assignees_user", "user_id", "users", "id")

	// Contraintes departments
	_ = addFK("departments", "fk_departments_office", "office_id", "offices", "id")

	// Contraintes role_permissions
	_ = addFK("role_permissions", "fk_role_permissions_role", "role_id", "roles", "id")
	_ = addFK("role_permissions", "fk_role_permissions_permission", "permission_id", "permissions", "id")

	// Contraintes roles (créateur et filiale)
	_ = addFK("roles", "fk_roles_created_by", "created_by_id", "users", "id")
	_ = addFK("roles", "fk_roles_filiale", "filiale_id", "filiales", "id")

	// Contraintes user_sessions
	_ = addFK("user_sessions", "fk_user_sessions_user", "user_id", "users", "id")

	// Contraintes projects (chef de projet, lead)
	_ = addFK("projects", "fk_projects_project_manager", "project_manager_id", "users", "id")
	_ = addFK("projects", "fk_projects_lead", "lead_id", "users", "id")

	// Contraintes project_phases
	_ = addFK("project_phases", "fk_project_phases_project", "project_id", "projects", "id")

	// Contraintes project_functions
	_ = addFK("project_functions", "fk_project_functions_project", "project_id", "projects", "id")

	// Contraintes project_members
	_ = addFK("project_members", "fk_project_members_project", "project_id", "projects", "id")
	_ = addFK("project_members", "fk_project_members_user", "user_id", "users", "id")
	_ = addFK("project_members", "fk_project_members_function", "project_function_id", "project_functions", "id")

	// Contraintes project_member_functions
	_ = addFK("project_member_functions", "fk_pmf_member", "project_member_id", "project_members", "id")
	_ = addFK("project_member_functions", "fk_pmf_function", "project_function_id", "project_functions", "id")

	// Contraintes project_phase_members
	_ = addFK("project_phase_members", "fk_project_phase_members_phase", "project_phase_id", "project_phases", "id")
	_ = addFK("project_phase_members", "fk_project_phase_members_user", "user_id", "users", "id")
	_ = addFK("project_phase_members", "fk_project_phase_members_function", "project_function_id", "project_functions", "id")

	// Contraintes project_tasks
	_ = addFK("project_tasks", "fk_project_tasks_project", "project_id", "projects", "id")
	_ = addFK("project_tasks", "fk_project_tasks_phase", "project_phase_id", "project_phases", "id")
	_ = addFK("project_tasks", "fk_project_tasks_assigned_to", "assigned_to_id", "users", "id")
	_ = addFK("project_tasks", "fk_project_tasks_created_by", "created_by_id", "users", "id")

	// Contraintes project_task_assignees
	_ = addFK("project_task_assignees", "fk_project_task_assignees_task", "project_task_id", "project_tasks", "id")
	_ = addFK("project_task_assignees", "fk_project_task_assignees_user", "user_id", "users", "id")

	// Contraintes project_task_comments
	_ = addFK("project_task_comments", "fk_project_task_comments_task", "project_task_id", "project_tasks", "id")
	_ = addFK("project_task_comments", "fk_project_task_comments_user", "user_id", "users", "id")

	// Contraintes project_task_attachments
	_ = addFK("project_task_attachments", "fk_project_task_attachments_task", "project_task_id", "project_tasks", "id")
	_ = addFK("project_task_attachments", "fk_project_task_attachments_user", "user_id", "users", "id")

	// Contraintes project_task_history
	_ = addFK("project_task_history", "fk_project_task_history_task", "project_task_id", "project_tasks", "id")
	_ = addFK("project_task_history", "fk_project_task_history_user", "user_id", "users", "id")

	// Contraintes project_budget_extensions
	_ = addFK("project_budget_extensions", "fk_project_budget_extensions_project", "project_id", "projects", "id")
	_ = addFK("project_budget_extensions", "fk_project_budget_extensions_created_by", "created_by_id", "users", "id")

	// Contraintes time_entries (project_task_id nullable)
	_ = addFK("time_entries", "fk_time_entries_project_task", "project_task_id", "project_tasks", "id")

	// Contraintes multi-filiales : filiales
	_ = addFK("filiale_software", "fk_filiale_software_filiale", "filiale_id", "filiales", "id")
	_ = addFK("filiale_software", "fk_filiale_software_software", "software_id", "software", "id")

	// Contraintes multi-filiales : users
	_ = addFK("users", "fk_users_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : departments
	_ = addFK("departments", "fk_departments_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : tickets
	_ = addFK("tickets", "fk_tickets_filiale", "filiale_id", "filiales", "id")
	_ = addFK("tickets", "fk_tickets_software", "software_id", "software", "id")
	_ = addFK("tickets", "fk_tickets_validated_by", "validated_by_user_id", "users", "id")

	// Contraintes multi-filiales : projects
	_ = addFK("projects", "fk_projects_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : knowledge
	_ = addFK("knowledge_articles", "fk_knowledge_articles_filiale", "filiale_id", "filiales", "id")
	_ = addFK("knowledge_categories", "fk_knowledge_categories_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : delays
	_ = addFK("delays", "fk_delays_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : declarations
	_ = addFK("daily_declarations", "fk_daily_declarations_filiale", "filiale_id", "filiales", "id")
	_ = addFK("weekly_declarations", "fk_weekly_declarations_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : assets
	_ = addFK("assets", "fk_assets_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : offices
	_ = addFK("offices", "fk_offices_filiale", "filiale_id", "filiales", "id")

	// Contraintes multi-filiales : sla
	_ = addFK("sla", "fk_sla_filiale", "filiale_id", "filiales", "id")

	log.Println("   ✅ Contraintes de clés étrangères ajoutées")
	return nil
}

// seedDefaultPermissions crée toutes les permissions disponibles dans le système
func seedDefaultPermissions() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🌱 Seeding des permissions par défaut...")

	// Liste complète des permissions du système
	permissions := []struct {
		Code        string
		Name        string
		Description string
		Module      string
	}{
		// Permissions Tickets
		{"tickets.view_all", "Voir tous les tickets", "Voir tous les tickets du système", "tickets"},
		{"tickets.view_filiale", "Voir tous les tickets de sa filiale", "Voir tous les tickets de sa filiale (DSI filiale)", "tickets"},
		{"tickets.view_team", "Voir tickets de son équipe", "Voir les tickets de son équipe/département", "tickets"},
		{"tickets.view_own", "Voir ses tickets", "Voir uniquement ses tickets assignés", "tickets"},
		{"tickets.create", "Créer un ticket", "Créer un nouveau ticket dans sa propre filiale", "tickets"},
		{"tickets.create_any_filiale", "Créer un ticket pour n'importe quelle filiale", "Créer un ticket pour n'importe quelle filiale (Département IT MCI CARE CI)", "tickets"},
		{"tickets.update", "Modifier un ticket", "Modifier un ticket", "tickets"},
		{"tickets.delete", "Supprimer un ticket", "Supprimer un ticket", "tickets"},
		{"tickets.assign", "Assigner un ticket", "Assigner un ticket à un utilisateur", "tickets"},
		{"tickets.reassign", "Réassigner un ticket", "Réassigner un ticket", "tickets"},
		{"tickets.close", "Clôturer un ticket", "Clôturer un ticket", "tickets"},
		{"tickets.resolve_all", "Résoudre tous les tickets", "Résoudre tous les tickets (IT MCI CARE CI)", "tickets"},
		{"tickets.resolve_own_filiale", "Résoudre tickets de sa filiale", "Résoudre les tickets de sa filiale uniquement", "tickets"},
		{"tickets.validate", "Valider les tickets résolus", "Valider les tickets résolus", "tickets"},
		{"tickets.validate_own", "Valider ses propres tickets", "Valider uniquement ses propres tickets créés", "tickets"},

		// Permissions Tickets internes (départements non-IT, scope département / filiale / global)
		{"tickets_internes.view_own", "Voir ses tickets internes", "Voir ses tickets internes (créés ou assignés)", "tickets_internes"},
		{"tickets_internes.view_department", "Voir tickets internes de son département", "Voir les tickets internes de son département", "tickets_internes"},
		{"tickets_internes.view_filiale", "Voir tickets internes de sa filiale", "Voir les tickets internes de sa filiale", "tickets_internes"},
		{"tickets_internes.view_all", "Voir tous les tickets internes", "Voir tous les tickets internes (vue DG/DGA/PDG)", "tickets_internes"},
		{"tickets_internes.create", "Créer un ticket interne", "Créer un ticket interne dans son périmètre", "tickets_internes"},
		{"tickets_internes.update", "Modifier un ticket interne", "Modifier un ticket interne", "tickets_internes"},
		{"tickets_internes.assign", "Assigner un ticket interne", "Assigner un ticket interne à un utilisateur", "tickets_internes"},
		{"tickets_internes.validate", "Valider un ticket interne résolu", "Valider un ticket interne (passer en résolu)", "tickets_internes"},
		{"tickets_internes.close", "Clôturer un ticket interne", "Clôturer un ticket interne", "tickets_internes"},
		{"tickets_internes.delete", "Supprimer un ticket interne", "Supprimer un ticket interne", "tickets_internes"},

		// Permissions Software
		{"software.view", "Voir les logiciels", "Voir les logiciels gérés", "software"},
		{"software.create", "Créer un logiciel", "Créer un nouveau logiciel (IT MCI CARE CI)", "software"},
		{"software.update", "Modifier un logiciel", "Modifier un logiciel (IT MCI CARE CI)", "software"},
		{"software.delete", "Supprimer un logiciel", "Supprimer un logiciel (IT MCI CARE CI)", "software"},
		{"software.deploy", "Déployer un logiciel", "Déployer un logiciel chez une filiale (IT MCI CARE CI)", "software"},
		{"software.manage_deployments", "Gérer les déploiements", "Gérer les déploiements de logiciels (IT MCI CARE CI)", "software"},

		// Permissions Filiales
		{"filiales.view", "Voir les filiales", "Voir les filiales (sa filiale uniquement sans view_all)", "filiales"},
		{"filiales.view_all", "Voir toutes les filiales", "Voir toutes les filiales du groupe", "filiales"},
		{"filiales.create", "Créer une filiale", "Créer une nouvelle filiale (Super Admin)", "filiales"},
		{"filiales.update", "Modifier une filiale", "Modifier une filiale (Super Admin)", "filiales"},
		{"filiales.manage", "Gestion complète des filiales", "Gestion complète des filiales (Super Admin)", "filiales"},

		// Permissions Notifications
		{"notifications.filter_by_filiale", "Filtrer les notifications par filiale", "Filtrer l'historique des notifications par filiale (résolveurs, développeurs)", "notifications"},

		// Permissions Timesheet
		{"timesheet.create_entry", "Saisir le temps", "Saisir le temps passé sur un ticket", "timesheet"},
		{"timesheet.view_all", "Voir toutes les déclarations", "Voir toutes les déclarations", "timesheet"},
		{"timesheet.view_team", "Voir déclarations de son équipe", "Voir les déclarations de son équipe", "timesheet"},
		{"timesheet.view_own", "Voir ses déclarations", "Voir uniquement ses propres déclarations de temps", "timesheet"},
		{"timesheet.validate", "Valider déclarations", "Valider les déclarations de temps", "timesheet"},
		{"timesheet.justify_delay", "Justifier un retard", "Justifier un retard", "timesheet"},
		{"timesheet.validate_justification", "Valider justifications", "Valider les justifications de retards", "timesheet"},
		{"timesheet.view_budget", "Voir le budget temps", "Accéder à l'onglet Budget temps (temps estimés par ticket, alertes budget)", "timesheet"},
		{"timesheet.create_daily", "Créer une déclaration journalière", "Créer ou modifier une déclaration journalière de temps", "timesheet"},
		{"timesheet.create_weekly", "Créer une déclaration hebdomadaire", "Créer ou modifier une déclaration hebdomadaire de temps", "timesheet"},

		// Permissions Users
		{"users.view_all", "Voir tous les utilisateurs", "Voir tous les utilisateurs", "users"},
		{"users.view_filiale", "Voir utilisateurs de sa filiale", "Voir les utilisateurs de sa propre filiale", "users"},
		{"users.view_team", "Voir utilisateurs de son équipe", "Voir les utilisateurs de son équipe", "users"},
		{"users.view_own", "Voir son propre profil", "Voir son propre profil", "users"},
		{"users.create", "Créer un utilisateur", "Créer un nouvel utilisateur dans sa propre filiale", "users"},
		{"users.create_any_filiale", "Créer un utilisateur dans n'importe quelle filiale", "Créer un utilisateur dans n'importe quelle filiale (admin principal)", "users"},
		{"users.update", "Modifier un utilisateur", "Modifier un utilisateur de sa propre filiale", "users"},
		{"users.update_any_filiale", "Modifier un utilisateur dans n'importe quelle filiale", "Modifier un utilisateur dans n'importe quelle filiale (admin principal)", "users"},
		{"users.delete", "Supprimer un utilisateur", "Supprimer un utilisateur", "users"},

		// Permissions Roles
		{"roles.view", "Voir les rôles", "Voir les rôles", "roles"},
		{"roles.view_filiale", "Voir les rôles de sa filiale", "Voir uniquement les rôles globaux et les rôles de sa filiale", "roles"},
		{"roles.view_department", "Voir les rôles de son département", "Voir uniquement les rôles utilisés par les utilisateurs de son département", "roles"},
		{"roles.view_assigned_only", "Voir uniquement les permissions assignées", "Voir uniquement les permissions actuellement assignées à un rôle (lecture seule)", "roles"},
		{"roles.create", "Créer un rôle", "Créer un nouveau rôle", "roles"},
		{"roles.update", "Modifier un rôle", "Modifier un rôle existant", "roles"},
		{"roles.delete", "Supprimer un rôle", "Supprimer un rôle", "roles"},
		{"roles.manage", "Gérer les rôles", "Créer, modifier, supprimer les rôles (permission globale)", "roles"},
		{"roles.delegate_permissions", "Déléguer des permissions", "Créer des rôles et leur assigner un sous-ensemble de ses propres permissions", "roles"},

		// Permissions Reports
		{"reports.view_global", "Rapports globaux groupe", "Voir les rapports globaux du groupe (IT MCI CARE CI)", "reports"},
		{"reports.view_filiale", "Rapports de sa filiale", "Voir les rapports de sa filiale", "reports"},
		{"reports.view_team", "Rapports d'équipe", "Voir les rapports de son équipe", "reports"},
		{"reports.view_own", "Rapports personnels", "Voir ses rapports personnels", "reports"},
		{"reports.view_departments", "Rapports par départements", "Voir les rapports par départements", "reports"},
		{"reports.view_employees", "Rapports par employé", "Voir les rapports par employé", "reports"},
		{"reports.compare_filiales", "Comparer entre filiales", "Comparer les rapports entre filiales (IT MCI CARE CI)", "reports"},

		// Permissions Assets
		{"assets.view_all", "Voir tous les actifs", "Voir tous les actifs IT", "assets"},
		{"assets.view_team", "Voir actifs de son équipe", "Voir les actifs de son équipe/département", "assets"},
		{"assets.view_own", "Voir ses actifs assignés", "Voir les actifs qui lui sont assignés", "assets"},
		{"assets.create", "Créer un actif", "Créer un actif IT", "assets"},
		{"assets.update", "Modifier un actif", "Modifier un actif IT", "assets"},
		{"assets.delete", "Supprimer un actif", "Supprimer un actif IT", "assets"},

		// Permissions Knowledge Base
		{"knowledge.view_all", "Voir tous les articles", "Voir tous les articles", "knowledge"},
		{"knowledge.view_published", "Voir les articles publiés", "Voir les articles publiés", "knowledge"},
		{"knowledge.view_own", "Voir ses propres articles", "Voir ses propres articles", "knowledge"},
		{"knowledge.create", "Créer un article", "Créer un article", "knowledge"},
		{"knowledge.update", "Modifier un article", "Modifier un article", "knowledge"},
		{"knowledge.delete", "Supprimer un article", "Supprimer un article", "knowledge"},
		{"knowledge.publish", "Publier un article", "Publier un article", "knowledge"},

		// Permissions Settings
		{"settings.view", "Voir les paramètres", "Voir les paramètres système", "settings"},
		{"settings.update", "Modifier les paramètres", "Modifier les paramètres système", "settings"},
		{"settings.manage", "Configuration système", "Gérer la configuration système (permission globale)", "settings"},

		// Permissions SLA
		{"sla.view", "Voir les SLA", "Voir les SLA", "sla"},
		{"sla.view_all", "Voir tous les SLA et violations", "Voir tous les SLA et violations", "sla"},
		{"sla.view_team", "Voir SLA de son équipe", "Voir les SLA/violations de son équipe", "sla"},
		{"sla.view_own", "Voir ses SLA", "Voir les SLA liés à ses tickets", "sla"},
		{"sla.create", "Créer un SLA", "Créer un SLA", "sla"},
		{"sla.update", "Modifier un SLA", "Modifier un SLA", "sla"},
		{"sla.delete", "Supprimer un SLA", "Supprimer un SLA", "sla"},
		{"sla.manage", "Gestion SLA", "Gérer les SLA (permission globale)", "sla"},

		// Permissions Audit
		{"audit.view_all", "Voir tous les logs", "Voir tous les logs d'audit", "audit"},
		{"audit.view_team", "Voir logs de son équipe", "Voir les logs de son équipe", "audit"},
		{"audit.view_own", "Voir ses propres logs", "Voir ses propres actions enregistrées", "audit"},

		// Permissions Offices (Sièges)
		{"offices.view", "Voir les sièges", "Voir les sièges (équivalent à view_filiale pour rétrocompat)", "offices"},
		{"offices.view_filiale", "Voir sièges de sa filiale", "Voir uniquement les sièges de sa propre filiale", "offices"},
		{"offices.view_all", "Voir tous les sièges", "Voir les sièges de toutes les filiales du système", "offices"},
		{"offices.create", "Créer un siège", "Créer un nouveau siège dans sa propre filiale", "offices"},
		{"offices.create_any_filiale", "Créer un siège dans n'importe quelle filiale", "Créer un siège dans n'importe quelle filiale (admin principal)", "offices"},
		{"offices.update", "Modifier un siège", "Modifier un siège de sa propre filiale", "offices"},
		{"offices.update_any_filiale", "Modifier un siège dans n'importe quelle filiale", "Modifier un siège dans n'importe quelle filiale (admin principal)", "offices"},
		{"offices.delete", "Supprimer un siège", "Supprimer un siège", "offices"},

		// Permissions Departments (Départements)
		{"departments.view", "Voir les départements", "Voir les départements (équivalent à view_filiale pour rétrocompat)", "departments"},
		{"departments.view_filiale", "Voir départements de sa filiale", "Voir uniquement les départements de sa propre filiale", "departments"},
		{"departments.view_all", "Voir tous les départements", "Voir les départements de toutes les filiales du système", "departments"},
		{"departments.create", "Créer un département", "Créer un nouveau département dans sa propre filiale", "departments"},
		{"departments.create_any_filiale", "Créer un département dans n'importe quelle filiale", "Créer un département dans n'importe quelle filiale (admin principal)", "departments"},
		{"departments.update", "Modifier un département", "Modifier un département de sa propre filiale", "departments"},
		{"departments.update_any_filiale", "Modifier un département dans n'importe quelle filiale", "Modifier un département dans n'importe quelle filiale (admin principal)", "departments"},
		{"departments.delete", "Supprimer un département", "Supprimer un département", "departments"},

		// Permissions Incidents
		{"incidents.view", "Voir les incidents", "Voir les incidents", "incidents"},
		{"incidents.view_all", "Voir tous les incidents", "Voir tous les incidents du système", "incidents"},
		{"incidents.view_team", "Voir incidents de son équipe", "Voir les incidents de son équipe/département", "incidents"},
		{"incidents.view_own", "Voir ses incidents", "Voir les incidents liés à ses tickets", "incidents"},
		{"incidents.create", "Créer un incident", "Créer un nouvel incident", "incidents"},
		{"incidents.update", "Modifier un incident", "Modifier un incident existant", "incidents"},
		{"incidents.delete", "Supprimer un incident", "Supprimer un incident", "incidents"},

		// Permissions Service Requests (Demandes de service)
		{"service_requests.view", "Voir les demandes de service", "Voir les demandes de service", "service_requests"},
		{"service_requests.view_all", "Voir toutes les demandes de service", "Voir toutes les demandes de service du système", "service_requests"},
		{"service_requests.view_team", "Voir demandes de son équipe", "Voir les demandes de service de son équipe/département", "service_requests"},
		{"service_requests.view_own", "Voir ses demandes de service", "Voir les demandes liées à ses tickets", "service_requests"},
		{"service_requests.create", "Créer une demande de service", "Créer une nouvelle demande de service", "service_requests"},
		{"service_requests.update", "Modifier une demande de service", "Modifier une demande de service existante", "service_requests"},
		{"service_requests.delete", "Supprimer une demande de service", "Supprimer une demande de service", "service_requests"},

		// Permissions Changes (Changements)
		{"changes.view", "Voir les changements", "Voir les changements", "changes"},
		{"changes.view_all", "Voir tous les changements", "Voir tous les changements du système", "changes"},
		{"changes.view_team", "Voir changements de son équipe", "Voir les changements de son équipe/département", "changes"},
		{"changes.view_own", "Voir ses changements", "Voir les changements liés à ses tickets", "changes"},
		{"changes.create", "Créer un changement", "Créer un nouveau changement", "changes"},
		{"changes.update", "Modifier un changement", "Modifier un changement existant", "changes"},
		{"changes.delete", "Supprimer un changement", "Supprimer un changement", "changes"},

		// Permissions Delays (Retards)
		{"delays.view", "Voir les retards", "Voir les retards", "delays"},
		{"delays.view_all", "Voir tous les retards", "Voir tous les retards du système", "delays"},
		{"delays.view_department", "Voir retards de son département", "Voir les retards de son département", "delays"},
		{"delays.view_own", "Voir ses propres retards", "Voir ses propres retards", "delays"},
		{"delays.validate", "Valider les retards", "Valider ou rejeter les justifications de retards", "delays"},

		// Permissions Projects (Projets) — entité principale
		{"projects.view", "Voir les projets", "Voir la liste des projets (selon scope)", "projects"},
		{"projects.view_all", "Voir tous les projets", "Voir tous les projets du système", "projects"},
		{"projects.view_team", "Voir projets de son équipe", "Voir les projets dont un membre est du même département", "projects"},
		{"projects.view_own", "Voir ses projets", "Voir les projets où l'utilisateur est membre ou membre d'une étape", "projects"},
		{"projects.create", "Créer un projet", "Créer un nouveau projet", "projects"},
		{"projects.update", "Modifier un projet", "Modifier les infos d'un projet (nom, description, dates, statut, budget)", "projects"},
		{"projects.delete", "Supprimer un projet", "Supprimer un projet", "projects"},
		{"projects.set_project_manager", "Désigner le chef de projet", "Désigner ou changer le chef de projet", "projects"},
		{"projects.set_lead", "Désigner le lead", "Désigner ou changer le lead technique ou fonctionnel", "projects"},

		// Permissions Projects — étapes (phases)
		{"projects.phases.view", "Voir les étapes", "Voir les étapes d'un projet", "projects"},
		{"projects.phases.create", "Créer une étape", "Créer une étape", "projects"},
		{"projects.phases.update", "Modifier une étape", "Modifier une étape (nom, ordre, dates, statut)", "projects"},
		{"projects.phases.delete", "Supprimer une étape", "Supprimer une étape", "projects"},
		{"projects.phases.reorder", "Réordonner les étapes", "Changer l'ordre des étapes", "projects"},

		// Permissions Projects — fonctions (au sens fonction projet)
		{"projects.functions.view", "Voir les fonctions", "Voir les fonctions d'un projet ou le catalogue global", "projects"},
		{"projects.functions.create", "Créer une fonction", "Créer une fonction (pour un projet ou globale)", "projects"},
		{"projects.functions.update", "Modifier une fonction", "Modifier une fonction", "projects"},
		{"projects.functions.delete", "Supprimer une fonction", "Supprimer une fonction", "projects"},

		// Permissions Projects — membres du projet
		{"projects.members.view", "Voir les membres", "Voir la liste des membres du projet", "projects"},
		{"projects.members.add", "Ajouter un membre", "Ajouter un membre au projet", "projects"},
		{"projects.members.remove", "Retirer un membre", "Retirer un membre du projet", "projects"},
		{"projects.members.assign_function", "Affecter une fonction", "Affecter ou modifier la fonction d'un membre", "projects"},
		{"projects.members.set_project_manager", "Désigner chef de projet (membre)", "Désigner un membre comme chef de projet", "projects"},
		{"projects.members.set_lead", "Désigner lead (membre)", "Désigner un membre comme lead", "projects"},

		// Permissions Projects — membres par étape
		{"projects.phase_members.view", "Voir les membres d'étape", "Voir les membres d'une étape", "projects"},
		{"projects.phase_members.add", "Ajouter un membre à une étape", "Ajouter un membre à une étape", "projects"},
		{"projects.phase_members.remove", "Retirer un membre d'étape", "Retirer un membre d'une étape", "projects"},
		{"projects.phase_members.assign_function", "Affecter fonction (membre d'étape)", "Affecter ou modifier la fonction d'un membre d'étape", "projects"},

		// Permissions Projects — tâches (project_tasks)
		{"projects.tasks.view", "Voir les tâches", "Voir les tâches (selon scope)", "projects"},
		{"projects.tasks.view_project", "Voir toutes les tâches du projet", "Voir toutes les tâches du projet", "projects"},
		{"projects.tasks.view_phase", "Voir les tâches de ses étapes", "Voir les tâches des étapes où l'utilisateur est membre", "projects"},
		{"projects.tasks.view_own", "Voir ses tâches", "Voir uniquement les tâches assignées à l'utilisateur", "projects"},
		{"projects.tasks.create", "Créer une tâche", "Créer une tâche dans une étape du projet", "projects"},
		{"projects.tasks.update", "Modifier une tâche", "Modifier une tâche (titre, description, statut, priorité, etc.)", "projects"},
		{"projects.tasks.delete", "Supprimer une tâche", "Supprimer une tâche", "projects"},
		{"projects.tasks.assign", "Assigner une tâche", "Assigner ou réassigner une tâche", "projects"},
		{"projects.tasks.close", "Clôturer une tâche", "Clôturer une tâche", "projects"},

		// Permissions Projects — commentaires des tâches
		{"projects.tasks.comments.view", "Voir les commentaires (tâches)", "Voir les commentaires d'une tâche", "projects"},
		{"projects.tasks.comments.create", "Créer un commentaire (tâche)", "Ajouter un commentaire à une tâche", "projects"},
		{"projects.tasks.comments.update", "Modifier un commentaire (tâche)", "Modifier son propre commentaire", "projects"},
		{"projects.tasks.comments.delete", "Supprimer un commentaire (tâche)", "Supprimer un commentaire", "projects"},

		// Permissions Projects — pièces jointes des tâches
		{"projects.tasks.attachments.view", "Voir les pièces jointes (tâches)", "Voir les pièces jointes d'une tâche", "projects"},
		{"projects.tasks.attachments.create", "Ajouter une pièce jointe (tâche)", "Ajouter une pièce jointe à une tâche", "projects"},
		{"projects.tasks.attachments.delete", "Supprimer une pièce jointe (tâche)", "Supprimer une pièce jointe d'une tâche", "projects"},

		// Permissions Projects — saisie de temps sur les tâches
		{"projects.tasks.time.view", "Voir le temps (tâches)", "Voir les saisies de temps des tâches du projet", "projects"},
		{"projects.tasks.time.create", "Saisir le temps (tâche)", "Saisir du temps sur une tâche", "projects"},
		{"projects.tasks.time.update", "Modifier une saisie (tâche)", "Modifier une saisie de temps sur une tâche", "projects"},
		{"projects.tasks.time.delete", "Supprimer une saisie (tâche)", "Supprimer une saisie de temps sur une tâche", "projects"},

		// Permissions Projects — budget et tableau de bord
		{"projects.budget.view", "Voir le budget projet", "Voir le budget temps du projet", "projects"},
		{"projects.budget.manage", "Gérer le budget projet", "Modifier le budget temps du projet", "projects"},
		{"projects.budget.extensions.update", "Modifier une extension de budget", "Modifier une extension de budget", "projects"},
		{"projects.budget.extensions.delete", "Supprimer une extension de budget", "Supprimer une extension de budget", "projects"},
		{"projects.dashboard.view", "Voir le tableau de bord projet", "Voir le tableau de bord (avancement, statistiques)", "projects"},

		// Permissions Asset Categories (Catégories d'actifs)
		{"asset_categories.view", "Voir les catégories d'actifs", "Voir les catégories d'actifs", "asset_categories"},
		{"asset_categories.create", "Créer une catégorie d'actif", "Créer une nouvelle catégorie d'actif", "asset_categories"},
		{"asset_categories.update", "Modifier une catégorie d'actif", "Modifier une catégorie d'actif existante", "asset_categories"},
		{"asset_categories.delete", "Supprimer une catégorie d'actif", "Supprimer une catégorie d'actif", "asset_categories"},

		// Permissions Knowledge Categories (Catégories de connaissances)
		{"knowledge_categories.view", "Voir les catégories de connaissances", "Voir les catégories de connaissances", "knowledge_categories"},
		{"knowledge_categories.create", "Créer une catégorie de connaissances", "Créer une nouvelle catégorie de connaissances", "knowledge_categories"},
		{"knowledge_categories.update", "Modifier une catégorie de connaissances", "Modifier une catégorie de connaissances existante", "knowledge_categories"},
		{"knowledge_categories.delete", "Supprimer une catégorie de connaissances", "Supprimer une catégorie de connaissances", "knowledge_categories"},

		// Permissions Ticket Categories (Catégories de tickets)
		{"ticket_categories.view", "Voir les catégories de tickets", "Voir les catégories de tickets", "ticket_categories"},
		{"ticket_categories.create", "Créer une catégorie de ticket", "Créer une nouvelle catégorie de ticket", "ticket_categories"},
		{"ticket_categories.update", "Modifier une catégorie de ticket", "Modifier une catégorie de ticket existante", "ticket_categories"},
		{"ticket_categories.delete", "Supprimer une catégorie de ticket", "Supprimer une catégorie de ticket", "ticket_categories"},
	}

	for _, perm := range permissions {
		var existing models.Permission
		result := DB.Where("code = ?", perm.Code).First(&existing)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newPerm := models.Permission{
				Code:        perm.Code,
				Name:        perm.Name,
				Description: perm.Description,
				Module:      perm.Module,
			}
			if err := DB.Create(&newPerm).Error; err != nil {
				log.Printf("   ⚠️  Erreur lors de la création de la permission %s: %v", perm.Code, err)
			}
		}
	}

	log.Println("   ✅ Permissions créées")
	return nil
}

// seedDefaultUserRole crée le rôle USER par défaut
func seedDefaultUserRole() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🌱 Seeding du rôle USER...")

	var existingRole models.Role
	result := DB.Where("name = ?", "USER").First(&existingRole)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		userRole := models.Role{
			Name:        "USER",
			Description: "Utilisateur standard",
			IsSystem:    true,
		}
		if err := DB.Create(&userRole).Error; err != nil {
			return fmt.Errorf("erreur lors de la création du rôle USER: %w", err)
		}
		log.Println("   ✅ Rôle USER créé")
	}

	return nil
}

// seedUserRoleProjectPermissions attribue au rôle USER les permissions nécessaires pour voir ses projets sur le tableau de bord (non-IT et IT)
func seedUserRoleProjectPermissions() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	var userRole models.Role
	if err := DB.Where("name = ?", "USER").First(&userRole).Error; err != nil {
		return nil // pas de rôle USER, rien à faire
	}

	codes := []string{"projects.view_own", "projects.tasks.view_own"}
	for _, code := range codes {
		var perm models.Permission
		if err := DB.Where("code = ?", code).First(&perm).Error; err != nil {
			continue
		}
		var exists int64
		DB.Model(&models.RolePermission{}).Where("role_id = ? AND permission_id = ?", userRole.ID, perm.ID).Count(&exists)
		if exists > 0 {
			continue
		}
		if err := DB.Create(&models.RolePermission{RoleID: userRole.ID, PermissionID: perm.ID, CreatedAt: time.Now()}).Error; err != nil {
			log.Printf("   ⚠️  Attribution %s au rôle USER: %v", code, err)
		}
	}
	log.Println("   ✅ Permissions projets (view_own, tasks.view_own) attribuées au rôle USER si besoin")
	return nil
}

// seedDefaultAdmin crée l'utilisateur admin par défaut
func seedDefaultAdmin() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🌱 Seeding de l'utilisateur admin...")

	// Vérifier si l'admin existe déjà
	var existingUser models.User
	result := DB.Where("username = ? OR email = ?", "admin", "admin@kronos.com").First(&existingUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Récupérer le rôle ADMIN
		var adminRole models.Role
		if err := DB.Where("name = ?", "ADMIN").First(&adminRole).Error; err != nil {
			log.Println("   ⚠️  Rôle ADMIN non trouvé, création...")
			adminRole = models.Role{
				Name:        "ADMIN",
				Description: "Administrateur système",
				IsSystem:    true,
			}
			if err := DB.Create(&adminRole).Error; err != nil {
				return fmt.Errorf("erreur lors de la création du rôle ADMIN: %w", err)
			}
		}

		// Hasher le mot de passe
		hashedPassword, err := utils.HashPassword("kronos12345")
		if err != nil {
			return fmt.Errorf("erreur lors du hashage du mot de passe: %w", err)
		}

		adminUser := models.User{
			Username:     "admin",
			Email:        "admin@kronos.com",
			PasswordHash: hashedPassword,
			FirstName:    "Admin",
			LastName:     "System",
			RoleID:       adminRole.ID,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := DB.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("erreur lors de la création de l'admin: %w", err)
		}
		log.Println("   ✅ Utilisateur admin créé (username: admin, email: admin@kronos.com, password: kronos12345)")
	} else {
		log.Println("   ℹ️  Utilisateur admin déjà existant")
	}

	// Attribuer toutes les permissions au rôle ADMIN
	if err := assignAllPermissionsToAdmin(); err != nil {
		log.Printf("   ⚠️  Erreur lors de l'attribution des permissions au rôle ADMIN: %v", err)
	}

	return nil
}

// assignAllPermissionsToAdmin attribue toutes les permissions au rôle ADMIN
func assignAllPermissionsToAdmin() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	// Récupérer le rôle ADMIN
	var adminRole models.Role
	if err := DB.Where("name = ?", "ADMIN").First(&adminRole).Error; err != nil {
		return fmt.Errorf("rôle ADMIN non trouvé: %w", err)
	}

	// Récupérer toutes les permissions
	var allPermissions []models.Permission
	if err := DB.Find(&allPermissions).Error; err != nil {
		return fmt.Errorf("erreur lors de la récupération des permissions: %w", err)
	}

	if len(allPermissions) == 0 {
		log.Println("   ℹ️  Aucune permission trouvée, attribution ignorée")
		return nil
	}

	// Vérifier quelles permissions sont déjà attribuées
	var existingRolePermissions []models.RolePermission
	if err := DB.Where("role_id = ?", adminRole.ID).Find(&existingRolePermissions).Error; err != nil {
		return fmt.Errorf("erreur lors de la vérification des permissions existantes: %w", err)
	}

	// Créer un map des permissions déjà attribuées
	existingPermIDs := make(map[uint]bool)
	for _, rp := range existingRolePermissions {
		existingPermIDs[rp.PermissionID] = true
	}

	// Ajouter les permissions manquantes
	newRolePermissions := []models.RolePermission{}
	for _, perm := range allPermissions {
		if !existingPermIDs[perm.ID] {
			newRolePermissions = append(newRolePermissions, models.RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: perm.ID,
				CreatedAt:    time.Now(),
			})
		}
	}

	if len(newRolePermissions) > 0 {
		if err := DB.Create(&newRolePermissions).Error; err != nil {
			return fmt.Errorf("erreur lors de l'attribution des permissions: %w", err)
		}
		log.Printf("   ✅ %d permission(s) attribuée(s) au rôle ADMIN", len(newRolePermissions))
	} else {
		log.Println("   ℹ️  Toutes les permissions sont déjà attribuées au rôle ADMIN")
	}

	return nil
}

// seedDefaultTicketCategories crée les catégories de tickets par défaut
func seedDefaultTicketCategories() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	log.Println("🌱 Seeding des catégories de tickets...")

	categories := []struct {
		Name        string
		Slug        string
		Description string
		Icon        string
		Color       string
	}{
		{"Incident", "incident", "Problème technique nécessitant une résolution", "alert-circle", "red"},
		{"Demande", "demande", "Demande de service ou d'assistance", "help-circle", "blue"},
		{"Changement", "changement", "Demande de modification ou d'évolution", "refresh-cw", "orange"},
	}

	for _, cat := range categories {
		var existing models.TicketCategory
		// Vérifier par slug OU par name (contrainte unique sur name : idx_ticket_categories_name)
		// Unscoped() pour inclure les lignes soft-deleted : l'index unique sur name s'applique à toutes les lignes
		result := DB.Unscoped().Where("slug = ? OR name = ?", cat.Slug, cat.Name).First(&existing)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newCat := models.TicketCategory{
				Name:        cat.Name,
				Slug:        cat.Slug,
				Description: cat.Description,
				Icon:        cat.Icon,
				Color:       cat.Color,
				IsActive:    true,
			}
			if err := DB.Create(&newCat).Error; err != nil {
				log.Printf("   ⚠️  Erreur lors de la création de la catégorie %s: %v", cat.Slug, err)
			}
		}
	}

	log.Println("   ✅ Catégories créées")
	return nil
}

// generateTicketCodes génère les codes pour les tickets existants qui n'en ont pas
func generateTicketCodes() error {
	// Fonction simplifiée - peut être complétée plus tard
	return nil
}

// migrateRequesterIDs migre les requester_id pour les tickets existants
func migrateRequesterIDs() error {
	// Fonction simplifiée - peut être complétée plus tard
	return nil
}

// migrateProjectFunctionTypesAndMemberFunctions : 1) colonne function_type (évite le mot réservé "type") ;
// 2) défaut function_type='execution' ; 3) copie project_members.project_function_id vers project_member_functions.
func migrateProjectFunctionTypesAndMemberFunctions() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("récupération sqlDB: %w", err)
	}
	// 1) S'assurer que la colonne function_type existe (évite le mot réservé MySQL "type")
	var hasFunctionType, hasType int
	_ = sqlDB.QueryRow(`
		SELECT COUNT(*) FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'project_functions' AND COLUMN_NAME = 'function_type'
	`).Scan(&hasFunctionType)
	_ = sqlDB.QueryRow(`
		SELECT COUNT(*) FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'project_functions' AND COLUMN_NAME = 'type'
	`).Scan(&hasType)
	if hasFunctionType == 0 {
		if hasType != 0 {
			// Renommer type -> function_type
			if _, err := sqlDB.Exec("ALTER TABLE project_functions CHANGE COLUMN `type` function_type VARCHAR(20) NOT NULL DEFAULT 'execution'"); err != nil {
				log.Printf("   ℹ️  project_functions type->function_type: %v", err)
			}
		} else {
			if _, err := sqlDB.Exec("ALTER TABLE project_functions ADD COLUMN function_type VARCHAR(20) NOT NULL DEFAULT 'execution' AFTER name"); err != nil {
				log.Printf("   ℹ️  project_functions ADD function_type: %v", err)
			}
		}
	}
	// 2) Mettre function_type='execution' pour les lignes sans valeur
	if err := DB.Exec("UPDATE project_functions SET function_type = ? WHERE function_type IS NULL OR function_type = ''", "execution").Error; err != nil {
		log.Printf("   ℹ️  project_functions.function_type UPDATE: %v", err)
	}
	// 3) Copier project_function_id vers project_member_functions
	var members []models.ProjectMember
	if err := DB.Where("project_function_id IS NOT NULL").Find(&members).Error; err != nil {
		return err
	}
	for _, m := range members {
		if m.ProjectFunctionID == nil {
			continue
		}
		var n int64
		DB.Model(&models.ProjectMemberFunction{}).Where("project_member_id = ? AND project_function_id = ?", m.ID, *m.ProjectFunctionID).Count(&n)
		if n == 0 {
			_ = DB.Create(&models.ProjectMemberFunction{ProjectMemberID: m.ID, ProjectFunctionID: *m.ProjectFunctionID}).Error
		}
	}
	return nil
}

// migrateEnsureDefaultDirectionFunctions ajoute « Chef de projet » et « Lead » (direction) pour chaque projet qui ne les a pas.
func migrateEnsureDefaultDirectionFunctions() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}
	var projectIDs []uint
	if err := DB.Model(&models.Project{}).Pluck("id", &projectIDs).Error; err != nil {
		return err
	}
	for _, pid := range projectIDs {
		// Chef de projet
		var n int64
		DB.Model(&models.ProjectFunction{}).Where("project_id = ? AND name = ?", pid, "Chef de projet").Count(&n)
		if n == 0 {
			if err := DB.Create(&models.ProjectFunction{ProjectID: &pid, Name: "Chef de projet", Type: "direction", DisplayOrder: 0}).Error; err != nil {
				log.Printf("   migrateEnsureDefaultDirectionFunctions project %d Chef de projet: %v", pid, err)
			}
		}
		// Lead
		DB.Model(&models.ProjectFunction{}).Where("project_id = ? AND name = ?", pid, "Lead").Count(&n)
		if n == 0 {
			if err := DB.Create(&models.ProjectFunction{ProjectID: &pid, Name: "Lead", Type: "direction", DisplayOrder: 1}).Error; err != nil {
				log.Printf("   migrateEnsureDefaultDirectionFunctions project %d Lead: %v", pid, err)
			}
		}
	}
	return nil
}

// migrateProjectTasksCodeUniquePerProject remplace l'index unique sur (code) par un index unique sur (project_id, code),
// afin que chaque projet puisse avoir ses propres TAP-YYYY-0001, TAP-YYYY-0002, etc.
func migrateProjectTasksCodeUniquePerProject() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}
	// Vérifier si l'ancien index unique sur (code) existe
	var n int
	if err := sqlDB.QueryRow(`
		SELECT COUNT(*) FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'project_tasks' AND INDEX_NAME = 'idx_project_tasks_code'
	`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		log.Println("   🔧 project_tasks: suppression de l'index unique idx_project_tasks_code (code global)...")
		if _, err := sqlDB.Exec("ALTER TABLE project_tasks DROP INDEX idx_project_tasks_code"); err != nil {
			return fmt.Errorf("DROP INDEX idx_project_tasks_code: %w", err)
		}
	}
	// Vérifier si le nouvel index composite (project_id, code) existe déjà
	if err := sqlDB.QueryRow(`
		SELECT COUNT(*) FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'project_tasks' AND INDEX_NAME = 'idx_project_tasks_project_code'
	`).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		log.Println("   🔧 project_tasks: création de l'index unique idx_project_tasks_project_code (project_id, code)...")
		if _, err := sqlDB.Exec("ALTER TABLE project_tasks ADD UNIQUE INDEX idx_project_tasks_project_code (project_id, code)"); err != nil {
			return fmt.Errorf("ADD UNIQUE INDEX idx_project_tasks_project_code: %w", err)
		}
	}
	return nil
}

// migrateSoftwareCodeVersionUnique remplace l'index unique sur (code) par un index unique sur (code, version),
// pour permettre plusieurs versions du même logiciel (ex. ISA 33 et ISA 35).
func migrateSoftwareCodeVersionUnique() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("récupération sqlDB: %w", err)
	}
	// Supprimer l'ancien index unique sur (code) s'il existe
	var n int
	if err := sqlDB.QueryRow(`
		SELECT COUNT(*) FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'software' AND INDEX_NAME = 'idx_software_code'
	`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		log.Println("   🔧 software: suppression de l'index unique idx_software_code (code seul)...")
		if _, err := sqlDB.Exec("ALTER TABLE software DROP INDEX idx_software_code"); err != nil {
			return fmt.Errorf("DROP INDEX idx_software_code: %w", err)
		}
	}
	// Créer l'index unique composite (code, version) s'il n'existe pas
	if err := sqlDB.QueryRow(`
		SELECT COUNT(*) FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'software' AND INDEX_NAME = 'idx_software_code_version'
	`).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		// Normaliser les versions NULL en '' pour que l'unicité (code, version) soit cohérente
		if _, err := sqlDB.Exec("UPDATE software SET version = '' WHERE version IS NULL"); err != nil {
			log.Printf("   ⚠️  software: mise à jour version NULL (ignoré): %v", err)
		}
		log.Println("   🔧 software: création de l'index unique idx_software_code_version (code, version)...")
		if _, err := sqlDB.Exec("ALTER TABLE software ADD UNIQUE INDEX idx_software_code_version (code, version)"); err != nil {
			return fmt.Errorf("ADD UNIQUE INDEX idx_software_code_version: %w", err)
		}
	}
	return nil
}

// migrateProjectsStartEndDates ajoute start_date et end_date à projects si les colonnes n'existent pas.
func migrateProjectsStartEndDates() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("récupération sqlDB: %w", err)
	}
	for _, col := range []string{"start_date", "end_date"} {
		var n int
		if err := sqlDB.QueryRow(`
			SELECT COUNT(*) FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'projects' AND COLUMN_NAME = ?
		`, col).Scan(&n); err != nil {
			return err
		}
		if n == 0 {
			log.Printf("   🔧 projects: ajout de la colonne %s (DATE NULL)", col)
			if _, err := sqlDB.Exec("ALTER TABLE projects ADD COLUMN " + col + " DATE NULL"); err != nil {
				return fmt.Errorf("ADD COLUMN projects."+col+": %w", err)
			}
		}
	}
	return nil
}

// migrateProjectBudgetExtensionsStartEndDates ajoute start_date et end_date à project_budget_extensions si absentes.
func migrateProjectBudgetExtensionsStartEndDates() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("récupération sqlDB: %w", err)
	}
	for _, col := range []string{"start_date", "end_date"} {
		var n int
		if err := sqlDB.QueryRow(`
			SELECT COUNT(*) FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'project_budget_extensions' AND COLUMN_NAME = ?
		`, col).Scan(&n); err != nil {
			return err
		}
		if n == 0 {
			log.Printf("   🔧 project_budget_extensions: ajout de la colonne %s (DATE NULL)", col)
			if _, err := sqlDB.Exec("ALTER TABLE project_budget_extensions ADD COLUMN " + col + " DATE NULL"); err != nil {
				return fmt.Errorf("ADD COLUMN project_budget_extensions."+col+": %w", err)
			}
		}
	}
	return nil
}

// makeAssetSoftwareAssetIDNullable rend la colonne asset_id de asset_software nullable
// Cela permet de créer des logiciels indépendamment des actifs
func makeAssetSoftwareAssetIDNullable() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}

	// Vérifier si la colonne existe et si elle est déjà nullable
	var isNullable string
	var columnType string
	err = sqlDB.QueryRow(`
		SELECT IS_NULLABLE, COLUMN_TYPE
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = 'asset_software'
		AND COLUMN_NAME = 'asset_id'
	`).Scan(&isNullable, &columnType)

	if err != nil {
		// La colonne n'existe pas encore, GORM la créera avec le bon type
		log.Println("   ℹ️  Colonne asset_software.asset_id n'existe pas encore, sera créée par GORM")
		return nil
	}

	// Si la colonne est déjà nullable, rien à faire
	if isNullable == "YES" {
		log.Println("   ℹ️  Colonne asset_software.asset_id est déjà nullable")
		return nil
	}

	// Modifier la colonne pour la rendre nullable
	log.Println("   🔧 Modification de asset_software.asset_id pour la rendre nullable...")

	// D'abord, supprimer la contrainte de clé étrangère si elle existe
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer func() {
		_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}()

	// Récupérer le nom de la contrainte de clé étrangère
	var constraintName string
	err = sqlDB.QueryRow(`
		SELECT CONSTRAINT_NAME
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = 'asset_software'
		AND COLUMN_NAME = 'asset_id'
		AND REFERENCED_TABLE_NAME IS NOT NULL
		LIMIT 1
	`).Scan(&constraintName)

	if err == nil && constraintName != "" {
		log.Printf("   🗑️  Suppression de la contrainte FK: %s", constraintName)
		_, _ = sqlDB.Exec(fmt.Sprintf("ALTER TABLE `asset_software` DROP FOREIGN KEY `%s`", constraintName))
	}

	// Modifier la colonne pour la rendre nullable
	_, err = sqlDB.Exec(`
		ALTER TABLE asset_software 
		MODIFY COLUMN asset_id INT UNSIGNED NULL
	`)
	if err != nil {
		return fmt.Errorf("erreur lors de la modification de la colonne asset_id: %w", err)
	}

	// Recréer la contrainte de clé étrangère avec SET NULL
	_, err = sqlDB.Exec(`
		ALTER TABLE asset_software
		ADD CONSTRAINT fk_asset_software_asset
		FOREIGN KEY (asset_id) REFERENCES assets(id)
		ON DELETE SET NULL
	`)
	if err != nil {
		// Si la contrainte existe déjà ou si elle ne peut pas être créée, ce n'est pas grave
		// Elle sera recréée par addAllForeignKeys()
		log.Printf("   ⚠️  Impossible de recréer la contrainte FK (sera gérée par addAllForeignKeys): %v", err)
	}

	log.Println("   ✅ Colonne asset_software.asset_id modifiée avec succès")
	return nil
}

// migrateMultiFiliales ajoute les colonnes nécessaires pour le support multi-filiales
func migrateMultiFiliales() error {
	if DB == nil {
		return fmt.Errorf("la base de données n'est pas initialisée")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'instance DB: %w", err)
	}

	log.Println("🔧 Migration multi-filiales: ajout des colonnes nécessaires...")

	// Fonction helper pour ajouter une colonne si elle n'existe pas
	addColumnIfNotExists := func(table, column, columnType string) error {
		var exists int
		err := sqlDB.QueryRow(`
			SELECT COUNT(*) FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?
		`, table, column).Scan(&exists)
		if err != nil {
			return err
		}
		if exists == 0 {
			log.Printf("   🔧 Ajout de la colonne %s.%s", table, column)
			_, err = sqlDB.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, columnType))
			if err != nil {
				return fmt.Errorf("erreur lors de l'ajout de %s.%s: %w", table, column, err)
			}
		}
		return nil
	}

	// Ajouter filiale_id aux tables
	tablesWithFilialeID := []struct {
		table      string
		columnType string
	}{
		{"users", "INT UNSIGNED NULL"},
		{"departments", "INT UNSIGNED NULL"},
		{"tickets", "INT UNSIGNED NULL"},
		{"projects", "INT UNSIGNED NULL"},
		{"knowledge_articles", "INT UNSIGNED NULL"},
		{"knowledge_categories", "INT UNSIGNED NULL"},
		{"delays", "INT UNSIGNED NULL"},
		{"daily_declarations", "INT UNSIGNED NULL"},
		{"weekly_declarations", "INT UNSIGNED NULL"},
		{"assets", "INT UNSIGNED NULL"},
		{"offices", "INT UNSIGNED NULL"},
		{"sla", "INT UNSIGNED NULL"},
	}

	for _, t := range tablesWithFilialeID {
		if err := addColumnIfNotExists(t.table, "filiale_id", t.columnType); err != nil {
			log.Printf("   ⚠️  Erreur pour %s.filiale_id: %v", t.table, err)
		}
	}

	// Ajouter software_id aux tickets
	if err := addColumnIfNotExists("tickets", "software_id", "INT UNSIGNED NULL"); err != nil {
		log.Printf("   ⚠️  Erreur pour tickets.software_id: %v", err)
	}

	// Ajouter validated_by_user_id et validated_at aux tickets
	if err := addColumnIfNotExists("tickets", "validated_by_user_id", "INT UNSIGNED NULL"); err != nil {
		log.Printf("   ⚠️  Erreur pour tickets.validated_by_user_id: %v", err)
	}
	if err := addColumnIfNotExists("tickets", "validated_at", "DATETIME NULL"); err != nil {
		log.Printf("   ⚠️  Erreur pour tickets.validated_at: %v", err)
	}

	// Ajouter is_it_department aux departments
	if err := addColumnIfNotExists("departments", "is_it_department", "BOOLEAN DEFAULT FALSE"); err != nil {
		log.Printf("   ⚠️  Erreur pour departments.is_it_department: %v", err)
	}

	// Tickets internes : colonnes time_entries et delays (pour temps et retards sur tickets internes)
	if err := addColumnIfNotExists("time_entries", "ticket_internal_id", "INT UNSIGNED NULL"); err != nil {
		log.Printf("   ⚠️  Erreur pour time_entries.ticket_internal_id: %v", err)
	}
	if err := addColumnIfNotExists("delays", "ticket_internal_id", "INT UNSIGNED NULL"); err != nil {
		log.Printf("   ⚠️  Erreur pour delays.ticket_internal_id: %v", err)
	}
	_, _ = sqlDB.Exec("ALTER TABLE time_entries MODIFY COLUMN ticket_id INT UNSIGNED NULL")
	_, _ = sqlDB.Exec("ALTER TABLE delays MODIFY COLUMN ticket_id INT UNSIGNED NULL")

	// Ajouter les index pour les nouvelles colonnes
	indexes := []struct {
		table  string
		column string
	}{
		{"users", "filiale_id"},
		{"departments", "filiale_id"},
		{"tickets", "filiale_id"},
		{"tickets", "software_id"},
		{"tickets", "validated_by_user_id"},
		{"projects", "filiale_id"},
		{"knowledge_articles", "filiale_id"},
		{"knowledge_categories", "filiale_id"},
		{"delays", "filiale_id"},
		{"daily_declarations", "filiale_id"},
		{"weekly_declarations", "filiale_id"},
		{"assets", "filiale_id"},
		{"offices", "filiale_id"},
		{"sla", "filiale_id"},
		{"departments", "is_it_department"},
		{"time_entries", "ticket_internal_id"},
		{"delays", "ticket_internal_id"},
	}

	for _, idx := range indexes {
		var exists int
		err := sqlDB.QueryRow(`
			SELECT COUNT(*) FROM information_schema.STATISTICS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?
		`, idx.table, idx.column).Scan(&exists)
		if err == nil && exists == 0 {
			indexName := fmt.Sprintf("idx_%s_%s", idx.table, idx.column)
			log.Printf("   🔧 Création de l'index %s sur %s.%s", indexName, idx.table, idx.column)
			_, _ = sqlDB.Exec(fmt.Sprintf("CREATE INDEX %s ON %s (%s)", indexName, idx.table, idx.column))
		}
	}

	// Index uniques delays : un retard par ticket normal et un par ticket interne (créés manuellement pour éviter le bug GORM uniqueIndex)
	for _, u := range []struct{ indexName, column string }{
		{"idx_delays_ticket_id", "ticket_id"},
		{"idx_delays_ticket_internal_id", "ticket_internal_id"},
	} {
		var count int
		var nonUnique int
		_ = sqlDB.QueryRow(`
			SELECT COUNT(*), COALESCE(MAX(NON_UNIQUE), 1) FROM information_schema.STATISTICS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'delays' AND INDEX_NAME = ?
		`, u.indexName).Scan(&count, &nonUnique)
		if count > 0 && nonUnique == 1 {
			log.Printf("   🔧 Suppression de l'index non-unique %s pour le remplacer par un index UNIQUE", u.indexName)
			_, _ = sqlDB.Exec(fmt.Sprintf("DROP INDEX %s ON delays", u.indexName))
		}
		if count == 0 || nonUnique == 1 {
			log.Printf("   🔧 Création de l'index UNIQUE %s sur delays.%s", u.indexName, u.column)
			_, _ = sqlDB.Exec(fmt.Sprintf("CREATE UNIQUE INDEX %s ON delays (%s)", u.indexName, u.column))
		}
	}

	log.Println("   ✅ Migration multi-filiales terminée")
	return nil
}
