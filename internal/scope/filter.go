package scope

import (
	"log"

	"gorm.io/gorm"
)

// ApplyFilialeScope applique le filtrage par filiale sur une requête
// Cette fonction détermine si l'utilisateur peut voir toutes les filiales ou uniquement la sienne
// tableName: nom de la table (ex: "tickets", "projects", "users")
// filialeColumn: nom de la colonne filiale_id dans la table (ex: "filiale_id")
func ApplyFilialeScope(db *gorm.DB, scope *QueryScope, tableName, filialeColumn string) *gorm.DB {
	query := db

	// Si un filtre de filiale spécifique est demandé (pour les rapports filtrés)
	if scope.FilterFilialeID != nil {
		query = query.Where(tableName+"."+filialeColumn+" = ?", *scope.FilterFilialeID)
		return query
	}

	// Si l'utilisateur a la permission de voir globalement (IT MCI CARE CI)
	// Permissions qui donnent accès à toutes les filiales :
	// - reports.view_global : rapports globaux groupe
	// - tickets.resolve_all : résoudre tous les tickets
	// - reports.compare_filiales : comparer entre filiales
	if scope.HasAnyPermission("reports.view_global", "tickets.resolve_all", "reports.compare_filiales") {
		// Pas de filtre, voir toutes les filiales
		return query
	}

	// Sinon, filtrer par la filiale de l'utilisateur
	if scope.FilialeID != nil {
		query = query.Where(tableName+"."+filialeColumn+" = ?", *scope.FilialeID)
		return query
	}

	// Si l'utilisateur n'a pas de filiale assignée et n'a pas les permissions globales,
	// ne rien retourner (sécurité par défaut)
	log.Printf("[scope] ApplyFilialeScope: FilialeID=nil pour user → 0 résultat. Vérifiez que l'utilisateur a filiale_id renseigné en base.")
	query = query.Where("1 = 0")
	return query
}

// assigneesTableChecker est une fonction qui vérifie si la table ticket_assignees existe
// Cette fonction peut être injectée depuis l'extérieur pour éviter les cycles d'importation
var assigneesTableChecker func() bool

// SetAssigneesTableChecker définit la fonction de vérification de la table ticket_assignees
// Cette fonction doit être appelée au démarrage de l'application
func SetAssigneesTableChecker(checker func() bool) {
	assigneesTableChecker = checker
}

// assigneesTableExists vérifie si la table ticket_assignees existe
func assigneesTableExists() bool {
	if assigneesTableChecker == nil {
		// Par défaut, supposer que la table n'existe pas si le checker n'est pas défini
		return false
	}
	return assigneesTableChecker()
}

// ApplyTicketScope applique les filtres de scope sur une requête de tickets
// Cette fonction détermine automatiquement quels tickets l'utilisateur peut voir
// selon ses permissions
func ApplyTicketScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Périmètre forcé pour le tableau de bord (department / filiale / global)
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			return ApplyFilialeScope(query, scope, "tickets", "filiale_id")
		}
	case "department":
		if scope.DepartmentID != nil {
			query = ApplyFilialeScope(query, scope, "tickets", "filiale_id")
			deptID := *scope.DepartmentID
			deptSub := "SELECT id FROM users WHERE department_id = ? AND is_active = 1"
			if assigneesTableExists() {
				query = query.Where(
					"(tickets.created_by_id IN ("+deptSub+") OR tickets.assigned_to_id IN ("+deptSub+") OR tickets.requester_id IN ("+deptSub+") OR EXISTS (SELECT 1 FROM ticket_assignees ta INNER JOIN users u_ta ON u_ta.id = ta.user_id AND u_ta.department_id = ? WHERE ta.ticket_id = tickets.id))",
					deptID, deptID, deptID, deptID,
				)
			} else {
				query = query.Where(
					"(tickets.created_by_id IN ("+deptSub+") OR tickets.assigned_to_id IN ("+deptSub+") OR tickets.requester_id IN ("+deptSub+"))",
					deptID, deptID, deptID,
				)
			}
			return query
		}
	}

	// Appliquer le filtrage par filiale en premier
	query = ApplyFilialeScope(query, scope, "tickets", "filiale_id")

	// Diagnostic : quel branche de permission est utilisée (pour debug liste vide)
	hasViewAll := scope.HasPermission("tickets.view_all")
	hasViewFiliale := scope.HasPermission("tickets.view_filiale")
	hasViewTeam := scope.HasPermission("tickets.view_team")
	hasViewOwn := scope.HasPermission("tickets.view_own")
	hasCreate := scope.HasPermission("tickets.create")
	log.Printf("[scope] ApplyTicketScope: user=%d FilialeID=%v view_all=%v view_filiale=%v view_team=%v view_own=%v create=%v",
		scope.UserID, scope.FilialeID, hasViewAll, hasViewFiliale, hasViewTeam, hasViewOwn, hasCreate)

	// Si l'utilisateur a la permission de voir tous les tickets, pas de filtre supplémentaire
	if hasViewAll {
		return query
	}

	// Si l'utilisateur peut voir tous les tickets de sa filiale (DSI filiale) : ApplyFilialeScope a déjà filtré par filiale
	if hasViewFiliale {
		return query
	}

	// Si l'utilisateur peut voir les tickets de son département (et éventuellement view_own)
	// Inclure aussi les tickets qu'il a créés ou qui lui sont assignés (ex. rôle délégué filiale)
	if scope.HasPermission("tickets.view_team") && scope.DepartmentID != nil {
		query = query.Joins("LEFT JOIN users ON users.id = tickets.requester_id")
		if assigneesTableExists() {
			query = query.Where(
				"users.department_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				*scope.DepartmentID, scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"users.department_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				*scope.DepartmentID, scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres tickets
	if scope.HasPermission("tickets.view_own") {
		// Voir les tickets créés par l'utilisateur, assignés à l'utilisateur,
		// ou où l'utilisateur est dans la liste des assignés (si la table existe)
		if assigneesTableExists() {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Si l'utilisateur a tickets.create mais pas de permission de vue explicite,
	// il peut au moins voir les tickets qu'il a créés (logique : si on peut créer, on peut voir ce qu'on crée)
	if scope.HasPermission("tickets.create") {
		query = query.Where("tickets.created_by_id = ?", scope.UserID)
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// applyTicketScopeTeamOrOwn applique le filtre view_team ou view_own sur la table tickets.
// Utilisé par ApplyTicketScopeForCategory pour incidents, service_requests et changes.
func applyTicketScopeTeamOrOwn(db *gorm.DB, scope *QueryScope, permTeam, permOwn string) *gorm.DB {
	query := db
	if scope.HasPermission(permTeam) && scope.DepartmentID != nil {
		query = query.Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}
	if scope.HasPermission(permOwn) {
		if assigneesTableExists() {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID,
			)
		}
		return query
	}
	return query.Where("1 = 0")
}

// ApplyTicketScopeForCategory applique le scope pour une requête de tickets filtrée par catégorie.
// Utilisé par FindByCategory pour prendre en charge :
// - tickets.view_* (comportement identique à ApplyTicketScope)
// - ticket_categories.view : voir tous les tickets de la catégorie
// - incidents.view_* quand category=incident
// - service_requests.view_* quand category=demande ou service_request
// - changes.view_* quand category=changement ou change
func ApplyTicketScopeForCategory(db *gorm.DB, scope *QueryScope, category string) *gorm.DB {
	query := db

	// 1) tickets.view_* : réutiliser la logique standard
	if scope.HasPermission("tickets.view_all") {
		return query
	}
	if scope.HasPermission("tickets.view_filiale") || scope.HasPermission("tickets.view_team") || scope.HasPermission("tickets.view_own") {
		return ApplyTicketScope(db, scope)
	}

	// 2) ticket_categories.view : accès à tous les tickets de la catégorie (filtre category déjà appliqué)
	if scope.HasPermission("ticket_categories.view") {
		return query
	}

	// 3) incident : incidents.view_all, .view, .view_team, .view_own
	if category == "incident" {
		if scope.HasPermission("incidents.view_all") || scope.HasPermission("incidents.view") {
			return query
		}
		if scope.HasPermission("incidents.view_team") || scope.HasPermission("incidents.view_own") {
			return applyTicketScopeTeamOrOwn(query, scope, "incidents.view_team", "incidents.view_own")
		}
	}

	// 4) demande / service_request : service_requests.view_*
	if category == "demande" || category == "service_request" {
		if scope.HasPermission("service_requests.view_all") || scope.HasPermission("service_requests.view") {
			return query
		}
		if scope.HasPermission("service_requests.view_team") || scope.HasPermission("service_requests.view_own") {
			return applyTicketScopeTeamOrOwn(query, scope, "service_requests.view_team", "service_requests.view_own")
		}
	}

	// 5) changement / change : changes.view_*
	if category == "changement" || category == "change" {
		if scope.HasPermission("changes.view_all") || scope.HasPermission("changes.view") {
			return query
		}
		if scope.HasPermission("changes.view_team") || scope.HasPermission("changes.view_own") {
			return applyTicketScopeTeamOrOwn(query, scope, "changes.view_team", "changes.view_own")
		}
	}

	// 6) autres catégories (developpement, assistance, support) : déjà couvertes par tickets.view_* ou ticket_categories.view ci‑dessus
	// Par défaut : rien
	query = query.Where("1 = 0")
	return query
}

// ApplyReportScope applique les filtres de scope sur une requête de rapports
func ApplyReportScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Périmètre forcé pour le tableau de bord
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			return query.Where("tickets.filiale_id = ?", *scope.FilialeID)
		}
	case "department":
		if scope.DepartmentID != nil {
			query = query.Joins("LEFT JOIN users ON users.id = tickets.requester_id").
				Where("users.department_id = ?", *scope.DepartmentID)
			if scope.FilialeID != nil {
				query = query.Where("tickets.filiale_id = ?", *scope.FilialeID)
			}
			return query
		}
	}

	// Appliquer le filtrage par filiale
	// Si l'utilisateur a la permission de voir les rapports globaux groupe (IT MCI CARE CI)
	if scope.HasPermission("reports.view_global") {
		// Si FilterFilialeID est défini, filtrer par cette filiale spécifique
		if scope.FilterFilialeID != nil {
			query = query.Where("tickets.filiale_id = ?", *scope.FilterFilialeID)
		}
		// Sinon, voir toutes les filiales (pas de filtre)
		return query
	}

	// Si l'utilisateur peut voir les rapports de sa filiale
	if scope.HasPermission("reports.view_filiale") && scope.FilialeID != nil {
		query = query.Where("tickets.filiale_id = ?", *scope.FilialeID)
		return query
	}

	// Si l'utilisateur peut voir les rapports de son département
	if scope.HasPermission("reports.view_team") && scope.DepartmentID != nil {
		// Pour les rapports basés sur les tickets, on filtre par département du demandeur
		// ET par filiale de l'utilisateur
		query = query.Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		if scope.FilialeID != nil {
			query = query.Where("tickets.filiale_id = ?", *scope.FilialeID)
		}
		return query
	}

	// Par défaut, ne rien retourner
	query = query.Where("1 = 0")
	return query
}

// ApplyUserScope applique les filtres de scope sur une requête d'utilisateurs
// Cette fonction détermine automatiquement quels utilisateurs l'utilisateur peut voir
// selon ses permissions
func ApplyUserScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Si l'utilisateur a la permission de voir tous les utilisateurs, pas de filtre supplémentaire
	if scope.HasPermission("users.view_all") {
		// Appliquer le filtrage par filiale seulement si un filtre spécifique est demandé
		if scope.FilterFilialeID != nil {
			query = query.Where("users.filiale_id = ?", *scope.FilterFilialeID)
		}
		return query
	}

	// Si l'utilisateur peut voir les utilisateurs de sa filiale (uniquement la permission explicite view_filiale)
	if scope.HasPermission("users.view_filiale") {
		if scope.FilialeID != nil {
			query = query.Where("users.filiale_id = ?", *scope.FilialeID)
			return query
		}
		query = query.Where("1 = 0")
		return query
	}

	// Si l'utilisateur peut voir les utilisateurs de son département (équipe)
	if scope.HasPermission("users.view_team") && scope.DepartmentID != nil {
		// Filtrer par filiale seulement si l'utilisateur en a une (évite 0 résultat quand FilialeID est nil)
		if scope.FilialeID != nil {
			query = ApplyFilialeScope(query, scope, "users", "filiale_id")
		}
		query = query.Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que son propre profil
	if scope.HasPermission("users.view_own") {
		query = query.Where("id = ?", scope.UserID)
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyTicketScopeToTable applique le scope sur une requête qui utilise Table("tickets")
// Cette fonction est utile pour les requêtes qui utilisent Table() au lieu de Model()
func ApplyTicketScopeToTable(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Périmètre forcé pour le tableau de bord
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			// Inclure les tickets sans filiale (filiale_id IS NULL) comme pour les projets
			return query.Where("(tickets.filiale_id = ? OR tickets.filiale_id IS NULL)", *scope.FilialeID)
		}
	case "department":
		if scope.DepartmentID != nil {
			deptID := *scope.DepartmentID
			deptSub := "SELECT id FROM users WHERE department_id = ? AND is_active = 1"
			if assigneesTableExists() {
				query = query.Where(
					"(tickets.created_by_id IN ("+deptSub+") OR tickets.assigned_to_id IN ("+deptSub+") OR tickets.requester_id IN ("+deptSub+") OR EXISTS (SELECT 1 FROM ticket_assignees ta INNER JOIN users u_ta ON u_ta.id = ta.user_id AND u_ta.department_id = ? WHERE ta.ticket_id = tickets.id))",
					deptID, deptID, deptID, deptID,
				)
			} else {
				query = query.Where(
					"(tickets.created_by_id IN ("+deptSub+") OR tickets.assigned_to_id IN ("+deptSub+") OR tickets.requester_id IN ("+deptSub+"))",
					deptID, deptID, deptID,
				)
			}
			if scope.FilialeID != nil {
				// Inclure les tickets sans filiale (filiale_id IS NULL) comme pour les projets
				query = query.Where("(tickets.filiale_id = ? OR tickets.filiale_id IS NULL)", *scope.FilialeID)
			}
			return query
		}
	}

	// Si l'utilisateur a la permission de voir tous les tickets, pas de filtre
	if scope.HasPermission("tickets.view_all") {
		return query
	}

	// Si l'utilisateur peut voir tous les tickets de sa filiale (DSI filiale)
	if scope.HasPermission("tickets.view_filiale") {
		return query
	}

	// Si l'utilisateur peut voir les tickets de son département (inclure aussi créés/assignés)
	if scope.HasPermission("tickets.view_team") && scope.DepartmentID != nil {
		query = query.Joins("LEFT JOIN users ON users.id = tickets.requester_id")
		if assigneesTableExists() {
			query = query.Where(
				"users.department_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				*scope.DepartmentID, scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"users.department_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				*scope.DepartmentID, scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres tickets
	if scope.HasPermission("tickets.view_own") {
		// Voir les tickets créés par l'utilisateur, assignés à l'utilisateur,
		// ou où l'utilisateur est dans la liste des assignés (si la table existe)
		if assigneesTableExists() {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Si l'utilisateur a tickets.create mais pas de permission de vue explicite,
	// il peut au moins voir les tickets qu'il a créés (logique : si on peut créer, on peut voir ce qu'on crée)
	if scope.HasPermission("tickets.create") {
		query = query.Where("tickets.created_by_id = ?", scope.UserID)
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyAssetScope applique les filtres de scope sur une requête d'actifs
// Cette fonction détermine automatiquement quels actifs l'utilisateur peut voir
// selon ses permissions
func ApplyAssetScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Périmètre forcé pour le tableau de bord
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			return query.Where("assets.filiale_id = ?", *scope.FilialeID)
		}
	case "department":
		if scope.DepartmentID != nil {
			query = ApplyFilialeScope(query, scope, "assets", "filiale_id")
			// Actifs assignés aux membres du département (sous-requête pour éviter les conflits de JOIN)
			query = query.Where("assets.assigned_to_id IN (SELECT id FROM users WHERE department_id = ? AND is_active = 1)", *scope.DepartmentID)
			return query
		}
		// scope=department mais utilisateur sans département → aucun actif
		return query.Where("1 = 0")
	}

	// Appliquer le filtrage par filiale en premier
	query = ApplyFilialeScope(query, scope, "assets", "filiale_id")

	// Si l'utilisateur a la permission de voir tous les actifs, pas de filtre supplémentaire
	if scope.HasPermission("assets.view_all") {
		return query
	}

	// Si l'utilisateur peut voir les actifs de son département
	if scope.HasPermission("assets.view_team") && scope.DepartmentID != nil {
		// Filtrer par département de l'utilisateur assigné (via la relation users.department_id)
		query = query.Joins("LEFT JOIN users ON users.id = assets.assigned_to_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que les actifs qui lui sont assignés
	if scope.HasPermission("assets.view_own") {
		query = query.Where("assigned_to_id = ?", scope.UserID)
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyIncidentScope applique les filtres de scope sur une requête d'incidents
// Cette fonction détermine automatiquement quels incidents l'utilisateur peut voir
// selon ses permissions. Les incidents sont filtrés via leurs tickets associés.
func ApplyIncidentScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Si l'utilisateur a la permission de voir tous les incidents (ou view legacy)
	if scope.HasPermission("incidents.view_all") || scope.HasPermission("incidents.view") {
		return query
	}

	// Si l'utilisateur peut voir les incidents de son équipe
	if scope.HasPermission("incidents.view_team") && scope.DepartmentID != nil {
		// Filtrer par département du demandeur via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = incidents.ticket_id").
			Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres incidents (liés à ses tickets)
	if scope.HasPermission("incidents.view_own") {
		// Filtrer via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = incidents.ticket_id")

		// Voir les incidents dont le ticket est créé par l'utilisateur, assigné à l'utilisateur,
		// ou où l'utilisateur est dans la liste des assignés (si la table existe)
		if assigneesTableExists() {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyServiceRequestScope applique les filtres de scope sur une requête de demandes de service
// Cette fonction détermine automatiquement quelles demandes de service l'utilisateur peut voir
// selon ses permissions. Les demandes de service sont filtrées via leurs tickets associés.
func ApplyServiceRequestScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Si l'utilisateur a la permission de voir toutes les demandes de service (ou view legacy)
	if scope.HasPermission("service_requests.view_all") || scope.HasPermission("service_requests.view") {
		return query
	}

	// Si l'utilisateur peut voir les demandes de service de son équipe
	if scope.HasPermission("service_requests.view_team") && scope.DepartmentID != nil {
		// Filtrer par département du demandeur via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = service_requests.ticket_id").
			Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres demandes de service (liées à ses tickets)
	if scope.HasPermission("service_requests.view_own") {
		// Filtrer via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = service_requests.ticket_id")

		// Voir les demandes de service dont le ticket est créé par l'utilisateur, assigné à l'utilisateur,
		// ou où l'utilisateur est dans la liste des assignés (si la table existe)
		if assigneesTableExists() {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyChangeScope applique les filtres de scope sur une requête de changements
// Cette fonction détermine automatiquement quels changements l'utilisateur peut voir
// selon ses permissions. Les changements sont filtrés via leurs tickets associés.
func ApplyChangeScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Si l'utilisateur a la permission de voir tous les changements (ou view legacy)
	if scope.HasPermission("changes.view_all") || scope.HasPermission("changes.view") {
		return query
	}

	// Si l'utilisateur peut voir les changements de son équipe
	if scope.HasPermission("changes.view_team") && scope.DepartmentID != nil {
		// Filtrer par département du demandeur via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = changes.ticket_id").
			Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres changements (liés à ses tickets)
	if scope.HasPermission("changes.view_own") {
		// Filtrer via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = changes.ticket_id")

		// Voir les changements dont le ticket est créé par l'utilisateur, assigné à l'utilisateur,
		// ou où l'utilisateur est dans la liste des assignés (si la table existe)
		if assigneesTableExists() {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyKnowledgeScope applique les filtres de scope sur une requête d'articles de la base de connaissances
// Cette fonction détermine automatiquement quels articles l'utilisateur peut voir selon ses permissions
func ApplyKnowledgeScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Périmètre forcé pour le tableau de bord
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			return query.Where("knowledge_articles.filiale_id IS NULL OR knowledge_articles.filiale_id = ?", *scope.FilialeID)
		}
	case "department":
		if scope.DepartmentID != nil {
			// Articles dont l'auteur appartient au département
			query = query.Where("knowledge_articles.author_id IN (SELECT id FROM users WHERE department_id = ?)", *scope.DepartmentID)
			if scope.FilialeID != nil {
				query = query.Where("(knowledge_articles.filiale_id IS NULL OR knowledge_articles.filiale_id = ?)", *scope.FilialeID)
			}
			return query
		}
	}

	// Si l'utilisateur a la permission de voir tous les articles (IT MCI CARE CI)
	if scope.HasPermission("knowledge.view_all") {
		// Voir tous les articles (globaux + de toutes les filiales)
		return query
	}

	// Appliquer le filtrage par filiale : voir les articles globaux (filiale_id IS NULL) + articles de sa filiale
	// Si l'utilisateur peut voir toutes les filiales (IT MCI CARE CI)
	if scope.HasAnyPermission("reports.view_global", "tickets.resolve_all") {
		// Pas de filtre par filiale pour les articles globaux, mais filtrer les articles par filiale si nécessaire
		// Les articles globaux sont toujours visibles
	} else if scope.FilialeID != nil {
		// Voir les articles globaux (filiale_id IS NULL) + articles de sa filiale
		query = query.Where("knowledge_articles.filiale_id IS NULL OR knowledge_articles.filiale_id = ?", *scope.FilialeID)
	} else {
		// Si pas de filiale assignée, voir uniquement les articles globaux
		query = query.Where("knowledge_articles.filiale_id IS NULL")
	}

	// Si l'utilisateur peut voir les articles publiés
	if scope.HasPermission("knowledge.view_published") {
		// Voir les articles publiés OU ses propres articles (publiés ou non) si l'utilisateur peut voir ses propres articles
		if scope.HasPermission("knowledge.view_own") {
			query = query.Where("is_published = ? OR author_id = ?", true, scope.UserID)
		} else {
			query = query.Where("is_published = ?", true)
		}
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres articles
	if scope.HasPermission("knowledge.view_own") {
		query = query.Where("author_id = ?", scope.UserID)
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyTimeEntryScope applique les filtres de scope sur une requête d'entrées de temps
// Cette fonction détermine automatiquement quelles entrées de temps l'utilisateur peut voir
// selon ses permissions. Les entrées de temps sont filtrées via timesheet ou tickets associés.
// Note: Le timesheet est réservé à IT MCI CARE CI uniquement
func ApplyTimeEntryScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Appliquer le filtrage par filiale via les tickets associés
	// Si l'utilisateur peut voir toutes les filiales (IT MCI CARE CI)
	if scope.HasAnyPermission("reports.view_global", "tickets.resolve_all") {
		// Pas de filtre par filiale, voir toutes les entrées
	} else if scope.FilialeID != nil {
		// Filtrer par filiale via le ticket
		query = query.Joins("INNER JOIN tickets ON tickets.id = time_entries.ticket_id").
			Where("tickets.filiale_id = ?", *scope.FilialeID)
	} else {
		// Si pas de filiale assignée et pas de permissions globales, ne rien retourner
		query = query.Where("1 = 0")
		return query
	}

	// Si l'utilisateur a la permission de voir toutes les entrées de temps
	if scope.HasPermission("timesheet.view_all") {
		return query
	}

	// Si l'utilisateur peut voir les entrées de temps de son équipe
	if scope.HasPermission("timesheet.view_team") && scope.DepartmentID != nil {
		// Filtrer par département du demandeur via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = time_entries.ticket_id").
			Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres entrées de temps
	if scope.HasPermission("timesheet.view_own") {
		// Filtrer via le ticket associé OU si l'entrée de temps appartient à l'utilisateur
		query = query.Joins("INNER JOIN tickets ON tickets.id = time_entries.ticket_id")

		// Voir les entrées de temps dont le ticket est créé par l'utilisateur, assigné à l'utilisateur,
		// ou où l'utilisateur est dans la liste des assignés (si la table existe)
		// OU si l'entrée de temps appartient directement à l'utilisateur
		if assigneesTableExists() {
			query = query.Where(
				"time_entries.user_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"time_entries.user_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyTimeEntryScopeForPendingValidation est utilisé pour la liste "en attente de validation".
// Il élargit le scope par rapport à ApplyTimeEntryScope afin que les utilisateurs avec
// timesheet.validate voient les entrées qu'ils sont habilités à valider (ex. entrées
// des membres de leur département, ou toutes si validateur sans département).
func ApplyTimeEntryScopeForPendingValidation(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	if scope.HasPermission("timesheet.view_all") {
		return query
	}

	if scope.HasPermission("timesheet.view_team") && scope.DepartmentID != nil {
		query = query.Joins("INNER JOIN tickets ON tickets.id = time_entries.ticket_id").
			Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("(users.department_id = ?) OR (EXISTS (SELECT 1 FROM users u_te ON u_te.id = time_entries.user_id AND u_te.department_id = ?))",
				*scope.DepartmentID, *scope.DepartmentID)
		return query
	}

	if scope.HasPermission("timesheet.view_own") {
		// Validateur avec timesheet.validate mais sans département ni view_team/view_all : voir toutes les entrées en attente
		if scope.HasPermission("timesheet.validate") &&
			!scope.HasPermission("timesheet.view_all") &&
			!scope.HasPermission("timesheet.view_team") &&
			scope.DepartmentID == nil {
			return query
		}

		query = query.Joins("INNER JOIN tickets ON tickets.id = time_entries.ticket_id")
		if assigneesTableExists() {
			if scope.HasPermission("timesheet.validate") && scope.DepartmentID != nil {
				query = query.Where(
					"(time_entries.user_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)) OR (EXISTS (SELECT 1 FROM users u_te ON u_te.id = time_entries.user_id AND u_te.department_id = ?))",
					scope.UserID, scope.UserID, scope.UserID, scope.UserID, *scope.DepartmentID)
			} else {
				query = query.Where(
					"time_entries.user_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
					scope.UserID, scope.UserID, scope.UserID, scope.UserID)
			}
		} else {
			if scope.HasPermission("timesheet.validate") && scope.DepartmentID != nil {
				query = query.Where(
					"(time_entries.user_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ?) OR (EXISTS (SELECT 1 FROM users u_te ON u_te.id = time_entries.user_id AND u_te.department_id = ?))",
					scope.UserID, scope.UserID, scope.UserID, *scope.DepartmentID)
			} else {
				query = query.Where(
					"time_entries.user_id = ? OR tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
					scope.UserID, scope.UserID, scope.UserID)
			}
		}
		return query
	}

	query = query.Where("1 = 0")
	return query
}

// ApplyDelayScope applique les filtres de scope sur une requête de retards.
// Trois niveaux de permission :
//   - delays.view_all : tous les retards ; FilterUserID optionnel pour limiter à un utilisateur.
//   - delays.view_department : retards des membres du département ; FilterUserID optionnel pour un membre.
//   - delays.view_own : uniquement ses propres retards (FilterUserID ignoré).
//
// Note: Les retards sont filtrés par filiale via le ticket associé
func ApplyDelayScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Appliquer le filtrage par filiale via le ticket associé
	// Si l'utilisateur peut voir toutes les filiales (IT MCI CARE CI)
	if scope.HasAnyPermission("reports.view_global", "tickets.resolve_all") {
		// Pas de filtre par filiale, voir tous les retards
	} else if scope.FilialeID != nil {
		// Filtrer par filiale via le ticket
		query = query.Joins("INNER JOIN tickets ON tickets.id = delays.ticket_id").
			Where("tickets.filiale_id = ?", *scope.FilialeID)
	} else {
		// Si pas de filiale assignée et pas de permissions globales, ne rien retourner
		query = query.Where("1 = 0")
		return query
	}

	// delays.view_all (ou delays.view legacy) : tous les retards
	if scope.HasPermission("delays.view_all") || scope.HasPermission("delays.view") {
		if scope.FilterUserID != nil {
			return query.Where("delays.user_id = ?", *scope.FilterUserID)
		}
		return query
	}

	// delays.view_department : retards des membres de son département
	if scope.HasPermission("delays.view_department") && scope.DepartmentID != nil {
		query = query.Joins("LEFT JOIN users u ON u.id = delays.user_id").
			Where("u.department_id = ?", *scope.DepartmentID)
		if scope.FilterUserID != nil {
			query = query.Where("delays.user_id = ?", *scope.FilterUserID)
		}
		return query
	}

	// delays.view_own : uniquement ses propres retards (FilterUserID ignoré)
	if scope.HasPermission("delays.view_own") {
		return query.Where("delays.user_id = ?", scope.UserID)
	}

	// Par défaut : ne rien retourner
	return query.Where("1 = 0")
}

// ApplyAuditScope applique les filtres de scope sur une requête de logs d'audit
// Cette fonction détermine automatiquement quels logs d'audit l'utilisateur peut voir
// selon ses permissions du module audit.
func ApplyAuditScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Si l'utilisateur a la permission de voir tous les logs d'audit
	if scope.HasPermission("audit.view_all") {
		return query
	}

	// Si l'utilisateur peut voir les logs de son équipe
	if scope.HasPermission("audit.view_team") && scope.DepartmentID != nil {
		// Filtrer par département de l'utilisateur qui a effectué l'action
		query = query.Joins("LEFT JOIN users ON users.id = audit_logs.user_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres logs
	if scope.HasPermission("audit.view_own") {
		query = query.Where("user_id = ?", scope.UserID)
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyProjectScope applique les filtres de scope sur une requête de projets
// Cette fonction détermine automatiquement quels projets l'utilisateur peut voir
// selon ses permissions. Les projets sont filtrés via leurs tickets associés.
func ApplyProjectScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Périmètre forcé pour le tableau de bord
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			return query.Where("projects.filiale_id = ?", *scope.FilialeID)
		}
	case "department":
		if scope.DepartmentID != nil {
			// Filiale : même filiale que l'utilisateur, ou pas de filiale (projets créés avant affectation filiale)
			if scope.FilialeID != nil {
				query = query.Where("(projects.filiale_id = ? OR projects.filiale_id IS NULL)", *scope.FilialeID)
			} else {
				query = query.Where("projects.filiale_id IS NULL")
			}
			// Projets créés par les membres du département
			query = query.Where("projects.created_by_id IN (SELECT id FROM users WHERE department_id = ? AND is_active = 1)", *scope.DepartmentID)
			return query
		}
	}

	// Appliquer le filtrage par filiale seulement si l'utilisateur a une filiale (évite 0 résultat pour "Mes projets" quand scope=own sans FilialeID)
	if scope.FilialeID != nil {
		query = ApplyFilialeScope(query, scope, "projects", "filiale_id")
	}

	// Si l'utilisateur a la permission de voir tous les projets (ou view legacy)
	if scope.HasPermission("projects.view_all") || scope.HasPermission("projects.view") {
		return query
	}

	// Si l'utilisateur peut voir les projets de son équipe (chef de département : projets où lui ou un membre du département est membre, assigné à une tâche, ou a des tickets liés)
	if scope.HasPermission("projects.view_team") && scope.DepartmentID != nil {
		deptID := *scope.DepartmentID
		// Projets créés par un membre du département, ou avec un membre du département comme membre du projet, assigné à une tâche, ou ayant un ticket lié
		byCreatedBy := "projects.created_by_id IN (SELECT id FROM users WHERE department_id = ? AND is_active = 1)"
		byMember := "EXISTS (SELECT 1 FROM project_members pm INNER JOIN users u ON u.id = pm.user_id WHERE pm.project_id = projects.id AND u.department_id = ? AND u.is_active = 1)"
		byTaskAssignee := "EXISTS (SELECT 1 FROM project_tasks pt INNER JOIN users u ON u.id = pt.assigned_to_id WHERE pt.project_id = projects.id AND u.department_id = ? AND u.is_active = 1) OR " +
			"EXISTS (SELECT 1 FROM project_task_assignees pta INNER JOIN project_tasks pt ON pt.id = pta.project_task_id INNER JOIN users u ON u.id = pta.user_id WHERE pt.project_id = projects.id AND u.department_id = ? AND u.is_active = 1)"
		byTicket := "EXISTS (SELECT 1 FROM ticket_projects tp INNER JOIN tickets t ON t.id = tp.ticket_id LEFT JOIN users u ON u.id = t.requester_id WHERE tp.project_id = projects.id AND u.department_id = ?)"
		query = query.Where("("+byCreatedBy+" OR "+byMember+" OR "+byTaskAssignee+" OR "+byTicket+")", deptID, deptID, deptID, deptID, deptID)
		return query
	}

	// Si l'utilisateur ne peut voir que ses propres projets (créateur, membre, assigné à une tâche, ou tickets liés)
	if scope.HasPermission("projects.view_own") {
		// Projets créés par l'utilisateur, ou où il est membre (project_members), assigné à une tâche, ou a des tickets liés
		byCreatedBy := "projects.created_by_id = ?"
		byMember := "EXISTS (SELECT 1 FROM project_members pm WHERE pm.project_id = projects.id AND pm.user_id = ?)"
		byTask := "EXISTS (SELECT 1 FROM project_tasks pt WHERE pt.project_id = projects.id AND pt.assigned_to_id = ?) OR " +
			"EXISTS (SELECT 1 FROM project_task_assignees pta INNER JOIN project_tasks pt ON pt.id = pta.project_task_id WHERE pt.project_id = projects.id AND pta.user_id = ?)"
		var byTicket string
		if assigneesTableExists() {
			byTicket = "EXISTS (SELECT 1 FROM ticket_projects tp " +
				"INNER JOIN tickets t ON t.id = tp.ticket_id " +
				"WHERE tp.project_id = projects.id AND " +
				"(t.created_by_id = ? OR t.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = t.id AND ta.user_id = ?)))"
		} else {
			byTicket = "EXISTS (SELECT 1 FROM ticket_projects tp " +
				"INNER JOIN tickets t ON t.id = tp.ticket_id " +
				"WHERE tp.project_id = projects.id AND (t.created_by_id = ? OR t.assigned_to_id = ?))"
		}
		if assigneesTableExists() {
			query = query.Where("("+byCreatedBy+" OR "+byMember+" OR "+byTask+" OR "+byTicket+")", scope.UserID, scope.UserID, scope.UserID, scope.UserID, scope.UserID, scope.UserID, scope.UserID)
		} else {
			query = query.Where("("+byCreatedBy+" OR "+byMember+" OR "+byTask+" OR "+byTicket+")", scope.UserID, scope.UserID, scope.UserID, scope.UserID, scope.UserID, scope.UserID)
		}
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplySLAScope applique les filtres de scope sur une requête de violations SLA
// Cette fonction détermine automatiquement quelles violations SLA l'utilisateur peut voir
// selon ses permissions du module sla. Les violations sont filtrées via leurs tickets associés.
func ApplySLAScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db

	// Périmètre forcé pour le tableau de bord
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			return query.Joins("INNER JOIN tickets ON tickets.id = ticket_sla.ticket_id").
				Where("tickets.filiale_id = ?", *scope.FilialeID)
		}
	case "department":
		if scope.DepartmentID != nil {
			return query.Joins("INNER JOIN tickets ON tickets.id = ticket_sla.ticket_id").
				Joins("LEFT JOIN users ON users.id = tickets.requester_id").
				Where("users.department_id = ?", *scope.DepartmentID)
		}
	}

	// Appliquer le filtrage par filiale via les tickets associés
	// Si l'utilisateur peut voir toutes les filiales (IT MCI CARE CI)
	if scope.HasAnyPermission("reports.view_global", "tickets.resolve_all") {
		// Pas de filtre par filiale, voir toutes les violations
	} else if scope.FilialeID != nil {
		// Filtrer par filiale via le ticket
		query = query.Joins("INNER JOIN tickets ON tickets.id = ticket_sla.ticket_id").
			Where("tickets.filiale_id = ?", *scope.FilialeID)
	} else {
		// Si pas de filiale assignée et pas de permissions globales, ne rien retourner
		query = query.Where("1 = 0")
		return query
	}

	// Si l'utilisateur a la permission de voir tous les SLA/violations (ou view legacy)
	if scope.HasPermission("sla.view_all") || scope.HasPermission("sla.view") {
		return query
	}

	// Si l'utilisateur peut voir les SLA de son équipe
	if scope.HasPermission("sla.view_team") && scope.DepartmentID != nil {
		// Filtrer par département du demandeur via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = ticket_sla.ticket_id").
			Joins("LEFT JOIN users ON users.id = tickets.requester_id").
			Where("users.department_id = ?", *scope.DepartmentID)
		return query
	}

	// Si l'utilisateur ne peut voir que ses SLA (liés à ses tickets)
	if scope.HasPermission("sla.view_own") {
		// Filtrer via le ticket associé
		query = query.Joins("INNER JOIN tickets ON tickets.id = ticket_sla.ticket_id")

		// Voir les violations dont le ticket est créé par l'utilisateur, assigné à l'utilisateur,
		// ou où l'utilisateur est dans la liste des assignés (si la table existe)
		if assigneesTableExists() {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
				scope.UserID, scope.UserID, scope.UserID,
			)
		} else {
			query = query.Where(
				"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
				scope.UserID, scope.UserID,
			)
		}
		return query
	}

	// Par défaut, si aucune permission de vue n'est trouvée, ne rien retourner
	// (sécurité par défaut : ne rien montrer)
	query = query.Where("1 = 0")
	return query
}

// ApplyTicketInternalScope applique les filtres de scope sur une requête de tickets internes (table ticket_internes).
// Permissions : tickets_internes.view_all, view_filiale, view_department, view_own.
func ApplyTicketInternalScope(db *gorm.DB, scope *QueryScope) *gorm.DB {
	return ApplyTicketInternalScopeToTable(db, scope)
}

// ApplyTicketInternalScopeToTable applique le scope sur une requête qui utilise Table("ticket_internes")
func ApplyTicketInternalScopeToTable(db *gorm.DB, scope *QueryScope) *gorm.DB {
	query := db
	if scope == nil {
		return query.Where("1 = 0")
	}
	// Périmètre forcé pour le tableau de bord (department / filiale / global) — même logique que les tickets
	switch scope.DashboardScopeHint {
	case "global":
		return query
	case "filiale":
		if scope.FilialeID != nil {
			return query.Where("ticket_internes.filiale_id = ?", *scope.FilialeID)
		}
	case "department":
		if scope.DepartmentID != nil {
			return query.Where("ticket_internes.department_id = ?", *scope.DepartmentID)
		}
	}

	if scope.HasPermission("tickets_internes.view_all") {
		return query
	}
	if scope.HasPermission("tickets_internes.view_filiale") && scope.FilialeID != nil {
		return query.Where("ticket_internes.filiale_id = ?", *scope.FilialeID)
	}
	if scope.HasPermission("tickets_internes.view_department") && scope.DepartmentID != nil {
		return query.Where("ticket_internes.department_id = ?", *scope.DepartmentID)
	}
	if scope.HasPermission("tickets_internes.view_own") {
		return query.Where("ticket_internes.created_by_id = ? OR ticket_internes.assigned_to_id = ?", scope.UserID, scope.UserID)
	}
	query = query.Where("1 = 0")
	return query
}
