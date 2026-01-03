package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// ServiceRequestService interface pour les opérations sur les demandes de service
type ServiceRequestService interface {
	Create(req dto.CreateServiceRequestRequest, createdByID uint) (*dto.ServiceRequestDTO, error)
	GetByID(id uint) (*dto.ServiceRequestDTO, error)
	GetByTicketID(ticketID uint) (*dto.ServiceRequestDTO, error)
	GetAll() ([]dto.ServiceRequestDTO, error)
	GetByType(typeID uint) ([]dto.ServiceRequestDTO, error)
	GetValidated() ([]dto.ServiceRequestDTO, error)
	GetPendingValidation() ([]dto.ServiceRequestDTO, error)
	Update(id uint, req dto.UpdateServiceRequestRequest, updatedByID uint) (*dto.ServiceRequestDTO, error)
	Validate(id uint, req dto.ValidateServiceRequestRequest, validatedByID uint) (*dto.ServiceRequestDTO, error)
	Delete(id uint) error
}

// ServiceRequestTypeService interface pour les opérations sur les types de demandes de service
type ServiceRequestTypeService interface {
	Create(req dto.CreateServiceRequestTypeRequest, createdByID uint) (*dto.ServiceRequestTypeDTO, error)
	GetByID(id uint) (*dto.ServiceRequestTypeDTO, error)
	GetAll() ([]dto.ServiceRequestTypeDTO, error)
	GetActive() ([]dto.ServiceRequestTypeDTO, error)
	Update(id uint, req dto.UpdateServiceRequestTypeRequest, updatedByID uint) (*dto.ServiceRequestTypeDTO, error)
	Delete(id uint) error
}

// serviceRequestService implémente ServiceRequestService
type serviceRequestService struct {
	serviceRequestRepo     repositories.ServiceRequestRepository
	serviceRequestTypeRepo repositories.ServiceRequestTypeRepository
	ticketRepo             repositories.TicketRepository
	userRepo               repositories.UserRepository
}

// serviceRequestTypeService implémente ServiceRequestTypeService
type serviceRequestTypeService struct {
	serviceRequestTypeRepo repositories.ServiceRequestTypeRepository
	userRepo               repositories.UserRepository
}

// NewServiceRequestService crée une nouvelle instance de ServiceRequestService
func NewServiceRequestService(
	serviceRequestRepo repositories.ServiceRequestRepository,
	serviceRequestTypeRepo repositories.ServiceRequestTypeRepository,
	ticketRepo repositories.TicketRepository,
	userRepo repositories.UserRepository,
) ServiceRequestService {
	return &serviceRequestService{
		serviceRequestRepo:     serviceRequestRepo,
		serviceRequestTypeRepo: serviceRequestTypeRepo,
		ticketRepo:             ticketRepo,
		userRepo:               userRepo,
	}
}

// NewServiceRequestTypeService crée une nouvelle instance de ServiceRequestTypeService
func NewServiceRequestTypeService(
	serviceRequestTypeRepo repositories.ServiceRequestTypeRepository,
	userRepo repositories.UserRepository,
) ServiceRequestTypeService {
	return &serviceRequestTypeService{
		serviceRequestTypeRepo: serviceRequestTypeRepo,
		userRepo:               userRepo,
	}
}

// Create crée une nouvelle demande de service à partir d'un ticket
func (s *serviceRequestService) Create(req dto.CreateServiceRequestRequest, createdByID uint) (*dto.ServiceRequestDTO, error) {
	// Vérifier que le ticket existe et est de catégorie "demande"
	ticket, err := s.ticketRepo.FindByID(req.TicketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	if ticket.Category != "demande" {
		return nil, errors.New("le ticket doit être de catégorie 'demande'")
	}

	// Vérifier qu'une demande de service n'existe pas déjà pour ce ticket
	existingRequest, _ := s.serviceRequestRepo.FindByTicketID(req.TicketID)
	if existingRequest != nil {
		return nil, errors.New("une demande de service existe déjà pour ce ticket")
	}

	// Vérifier que le type existe
	typeModel, err := s.serviceRequestTypeRepo.FindByID(req.TypeID)
	if err != nil {
		return nil, errors.New("type de demande de service introuvable")
	}

	// Calculer le délai si fourni ou utiliser le délai par défaut du type
	var deadline *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		parsedDeadline, err := time.Parse("2006-01-02", *req.Deadline)
		if err == nil {
			deadline = &parsedDeadline
		}
	}

	// Si pas de deadline fournie, utiliser le délai par défaut du type
	if deadline == nil && typeModel.DefaultDeadline > 0 {
		deadlineTime := time.Now().Add(time.Duration(typeModel.DefaultDeadline) * time.Hour)
		deadline = &deadlineTime
	}

	// Créer la demande de service
	serviceRequest := &models.ServiceRequest{
		TicketID:  req.TicketID,
		TypeID:    req.TypeID,
		Deadline:  deadline,
		Validated: false,
	}

	if err := s.serviceRequestRepo.Create(serviceRequest); err != nil {
		return nil, errors.New("erreur lors de la création de la demande de service")
	}

	// Récupérer la demande créée avec ses relations
	createdRequest, err := s.serviceRequestRepo.FindByID(serviceRequest.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la demande créée")
	}

	// Convertir en DTO
	requestDTO := s.serviceRequestToDTO(createdRequest)
	return &requestDTO, nil
}

// GetByID récupère une demande de service par son ID
func (s *serviceRequestService) GetByID(id uint) (*dto.ServiceRequestDTO, error) {
	serviceRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("demande de service introuvable")
	}

	requestDTO := s.serviceRequestToDTO(serviceRequest)
	return &requestDTO, nil
}

// GetByTicketID récupère une demande de service par l'ID du ticket
func (s *serviceRequestService) GetByTicketID(ticketID uint) (*dto.ServiceRequestDTO, error) {
	serviceRequest, err := s.serviceRequestRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("demande de service introuvable")
	}

	requestDTO := s.serviceRequestToDTO(serviceRequest)
	return &requestDTO, nil
}

// GetAll récupère toutes les demandes de service
func (s *serviceRequestService) GetAll() ([]dto.ServiceRequestDTO, error) {
	serviceRequests, err := s.serviceRequestRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des demandes de service")
	}

	requestDTOs := make([]dto.ServiceRequestDTO, len(serviceRequests))
	for i, req := range serviceRequests {
		requestDTOs[i] = s.serviceRequestToDTO(&req)
	}

	return requestDTOs, nil
}

// GetByType récupère les demandes de service par type
func (s *serviceRequestService) GetByType(typeID uint) ([]dto.ServiceRequestDTO, error) {
	serviceRequests, err := s.serviceRequestRepo.FindByType(typeID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des demandes de service")
	}

	requestDTOs := make([]dto.ServiceRequestDTO, len(serviceRequests))
	for i, req := range serviceRequests {
		requestDTOs[i] = s.serviceRequestToDTO(&req)
	}

	return requestDTOs, nil
}

// GetValidated récupère les demandes de service validées
func (s *serviceRequestService) GetValidated() ([]dto.ServiceRequestDTO, error) {
	serviceRequests, err := s.serviceRequestRepo.FindValidated()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des demandes de service")
	}

	requestDTOs := make([]dto.ServiceRequestDTO, len(serviceRequests))
	for i, req := range serviceRequests {
		requestDTOs[i] = s.serviceRequestToDTO(&req)
	}

	return requestDTOs, nil
}

// GetPendingValidation récupère les demandes de service en attente de validation
func (s *serviceRequestService) GetPendingValidation() ([]dto.ServiceRequestDTO, error) {
	serviceRequests, err := s.serviceRequestRepo.FindPendingValidation()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des demandes de service")
	}

	requestDTOs := make([]dto.ServiceRequestDTO, len(serviceRequests))
	for i, req := range serviceRequests {
		requestDTOs[i] = s.serviceRequestToDTO(&req)
	}

	return requestDTOs, nil
}

// Update met à jour une demande de service
func (s *serviceRequestService) Update(id uint, req dto.UpdateServiceRequestRequest, updatedByID uint) (*dto.ServiceRequestDTO, error) {
	// Récupérer la demande existante
	serviceRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("demande de service introuvable")
	}

	// Mettre à jour les champs fournis
	if req.TypeID != nil {
		// Vérifier que le type existe
		_, err := s.serviceRequestTypeRepo.FindByID(*req.TypeID)
		if err != nil {
			return nil, errors.New("type de demande de service introuvable")
		}
		serviceRequest.TypeID = *req.TypeID
	}

	if req.Deadline != nil && *req.Deadline != "" {
		parsedDeadline, err := time.Parse("2006-01-02", *req.Deadline)
		if err == nil {
			serviceRequest.Deadline = &parsedDeadline
		}
	}

	// Sauvegarder
	if err := s.serviceRequestRepo.Update(serviceRequest); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la demande de service")
	}

	// Récupérer la demande mise à jour
	updatedRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la demande mise à jour")
	}

	requestDTO := s.serviceRequestToDTO(updatedRequest)
	return &requestDTO, nil
}

// Validate valide ou rejette une demande de service
func (s *serviceRequestService) Validate(id uint, req dto.ValidateServiceRequestRequest, validatedByID uint) (*dto.ServiceRequestDTO, error) {
	// Récupérer la demande
	serviceRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("demande de service introuvable")
	}

	// Vérifier que la demande n'est pas déjà validée
	if serviceRequest.Validated && req.Validated {
		return nil, errors.New("la demande de service est déjà validée")
	}

	// Mettre à jour
	serviceRequest.Validated = req.Validated
	serviceRequest.ValidatedByID = &validatedByID
	now := time.Now()
	serviceRequest.ValidatedAt = &now
	serviceRequest.ValidationComment = req.Comment

	// Sauvegarder
	if err := s.serviceRequestRepo.Update(serviceRequest); err != nil {
		return nil, errors.New("erreur lors de la validation de la demande de service")
	}

	// Récupérer la demande mise à jour
	updatedRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la demande mise à jour")
	}

	requestDTO := s.serviceRequestToDTO(updatedRequest)
	return &requestDTO, nil
}

// Delete supprime une demande de service
func (s *serviceRequestService) Delete(id uint) error {
	// Vérifier que la demande existe
	_, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return errors.New("demande de service introuvable")
	}

	// Supprimer
	if err := s.serviceRequestRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la demande de service")
	}

	return nil
}

// Create crée un nouveau type de demande de service
func (s *serviceRequestTypeService) Create(req dto.CreateServiceRequestTypeRequest, createdByID uint) (*dto.ServiceRequestTypeDTO, error) {
	// Créer le type
	serviceRequestType := &models.ServiceRequestType{
		Name:            req.Name,
		Description:     req.Description,
		DefaultDeadline: req.DefaultDeadline,
		IsActive:        true,
		CreatedByID:     &createdByID,
	}

	if err := s.serviceRequestTypeRepo.Create(serviceRequestType); err != nil {
		return nil, errors.New("erreur lors de la création du type de demande de service")
	}

	// Récupérer le type créé
	createdType, err := s.serviceRequestTypeRepo.FindByID(serviceRequestType.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du type créé")
	}

	// Convertir en DTO
	typeDTO := s.serviceRequestTypeToDTO(createdType)
	return &typeDTO, nil
}

// GetByID récupère un type par son ID
func (s *serviceRequestTypeService) GetByID(id uint) (*dto.ServiceRequestTypeDTO, error) {
	serviceRequestType, err := s.serviceRequestTypeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("type de demande de service introuvable")
	}

	typeDTO := s.serviceRequestTypeToDTO(serviceRequestType)
	return &typeDTO, nil
}

// GetAll récupère tous les types
func (s *serviceRequestTypeService) GetAll() ([]dto.ServiceRequestTypeDTO, error) {
	types, err := s.serviceRequestTypeRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des types")
	}

	typeDTOs := make([]dto.ServiceRequestTypeDTO, len(types))
	for i, t := range types {
		typeDTOs[i] = s.serviceRequestTypeToDTO(&t)
	}

	return typeDTOs, nil
}

// GetActive récupère tous les types actifs
func (s *serviceRequestTypeService) GetActive() ([]dto.ServiceRequestTypeDTO, error) {
	types, err := s.serviceRequestTypeRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des types actifs")
	}

	typeDTOs := make([]dto.ServiceRequestTypeDTO, len(types))
	for i, t := range types {
		typeDTOs[i] = s.serviceRequestTypeToDTO(&t)
	}

	return typeDTOs, nil
}

// Update met à jour un type
func (s *serviceRequestTypeService) Update(id uint, req dto.UpdateServiceRequestTypeRequest, updatedByID uint) (*dto.ServiceRequestTypeDTO, error) {
	// Récupérer le type existant
	serviceRequestType, err := s.serviceRequestTypeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("type de demande de service introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Name != "" {
		serviceRequestType.Name = req.Name
	}

	if req.Description != "" {
		serviceRequestType.Description = req.Description
	}

	if req.DefaultDeadline != nil {
		serviceRequestType.DefaultDeadline = *req.DefaultDeadline
	}

	if req.IsActive != nil {
		serviceRequestType.IsActive = *req.IsActive
	}

	// Sauvegarder
	if err := s.serviceRequestTypeRepo.Update(serviceRequestType); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du type")
	}

	// Récupérer le type mis à jour
	updatedType, err := s.serviceRequestTypeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du type mis à jour")
	}

	typeDTO := s.serviceRequestTypeToDTO(updatedType)
	return &typeDTO, nil
}

// Delete supprime un type
func (s *serviceRequestTypeService) Delete(id uint) error {
	// Vérifier que le type existe
	_, err := s.serviceRequestTypeRepo.FindByID(id)
	if err != nil {
		return errors.New("type de demande de service introuvable")
	}

	// Supprimer
	if err := s.serviceRequestTypeRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du type")
	}

	return nil
}

// serviceRequestToDTO convertit un modèle ServiceRequest en DTO ServiceRequestDTO
func (s *serviceRequestService) serviceRequestToDTO(serviceRequest *models.ServiceRequest) dto.ServiceRequestDTO {
	// Convertir le ticket en DTO
	ticketDTO := s.ticketToDTO(&serviceRequest.Ticket)

	// Convertir le type en DTO
	typeDTO := s.serviceRequestTypeToDTO(&serviceRequest.Type)

	var validatedByID *uint
	if serviceRequest.ValidatedByID != nil {
		validatedByID = serviceRequest.ValidatedByID
	}

	return dto.ServiceRequestDTO{
		ID:                serviceRequest.ID,
		TicketID:          serviceRequest.TicketID,
		Ticket:            &ticketDTO,
		TypeID:            serviceRequest.TypeID,
		Type:              &typeDTO,
		Deadline:          serviceRequest.Deadline,
		Validated:         serviceRequest.Validated,
		ValidatedBy:       validatedByID,
		ValidatedAt:       serviceRequest.ValidatedAt,
		ValidationComment: serviceRequest.ValidationComment,
		CreatedAt:         serviceRequest.CreatedAt,
		UpdatedAt:         serviceRequest.UpdatedAt,
	}
}

// serviceRequestTypeToDTO convertit un modèle ServiceRequestType en DTO ServiceRequestTypeDTO
func (s *serviceRequestService) serviceRequestTypeToDTO(serviceRequestType *models.ServiceRequestType) dto.ServiceRequestTypeDTO {
	return dto.ServiceRequestTypeDTO{
		ID:              serviceRequestType.ID,
		Name:            serviceRequestType.Name,
		Description:     serviceRequestType.Description,
		DefaultDeadline: serviceRequestType.DefaultDeadline,
		IsActive:        serviceRequestType.IsActive,
	}
}

// serviceRequestTypeToDTO convertit un modèle ServiceRequestType en DTO ServiceRequestTypeDTO (pour serviceRequestTypeService)
func (s *serviceRequestTypeService) serviceRequestTypeToDTO(serviceRequestType *models.ServiceRequestType) dto.ServiceRequestTypeDTO {
	return dto.ServiceRequestTypeDTO{
		ID:              serviceRequestType.ID,
		Name:            serviceRequestType.Name,
		Description:     serviceRequestType.Description,
		DefaultDeadline: serviceRequestType.DefaultDeadline,
		IsActive:        serviceRequestType.IsActive,
	}
}

// ticketToDTO convertit un modèle Ticket en DTO TicketDTO (méthode utilitaire)
func (s *serviceRequestService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
	var assignedToDTO *dto.UserDTO
	if ticket.AssignedTo != nil {
		assignedDTO := s.userToDTO(ticket.AssignedTo)
		assignedToDTO = &assignedDTO
	}

	createdByDTO := s.userToDTO(&ticket.CreatedBy)

	return dto.TicketDTO{
		ID:            ticket.ID,
		Title:         ticket.Title,
		Description:   ticket.Description,
		Category:      ticket.Category,
		Source:        ticket.Source,
		Status:        ticket.Status,
		Priority:      ticket.Priority,
		AssignedTo:    assignedToDTO,
		CreatedBy:     createdByDTO,
		EstimatedTime: ticket.EstimatedTime,
		ActualTime:    ticket.ActualTime,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
		ClosedAt:      ticket.ClosedAt,
	}
}

// userToDTO convertit un modèle User en DTO UserDTO (méthode utilitaire)
func (s *serviceRequestService) userToDTO(user *models.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role.Name,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
