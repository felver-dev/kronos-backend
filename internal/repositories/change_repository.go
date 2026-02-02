package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// ChangeRepository interface pour les opérations sur les changements
type ChangeRepository interface {
	Create(change *models.Change) error
	FindByID(id uint) (*models.Change, error)
	FindByTicketID(ticketID uint) (*models.Change, error)
	FindAll(scope interface{}) ([]models.Change, error) // scope peut être *scope.QueryScope ou nil
	FindByRisk(scope interface{}, risk string) ([]models.Change, error)
	FindByResponsible(scope interface{}, responsibleID uint) ([]models.Change, error) // scope peut être *scope.QueryScope ou nil
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
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *changeRepository) FindAll(scopeParam interface{}) ([]models.Change, error) {
	var changes []models.Change
	
	// Construire la requête de base
	query := database.DB.Model(&models.Change{}).
		Preload("Ticket").Preload("Responsible")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyChangeScope(query, queryScope)
		}
	}
	
	err := query.Find(&changes).Error
	return changes, err
}

// FindByRisk récupère les changements par niveau de risque
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *changeRepository) FindByRisk(scopeParam interface{}, risk string) ([]models.Change, error) {
	var changes []models.Change
	
	// Construire la requête de base
	query := database.DB.Model(&models.Change{}).
		Preload("Ticket").Preload("Responsible").
		Where("changes.risk = ?", risk)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyChangeScope(query, queryScope)
		}
	}
	
	err := query.Find(&changes).Error
	return changes, err
}

// FindByResponsible récupère les changements par responsable
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *changeRepository) FindByResponsible(scopeParam interface{}, responsibleID uint) ([]models.Change, error) {
	var changes []models.Change
	
	// Construire la requête de base
	query := database.DB.Model(&models.Change{}).
		Preload("Ticket").Preload("Responsible").
		Where("changes.responsible_id = ?", responsibleID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyChangeScope(query, queryScope)
		}
	}
	
	err := query.Find(&changes).Error
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
