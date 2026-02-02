package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// ChangeService interface pour les opérations sur les changements
type ChangeService interface {
	Create(req dto.CreateChangeRequest, createdByID uint) (*dto.ChangeDTO, error)
	GetByID(id uint) (*dto.ChangeDTO, error)
	GetByTicketID(ticketID uint) (*dto.ChangeDTO, error)
	GetAll(scope interface{}) ([]dto.ChangeDTO, error) // scope peut être *scope.QueryScope ou nil
	GetByRisk(scope interface{}, risk string) ([]dto.ChangeDTO, error)
	GetByResponsible(scope interface{}, responsibleID uint) ([]dto.ChangeDTO, error) // scope peut être *scope.QueryScope ou nil
	Update(id uint, req dto.UpdateChangeRequest, updatedByID uint) (*dto.ChangeDTO, error)
	AssignResponsible(id uint, req dto.AssignResponsibleRequest, assignedByID uint) (*dto.ChangeDTO, error)
	UpdateRisk(id uint, req dto.UpdateRiskRequest, updatedByID uint) (*dto.ChangeDTO, error)
	RecordResult(id uint, req dto.RecordChangeResultRequest, recordedByID uint) (*dto.ChangeDTO, error)
	Delete(id uint) error
}

// changeService implémente ChangeService
type changeService struct {
	changeRepo repositories.ChangeRepository
	ticketRepo repositories.TicketRepository
	userRepo   repositories.UserRepository
}

// NewChangeService crée une nouvelle instance de ChangeService
func NewChangeService(
	changeRepo repositories.ChangeRepository,
	ticketRepo repositories.TicketRepository,
	userRepo repositories.UserRepository,
) ChangeService {
	return &changeService{
		changeRepo: changeRepo,
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

// Create crée un nouveau changement à partir d'un ticket
func (s *changeService) Create(req dto.CreateChangeRequest, createdByID uint) (*dto.ChangeDTO, error) {
	// Vérifier que le ticket existe et est de catégorie "changement"
	ticket, err := s.ticketRepo.FindByID(req.TicketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	if ticket.Category != "changement" {
		return nil, errors.New("le ticket doit être de catégorie 'changement'")
	}

	// Vérifier qu'un changement n'existe pas déjà pour ce ticket
	existingChange, _ := s.changeRepo.FindByTicketID(req.TicketID)
	if existingChange != nil {
		return nil, errors.New("un changement existe déjà pour ce ticket")
	}

	// Créer le changement
	change := &models.Change{
		TicketID:        req.TicketID,
		Risk:            req.Risk,
		RiskDescription: req.RiskDescription,
	}

	if err := s.changeRepo.Create(change); err != nil {
		return nil, errors.New("erreur lors de la création du changement")
	}

	// Récupérer le changement créé avec ses relations
	createdChange, err := s.changeRepo.FindByID(change.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du changement créé")
	}

	// Convertir en DTO
	changeDTO := s.changeToDTO(createdChange)
	return &changeDTO, nil
}

// GetByID récupère un changement par son ID
func (s *changeService) GetByID(id uint) (*dto.ChangeDTO, error) {
	change, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("changement introuvable")
	}

	changeDTO := s.changeToDTO(change)
	return &changeDTO, nil
}

// GetByTicketID récupère un changement par l'ID du ticket
func (s *changeService) GetByTicketID(ticketID uint) (*dto.ChangeDTO, error) {
	change, err := s.changeRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("changement introuvable")
	}

	changeDTO := s.changeToDTO(change)
	return &changeDTO, nil
}

// GetAll récupère tous les changements
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *changeService) GetAll(scopeParam interface{}) ([]dto.ChangeDTO, error) {
	changes, err := s.changeRepo.FindAll(scopeParam)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des changements")
	}

	var changeDTOs []dto.ChangeDTO
	for _, change := range changes {
		changeDTOs = append(changeDTOs, s.changeToDTO(&change))
	}

	return changeDTOs, nil
}

// GetByRisk récupère les changements par niveau de risque
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *changeService) GetByRisk(scopeParam interface{}, risk string) ([]dto.ChangeDTO, error) {
	changes, err := s.changeRepo.FindByRisk(scopeParam, risk)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des changements")
	}

	var changeDTOs []dto.ChangeDTO
	for _, change := range changes {
		changeDTOs = append(changeDTOs, s.changeToDTO(&change))
	}

	return changeDTOs, nil
}

// GetByResponsible récupère les changements par responsable
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *changeService) GetByResponsible(scopeParam interface{}, responsibleID uint) ([]dto.ChangeDTO, error) {
	changes, err := s.changeRepo.FindByResponsible(scopeParam, responsibleID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des changements")
	}

	var changeDTOs []dto.ChangeDTO
	for _, change := range changes {
		changeDTOs = append(changeDTOs, s.changeToDTO(&change))
	}

	return changeDTOs, nil
}

// Update met à jour un changement
func (s *changeService) Update(id uint, req dto.UpdateChangeRequest, updatedByID uint) (*dto.ChangeDTO, error) {
	change, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("changement introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Risk != "" {
		change.Risk = req.Risk
	}
	if req.RiskDescription != "" {
		change.RiskDescription = req.RiskDescription
	}

	if err := s.changeRepo.Update(change); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du changement")
	}

	// Récupérer le changement mis à jour
	updatedChange, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du changement mis à jour")
	}

	changeDTO := s.changeToDTO(updatedChange)
	return &changeDTO, nil
}

// AssignResponsible assigne un responsable au changement
func (s *changeService) AssignResponsible(id uint, req dto.AssignResponsibleRequest, assignedByID uint) (*dto.ChangeDTO, error) {
	change, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("changement introuvable")
	}

	// Vérifier que l'utilisateur responsable existe
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, errors.New("utilisateur responsable introuvable")
	}

	change.ResponsibleID = &req.UserID

	if err := s.changeRepo.Update(change); err != nil {
		return nil, errors.New("erreur lors de l'assignation du responsable")
	}

	// Récupérer le changement mis à jour
	updatedChange, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du changement mis à jour")
	}

	changeDTO := s.changeToDTO(updatedChange)
	return &changeDTO, nil
}

// UpdateRisk met à jour le risque d'un changement
func (s *changeService) UpdateRisk(id uint, req dto.UpdateRiskRequest, updatedByID uint) (*dto.ChangeDTO, error) {
	change, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("changement introuvable")
	}

	change.Risk = req.Risk
	if req.RiskDescription != "" {
		change.RiskDescription = req.RiskDescription
	}

	if err := s.changeRepo.Update(change); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du risque")
	}

	// Récupérer le changement mis à jour
	updatedChange, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du changement mis à jour")
	}

	changeDTO := s.changeToDTO(updatedChange)
	return &changeDTO, nil
}

// RecordResult enregistre le résultat post-changement
func (s *changeService) RecordResult(id uint, req dto.RecordChangeResultRequest, recordedByID uint) (*dto.ChangeDTO, error) {
	change, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("changement introuvable")
	}

	now := time.Now()
	change.Result = req.Result
	change.ResultDescription = req.Description
	change.ResultDate = &now

	if err := s.changeRepo.Update(change); err != nil {
		return nil, errors.New("erreur lors de l'enregistrement du résultat")
	}

	// Récupérer le changement mis à jour
	updatedChange, err := s.changeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du changement mis à jour")
	}

	changeDTO := s.changeToDTO(updatedChange)
	return &changeDTO, nil
}

// Delete supprime un changement
func (s *changeService) Delete(id uint) error {
	_, err := s.changeRepo.FindByID(id)
	if err != nil {
		return errors.New("changement introuvable")
	}

	if err := s.changeRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du changement")
	}

	return nil
}

// changeToDTO convertit un modèle Change en DTO
func (s *changeService) changeToDTO(change *models.Change) dto.ChangeDTO {
	changeDTO := dto.ChangeDTO{
		ID:                change.ID,
		TicketID:          change.TicketID,
		Risk:              change.Risk,
		RiskDescription:   change.RiskDescription,
		Result:            change.Result,
		ResultDescription: change.ResultDescription,
		ResultDate:        change.ResultDate,
		CreatedAt:         change.CreatedAt,
		UpdatedAt:         change.UpdatedAt,
	}

	// Convertir le ticket si présent
	if change.Ticket.ID != 0 {
		ticketDTO := s.ticketToDTO(&change.Ticket)
		changeDTO.Ticket = &ticketDTO
	}

	// Convertir le responsable si présent
	if change.ResponsibleID != nil {
		changeDTO.ResponsibleID = change.ResponsibleID
		if change.Responsible != nil && change.Responsible.ID != 0 {
			userDTO := s.userToDTO(change.Responsible)
			changeDTO.Responsible = &userDTO
		}
	}

	return changeDTO
}

// ticketToDTO convertit un modèle Ticket en DTO (méthode helper)
func (s *changeService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
	ticketDTO := dto.TicketDTO{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Category:    ticket.Category,
		Source:      ticket.Source,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	if ticket.EstimatedTime != nil {
		ticketDTO.EstimatedTime = ticket.EstimatedTime
	}
	if ticket.ActualTime != nil {
		ticketDTO.ActualTime = ticket.ActualTime
	}
	if ticket.ClosedAt != nil {
		ticketDTO.ClosedAt = ticket.ClosedAt
	}

	// Convertir CreatedBy
	if ticket.CreatedBy.ID != 0 {
		userDTO := s.userToDTO(&ticket.CreatedBy)
		ticketDTO.CreatedBy = userDTO
	}

	// Convertir AssignedTo
	if ticket.AssignedTo != nil && ticket.AssignedTo.ID != 0 {
		userDTO := s.userToDTO(ticket.AssignedTo)
		ticketDTO.AssignedTo = &userDTO
	}

	return ticketDTO
}

// userToDTO convertit un modèle User en DTO (méthode helper)
func (s *changeService) userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		DepartmentID: user.DepartmentID,
		Avatar:     user.Avatar,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	if user.RoleID != 0 {
		userDTO.Role = user.Role.Name
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	// Inclure le département complet si présent
	if user.Department != nil {
		departmentDTO := dto.DepartmentDTO{
			ID:          user.Department.ID,
			Name:        user.Department.Name,
			Code:        user.Department.Code,
			Description: user.Department.Description,
			OfficeID:    user.Department.OfficeID,
			IsActive:    user.Department.IsActive,
			CreatedAt:   user.Department.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   user.Department.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if user.Department.Office != nil {
			departmentDTO.Office = &dto.OfficeDTO{
				ID:        user.Department.Office.ID,
				Name:      user.Department.Office.Name,
				Country:   user.Department.Office.Country,
				City:      user.Department.Office.City,
				Commune:   user.Department.Office.Commune,
				Address:   user.Department.Office.Address,
				Longitude: user.Department.Office.Longitude,
				Latitude:  user.Department.Office.Latitude,
				IsActive:  user.Department.Office.IsActive,
				CreatedAt: user.Department.Office.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt: user.Department.Office.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
		userDTO.Department = &departmentDTO
	}

	return userDTO
}
