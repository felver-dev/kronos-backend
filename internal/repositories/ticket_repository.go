package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TicketRepository interface pour les opérations sur les tickets
type TicketRepository interface {
	Create(ticket *models.Ticket) error
	FindByID(id uint) (*models.Ticket, error)
	FindAll() ([]models.Ticket, error)
	FindByStatus(status string) ([]models.Ticket, error)
	FindByCategory(category string) ([]models.Ticket, error)
	FindByPriority(priority string) ([]models.Ticket, error)
	FindByAssignedTo(userID uint) ([]models.Ticket, error)
	FindByCreatedBy(userID uint) ([]models.Ticket, error)
	Update(ticket *models.Ticket) error
	Delete(id uint) error
	CountByStatus(status string) (int64, error)
	CountByCategory(category string) (int64, error)
}

// ticketRepository implémente TicketRepository
type ticketRepository struct{}

// NewTicketRepository crée une nouvelle instance de TicketRepository
func NewTicketRepository() TicketRepository {
	return &ticketRepository{}
}

// Create crée un nouveau ticket
func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return database.DB.Create(ticket).Error
}

// FindByID trouve un ticket par son ID avec ses relations
func (r *ticketRepository) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("CreatedBy.Role").Preload("AssignedTo").Preload("AssignedTo.Role").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// FindAll récupère tous les tickets avec leurs relations
func (r *ticketRepository) FindAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("AssignedTo").Find(&tickets).Error
	return tickets, err
}

// FindByStatus récupère les tickets par statut
func (r *ticketRepository) FindByStatus(status string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("AssignedTo").Where("status = ?", status).Find(&tickets).Error
	return tickets, err
}

// FindByCategory récupère les tickets par catégorie
func (r *ticketRepository) FindByCategory(category string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("AssignedTo").Where("category = ?", category).Find(&tickets).Error
	return tickets, err
}

// FindByPriority récupère les tickets par priorité
func (r *ticketRepository) FindByPriority(priority string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("AssignedTo").Where("priority = ?", priority).Find(&tickets).Error
	return tickets, err
}

// FindByAssignedTo récupère les tickets assignés à un utilisateur
func (r *ticketRepository) FindByAssignedTo(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("AssignedTo").Where("assigned_to_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

// FindByCreatedBy récupère les tickets créés par un utilisateur
func (r *ticketRepository) FindByCreatedBy(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := database.DB.Preload("CreatedBy").Preload("AssignedTo").Where("created_by_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

// Update met à jour un ticket
func (r *ticketRepository) Update(ticket *models.Ticket) error {
	return database.DB.Save(ticket).Error
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
