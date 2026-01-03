package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// ServiceRequestRepository interface pour les opérations sur les demandes de service
type ServiceRequestRepository interface {
	Create(serviceRequest *models.ServiceRequest) error
	FindByID(id uint) (*models.ServiceRequest, error)
	FindByTicketID(ticketID uint) (*models.ServiceRequest, error)
	FindAll() ([]models.ServiceRequest, error)
	FindByType(typeID uint) ([]models.ServiceRequest, error)
	FindValidated() ([]models.ServiceRequest, error)
	FindPendingValidation() ([]models.ServiceRequest, error)
	Update(serviceRequest *models.ServiceRequest) error
	Delete(id uint) error
}

// ServiceRequestTypeRepository interface pour les opérations sur les types de demandes de service
type ServiceRequestTypeRepository interface {
	Create(serviceRequestType *models.ServiceRequestType) error
	FindByID(id uint) (*models.ServiceRequestType, error)
	FindAll() ([]models.ServiceRequestType, error)
	FindActive() ([]models.ServiceRequestType, error)
	Update(serviceRequestType *models.ServiceRequestType) error
	Delete(id uint) error
}

// serviceRequestRepository implémente ServiceRequestRepository
type serviceRequestRepository struct{}

// serviceRequestTypeRepository implémente ServiceRequestTypeRepository
type serviceRequestTypeRepository struct{}

// NewServiceRequestRepository crée une nouvelle instance de ServiceRequestRepository
func NewServiceRequestRepository() ServiceRequestRepository {
	return &serviceRequestRepository{}
}

// NewServiceRequestTypeRepository crée une nouvelle instance de ServiceRequestTypeRepository
func NewServiceRequestTypeRepository() ServiceRequestTypeRepository {
	return &serviceRequestTypeRepository{}
}

// Create crée une nouvelle demande de service
func (r *serviceRequestRepository) Create(serviceRequest *models.ServiceRequest) error {
	return database.DB.Create(serviceRequest).Error
}

// FindByID trouve une demande de service par son ID
func (r *serviceRequestRepository) FindByID(id uint) (*models.ServiceRequest, error) {
	var serviceRequest models.ServiceRequest
	err := database.DB.Preload("Ticket").Preload("Ticket.CreatedBy").Preload("Ticket.AssignedTo").Preload("Type").Preload("Validator").First(&serviceRequest, id).Error
	if err != nil {
		return nil, err
	}
	return &serviceRequest, nil
}

// FindByTicketID trouve une demande de service par l'ID du ticket
func (r *serviceRequestRepository) FindByTicketID(ticketID uint) (*models.ServiceRequest, error) {
	var serviceRequest models.ServiceRequest
	err := database.DB.Preload("Ticket").Preload("Type").Preload("Validator").Where("ticket_id = ?", ticketID).First(&serviceRequest).Error
	if err != nil {
		return nil, err
	}
	return &serviceRequest, nil
}

// FindAll récupère toutes les demandes de service
func (r *serviceRequestRepository) FindAll() ([]models.ServiceRequest, error) {
	var serviceRequests []models.ServiceRequest
	err := database.DB.Preload("Ticket").Preload("Type").Preload("Validator").Find(&serviceRequests).Error
	return serviceRequests, err
}

// FindByType récupère les demandes de service par type
func (r *serviceRequestRepository) FindByType(typeID uint) ([]models.ServiceRequest, error) {
	var serviceRequests []models.ServiceRequest
	err := database.DB.Preload("Ticket").Preload("Type").Where("type_id = ?", typeID).Find(&serviceRequests).Error
	return serviceRequests, err
}

// FindValidated récupère les demandes de service validées
func (r *serviceRequestRepository) FindValidated() ([]models.ServiceRequest, error) {
	var serviceRequests []models.ServiceRequest
	err := database.DB.Preload("Ticket").Preload("Type").Preload("Validator").Where("validated = ?", true).Find(&serviceRequests).Error
	return serviceRequests, err
}

// FindPendingValidation récupère les demandes de service en attente de validation
func (r *serviceRequestRepository) FindPendingValidation() ([]models.ServiceRequest, error) {
	var serviceRequests []models.ServiceRequest
	err := database.DB.Preload("Ticket").Preload("Type").Where("validated = ?", false).Find(&serviceRequests).Error
	return serviceRequests, err
}

// Update met à jour une demande de service
func (r *serviceRequestRepository) Update(serviceRequest *models.ServiceRequest) error {
	return database.DB.Save(serviceRequest).Error
}

// Delete supprime une demande de service
func (r *serviceRequestRepository) Delete(id uint) error {
	return database.DB.Delete(&models.ServiceRequest{}, id).Error
}

// Create crée un nouveau type de demande de service
func (r *serviceRequestTypeRepository) Create(serviceRequestType *models.ServiceRequestType) error {
	return database.DB.Create(serviceRequestType).Error
}

// FindByID trouve un type de demande de service par son ID
func (r *serviceRequestTypeRepository) FindByID(id uint) (*models.ServiceRequestType, error) {
	var serviceRequestType models.ServiceRequestType
	err := database.DB.First(&serviceRequestType, id).Error
	if err != nil {
		return nil, err
	}
	return &serviceRequestType, nil
}

// FindAll récupère tous les types de demandes de service
func (r *serviceRequestTypeRepository) FindAll() ([]models.ServiceRequestType, error) {
	var serviceRequestTypes []models.ServiceRequestType
	err := database.DB.Find(&serviceRequestTypes).Error
	return serviceRequestTypes, err
}

// FindActive récupère tous les types de demandes de service actifs
func (r *serviceRequestTypeRepository) FindActive() ([]models.ServiceRequestType, error) {
	var serviceRequestTypes []models.ServiceRequestType
	err := database.DB.Where("is_active = ?", true).Find(&serviceRequestTypes).Error
	return serviceRequestTypes, err
}

// Update met à jour un type de demande de service
func (r *serviceRequestTypeRepository) Update(serviceRequestType *models.ServiceRequestType) error {
	return database.DB.Save(serviceRequestType).Error
}

// Delete supprime un type de demande de service
func (r *serviceRequestTypeRepository) Delete(id uint) error {
	return database.DB.Delete(&models.ServiceRequestType{}, id).Error
}
