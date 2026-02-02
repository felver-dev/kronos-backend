package repositories

import (
	"strings"
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// TimeEntryRepository interface pour les opérations sur les entrées de temps
type TimeEntryRepository interface {
	Create(timeEntry *models.TimeEntry) error
	FindByID(id uint) (*models.TimeEntry, error)
	FindAll(scope interface{}) ([]models.TimeEntry, error) // scope peut être *scope.QueryScope ou nil
	FindByTicketID(ticketID uint) ([]models.TimeEntry, error)
	FindByUserID(userID uint) ([]models.TimeEntry, error)
	FindByDateRange(userID uint, startDate, endDate time.Time) ([]models.TimeEntry, error)
	FindValidated(scope interface{}) ([]models.TimeEntry, error)
	FindPendingValidation(scope interface{}) ([]models.TimeEntry, error)
	Search(scope interface{}, query string, limit int) ([]models.TimeEntry, error)
	Update(timeEntry *models.TimeEntry) error
	Delete(id uint) error
	SumByTicketID(ticketID uint) (int, error)
	SumByUserID(userID uint) (int, error)
	// ValidateByTicketID marque comme validées toutes les entrées de temps non encore validées du ticket (ex. après validation du ticket par le demandeur)
	ValidateByTicketID(ticketID uint, validatedByID uint) error
}

// timeEntryRepository implémente TimeEntryRepository
type timeEntryRepository struct{}

// NewTimeEntryRepository crée une nouvelle instance de TimeEntryRepository
func NewTimeEntryRepository() TimeEntryRepository {
	return &timeEntryRepository{}
}

// Create crée une nouvelle entrée de temps
func (r *timeEntryRepository) Create(timeEntry *models.TimeEntry) error {
	return database.DB.Create(timeEntry).Error
}

// FindByID trouve une entrée de temps par son ID
func (r *timeEntryRepository) FindByID(id uint) (*models.TimeEntry, error) {
	var timeEntry models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("User").Preload("ValidatedBy").First(&timeEntry, id).Error
	if err != nil {
		return nil, err
	}
	return &timeEntry, nil
}

// FindAll récupère toutes les entrées de temps
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *timeEntryRepository) FindAll(scopeParam interface{}) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry

	// Construire la requête de base
	query := database.DB.Model(&models.TimeEntry{}).
		Preload("Ticket").
		Preload("User").
		Preload("ValidatedBy")

	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTimeEntryScope(query, queryScope)
		}
	}

	err := query.Order("time_entries.created_at DESC").Find(&timeEntries).Error
	if err != nil {
		return nil, err
	}
	return timeEntries, nil
}

// FindByTicketID récupère les entrées de temps d'un ticket
func (r *timeEntryRepository) FindByTicketID(ticketID uint) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("User").Preload("ValidatedBy").Where("ticket_id = ?", ticketID).Find(&timeEntries).Error
	return timeEntries, err
}

// FindByUserID récupère les entrées de temps d'un utilisateur
func (r *timeEntryRepository) FindByUserID(userID uint) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("ValidatedBy").Where("user_id = ?", userID).Find(&timeEntries).Error
	return timeEntries, err
}

// FindByDateRange récupère les entrées de temps d'un utilisateur dans une plage de dates
func (r *timeEntryRepository) FindByDateRange(userID uint, startDate, endDate time.Time) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("ValidatedBy").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Find(&timeEntries).Error
	return timeEntries, err
}

// FindValidated récupère les entrées de temps validées
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *timeEntryRepository) FindValidated(scopeParam interface{}) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry

	// Construire la requête de base
	query := database.DB.Model(&models.TimeEntry{}).
		Preload("Ticket").Preload("User").Preload("ValidatedBy").
		Where("time_entries.validated = ?", true)

	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTimeEntryScope(query, queryScope)
		}
	}

	err := query.Find(&timeEntries).Error
	return timeEntries, err
}

// FindPendingValidation récupère les entrées de temps en attente de validation.
// Utilise ApplyTimeEntryScopeForPendingValidation pour que les validateurs (timesheet.validate)
// voient les entrées de leur département ou toutes si validateur sans département.
func (r *timeEntryRepository) FindPendingValidation(scopeParam interface{}) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry

	query := database.DB.Model(&models.TimeEntry{}).
		Preload("Ticket").Preload("User").
		Where("time_entries.validated = ?", false)

	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTimeEntryScopeForPendingValidation(query, queryScope)
		}
	}

	err := query.Find(&timeEntries).Error
	return timeEntries, err
}

// Search recherche des entrées de temps par description, ticket ou utilisateur
func (r *timeEntryRepository) Search(scopeParam interface{}, searchQuery string, limit int) ([]models.TimeEntry, error) {
	if limit <= 0 {
		limit = 20
	}
	like := "%" + strings.ToLower(searchQuery) + "%"

	// Construire la requête de base
	query := database.DB.Model(&models.TimeEntry{}).
		Select("time_entries.*").
		Joins("LEFT JOIN tickets ON tickets.id = time_entries.ticket_id").
		Joins("LEFT JOIN users ON users.id = time_entries.user_id").
		Where(
			"LOWER(time_entries.description) LIKE ? OR LOWER(tickets.code) LIKE ? OR LOWER(tickets.title) LIKE ? OR LOWER(users.username) LIKE ? OR LOWER(users.first_name) LIKE ? OR LOWER(users.last_name) LIKE ?",
			like, like, like, like, like, like,
		).
		Preload("Ticket").
		Preload("User").
		Preload("ValidatedBy")

	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyTimeEntryScope(query, queryScope)
		}
	}

	var timeEntries []models.TimeEntry
	err := query.Order("time_entries.created_at DESC").
		Limit(limit).
		Find(&timeEntries).Error
	return timeEntries, err
}

// Update met à jour une entrée de temps
func (r *timeEntryRepository) Update(timeEntry *models.TimeEntry) error {
	return database.DB.Save(timeEntry).Error
}

// Delete supprime une entrée de temps (soft delete)
func (r *timeEntryRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TimeEntry{}, id).Error
}

// SumByTicketID calcule la somme des temps passés sur un ticket
func (r *timeEntryRepository) SumByTicketID(ticketID uint) (int, error) {
	var total int
	err := database.DB.Model(&models.TimeEntry{}).
		Where("ticket_id = ?", ticketID).
		Select("COALESCE(SUM(time_spent), 0)").
		Scan(&total).Error
	return total, err
}

// SumByUserID calcule la somme des temps passés par un utilisateur
func (r *timeEntryRepository) SumByUserID(userID uint) (int, error) {
	var total int
	err := database.DB.Model(&models.TimeEntry{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(time_spent), 0)").
		Scan(&total).Error
	return total, err
}

// ValidateByTicketID marque comme validées toutes les entrées de temps non encore validées du ticket
func (r *timeEntryRepository) ValidateByTicketID(ticketID uint, validatedByID uint) error {
	now := time.Now()
	return database.DB.Model(&models.TimeEntry{}).
		Where("ticket_id = ? AND validated = ?", ticketID, false).
		Updates(map[string]interface{}{
			"validated":       true,
			"validated_by_id": validatedByID,
			"validated_at":    now,
			"updated_at":      now,
		}).Error
}
