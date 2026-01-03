# Migrations de Base de Données

Ce dossier contient les migrations pour créer toutes les tables de la base de données.

## Utilisation

### Exécuter les migrations

Pour créer toutes les tables dans la base de données :

```bash
go run cmd/migrate/main.go
```

### Exécuter les migrations avec seeding

Pour créer les tables et insérer les données initiales (rôles système) :

```bash
go run cmd/migrate/main.go -seed
```

## Structure

- `migrate.go` : Contient la fonction `RunMigrations()` qui crée toutes les tables et `SeedData()` qui insère les données initiales
- `cmd/migrate/main.go` : Point d'entrée pour exécuter les migrations

## Tables créées

Les migrations créent les tables suivantes (dans l'ordre) :

1. **Authentification et utilisateurs**
   - `roles`
   - `permissions`
   - `role_permissions`
   - `users`
   - `user_sessions`

2. **Tickets**
   - `tickets`
   - `ticket_comments`
   - `ticket_history`
   - `ticket_attachments`
   - `ticket_tags`
   - `ticket_tag_assignments`

3. **Incidents**
   - `incidents`
   - `incident_assets`

4. **Demandes de service**
   - `service_request_types`
   - `service_requests`

5. **Changements**
   - `changes`

6. **Gestion du temps**
   - `time_entries`
   - `daily_declarations`
   - `daily_declaration_tasks`
   - `weekly_declarations`
   - `weekly_declaration_tasks`

7. **Retards**
   - `delays`
   - `delay_justifications`

8. **Actifs IT**
   - `asset_categories`
   - `assets`
   - `ticket_assets`

9. **SLA**
   - `slas`
   - `ticket_slas`

10. **Notifications**
    - `notifications`

11. **Base de connaissances**
    - `knowledge_categories`
    - `knowledge_articles`
    - `knowledge_article_attachments`

12. **Projets**
    - `projects`
    - `ticket_projects`

13. **Paramétrage**
    - `settings`
    - `request_sources`

14. **Audit et sauvegarde**
    - `audit_logs`
    - `backup_configurations`
    - `backups`

## Données initiales

Le seeding insère les rôles système suivants :
- `DSI` : Directeur des Systèmes d'Information
- `RESPONSABLE_IT` : Responsable IT
- `TECHNICIEN_IT` : Technicien IT

## Notes

- Les migrations utilisent GORM `AutoMigrate` qui crée automatiquement les tables basées sur les modèles
- Les relations (clés étrangères) sont créées automatiquement par GORM
- Les migrations sont idempotentes : elles peuvent être exécutées plusieurs fois sans problème
- Le seeding vérifie si les données existent déjà avant d'insérer

