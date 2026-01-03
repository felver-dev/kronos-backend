package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// DelayService interface pour les opérations sur les retards
type DelayService interface {
	GetByID(id uint) (*dto.DelayDTO, error)
	GetByTicketID(ticketID uint) (*dto.DelayDTO, error)
	GetAll() ([]dto.DelayDTO, error)
	GetByUserID(userID uint) ([]dto.DelayDTO, error)
	GetByStatus(status string) ([]dto.DelayDTO, error)
	GetUnjustified() ([]dto.DelayDTO, error)
	Delete(id uint) error
	CreateJustification(delayID uint, req dto.CreateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error)
	UpdateJustification(id uint, req dto.UpdateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error)
	ValidateJustification(id uint, req dto.ValidateDelayJustificationRequest, validatedByID uint) (*dto.DelayJustificationDTO, error)
	GetJustificationByDelayID(delayID uint) (*dto.DelayJustificationDTO, error)
	GetStatusStats() (*dto.DelayStatusStatsDTO, error)
}

// DelayJustificationService interface pour les opérations sur les justifications de retards
type DelayJustificationService interface {
	Create(delayID uint, req dto.CreateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error)
	GetByID(id uint) (*dto.DelayJustificationDTO, error)
	GetByDelayID(delayID uint) (*dto.DelayJustificationDTO, error)
	GetAll() ([]dto.DelayJustificationDTO, error)
	GetPending() ([]dto.DelayJustificationDTO, error)
	Update(id uint, req dto.UpdateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error)
	Validate(id uint, req dto.ValidateDelayJustificationRequest, validatedByID uint) (*dto.DelayJustificationDTO, error)
	Delete(id uint) error
}

// delayService implémente DelayService
type delayService struct {
	delayRepo              repositories.DelayRepository
	delayJustificationRepo repositories.DelayJustificationRepository
	userRepo               repositories.UserRepository
}

// NewDelayService crée une nouvelle instance de DelayService
func NewDelayService(
	delayRepo repositories.DelayRepository,
	delayJustificationRepo repositories.DelayJustificationRepository,
	userRepo repositories.UserRepository,
) DelayService {
	return &delayService{
		delayRepo:              delayRepo,
		delayJustificationRepo: delayJustificationRepo,
		userRepo:               userRepo,
	}
}

// GetByID récupère un retard par son ID
func (s *delayService) GetByID(id uint) (*dto.DelayDTO, error) {
	delay, err := s.delayRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("retard introuvable")
	}

	delayDTO := s.delayToDTO(delay)
	return &delayDTO, nil
}

// GetByTicketID récupère un retard par l'ID du ticket
func (s *delayService) GetByTicketID(ticketID uint) (*dto.DelayDTO, error) {
	delay, err := s.delayRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("retard introuvable")
	}

	delayDTO := s.delayToDTO(delay)
	return &delayDTO, nil
}

// GetAll récupère tous les retards
func (s *delayService) GetAll() ([]dto.DelayDTO, error) {
	delays, err := s.delayRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des retards")
	}

	var delayDTOs []dto.DelayDTO
	for _, delay := range delays {
		delayDTOs = append(delayDTOs, s.delayToDTO(&delay))
	}

	return delayDTOs, nil
}

// GetByUserID récupère les retards d'un utilisateur
func (s *delayService) GetByUserID(userID uint) ([]dto.DelayDTO, error) {
	delays, err := s.delayRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des retards")
	}

	var delayDTOs []dto.DelayDTO
	for _, delay := range delays {
		delayDTOs = append(delayDTOs, s.delayToDTO(&delay))
	}

	return delayDTOs, nil
}

// GetByStatus récupère les retards par statut
func (s *delayService) GetByStatus(status string) ([]dto.DelayDTO, error) {
	delays, err := s.delayRepo.FindByStatus(status)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des retards")
	}

	var delayDTOs []dto.DelayDTO
	for _, delay := range delays {
		delayDTOs = append(delayDTOs, s.delayToDTO(&delay))
	}

	return delayDTOs, nil
}

// GetUnjustified récupère les retards non justifiés
func (s *delayService) GetUnjustified() ([]dto.DelayDTO, error) {
	delays, err := s.delayRepo.FindUnjustified()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des retards")
	}

	var delayDTOs []dto.DelayDTO
	for _, delay := range delays {
		delayDTOs = append(delayDTOs, s.delayToDTO(&delay))
	}

	return delayDTOs, nil
}

// Delete supprime un retard
func (s *delayService) Delete(id uint) error {
	_, err := s.delayRepo.FindByID(id)
	if err != nil {
		return errors.New("retard introuvable")
	}

	if err := s.delayRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du retard")
	}

	return nil
}

// CreateJustification crée une justification pour un retard
func (s *delayService) CreateJustification(delayID uint, req dto.CreateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error) {
	// Vérifier que le retard existe
	delay, err := s.delayRepo.FindByID(delayID)
	if err != nil {
		return nil, errors.New("retard introuvable")
	}

	// Vérifier que l'utilisateur est le technicien du retard
	if delay.UserID != userID {
		return nil, errors.New("vous n'êtes pas autorisé à justifier ce retard")
	}

	// Vérifier qu'une justification n'existe pas déjà
	existingJustification, _ := s.delayJustificationRepo.FindByDelayID(delayID)
	if existingJustification != nil {
		return nil, errors.New("une justification existe déjà pour ce retard")
	}

	// Créer la justification
	justification := &models.DelayJustification{
		DelayID:       delayID,
		UserID:        userID,
		Justification: req.Justification,
		Status:        "pending",
	}

	if err := s.delayJustificationRepo.Create(justification); err != nil {
		return nil, errors.New("erreur lors de la création de la justification")
	}

	// Mettre à jour le statut du retard
	delay.Status = "pending"
	if err := s.delayRepo.Update(delay); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du retard")
	}

	// Récupérer la justification créée
	createdJustification, err := s.delayJustificationRepo.FindByID(justification.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la justification créée")
	}

	justificationDTO := s.justificationToDTO(createdJustification)
	return &justificationDTO, nil
}

// UpdateJustification met à jour une justification (avant validation)
func (s *delayService) UpdateJustification(id uint, req dto.UpdateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error) {
	justification, err := s.delayJustificationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("justification introuvable")
	}

	// Vérifier que l'utilisateur est le créateur de la justification
	if justification.UserID != userID {
		return nil, errors.New("vous n'êtes pas autorisé à modifier cette justification")
	}

	// Vérifier que la justification n'est pas déjà validée ou rejetée
	if justification.Status != "pending" {
		return nil, errors.New("impossible de modifier une justification déjà validée ou rejetée")
	}

	justification.Justification = req.Justification

	if err := s.delayJustificationRepo.Update(justification); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la justification")
	}

	// Récupérer la justification mise à jour
	updatedJustification, err := s.delayJustificationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la justification mise à jour")
	}

	justificationDTO := s.justificationToDTO(updatedJustification)
	return &justificationDTO, nil
}

// ValidateJustification valide ou rejette une justification
func (s *delayService) ValidateJustification(id uint, req dto.ValidateDelayJustificationRequest, validatedByID uint) (*dto.DelayJustificationDTO, error) {
	justification, err := s.delayJustificationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("justification introuvable")
	}

	// Vérifier que la justification est en attente
	if justification.Status != "pending" {
		return nil, errors.New("la justification a déjà été traitée")
	}

	now := time.Now()
	justification.ValidatedByID = &validatedByID
	justification.ValidatedAt = &now
	justification.ValidationComment = req.Comment

	if req.Validated {
		justification.Status = "validated"
	} else {
		justification.Status = "rejected"
	}

	if err := s.delayJustificationRepo.Update(justification); err != nil {
		return nil, errors.New("erreur lors de la validation de la justification")
	}

	// Mettre à jour le statut du retard
	delay, err := s.delayRepo.FindByID(justification.DelayID)
	if err == nil {
		if req.Validated {
			delay.Status = "justified"
		} else {
			delay.Status = "unjustified"
		}
		s.delayRepo.Update(delay)
	}

	// Récupérer la justification mise à jour
	updatedJustification, err := s.delayJustificationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la justification mise à jour")
	}

	justificationDTO := s.justificationToDTO(updatedJustification)
	return &justificationDTO, nil
}

// GetJustificationByDelayID récupère la justification d'un retard
func (s *delayService) GetJustificationByDelayID(delayID uint) (*dto.DelayJustificationDTO, error) {
	justification, err := s.delayJustificationRepo.FindByDelayID(delayID)
	if err != nil {
		return nil, errors.New("justification introuvable")
	}

	justificationDTO := s.justificationToDTO(justification)
	return &justificationDTO, nil
}

// GetStatusStats récupère les statistiques de retards par statut
func (s *delayService) GetStatusStats() (*dto.DelayStatusStatsDTO, error) {
	unjustified, _ := s.delayRepo.FindByStatus("unjustified")
	pending, _ := s.delayRepo.FindByStatus("pending")
	justified, _ := s.delayRepo.FindByStatus("justified")
	rejected, _ := s.delayRepo.FindByStatus("rejected")

	return &dto.DelayStatusStatsDTO{
		Unjustified: len(unjustified),
		Pending:     len(pending),
		Justified:   len(justified),
		Rejected:    len(rejected),
	}, nil
}

// delayToDTO convertit un modèle Delay en DTO
func (s *delayService) delayToDTO(delay *models.Delay) dto.DelayDTO {
	delayDTO := dto.DelayDTO{
		ID:              delay.ID,
		TicketID:        delay.TicketID,
		UserID:          delay.UserID,
		EstimatedTime:   delay.EstimatedTime,
		ActualTime:      delay.ActualTime,
		DelayTime:       delay.DelayTime,
		DelayPercentage: delay.DelayPercentage,
		Status:          delay.Status,
		DetectedAt:      delay.DetectedAt,
		CreatedAt:       delay.CreatedAt,
		UpdatedAt:       delay.UpdatedAt,
	}

	// Convertir le ticket si présent
	if delay.Ticket.ID != 0 {
		ticketDTO := s.ticketToDTO(&delay.Ticket)
		delayDTO.Ticket = &ticketDTO
	}

	// Convertir l'utilisateur si présent
	if delay.User.ID != 0 {
		userDTO := s.userToDTO(&delay.User)
		delayDTO.User = &userDTO
	}

	// Convertir la justification si présente
	if delay.Justification != nil && delay.Justification.ID != 0 {
		justificationDTO := s.justificationToDTO(delay.Justification)
		delayDTO.Justification = &justificationDTO
	}

	return delayDTO
}

// justificationToDTO convertit un modèle DelayJustification en DTO
func (s *delayService) justificationToDTO(justification *models.DelayJustification) dto.DelayJustificationDTO {
	justificationDTO := dto.DelayJustificationDTO{
		ID:            justification.ID,
		DelayID:       justification.DelayID,
		UserID:        justification.UserID,
		Justification: justification.Justification,
		Status:        justification.Status,
		CreatedAt:     justification.CreatedAt,
		UpdatedAt:     justification.UpdatedAt,
	}

	if justification.ValidatedByID != nil {
		justificationDTO.ValidatedBy = justification.ValidatedByID
	}
	if justification.ValidatedAt != nil {
		justificationDTO.ValidatedAt = justification.ValidatedAt
	}
	if justification.ValidationComment != "" {
		justificationDTO.ValidationComment = justification.ValidationComment
	}

	// Convertir l'utilisateur si présent
	if justification.User.ID != 0 {
		userDTO := s.userToDTO(&justification.User)
		justificationDTO.User = &userDTO
	}

	return justificationDTO
}

// ticketToDTO convertit un modèle Ticket en DTO (méthode helper)
func (s *delayService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {
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
func (s *delayService) userToDTO(user *models.User) dto.UserDTO {
	userDTO := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.RoleID != 0 {
		userDTO.Role = user.Role.Name
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	return userDTO
}
