package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// IncidentService interface pour les opérations sur les incidents
type IncidentService interface {
	Create(req dto.CreateIncidentRequest, createdByID uint) (*dto.IncidentDTO, error)
	GetByID(id uint) (*dto.IncidentDTO, error)
	GetByTicketID(ticketID uint) (*dto.IncidentDTO, error)
	GetAll() ([]dto.IncidentDTO, error)
	GetByImpact(impact string) ([]dto.IncidentDTO, error)
	GetByUrgency(urgency string) ([]dto.IncidentDTO, error)
	Update(id uint, req dto.UpdateIncidentRequest, updatedByID uint) (*dto.IncidentDTO, error)
	Qualify(id uint, req dto.QualifyIncidentRequest, qualifiedByID uint) (*dto.IncidentDTO, error)
	LinkAsset(incidentID uint, assetID uint, linkedByID uint) error
	UnlinkAsset(incidentID uint, assetID uint) error
	Resolve(id uint, resolvedByID uint) (*dto.IncidentDTO, error)
	Delete(id uint) error
}

// incidentService implémente IncidentService
type incidentService struct {
	incidentRepo    repositories.IncidentRepository
	ticketRepo      repositories.TicketRepository
	ticketAssetRepo repositories.TicketAssetRepository
	assetRepo       repositories.AssetRepository
}

// NewIncidentService crée une nouvelle instance de IncidentService
func NewIncidentService(
	incidentRepo repositories.IncidentRepository,
	ticketRepo repositories.TicketRepository,
	ticketAssetRepo repositories.TicketAssetRepository,
	assetRepo repositories.AssetRepository,
) IncidentService {
	return &incidentService{
		incidentRepo:    incidentRepo,
		ticketRepo:      ticketRepo,
		ticketAssetRepo: ticketAssetRepo,
		assetRepo:       assetRepo,
	}
}

// Create crée un nouvel incident à partir d'un ticket
func (s *incidentService) Create(req dto.CreateIncidentRequest, createdByID uint) (*dto.IncidentDTO, error) {
	// Vérifier que le ticket existe et est de catégorie "incident"
	ticket, err := s.ticketRepo.FindByID(req.TicketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	if ticket.Category != "incident" {
		return nil, errors.New("le ticket doit être de catégorie 'incident'")
	}

	// Vérifier qu'un incident n'existe pas déjà pour ce ticket
	existingIncident, _ := s.incidentRepo.FindByTicketID(req.TicketID)
	if existingIncident != nil {
		return nil, errors.New("un incident existe déjà pour ce ticket")
	}

	// Valider l'impact et l'urgence
	validImpacts := []string{"low", "medium", "high", "critical"}
	validUrgencies := []string{"low", "medium", "high", "critical"}

	impactValid := false
	for _, vi := range validImpacts {
		if req.Impact == vi {
			impactValid = true
			break
		}
	}
	if !impactValid {
		return nil, errors.New("impact invalide")
	}

	urgencyValid := false
	for _, vu := range validUrgencies {
		if req.Urgency == vu {
			urgencyValid = true
			break
		}
	}
	if !urgencyValid {
		return nil, errors.New("urgence invalide")
	}

	// Créer l'incident
	incident := &models.Incident{
		TicketID: req.TicketID,
		Impact:   req.Impact,
		Urgency:  req.Urgency,
	}

	if err := s.incidentRepo.Create(incident); err != nil {
		return nil, errors.New("erreur lors de la création de l'incident")
	}

	// Récupérer l'incident créé avec ses relations
	createdIncident, err := s.incidentRepo.FindByID(incident.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'incident créé")
	}

	// Convertir en DTO
	incidentDTO := s.incidentToDTO(createdIncident)
	return &incidentDTO, nil
}

// GetByID récupère un incident par son ID
func (s *incidentService) GetByID(id uint) (*dto.IncidentDTO, error) {
	incident, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("incident introuvable")
	}

	incidentDTO := s.incidentToDTO(incident)
	return &incidentDTO, nil
}

// GetByTicketID récupère un incident par l'ID du ticket
func (s *incidentService) GetByTicketID(ticketID uint) (*dto.IncidentDTO, error) {
	incident, err := s.incidentRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("incident introuvable")
	}

	incidentDTO := s.incidentToDTO(incident)
	return &incidentDTO, nil
}

// GetAll récupère tous les incidents
func (s *incidentService) GetAll() ([]dto.IncidentDTO, error) {
	incidents, err := s.incidentRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des incidents")
	}

	incidentDTOs := make([]dto.IncidentDTO, len(incidents))
	for i, incident := range incidents {
		incidentDTOs[i] = s.incidentToDTO(&incident)
	}

	return incidentDTOs, nil
}

// GetByImpact récupère les incidents par impact
func (s *incidentService) GetByImpact(impact string) ([]dto.IncidentDTO, error) {
	incidents, err := s.incidentRepo.FindByImpact(impact)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des incidents")
	}

	incidentDTOs := make([]dto.IncidentDTO, len(incidents))
	for i, incident := range incidents {
		incidentDTOs[i] = s.incidentToDTO(&incident)
	}

	return incidentDTOs, nil
}

// GetByUrgency récupère les incidents par urgence
func (s *incidentService) GetByUrgency(urgency string) ([]dto.IncidentDTO, error) {
	incidents, err := s.incidentRepo.FindByUrgency(urgency)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des incidents")
	}

	incidentDTOs := make([]dto.IncidentDTO, len(incidents))
	for i, incident := range incidents {
		incidentDTOs[i] = s.incidentToDTO(&incident)
	}

	return incidentDTOs, nil
}

// Update met à jour un incident
func (s *incidentService) Update(id uint, req dto.UpdateIncidentRequest, updatedByID uint) (*dto.IncidentDTO, error) {
	// Récupérer l'incident existant
	incident, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("incident introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Impact != "" {
		// Valider l'impact
		validImpacts := []string{"low", "medium", "high", "critical"}
		valid := false
		for _, vi := range validImpacts {
			if req.Impact == vi {
				valid = true
				break
			}
		}
		if !valid {
			return nil, errors.New("impact invalide")
		}
		incident.Impact = req.Impact
	}

	if req.Urgency != "" {
		// Valider l'urgence
		validUrgencies := []string{"low", "medium", "high", "critical"}
		valid := false
		for _, vu := range validUrgencies {
			if req.Urgency == vu {
				valid = true
				break
			}
		}
		if !valid {
			return nil, errors.New("urgence invalide")
		}
		incident.Urgency = req.Urgency
	}

	// Sauvegarder
	if err := s.incidentRepo.Update(incident); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de l'incident")
	}

	// Récupérer l'incident mis à jour
	updatedIncident, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'incident mis à jour")
	}

	incidentDTO := s.incidentToDTO(updatedIncident)
	return &incidentDTO, nil
}

// Qualify qualifie un incident (met à jour impact et urgence)
func (s *incidentService) Qualify(id uint, req dto.QualifyIncidentRequest, qualifiedByID uint) (*dto.IncidentDTO, error) {
	// Récupérer l'incident
	incident, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("incident introuvable")
	}

	// Valider l'impact et l'urgence
	validImpacts := []string{"low", "medium", "high", "critical"}
	validUrgencies := []string{"low", "medium", "high", "critical"}

	impactValid := false
	for _, vi := range validImpacts {
		if req.Impact == vi {
			impactValid = true
			break
		}
	}
	if !impactValid {
		return nil, errors.New("impact invalide")
	}

	urgencyValid := false
	for _, vu := range validUrgencies {
		if req.Urgency == vu {
			urgencyValid = true
			break
		}
	}
	if !urgencyValid {
		return nil, errors.New("urgence invalide")
	}

	// Mettre à jour
	incident.Impact = req.Impact
	incident.Urgency = req.Urgency

	if err := s.incidentRepo.Update(incident); err != nil {
		return nil, errors.New("erreur lors de la qualification de l'incident")
	}

	// Récupérer l'incident mis à jour
	updatedIncident, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'incident mis à jour")
	}

	incidentDTO := s.incidentToDTO(updatedIncident)
	return &incidentDTO, nil
}

// LinkAsset lie un actif à un incident
func (s *incidentService) LinkAsset(incidentID uint, assetID uint, linkedByID uint) error {
	// Vérifier que l'incident existe
	incident, err := s.incidentRepo.FindByID(incidentID)
	if err != nil {
		return errors.New("incident introuvable")
	}

	// Vérifier que l'actif existe
	_, err = s.assetRepo.FindByID(assetID)
	if err != nil {
		return errors.New("actif introuvable")
	}

	// Créer la liaison via ticket_assets (car incident est lié à un ticket)
	ticketAsset := &models.TicketAsset{
		TicketID: incident.TicketID,
		AssetID:  assetID,
	}

	if err := s.ticketAssetRepo.Create(ticketAsset); err != nil {
		return errors.New("erreur lors de la liaison de l'actif")
	}

	return nil
}

// UnlinkAsset supprime la liaison entre un incident et un actif
func (s *incidentService) UnlinkAsset(incidentID uint, assetID uint) error {
	// Vérifier que l'incident existe
	incident, err := s.incidentRepo.FindByID(incidentID)
	if err != nil {
		return errors.New("incident introuvable")
	}

	// Supprimer la liaison directement
	if err := s.ticketAssetRepo.DeleteByTicketAndAsset(incident.TicketID, assetID); err != nil {
		return errors.New("erreur lors de la suppression de la liaison")
	}

	return nil
}

// Resolve résout un incident
func (s *incidentService) Resolve(id uint, resolvedByID uint) (*dto.IncidentDTO, error) {
	// Récupérer l'incident
	incident, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("incident introuvable")
	}

	// Calculer le temps de résolution si le ticket est clôturé
	ticket, err := s.ticketRepo.FindByID(incident.TicketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	if ticket.ClosedAt != nil && ticket.CreatedAt.Before(*ticket.ClosedAt) {
		duration := ticket.ClosedAt.Sub(ticket.CreatedAt)
		resolutionTime := int(duration.Minutes())
		incident.ResolutionTime = &resolutionTime
		now := time.Now()
		incident.ResolvedAt = &now
	}

	// Sauvegarder
	if err := s.incidentRepo.Update(incident); err != nil {
		return nil, errors.New("erreur lors de la résolution de l'incident")
	}

	// Récupérer l'incident mis à jour
	updatedIncident, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'incident mis à jour")
	}

	incidentDTO := s.incidentToDTO(updatedIncident)
	return &incidentDTO, nil
}

// Delete supprime un incident
func (s *incidentService) Delete(id uint) error {
	// Vérifier que l'incident existe
	_, err := s.incidentRepo.FindByID(id)
	if err != nil {
		return errors.New("incident introuvable")
	}

	// Supprimer
	if err := s.incidentRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de l'incident")
	}

	return nil
}

// incidentToDTO convertit un modèle Incident en DTO IncidentDTO
func (s *incidentService) incidentToDTO(incident *models.Incident) dto.IncidentDTO {
	// Convertir le ticket en DTO
	ticketDTO := s.ticketToDTO(&incident.Ticket)

	// Récupérer les actifs liés
	ticketAssets, _ := s.ticketAssetRepo.FindByTicketID(incident.TicketID)
	assetDTOs := make([]dto.AssetDTO, 0)
	for _, ta := range ticketAssets {
		asset, _ := s.assetRepo.FindByID(ta.AssetID)
		if asset != nil {
			assetDTO := s.assetToDTO(asset)
			assetDTOs = append(assetDTOs, assetDTO)
		}
	}

	return dto.IncidentDTO{
		ID:             incident.ID,
		TicketID:       incident.TicketID,
		Ticket:         &ticketDTO,
		Impact:         incident.Impact,
		Urgency:        incident.Urgency,
		ResolutionTime: incident.ResolutionTime,
		ResolvedAt:     incident.ResolvedAt,
		LinkedAssets:   assetDTOs,
		CreatedAt:      incident.CreatedAt,
		UpdatedAt:      incident.UpdatedAt,
	}
}

// ticketToDTO convertit un modèle Ticket en DTO TicketDTO (méthode utilitaire)
func (s *incidentService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
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
func (s *incidentService) userToDTO(user *models.User) dto.UserDTO {
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

// assetToDTO convertit un modèle Asset en DTO AssetDTO (méthode utilitaire)
func (s *incidentService) assetToDTO(asset *models.Asset) dto.AssetDTO {
	var assignedToDTO *dto.UserDTO
	if asset.AssignedTo != nil {
		assignedDTO := s.userToDTO(asset.AssignedTo)
		assignedToDTO = &assignedDTO
	}

	var assignedToID *uint
	if asset.AssignedToID != nil {
		assignedToID = asset.AssignedToID
	}

	return dto.AssetDTO{
		ID:           asset.ID,
		Name:         asset.Name,
		SerialNumber: asset.SerialNumber,
		Model:        asset.Model,
		Manufacturer: asset.Manufacturer,
		CategoryID:   asset.CategoryID,
		AssignedTo:   assignedToID,
		AssignedUser: assignedToDTO,
		Status:       asset.Status,
		CreatedAt:    asset.CreatedAt,
		UpdatedAt:    asset.UpdatedAt,
	}
}
