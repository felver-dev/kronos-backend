package repositories

import (
	"time"

	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TimeEntryRepository interface pour les opérations sur les entrées de temps
type TimeEntryRepository interface {
	Create(timeEntry *models.TimeEntry) error
	FindByID(id uint) (*models.TimeEntry, error)
	FindAll() ([]models.TimeEntry, error)
	FindByTicketID(ticketID uint) ([]models.TimeEntry, error)
	FindByUserID(userID uint) ([]models.TimeEntry, error)
	FindByDateRange(userID uint, startDate, endDate time.Time) ([]models.TimeEntry, error)
	FindValidated() ([]models.TimeEntry, error)
	FindPendingValidation() ([]models.TimeEntry, error)
	Update(timeEntry *models.TimeEntry) error
	Delete(id uint) error
	SumByTicketID(ticketID uint) (int, error)
	SumByUserID(userID uint) (int, error)
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
	err := database.DB.Preload("Ticket").Preload("User").Preload("Validator").First(&timeEntry, id).Error
	if err != nil {
		return nil, err
	}
	return &timeEntry, nil
}

// FindAll récupère toutes les entrées de temps
func (r *timeEntryRepository) FindAll() ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("User").Preload("Validator").Find(&timeEntries).Error
	return timeEntries, err
}

// FindByTicketID récupère les entrées de temps d'un ticket
func (r *timeEntryRepository) FindByTicketID(ticketID uint) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("User").Preload("Validator").Where("ticket_id = ?", ticketID).Find(&timeEntries).Error
	return timeEntries, err
}

// FindByUserID récupère les entrées de temps d'un utilisateur
func (r *timeEntryRepository) FindByUserID(userID uint) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("Validator").Where("user_id = ?", userID).Find(&timeEntries).Error
	return timeEntries, err
}

// FindByDateRange récupère les entrées de temps d'un utilisateur dans une plage de dates
func (r *timeEntryRepository) FindByDateRange(userID uint, startDate, endDate time.Time) ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("Validator").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Find(&timeEntries).Error
	return timeEntries, err
}

// FindValidated récupère les entrées de temps validées
func (r *timeEntryRepository) FindValidated() ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("User").Preload("Validator").Where("validated = ?", true).Find(&timeEntries).Error
	return timeEntries, err
}

// FindPendingValidation récupère les entrées de temps en attente de validation
func (r *timeEntryRepository) FindPendingValidation() ([]models.TimeEntry, error) {
	var timeEntries []models.TimeEntry
	err := database.DB.Preload("Ticket").Preload("User").Where("validated = ?", false).Find(&timeEntries).Error
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

