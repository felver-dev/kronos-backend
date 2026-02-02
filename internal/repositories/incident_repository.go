package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// IncidentRepository interface pour les opérations sur les incidents
type IncidentRepository interface {
	Create(incident *models.Incident) error
	FindByID(id uint) (*models.Incident, error)
	FindByTicketID(ticketID uint) (*models.Incident, error)
	FindAll(scope interface{}) ([]models.Incident, error) // scope peut être *scope.QueryScope ou nil
	FindByImpact(scope interface{}, impact string) ([]models.Incident, error)
	FindByUrgency(scope interface{}, urgency string) ([]models.Incident, error)
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
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *incidentRepository) FindAll(scopeParam interface{}) ([]models.Incident, error) {
	var incidents []models.Incident
	
	// Construire la requête de base
	query := database.DB.Model(&models.Incident{}).
		Preload("Ticket").Preload("Ticket.CreatedBy").Preload("Ticket.AssignedTo")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyIncidentScope(query, queryScope)
		}
	}
	
	err := query.Find(&incidents).Error
	return incidents, err
}

// FindByImpact récupère les incidents par impact
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *incidentRepository) FindByImpact(scopeParam interface{}, impact string) ([]models.Incident, error) {
	var incidents []models.Incident
	
	// Construire la requête de base
	query := database.DB.Model(&models.Incident{}).
		Preload("Ticket").
		Where("incidents.impact = ?", impact)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyIncidentScope(query, queryScope)
		}
	}
	
	err := query.Find(&incidents).Error
	return incidents, err
}

// FindByUrgency récupère les incidents par urgence
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *incidentRepository) FindByUrgency(scopeParam interface{}, urgency string) ([]models.Incident, error) {
	var incidents []models.Incident
	
	// Construire la requête de base
	query := database.DB.Model(&models.Incident{}).
		Preload("Ticket").
		Where("incidents.urgency = ?", urgency)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyIncidentScope(query, queryScope)
		}
	}
	
	err := query.Find(&incidents).Error
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
