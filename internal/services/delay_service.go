package services

import (
	"errors"
	"sync"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// DelayService interface pour les opérations sur les retards
type DelayService interface {
	GetByID(id uint) (*dto.DelayDTO, error)
	GetByTicketID(ticketID uint) (*dto.DelayDTO, error)
	GetAll(scope interface{}) ([]dto.DelayDTO, error) // scope peut être *scope.QueryScope ou nil
	GetByUserID(scope interface{}, userID uint) ([]dto.DelayDTO, error)
	GetByStatus(scope interface{}, status string) ([]dto.DelayDTO, error)
	GetUnjustified(scope interface{}) ([]dto.DelayDTO, error)
	Delete(id uint) error
	CreateJustification(delayID uint, req dto.CreateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error)
	UpdateJustification(id uint, req dto.UpdateDelayJustificationRequest, userID uint) (*dto.DelayJustificationDTO, error)
	ValidateJustification(id uint, req dto.ValidateDelayJustificationRequest, validatedByID uint) (*dto.DelayJustificationDTO, error)
	GetJustificationByDelayID(delayID uint) (*dto.DelayJustificationDTO, error)
	DeleteJustification(delayID uint, userID uint) error
	GetJustificationsByUserID(userID uint) ([]dto.DelayJustificationDTO, error)
	GetJustificationByTicketID(ticketID uint) (*dto.DelayJustificationDTO, error)
	GetValidatedJustifications() ([]dto.DelayJustificationDTO, error)
	GetRejectedJustifications() ([]dto.DelayJustificationDTO, error)
	GetJustificationsHistory() ([]dto.DelayJustificationDTO, error)
	RejectJustification(delayID uint, req dto.ValidateDelayJustificationRequest, rejectedByID uint) (*dto.DelayJustificationDTO, error)
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
	ticketRepo             repositories.TicketRepository
	syncMu                 sync.Mutex
	lastSync               time.Time
	syncing                bool
}

// NewDelayService crée une nouvelle instance de DelayService
func NewDelayService(
	delayRepo repositories.DelayRepository,
	delayJustificationRepo repositories.DelayJustificationRepository,
	userRepo repositories.UserRepository,
	ticketRepo repositories.TicketRepository,
) DelayService {
	return &delayService{
		delayRepo:              delayRepo,
		delayJustificationRepo: delayJustificationRepo,
		userRepo:               userRepo,
		ticketRepo:             ticketRepo,
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
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *delayService) GetAll(scopeParam interface{}) ([]dto.DelayDTO, error) {
	s.triggerSyncDelaysFromTickets()
	delays, err := s.delayRepo.FindAll(scopeParam)
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
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *delayService) GetByUserID(scopeParam interface{}, userID uint) ([]dto.DelayDTO, error) {
	s.triggerSyncDelaysFromTickets()
	delays, err := s.delayRepo.FindByUserID(scopeParam, userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des retards")
	}

	var delayDTOs []dto.DelayDTO
	for _, delay := range delays {
		delayDTOs = append(delayDTOs, s.delayToDTO(&delay))
	}

	return delayDTOs, nil
}

func (s *delayService) triggerSyncDelaysFromTickets() {
	s.syncMu.Lock()
	if s.syncing || (!s.lastSync.IsZero() && time.Since(s.lastSync) < 2*time.Minute) {
		s.syncMu.Unlock()
		return
	}
	s.syncing = true
	s.lastSync = time.Now()
	s.syncMu.Unlock()

	go func() {
		s.syncDelaysFromTickets()
		s.syncMu.Lock()
		s.syncing = false
		s.syncMu.Unlock()
	}()
}

func (s *delayService) syncDelaysFromTickets() {
	page := 1
	limit := 500
	for {
		tickets, total, err := s.ticketRepo.FindAll(nil, page, limit, nil) // nil scope = pas de filtre (utilisé en interne)
		if err != nil {
			return
		}
		for _, ticket := range tickets {
			if ticket.EstimatedTime == nil || *ticket.EstimatedTime <= 0 || ticket.ActualTime == nil {
				continue
			}
			estimated := *ticket.EstimatedTime
			actual := *ticket.ActualTime
			delayTime := actual - estimated

			if delayTime <= 0 {
				if existing, err := s.delayRepo.FindByTicketID(ticket.ID); err == nil && existing != nil {
					if existing.Status == "unjustified" {
						_ = s.delayRepo.Delete(existing.ID)
					}
				}
				continue
			}

			percentage := float64(delayTime) / float64(estimated) * 100
			if percentage > 999.99 {
				percentage = 999.99
			}
			existing, err := s.delayRepo.FindByTicketID(ticket.ID)
			ownerID := ticket.CreatedByID
			if ticket.AssignedToID != nil {
				ownerID = *ticket.AssignedToID
			}
			if err != nil || existing == nil {
				tid := ticket.ID
				delay := &models.Delay{
					TicketID:        &tid,
					UserID:          ownerID,
					EstimatedTime:   estimated,
					ActualTime:      actual,
					DelayTime:       delayTime,
					DelayPercentage: percentage,
					Status:          "unjustified",
					DetectedAt:      time.Now(),
				}
				_ = s.delayRepo.Create(delay)
				continue
			}

			existing.EstimatedTime = estimated
			existing.ActualTime = actual
			existing.DelayTime = delayTime
			existing.DelayPercentage = percentage
			if ownerID != 0 {
				existing.UserID = ownerID
			}
			if existing.Status == "rejected" {
				existing.Status = "unjustified"
			}
			_ = s.delayRepo.Update(existing)
		}
		if page*limit >= int(total) {
			break
		}
		page++
	}
}

// GetByStatus récupère les retards par statut
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (s *delayService) GetByStatus(scopeParam interface{}, status string) ([]dto.DelayDTO, error) {
	delays, err := s.delayRepo.FindByStatus(scopeParam, status)
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
func (s *delayService) GetUnjustified(scopeParam interface{}) ([]dto.DelayDTO, error) {
	delays, err := s.delayRepo.FindUnjustified(scopeParam)
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

	if req.Validated != nil && *req.Validated {
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
		if req.Validated != nil && *req.Validated {
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

// DeleteJustification supprime une justification
func (s *delayService) DeleteJustification(delayID uint, userID uint) error {
	justification, err := s.delayJustificationRepo.FindByDelayID(delayID)
	if err != nil {
		return errors.New("justification introuvable")
	}

	// Vérifier que l'utilisateur est le créateur de la justification
	if justification.UserID != userID {
		return errors.New("vous n'êtes pas autorisé à supprimer cette justification")
	}

	// Vérifier que la justification n'est pas déjà validée ou rejetée
	if justification.Status != "pending" {
		return errors.New("impossible de supprimer une justification déjà validée ou rejetée")
	}

	if err := s.delayJustificationRepo.Delete(justification.ID); err != nil {
		return errors.New("erreur lors de la suppression de la justification")
	}

	// Remettre le statut du retard à unjustified
	delay, err := s.delayRepo.FindByID(delayID)
	if err == nil {
		delay.Status = "unjustified"
		s.delayRepo.Update(delay)
	}

	return nil
}

// GetJustificationsByUserID récupère les justifications d'un utilisateur
func (s *delayService) GetJustificationsByUserID(userID uint) ([]dto.DelayJustificationDTO, error) {
	justifications, err := s.delayJustificationRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des justifications")
	}

	justificationDTOs := make([]dto.DelayJustificationDTO, len(justifications))
	for i, justification := range justifications {
		justificationDTOs[i] = s.justificationToDTO(&justification)
	}

	return justificationDTOs, nil
}

// GetJustificationByTicketID récupère la justification d'un ticket
func (s *delayService) GetJustificationByTicketID(ticketID uint) (*dto.DelayJustificationDTO, error) {
	delay, err := s.delayRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("retard introuvable pour ce ticket")
	}

	if delay.Justification == nil {
		return nil, errors.New("aucune justification trouvée pour ce ticket")
	}

	justificationDTO := s.justificationToDTO(delay.Justification)
	return &justificationDTO, nil
}

// GetValidatedJustifications récupère les justifications validées
func (s *delayService) GetValidatedJustifications() ([]dto.DelayJustificationDTO, error) {
	justifications, err := s.delayJustificationRepo.FindValidated()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des justifications validées")
	}

	justificationDTOs := make([]dto.DelayJustificationDTO, len(justifications))
	for i, justification := range justifications {
		justificationDTOs[i] = s.justificationToDTO(&justification)
	}

	return justificationDTOs, nil
}

// GetRejectedJustifications récupère les justifications rejetées
func (s *delayService) GetRejectedJustifications() ([]dto.DelayJustificationDTO, error) {
	justifications, err := s.delayJustificationRepo.FindRejected()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des justifications rejetées")
	}

	justificationDTOs := make([]dto.DelayJustificationDTO, len(justifications))
	for i, justification := range justifications {
		justificationDTOs[i] = s.justificationToDTO(&justification)
	}

	return justificationDTOs, nil
}

// GetJustificationsHistory récupère l'historique de toutes les justifications
func (s *delayService) GetJustificationsHistory() ([]dto.DelayJustificationDTO, error) {
	justifications, err := s.delayJustificationRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'historique")
	}

	justificationDTOs := make([]dto.DelayJustificationDTO, len(justifications))
	for i, justification := range justifications {
		justificationDTOs[i] = s.justificationToDTO(&justification)
	}

	return justificationDTOs, nil
}

// RejectJustification rejette une justification
func (s *delayService) RejectJustification(delayID uint, req dto.ValidateDelayJustificationRequest, rejectedByID uint) (*dto.DelayJustificationDTO, error) {
	justification, err := s.delayJustificationRepo.FindByDelayID(delayID)
	if err != nil {
		return nil, errors.New("justification introuvable")
	}

	// Vérifier que la justification est en attente
	if justification.Status != "pending" {
		return nil, errors.New("la justification a déjà été traitée")
	}

	// Rejeter la justification
	validated := false
	rejectReq := dto.ValidateDelayJustificationRequest{
		Validated: &validated,
		Comment:   req.Comment,
	}

	return s.ValidateJustification(justification.ID, rejectReq, rejectedByID)
}

// GetStatusStats récupère les statistiques de retards par statut
// Note: Cette méthode est utilisée en interne, donc on passe nil pour le scope
func (s *delayService) GetStatusStats() (*dto.DelayStatusStatsDTO, error) {
	unjustified, _ := s.delayRepo.FindByStatus(nil, "unjustified")
	pending, _ := s.delayRepo.FindByStatus(nil, "pending")
	justified, _ := s.delayRepo.FindByStatus(nil, "justified")
	rejected, _ := s.delayRepo.FindByStatus(nil, "rejected")

	return &dto.DelayStatusStatsDTO{
		Unjustified: len(unjustified),
		Pending:     len(pending),
		Justified:   len(justified),
		Rejected:    len(rejected),
	}, nil
}

// delayToDTO convertit un modèle Delay en DTO
func (s *delayService) delayToDTO(delay *models.Delay) dto.DelayDTO {
	var ticketID uint
	if delay.TicketID != nil {
		ticketID = *delay.TicketID
	}
	delayDTO := dto.DelayDTO{
		ID:              delay.ID,
		TicketID:        ticketID,
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

	if delay.Ticket != nil && delay.Ticket.ID != 0 {
		ticketDTO := s.ticketToDTO(delay.Ticket)
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

	if justification.Delay.ID != 0 {
		if justification.Delay.TicketID != nil {
			justificationDTO.TicketID = justification.Delay.TicketID
		}
		if justification.Delay.Ticket != nil && justification.Delay.Ticket.ID != 0 {
			justificationDTO.TicketCode = justification.Delay.Ticket.Code
			justificationDTO.TicketTitle = justification.Delay.Ticket.Title
		}
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
		Code:        ticket.Code,
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
