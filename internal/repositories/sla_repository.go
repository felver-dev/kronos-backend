package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// SLARepository interface pour les opérations sur les SLA
type SLARepository interface {
	Create(sla *models.SLA) error
	FindByID(id uint) (*models.SLA, error)
	FindAll() ([]models.SLA, error)
	FindActive() ([]models.SLA, error)
	FindByCategory(category string) ([]models.SLA, error)
	FindByCategoryAndPriority(category, priority string) (*models.SLA, error)
	Update(sla *models.SLA) error
	Delete(id uint) error
}

// TicketSLARepository interface pour les opérations sur les associations ticket-SLA
type TicketSLARepository interface {
	Create(ticketSLA *models.TicketSLA) error
	FindByID(id uint) (*models.TicketSLA, error)
	FindByTicketID(ticketID uint) (*models.TicketSLA, error)
	FindBySLAID(slaID uint) ([]models.TicketSLA, error)
	FindAll() ([]models.TicketSLA, error)
	FindByStatus(status string) ([]models.TicketSLA, error)
	FindViolated() ([]models.TicketSLA, error)
	Update(ticketSLA *models.TicketSLA) error
	Delete(id uint) error
}

// slaRepository implémente SLARepository
type slaRepository struct{}

// ticketSLARepository implémente TicketSLARepository
type ticketSLARepository struct{}

// NewSLARepository crée une nouvelle instance de SLARepository
func NewSLARepository() SLARepository {
	return &slaRepository{}
}

// NewTicketSLARepository crée une nouvelle instance de TicketSLARepository
func NewTicketSLARepository() TicketSLARepository {
	return &ticketSLARepository{}
}

// Create crée un nouveau SLA
func (r *slaRepository) Create(sla *models.SLA) error {
	return database.DB.Create(sla).Error
}

// FindByID trouve un SLA par son ID
func (r *slaRepository) FindByID(id uint) (*models.SLA, error) {
	var sla models.SLA
	err := database.DB.Preload("CreatedBy").First(&sla, id).Error
	if err != nil {
		return nil, err
	}
	return &sla, nil
}

// FindAll récupère tous les SLA
func (r *slaRepository) FindAll() ([]models.SLA, error) {
	var slas []models.SLA
	err := database.DB.Preload("CreatedBy").Find(&slas).Error
	return slas, err
}

// FindActive récupère tous les SLA actifs
func (r *slaRepository) FindActive() ([]models.SLA, error) {
	var slas []models.SLA
	err := database.DB.Preload("CreatedBy").Where("is_active = ?", true).Find(&slas).Error
	return slas, err
}

// FindByCategory récupère les SLA par catégorie de ticket
func (r *slaRepository) FindByCategory(category string) ([]models.SLA, error) {
	var slas []models.SLA
	err := database.DB.Preload("CreatedBy").Where("ticket_category = ? AND is_active = ?", category, true).Find(&slas).Error
	return slas, err
}

// FindByCategoryAndPriority trouve un SLA par catégorie et priorité
func (r *slaRepository) FindByCategoryAndPriority(category, priority string) (*models.SLA, error) {
	var sla models.SLA
	query := database.DB.Where("ticket_category = ? AND is_active = ?", category, true)
	if priority != "" {
		query = query.Where("priority = ?", priority)
	} else {
		query = query.Where("priority IS NULL")
	}
	err := query.Preload("CreatedBy").First(&sla).Error
	if err != nil {
		return nil, err
	}
	return &sla, nil
}

// Update met à jour un SLA
func (r *slaRepository) Update(sla *models.SLA) error {
	return database.DB.Save(sla).Error
}

// Delete supprime un SLA
func (r *slaRepository) Delete(id uint) error {
	return database.DB.Delete(&models.SLA{}, id).Error
}

// Create crée une nouvelle association ticket-SLA
func (r *ticketSLARepository) Create(ticketSLA *models.TicketSLA) error {
	return database.DB.Create(ticketSLA).Error
}

// FindByID trouve une association ticket-SLA par son ID
func (r *ticketSLARepository) FindByID(id uint) (*models.TicketSLA, error) {
	var ticketSLA models.TicketSLA
	err := database.DB.Preload("Ticket").Preload("SLA").First(&ticketSLA, id).Error
	if err != nil {
		return nil, err
	}
	return &ticketSLA, nil
}

// FindByTicketID trouve une association ticket-SLA par l'ID du ticket
func (r *ticketSLARepository) FindByTicketID(ticketID uint) (*models.TicketSLA, error) {
	var ticketSLA models.TicketSLA
	err := database.DB.Preload("Ticket").Preload("SLA").Where("ticket_id = ?", ticketID).First(&ticketSLA).Error
	if err != nil {
		return nil, err
	}
	return &ticketSLA, nil
}

// FindBySLAID trouve toutes les associations ticket-SLA par l'ID du SLA
func (r *ticketSLARepository) FindBySLAID(slaID uint) ([]models.TicketSLA, error) {
	var ticketSLAs []models.TicketSLA
	err := database.DB.Preload("Ticket").Preload("SLA").Where("sla_id = ?", slaID).Find(&ticketSLAs).Error
	return ticketSLAs, err
}

// FindAll récupère toutes les associations ticket-SLA
func (r *ticketSLARepository) FindAll() ([]models.TicketSLA, error) {
	var ticketSLAs []models.TicketSLA
	err := database.DB.Preload("Ticket").Preload("SLA").Find(&ticketSLAs).Error
	return ticketSLAs, err
}

// FindByStatus récupère les associations ticket-SLA par statut
func (r *ticketSLARepository) FindByStatus(status string) ([]models.TicketSLA, error) {
	var ticketSLAs []models.TicketSLA
	err := database.DB.Preload("Ticket").Preload("SLA").Where("status = ?", status).Find(&ticketSLAs).Error
	return ticketSLAs, err
}

// FindViolated récupère les associations ticket-SLA violées
func (r *ticketSLARepository) FindViolated() ([]models.TicketSLA, error) {
	var ticketSLAs []models.TicketSLA
	err := database.DB.Preload("Ticket").Preload("SLA").Where("status = ?", "violated").Find(&ticketSLAs).Error
	return ticketSLAs, err
}

// Update met à jour une association ticket-SLA
func (r *ticketSLARepository) Update(ticketSLA *models.TicketSLA) error {
	return database.DB.Save(ticketSLA).Error
}

// Delete supprime une association ticket-SLA
func (r *ticketSLARepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketSLA{}, id).Error
}
