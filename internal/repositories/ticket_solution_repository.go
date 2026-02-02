package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// TicketSolutionRepository interface pour les opérations sur les solutions de tickets
type TicketSolutionRepository interface {
	Create(solution *models.TicketSolution) error
	FindByID(id uint) (*models.TicketSolution, error)
	FindByTicketID(ticketID uint) ([]models.TicketSolution, error)
	Update(solution *models.TicketSolution) error
	Delete(id uint) error
}

// ticketSolutionRepository implémente TicketSolutionRepository
type ticketSolutionRepository struct{}

// NewTicketSolutionRepository crée une nouvelle instance de TicketSolutionRepository
func NewTicketSolutionRepository() TicketSolutionRepository {
	return &ticketSolutionRepository{}
}

// Create crée une nouvelle solution
func (r *ticketSolutionRepository) Create(solution *models.TicketSolution) error {
	return database.DB.Create(solution).Error
}

// FindByID trouve une solution par son ID
func (r *ticketSolutionRepository) FindByID(id uint) (*models.TicketSolution, error) {
	var solution models.TicketSolution
	err := database.DB.Preload("CreatedBy").Preload("Ticket").First(&solution, id).Error
	if err != nil {
		return nil, err
	}
	return &solution, nil
}

// FindByTicketID trouve toutes les solutions d'un ticket
func (r *ticketSolutionRepository) FindByTicketID(ticketID uint) ([]models.TicketSolution, error) {
	var solutions []models.TicketSolution
	err := database.DB.Preload("CreatedBy").Where("ticket_id = ?", ticketID).Order("created_at DESC").Find(&solutions).Error
	return solutions, err
}

// Update met à jour une solution
func (r *ticketSolutionRepository) Update(solution *models.TicketSolution) error {
	return database.DB.Save(solution).Error
}

// Delete supprime une solution
func (r *ticketSolutionRepository) Delete(id uint) error {
	return database.DB.Delete(&models.TicketSolution{}, id).Error
}
