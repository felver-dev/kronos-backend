package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ChangeRepository interface pour les opérations sur les changements
type ChangeRepository interface {
	Create(change *models.Change) error
	FindByID(id uint) (*models.Change, error)
	FindByTicketID(ticketID uint) (*models.Change, error)
	FindAll() ([]models.Change, error)
	FindByRisk(risk string) ([]models.Change, error)
	FindByResponsible(responsibleID uint) ([]models.Change, error)
	Update(change *models.Change) error
	Delete(id uint) error
}

// changeRepository implémente ChangeRepository
type changeRepository struct{}

// NewChangeRepository crée une nouvelle instance de ChangeRepository
func NewChangeRepository() ChangeRepository {
	return &changeRepository{}
}

// Create crée un nouveau changement
func (r *changeRepository) Create(change *models.Change) error {
	return database.DB.Create(change).Error
}

// FindByID trouve un changement par son ID
func (r *changeRepository) FindByID(id uint) (*models.Change, error) {
	var change models.Change
	err := database.DB.Preload("Ticket").Preload("Ticket.CreatedBy").Preload("Ticket.AssignedTo").Preload("Responsible").First(&change, id).Error
	if err != nil {
		return nil, err
	}
	return &change, nil
}

// FindByTicketID trouve un changement par l'ID du ticket
func (r *changeRepository) FindByTicketID(ticketID uint) (*models.Change, error) {
	var change models.Change
	err := database.DB.Preload("Ticket").Preload("Responsible").Where("ticket_id = ?", ticketID).First(&change).Error
	if err != nil {
		return nil, err
	}
	return &change, nil
}

// FindAll récupère tous les changements
func (r *changeRepository) FindAll() ([]models.Change, error) {
	var changes []models.Change
	err := database.DB.Preload("Ticket").Preload("Responsible").Find(&changes).Error
	return changes, err
}

// FindByRisk récupère les changements par niveau de risque
func (r *changeRepository) FindByRisk(risk string) ([]models.Change, error) {
	var changes []models.Change
	err := database.DB.Preload("Ticket").Preload("Responsible").Where("risk = ?", risk).Find(&changes).Error
	return changes, err
}

// FindByResponsible récupère les changements par responsable
func (r *changeRepository) FindByResponsible(responsibleID uint) ([]models.Change, error) {
	var changes []models.Change
	err := database.DB.Preload("Ticket").Preload("Responsible").Where("responsible_id = ?", responsibleID).Find(&changes).Error
	return changes, err
}

// Update met à jour un changement
func (r *changeRepository) Update(change *models.Change) error {
	return database.DB.Save(change).Error
}

// Delete supprime un changement
func (r *changeRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Change{}, id).Error
}
