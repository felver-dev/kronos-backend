package repositories

import (
	"fmt"
	"sync"
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
	"gorm.io/gorm"
)

var (
	assigneesTableOnce sync.Once
	hasAssigneesTable  bool
)

func assigneesTableExists() bool {
	assigneesTableOnce.Do(func() {
		hasAssigneesTable = database.DB.Migrator().HasTable(&models.TicketAssignee{})
	})
	if hasAssigneesTable {
		return true
	}
	// Ne pas faire confiance au cache false indéfiniment (migration retardée) : revérifier
	return database.DB.Migrator().HasTable(&models.TicketAssignee{})
}

// TicketRepository interface pour les opérations sur les tickets
type TicketRepository interface {
	Create(ticket *models.Ticket) error
	ExistsByID(id uint) (bool, error)
	FindByID(id uint) (*models.Ticket, error)
	FindByIDLean(id uint) (*models.Ticket, error)
	FindByIDForAssign(id uint) (*models.Ticket, error)
	FindByIDForUpdate(id uint) (*models.Ticket, error)
	FindAll(scope interface{}, page, limit int, filterFilialeID *uint) ([]models.Ticket, int64, error) // scope peut être *scope.QueryScope ou nil; filterFilialeID = filtre par filiale du ticket (envoyée par)
	FindWithFilters(scope interface{}, page, limit int, status string, filterFilialeID *uint, assigneeUserID *uint) ([]models.Ticket, int64, error)
	FindByStatus(scope interface{}, status string, page, limit int) ([]models.Ticket, int64, error)
	FindByCategory(scope interface{}, category string, page, limit int, status, priority string) ([]models.Ticket, int64, error)
	FindByPriority(priority string) ([]models.Ticket, error)
	FindByAssignedTo(userID uint, page, limit int) ([]models.Ticket, int64, error)
	FindPanierByUser(userID uint, page, limit int) ([]models.Ticket, int64, error) // Tickets assignés à l'utilisateur, non clôturés
	FindByCreatedBy(userID uint, page, limit int) ([]models.Ticket, int64, error)
	FindByUser(userID uint, page, limit int, status string) ([]models.Ticket, int64, error)
	FindBySource(scope interface{}, source string, page, limit int) ([]models.Ticket, int64, error)
	FindByDepartment(departmentID uint, page, limit int) ([]models.Ticket, int64, error)
	Update(ticket *models.Ticket) error
	UpdateFields(id uint, updates map[string]interface{}) error
	UpdateAssignFields(id uint, assignedToID *uint, estimatedTime *int, status string) error
	UpdateRequesterNameByCreatedBy(createdByID uint, requesterName string) error
	UpdateRequesterNameByName(oldName string, newName string) error
	UpdateRequesterNameByRequesterID(requesterID uint, requesterName string) error
	Delete(id uint) error
	CountByStatus(status string) (int64, error)
	CountByCategory(category string) (int64, error)
	Search(scope interface{}, query string, status string, limit int) ([]models.Ticket, error) // scope peut être *scope.QueryScope ou nil
	GetNextSequenceNumber(year int) (int, error) // Obtient le prochain numéro séquentiel pour une année donnée
	CodeExists(code string) (bool, error)        // Vérifie si un code existe déjà
}

// ticketRepository implémente TicketRepository
type ticketRepository struct{}

// NewTicketRepository crée une nouvelle instance de TicketRepository
func NewTicketRepository() TicketRepository {
	return &ticketRepository{}
}

// applyTicketPreloads applique les Preloads standards pour les tickets (relations communes)
func applyTicketPreloads(query *gorm.DB) *gorm.DB {
	return query.Preload("CreatedBy").Preload("CreatedBy.Department").
		Preload("AssignedTo").
		Preload("Requester").
		Preload("ValidatedBy").
		Preload("Filiale").
		Preload("Software").
		Preload("Assignees").Preload("Assignees.User")
}

// applyTicketPreloadsBasic applique les Preloads de base (sans toutes les relations)
func applyTicketPreloadsBasic(query *gorm.DB) *gorm.DB {
	return query.Preload("CreatedBy").Preload("AssignedTo").
		Preload("Filiale").Preload("Software")
}

// Create crée un nouveau ticket
func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return database.DB.Create(ticket).Error
}

// ExistsByID vérifie si un ticket existe (requête légère)
func (r *ticketRepository) ExistsByID(id uint) (bool, error) {
	var count int64
	if err := database.DB.Model(&models.Ticket{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByID trouve un ticket par son ID avec ses relations (incluant les attachments)
func (r *ticketRepository) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := applyTicketPreloads(database.DB).
		First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// FindByIDLean charge un ticket avec un minimum de relations pour la vue détail
func (r *ticketRepository) FindByIDLean(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.Preload("CreatedBy").
		Preload("AssignedTo").
		Preload("Requester").
		Preload("Filiale").
		Preload("Software").
		Preload("Assignees").Preload("Assignees.User").
		First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// FindByIDForAssign charge uniquement les champs nécessaires pour l'assignation
func (r *ticketRepository) FindByIDForAssign(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.Select("id", "assigned_to_id", "status", "estimated_time").
		First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// FindByIDForUpdate charge uniquement les champs nécessaires pour la mise à jour
func (r *ticketRepository) FindByIDForUpdate(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.Select(
		"id",
		"code",
		"title",
		"description",
		"category",
		"source",
		"status",
		"priority",
		"requester_id",
		"requester_name",
		"requester_department",
		"assigned_to_id",
		"created_by_id",
		"parent_id",
	).
		First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// FindAll récupère tous les tickets avec leurs relations (avec pagination)
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *ticketRepository) FindAll(scopeParam interface{}, page, limit int, filterFilialeID *uint) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	// Construire la requête de base
	baseQuery := database.DB.Model(&models.Ticket{})

	if filterFilialeID != nil {
		baseQuery = baseQuery.Where("filiale_id = ?", *filterFilialeID)
	}

	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScope(baseQuery, queryScope)
		}
	}

	// Compter le total avec le scope appliqué
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les tickets avec pagination
	// Note: On doit réappliquer le scope sur la requête de récupération
	query := database.DB.Model(&models.Ticket{})
	if filterFilialeID != nil {
		query = query.Where("filiale_id = ?", *filterFilialeID)
	}
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketScope(query, queryScope)
		}
	}

	err := applyTicketPreloadsBasic(query).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error

	return tickets, total, err
}

// FindWithFilters récupère les tickets avec filtres optionnels (statut, filiale, assigné)
func (r *ticketRepository) FindWithFilters(scopeParam interface{}, page, limit int, status string, filterFilialeID *uint, assigneeUserID *uint) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	baseQuery := database.DB.Model(&models.Ticket{})
	if status != "" {
		baseQuery = baseQuery.Where("status = ?", status)
	}
	if filterFilialeID != nil {
		baseQuery = baseQuery.Where("filiale_id = ?", *filterFilialeID)
	}
	if assigneeUserID != nil {
		baseQuery = baseQuery.Where("assigned_to_id = ? OR id IN (SELECT ticket_id FROM ticket_assignees WHERE user_id = ?)", *assigneeUserID, *assigneeUserID)
	}
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScope(baseQuery, queryScope)
		}
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Ticket{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if filterFilialeID != nil {
		query = query.Where("filiale_id = ?", *filterFilialeID)
	}
	if assigneeUserID != nil {
		query = query.Where("assigned_to_id = ? OR id IN (SELECT ticket_id FROM ticket_assignees WHERE user_id = ?)", *assigneeUserID, *assigneeUserID)
	}
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketScope(query, queryScope)
		}
	}

	err := applyTicketPreloadsBasic(query).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error
	return tickets, total, err
}

// FindByStatus récupère les tickets par statut (avec pagination)
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *ticketRepository) FindByStatus(scopeParam interface{}, status string, page, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	// Construire la requête de base avec le filtre de statut
	baseQuery := database.DB.Model(&models.Ticket{}).Where("status = ?", status)

	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScope(baseQuery, queryScope)
		}
	}

	// Compter le total avec le scope appliqué
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les tickets avec pagination
	query := database.DB.Model(&models.Ticket{}).Where("status = ?", status)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketScope(query, queryScope)
		}
	}

	err := applyTicketPreloadsBasic(query).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error

	return tickets, total, err
}

// FindByCategory récupère les tickets par catégorie (avec pagination et filtres optionnels)
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *ticketRepository) FindByCategory(scopeParam interface{}, category string, page, limit int, status, priority string) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	fmt.Printf("DEBUG: FindByCategory - category: %s, status: %s, priority: %s\n", category, status, priority)

	// Construire la requête de base
	query := database.DB.Model(&models.Ticket{}).Where("category = ?", category)
	countQuery := database.DB.Model(&models.Ticket{}).Where("category = ?", category)

	// Ajouter les filtres optionnels
	if status != "" && status != "all" {
		fmt.Printf("DEBUG: Ajout du filtre status: %s\n", status)
		query = query.Where("status = ?", status)
		countQuery = countQuery.Where("status = ?", status)
	}
	if priority != "" && priority != "all" {
		fmt.Printf("DEBUG: Ajout du filtre priority: %s\n", priority)
		query = query.Where("priority = ?", priority)
		countQuery = countQuery.Where("priority = ?", priority)
	} else {
		fmt.Printf("DEBUG: Pas de filtre priority (valeur: '%s')\n", priority)
	}

	// Appliquer le scope si fourni (scope dépendant de la catégorie : incidents.*, service_requests.*, changes.*, ticket_categories.view, tickets.view_*)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketScopeForCategory(query, queryScope, category)
			countQuery = scope.ApplyTicketScopeForCategory(countQuery, queryScope, category)
		}
	}

	// Compter le total avec les filtres et le scope
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les tickets avec pagination
	err := applyTicketPreloadsBasic(query).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error

	return tickets, total, err
}

// FindByPriority récupère les tickets par priorité
func (r *ticketRepository) FindByPriority(priority string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("AssignedTo").
		Preload("Filiale").Preload("Software").
		Preload("Assignees").Preload("Assignees.User").
		Where("priority = ?", priority).Find(&tickets).Error
	return tickets, err
}

// FindByAssignedTo récupère les tickets assignés à un utilisateur (avec pagination)
func (r *ticketRepository) FindByAssignedTo(userID uint, page, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64
	start := time.Now()

	baseQuery := database.DB.Model(&models.Ticket{})
	if assigneesTableExists() {
		baseQuery = baseQuery.Where(
			"tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
			userID, userID,
		)
	} else {
		baseQuery = baseQuery.Where("tickets.assigned_to_id = ?", userID)
	}

	// Compter le total
	countStart := time.Now()
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	countDur := time.Since(countStart)

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les tickets avec pagination
	queryStart := time.Now()
	err := applyTicketPreloadsBasic(baseQuery).
		Order("tickets.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error
	queryDur := time.Since(queryStart)
	fmt.Printf("PERF FindByAssignedTo user=%d count=%d countDur=%s queryDur=%s totalDur=%s\n", userID, total, countDur, queryDur, time.Since(start))

	return tickets, total, err
}

// FindPanierByUser récupère les tickets assignés à l'utilisateur et non clôturés (panier / file de travail)
func (r *ticketRepository) FindPanierByUser(userID uint, page, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	if userID == 0 {
		return tickets, 0, nil
	}

	baseQuery := database.DB.Model(&models.Ticket{}).Where("tickets.status NOT IN (?, ?)", "cloture", "closed")
	if assigneesTableExists() {
		baseQuery = baseQuery.Where(
			"tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
			userID, userID,
		)
	} else {
		baseQuery = baseQuery.Where("tickets.assigned_to_id = ?", userID)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := baseQuery.
		Preload("CreatedBy").Preload("AssignedTo").
		Preload("Requester").
		Preload("Filiale").Preload("Software").
		Order("tickets.updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error
	return tickets, total, err
}

// FindByCreatedBy récupère les tickets créés par un utilisateur (avec pagination)
func (r *ticketRepository) FindByCreatedBy(userID uint, page, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	// Compter le total
	if err := database.DB.Model(&models.Ticket{}).Where("created_by_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les tickets avec pagination
	err := applyTicketPreloadsBasic(database.DB).
		Where("created_by_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error

	return tickets, total, err
}

// FindByUser récupère les tickets créés par l'utilisateur ou qui lui sont assignés (avec pagination)
func (r *ticketRepository) FindByUser(userID uint, page, limit int, status string) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64
	start := time.Now()

	baseQuery := database.DB.Model(&models.Ticket{})
	if assigneesTableExists() {
		baseQuery = baseQuery.Where(
			"tickets.created_by_id = ? OR tickets.assigned_to_id = ? OR EXISTS (SELECT 1 FROM ticket_assignees ta WHERE ta.ticket_id = tickets.id AND ta.user_id = ?)",
			userID, userID, userID,
		)
	} else {
		baseQuery = baseQuery.Where(
			"tickets.created_by_id = ? OR tickets.assigned_to_id = ?",
			userID, userID,
		)
	}

	if status != "" && status != "all" {
		// "cloture" inclut aussi les tickets validés (resolu) pour les stats et listes "résolus"
		if status == "cloture" {
			baseQuery = baseQuery.Where("tickets.status IN ?", []string{"cloture", "resolu"})
		} else {
			baseQuery = baseQuery.Where("tickets.status = ?", status)
		}
	}

	countStart := time.Now()
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	countDur := time.Since(countStart)

	offset := (page - 1) * limit

	queryStart := time.Now()
	err := applyTicketPreloadsBasic(baseQuery).
		Order("tickets.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error
	queryDur := time.Since(queryStart)
	fmt.Printf("PERF FindByUser user=%d status=%s count=%d countDur=%s queryDur=%s totalDur=%s\n", userID, status, total, countDur, queryDur, time.Since(start))

	return tickets, total, err
}

// FindBySource récupère les tickets par source (avec pagination)
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *ticketRepository) FindBySource(scopeParam interface{}, source string, page, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	// Construire la requête de base avec le filtre de source
	baseQuery := database.DB.Model(&models.Ticket{}).Where("source = ?", source)

	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			baseQuery = scope.ApplyTicketScope(baseQuery, queryScope)
		}
	}

	// Compter le total avec le scope appliqué
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les tickets avec pagination
	query := database.DB.Model(&models.Ticket{}).Where("source = ?", source)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTicketScope(query, queryScope)
		}
	}

	err := applyTicketPreloadsBasic(query).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error

	return tickets, total, err
}

// FindByDepartment récupère les tickets par département du demandeur (avec pagination)
// Le filtrage se fait principalement sur le département du Requester (relation users.department_id).
// Cela permet aux chefs de service de voir les tickets dont le demandeur appartient à leur département.
func (r *ticketRepository) FindByDepartment(departmentID uint, page, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	// Construire la requête de base en joignant la table users sur requester_id
	baseQuery := database.DB.
		Model(&models.Ticket{}).
		Joins("LEFT JOIN users ON users.id = tickets.requester_id").
		Where("users.department_id = ?", departmentID)

	// Compter le total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les tickets avec pagination
	err := applyTicketPreloadsBasic(baseQuery).
		Order("tickets.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&tickets).Error

	return tickets, total, err
}

// Update met à jour un ticket
func (r *ticketRepository) Update(ticket *models.Ticket) error {
	return database.DB.Save(ticket).Error
}

// UpdateFields met à jour uniquement les champs fournis pour un ticket
func (r *ticketRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	return database.DB.Model(&models.Ticket{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateAssignFields met à jour uniquement les champs d'assignation d'un ticket
func (r *ticketRepository) UpdateAssignFields(id uint, assignedToID *uint, estimatedTime *int, status string) error {
	updates := map[string]interface{}{
		"assigned_to_id": assignedToID,
	}
	if estimatedTime != nil {
		updates["estimated_time"] = estimatedTime
	}
	if status != "" {
		updates["status"] = status
	}
	return database.DB.Model(&models.Ticket{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateRequesterNameByCreatedBy met à jour le nom du demandeur pour tous les tickets créés par un utilisateur
func (r *ticketRepository) UpdateRequesterNameByCreatedBy(createdByID uint, requesterName string) error {
	return database.DB.Model(&models.Ticket{}).
		Where("created_by_id = ?", createdByID).
		Update("requester_name", requesterName).Error
}

// UpdateRequesterNameByName met à jour le nom du demandeur pour tous les tickets où le requester_name correspond à l'ancien nom
func (r *ticketRepository) UpdateRequesterNameByName(oldName string, newName string) error {
	return database.DB.Model(&models.Ticket{}).
		Where("requester_name = ?", oldName).
		Update("requester_name", newName).Error
}

// UpdateRequesterNameByRequesterID met à jour le nom du demandeur pour tous les tickets où le requester_id correspond à l'utilisateur
func (r *ticketRepository) UpdateRequesterNameByRequesterID(requesterID uint, requesterName string) error {
	return database.DB.Model(&models.Ticket{}).
		Where("requester_id = ?", requesterID).
		Update("requester_name", requesterName).Error
}

// Delete supprime un ticket (soft delete)
func (r *ticketRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Ticket{}, id).Error
}

// CountByStatus compte les tickets par statut
func (r *ticketRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Ticket{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountByCategory compte les tickets par catégorie
func (r *ticketRepository) CountByCategory(category string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Ticket{}).Where("category = ?", category).Count(&count).Error
	return count, err
}

// Search recherche des tickets par titre ou description
func (r *ticketRepository) Search(scopeParam interface{}, query string, status string, limit int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	searchPattern := "%" + query + "%"

	// Construire la requête de base
	db := applyTicketPreloadsBasic(database.DB.Model(&models.Ticket{})).
		Where("tickets.title LIKE ? OR tickets.description LIKE ?", searchPattern, searchPattern)

	// Appliquer le scope si fourni (doit être fait avant les autres filtres)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			db = scope.ApplyTicketScope(db, queryScope)
		}
	}

	if status != "" {
		db = db.Where("tickets.status = ?", status)
	}

	if limit > 0 {
		db = db.Limit(limit)
	}

	err := db.Order("tickets.created_at DESC").Find(&tickets).Error
	return tickets, err
}

// GetNextSequenceNumber obtient le prochain numéro séquentiel pour une année donnée
// Le format est TKT-YYYY-NNNN, donc on trouve le numéro maximum existant pour cette année
func (r *ticketRepository) GetNextSequenceNumber(year int) (int, error) {
	var tickets []models.Ticket
	codePattern := fmt.Sprintf("TKT-%d-%%", year)

	// Récupérer tous les tickets avec un code correspondant au pattern de l'année
	// Inclure même les tickets supprimés (soft delete) pour éviter les collisions
	err := database.DB.Unscoped().Model(&models.Ticket{}).
		Where("code LIKE ?", codePattern).
		Select("code").
		Find(&tickets).Error

	if err != nil {
		return 0, err
	}

	// Trouver le numéro de séquence maximum
	maxSequence := 0
	for _, ticket := range tickets {
		// Extraire le numéro de séquence du code (format: TKT-YYYY-NNNN)
		if len(ticket.Code) >= 13 {
			var ticketYear, seq int
			if _, err := fmt.Sscanf(ticket.Code, "TKT-%d-%d", &ticketYear, &seq); err == nil && ticketYear == year {
				if seq > maxSequence {
					maxSequence = seq
				}
			}
		}
	}

	// Le prochain numéro est maxSequence + 1
	return maxSequence + 1, nil
}

// CodeExists vérifie si un code de ticket existe déjà (y compris les tickets supprimés)
func (r *ticketRepository) CodeExists(code string) (bool, error) {
	var count int64
	// Utiliser Unscoped() pour inclure les tickets supprimés (soft delete)
	err := database.DB.Unscoped().Model(&models.Ticket{}).
		Where("code = ?", code).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
