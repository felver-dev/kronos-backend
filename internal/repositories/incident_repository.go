package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// IncidentRepository interface pour les opérations sur les incidents
type IncidentRepository interface {
	Create(incident *models.Incident) error
	FindByID(id uint) (*models.Incident, error)
	FindByTicketID(ticketID uint) (*models.Incident, error)
	FindAll() ([]models.Incident, error)
	FindByImpact(impact string) ([]models.Incident, error)
	FindByUrgency(urgency string) ([]models.Incident, error)
	Update(incident *models.Incident) error
	Delete(id uint) error
}

// incidentRepository implémente IncidentRepository
type incidentRepository struct{}

// NewIncidentRepository crée une nouvelle instance de IncidentRepository
func NewIncidentRepository() IncidentRepository {
	return &incidentRepository{}
}

// Create crée un nouvel incident
func (r *incidentRepository) Create(incident *models.Incident) error {
	return database.DB.Create(incident).Error
}

// FindByID trouve un incident par son ID avec son ticket
func (r *incidentRepository) FindByID(id uint) (*models.Incident, error) {
	var incident models.Incident
	err := database.DB.Preload("Ticket").Preload("Ticket.CreatedBy").Preload("Ticket.AssignedTo").First(&incident, id).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// FindByTicketID trouve un incident par l'ID du ticket
func (r *incidentRepository) FindByTicketID(ticketID uint) (*models.Incident, error) {
	var incident models.Incident
	err := database.DB.Preload("Ticket").Preload("Ticket.CreatedBy").Preload("Ticket.AssignedTo").Where("ticket_id = ?", ticketID).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// FindAll récupère tous les incidents avec leurs tickets
func (r *incidentRepository) FindAll() ([]models.Incident, error) {
	var incidents []models.Incident
	err := database.DB.Preload("Ticket").Preload("Ticket.CreatedBy").Preload("Ticket.AssignedTo").Find(&incidents).Error
	return incidents, err
}

// FindByImpact récupère les incidents par impact
func (r *incidentRepository) FindByImpact(impact string) ([]models.Incident, error) {
	var incidents []models.Incident
	err := database.DB.Preload("Ticket").Where("impact = ?", impact).Find(&incidents).Error
	return incidents, err
}

// FindByUrgency récupère les incidents par urgence
func (r *incidentRepository) FindByUrgency(urgency string) ([]models.Incident, error) {
	var incidents []models.Incident
	err := database.DB.Preload("Ticket").Where("urgency = ?", urgency).Find(&incidents).Error
	return incidents, err
}

// Update met à jour un incident
func (r *incidentRepository) Update(incident *models.Incident) error {
	return database.DB.Save(incident).Error
}

// Delete supprime un incident
func (r *incidentRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Incident{}, id).Error
}
