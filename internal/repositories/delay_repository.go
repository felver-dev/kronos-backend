package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// DelayRepository interface pour les opérations sur les retards
type DelayRepository interface {
	Create(delay *models.Delay) error
	FindByID(id uint) (*models.Delay, error)
	FindByTicketID(ticketID uint) (*models.Delay, error)
	FindAll(scope interface{}) ([]models.Delay, error) // scope peut être *scope.QueryScope ou nil
	FindByUserID(scope interface{}, userID uint) ([]models.Delay, error)
	FindByStatus(scope interface{}, status string) ([]models.Delay, error)
	FindUnjustified(scope interface{}) ([]models.Delay, error)
	Update(delay *models.Delay) error
	Delete(id uint) error
}

// DelayJustificationRepository interface pour les opérations sur les justifications de retards
type DelayJustificationRepository interface {
	Create(justification *models.DelayJustification) error
	FindByID(id uint) (*models.DelayJustification, error)
	FindByDelayID(delayID uint) (*models.DelayJustification, error)
	FindAll() ([]models.DelayJustification, error)
	FindByStatus(status string) ([]models.DelayJustification, error)
	FindPending() ([]models.DelayJustification, error)
	FindByUserID(userID uint) ([]models.DelayJustification, error)
	FindValidated() ([]models.DelayJustification, error)
	FindRejected() ([]models.DelayJustification, error)
	Update(justification *models.DelayJustification) error
	Delete(id uint) error
}

// delayRepository implémente DelayRepository
type delayRepository struct{}

// delayJustificationRepository implémente DelayJustificationRepository
type delayJustificationRepository struct{}

// NewDelayRepository crée une nouvelle instance de DelayRepository
func NewDelayRepository() DelayRepository {
	return &delayRepository{}
}

// NewDelayJustificationRepository crée une nouvelle instance de DelayJustificationRepository
func NewDelayJustificationRepository() DelayJustificationRepository {
	return &delayJustificationRepository{}
}

// Create crée un nouveau retard
func (r *delayRepository) Create(delay *models.Delay) error {
	return database.DB.Create(delay).Error
}

// FindByID trouve un retard par son ID
func (r *delayRepository) FindByID(id uint) (*models.Delay, error) {
	var delay models.Delay
	err := database.DB.Preload("Ticket").Preload("User").Preload("Justification").Preload("Justification.User").Preload("Justification.ValidatedBy").First(&delay, id).Error
	if err != nil {
		return nil, err
	}
	return &delay, nil
}

// FindByTicketID trouve un retard par l'ID du ticket
func (r *delayRepository) FindByTicketID(ticketID uint) (*models.Delay, error) {
	var delay models.Delay
	err := database.DB.Preload("Ticket").Preload("User").Preload("Justification").Where("ticket_id = ?", ticketID).First(&delay).Error
	if err != nil {
		return nil, err
	}
	return &delay, nil
}

// FindAll récupère tous les retards
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *delayRepository) FindAll(scopeParam interface{}) ([]models.Delay, error) {
	var delays []models.Delay
	
	// Construire la requête de base
	query := database.DB.Model(&models.Delay{}).
		Preload("Ticket").Preload("User").Preload("Justification")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyDelayScope(query, queryScope)
		}
	}
	
	err := query.Find(&delays).Error
	return delays, err
}

// FindByUserID récupère les retards d'un utilisateur
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *delayRepository) FindByUserID(scopeParam interface{}, userID uint) ([]models.Delay, error) {
	var delays []models.Delay
	
	// Construire la requête de base
	query := database.DB.Model(&models.Delay{}).
		Preload("Ticket").Preload("User").Preload("Justification").
		Where("delays.user_id = ?", userID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyDelayScope(query, queryScope)
		}
	}
	
	err := query.Find(&delays).Error
	return delays, err
}

// FindByStatus récupère les retards par statut
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *delayRepository) FindByStatus(scopeParam interface{}, status string) ([]models.Delay, error) {
	var delays []models.Delay
	
	// Construire la requête de base
	query := database.DB.Model(&models.Delay{}).
		Preload("Ticket").Preload("User").Preload("Justification").
		Where("delays.status = ?", status)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyDelayScope(query, queryScope)
		}
	}
	
	err := query.Find(&delays).Error
	return delays, err
}

// FindUnjustified récupère les retards non justifiés
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *delayRepository) FindUnjustified(scopeParam interface{}) ([]models.Delay, error) {
	var delays []models.Delay
	
	// Construire la requête de base
	query := database.DB.Model(&models.Delay{}).
		Preload("Ticket").Preload("User").Preload("Justification").
		Where("delays.status = ?", "unjustified")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyDelayScope(query, queryScope)
		}
	}
	
	err := query.Find(&delays).Error
	return delays, err
}

// Update met à jour un retard
func (r *delayRepository) Update(delay *models.Delay) error {
	return database.DB.Save(delay).Error
}

// Delete supprime un retard (soft delete)
func (r *delayRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Delay{}, id).Error
}

// Create crée une nouvelle justification de retard
func (r *delayJustificationRepository) Create(justification *models.DelayJustification) error {
	return database.DB.Create(justification).Error
}

// FindByID trouve une justification par son ID
func (r *delayJustificationRepository) FindByID(id uint) (*models.DelayJustification, error) {
	var justification models.DelayJustification
	err := database.DB.Preload("Delay").Preload("Delay.Ticket").Preload("User").Preload("ValidatedBy").First(&justification, id).Error
	if err != nil {
		return nil, err
	}
	return &justification, nil
}

// FindByDelayID trouve une justification par l'ID du retard
func (r *delayJustificationRepository) FindByDelayID(delayID uint) (*models.DelayJustification, error) {
	var justification models.DelayJustification
	err := database.DB.Preload("Delay").Preload("User").Preload("ValidatedBy").Where("delay_id = ?", delayID).First(&justification).Error
	if err != nil {
		return nil, err
	}
	return &justification, nil
}

// FindAll récupère toutes les justifications
func (r *delayJustificationRepository) FindAll() ([]models.DelayJustification, error) {
	var justifications []models.DelayJustification
	err := database.DB.Preload("Delay").Preload("User").Preload("ValidatedBy").Find(&justifications).Error
	return justifications, err
}

// FindByStatus récupère les justifications par statut
func (r *delayJustificationRepository) FindByStatus(status string) ([]models.DelayJustification, error) {
	var justifications []models.DelayJustification
	err := database.DB.Preload("Delay").Preload("User").Preload("ValidatedBy").Where("status = ?", status).Find(&justifications).Error
	return justifications, err
}

// FindPending récupère les justifications en attente de validation
func (r *delayJustificationRepository) FindPending() ([]models.DelayJustification, error) {
	var justifications []models.DelayJustification
	err := database.DB.Preload("Delay").Preload("User").Where("status = ?", "pending").Find(&justifications).Error
	return justifications, err
}

// FindByUserID récupère les justifications d'un utilisateur
func (r *delayJustificationRepository) FindByUserID(userID uint) ([]models.DelayJustification, error) {
	var justifications []models.DelayJustification
	err := database.DB.Preload("Delay").Preload("Delay.Ticket").Preload("User").Preload("ValidatedBy").Where("user_id = ?", userID).Order("created_at DESC").Find(&justifications).Error
	return justifications, err
}

// FindValidated récupère les justifications validées
func (r *delayJustificationRepository) FindValidated() ([]models.DelayJustification, error) {
	var justifications []models.DelayJustification
	err := database.DB.Preload("Delay").Preload("Delay.Ticket").Preload("User").Preload("ValidatedBy").Where("status = ?", "validated").Order("validated_at DESC").Find(&justifications).Error
	return justifications, err
}

// FindRejected récupère les justifications rejetées
func (r *delayJustificationRepository) FindRejected() ([]models.DelayJustification, error) {
	var justifications []models.DelayJustification
	err := database.DB.Preload("Delay").Preload("Delay.Ticket").Preload("User").Preload("ValidatedBy").Where("status = ?", "rejected").Order("validated_at DESC").Find(&justifications).Error
	return justifications, err
}

// Update met à jour une justification
func (r *delayJustificationRepository) Update(justification *models.DelayJustification) error {
	return database.DB.Save(justification).Error
}

// Delete supprime une justification (soft delete)
func (r *delayJustificationRepository) Delete(id uint) error {
	return database.DB.Delete(&models.DelayJustification{}, id).Error
}
