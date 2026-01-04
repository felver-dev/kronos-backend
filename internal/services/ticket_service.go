package services

import (
	"errors"
	"time"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TicketService interface pour les opérations sur les tickets
type TicketService interface {
	Create(req dto.CreateTicketRequest, createdByID uint) (*dto.TicketDTO, error)
	GetByID(id uint) (*dto.TicketDTO, error)
	GetAll(page, limit int) (*dto.TicketListResponse, error)
	GetByStatus(status string, page, limit int) (*dto.TicketListResponse, error)
	GetByCategory(category string, page, limit int) (*dto.TicketListResponse, error)
	GetBySource(source string, page, limit int) (*dto.TicketListResponse, error)
	GetByAssignedTo(userID uint, page, limit int) (*dto.TicketListResponse, error)
	GetByCreatedBy(userID uint, page, limit int) (*dto.TicketListResponse, error)
	GetHistory(ticketID uint) ([]dto.TicketHistoryDTO, error)
	Update(id uint, req dto.UpdateTicketRequest, updatedByID uint) (*dto.TicketDTO, error)
	Assign(id uint, req dto.AssignTicketRequest, assignedByID uint) (*dto.TicketDTO, error)
	ChangeStatus(id uint, status string, changedByID uint) (*dto.TicketDTO, error)
	Close(id uint, closedByID uint) (*dto.TicketDTO, error)
	Delete(id uint) error
	AddComment(ticketID uint, req dto.CreateTicketCommentRequest, userID uint) (*dto.TicketCommentDTO, error)
	GetComments(ticketID uint) ([]dto.TicketCommentDTO, error)
}

// ticketService implémente TicketService
type ticketService struct {
	ticketRepo  repositories.TicketRepository
	userRepo    repositories.UserRepository
	commentRepo repositories.TicketCommentRepository
	historyRepo repositories.TicketHistoryRepository
}

// NewTicketService crée une nouvelle instance de TicketService
func NewTicketService(
	ticketRepo repositories.TicketRepository,
	userRepo repositories.UserRepository,
	commentRepo repositories.TicketCommentRepository,
	historyRepo repositories.TicketHistoryRepository,
) TicketService {
	return &ticketService{
		ticketRepo:  ticketRepo,
		userRepo:    userRepo,
		commentRepo: commentRepo,
		historyRepo: historyRepo,
	}
}

// Create crée un nouveau ticket
func (s *ticketService) Create(req dto.CreateTicketRequest, createdByID uint) (*dto.TicketDTO, error) {
	// Vérifier que l'utilisateur créateur existe
	_, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur créateur introuvable")
	}

	// Créer le ticket
	ticket := &models.Ticket{
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Source:        req.Source,
		Status:        "ouvert", // Statut par défaut
		Priority:      req.Priority,
		CreatedByID:   createdByID,
		EstimatedTime: req.EstimatedTime,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, errors.New("erreur lors de la création du ticket")
	}

	// Créer une entrée d'historique
	s.createHistory(ticket.ID, createdByID, "created", "", "", "Ticket créé")

	// Récupérer le ticket créé avec ses relations
	createdTicket, err := s.ticketRepo.FindByID(ticket.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket créé")
	}

	// Convertir en DTO
	ticketDTO := s.ticketToDTO(createdTicket)
	return &ticketDTO, nil
}

// GetByID récupère un ticket par son ID
func (s *ticketService) GetByID(id uint) (*dto.TicketDTO, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	ticketDTO := s.ticketToDTO(ticket)
	return &ticketDTO, nil
}

// GetAll récupère tous les tickets avec pagination
func (s *ticketService) GetAll(page, limit int) (*dto.TicketListResponse, error) {
	// TODO: Implémenter la pagination dans le repository
	tickets, err := s.ticketRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: int64(len(tickets)),
		},
	}, nil
}

// GetByStatus récupère les tickets par statut
func (s *ticketService) GetByStatus(status string, page, limit int) (*dto.TicketListResponse, error) {
	tickets, err := s.ticketRepo.FindByStatus(status)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: int64(len(tickets)),
		},
	}, nil
}

// GetByCategory récupère les tickets par catégorie
func (s *ticketService) GetByCategory(category string, page, limit int) (*dto.TicketListResponse, error) {
	tickets, err := s.ticketRepo.FindByCategory(category)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: int64(len(tickets)),
		},
	}, nil
}

// GetByAssignedTo récupère les tickets assignés à un utilisateur
func (s *ticketService) GetByAssignedTo(userID uint, page, limit int) (*dto.TicketListResponse, error) {
	tickets, err := s.ticketRepo.FindByAssignedTo(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: int64(len(tickets)),
		},
	}, nil
}

// GetByCreatedBy récupère les tickets créés par un utilisateur
func (s *ticketService) GetByCreatedBy(userID uint, page, limit int) (*dto.TicketListResponse, error) {
	tickets, err := s.ticketRepo.FindByCreatedBy(userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: int64(len(tickets)),
		},
	}, nil
}

// GetBySource récupère les tickets par source
func (s *ticketService) GetBySource(source string, page, limit int) (*dto.TicketListResponse, error) {
	tickets, err := s.ticketRepo.FindBySource(source)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des tickets")
	}

	ticketDTOs := make([]dto.TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = s.ticketToDTO(&ticket)
	}

	return &dto.TicketListResponse{
		Tickets: ticketDTOs,
		Pagination: dto.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: int64(len(tickets)),
		},
	}, nil
}

// GetHistory récupère l'historique d'un ticket
func (s *ticketService) GetHistory(ticketID uint) ([]dto.TicketHistoryDTO, error) {
	histories, err := s.historyRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'historique")
	}

	historyDTOs := make([]dto.TicketHistoryDTO, len(histories))
	for i, history := range histories {
		userDTO := s.userToDTO(&history.User)
		historyDTOs[i] = dto.TicketHistoryDTO{
			ID:          history.ID,
			TicketID:    history.TicketID,
			User:        userDTO,
			Action:      history.Action,
			FieldName:   history.FieldName,
			OldValue:    history.OldValue,
			NewValue:    history.NewValue,
			Description: history.Description,
			CreatedAt:   history.CreatedAt,
		}
	}

	return historyDTOs, nil
}

// Update met à jour un ticket
func (s *ticketService) Update(id uint, req dto.UpdateTicketRequest, updatedByID uint) (*dto.TicketDTO, error) {
	// Récupérer le ticket existant
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Title != "" {
		s.createHistory(id, updatedByID, "updated", "title", ticket.Title, req.Title)
		ticket.Title = req.Title
	}

	if req.Description != "" {
		s.createHistory(id, updatedByID, "updated", "description", ticket.Description, req.Description)
		ticket.Description = req.Description
	}

	if req.Priority != "" {
		s.createHistory(id, updatedByID, "updated", "priority", ticket.Priority, req.Priority)
		ticket.Priority = req.Priority
	}

	// EstimatedTime n'est pas dans UpdateTicketRequest, il faut utiliser AssignTicketRequest

	// Sauvegarder
	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du ticket")
	}

	// Récupérer le ticket mis à jour
	updatedTicket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket mis à jour")
	}

	ticketDTO := s.ticketToDTO(updatedTicket)
	return &ticketDTO, nil
}

// Assign assigne un ticket à un utilisateur
func (s *ticketService) Assign(id uint, req dto.AssignTicketRequest, assignedByID uint) (*dto.TicketDTO, error) {
	// Récupérer le ticket
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Vérifier que l'utilisateur assigné existe
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, errors.New("utilisateur assigné introuvable")
	}

	// Enregistrer l'ancien assigné pour l'historique
	oldAssignedID := ticket.AssignedToID
	assignedToID := req.UserID
	ticket.AssignedToID = &assignedToID

	// Mettre à jour le temps estimé si fourni
	if req.EstimatedTime != nil {
		ticket.EstimatedTime = req.EstimatedTime
	}

	// Changer le statut si assigné
	if ticket.Status == "ouvert" {
		ticket.Status = "en_cours"
	}

	// Sauvegarder
	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, errors.New("erreur lors de l'assignation du ticket")
	}

	// Créer une entrée d'historique
	oldValue := ""
	newValue := ""
	if oldAssignedID != nil {
		oldUser, _ := s.userRepo.FindByID(*oldAssignedID)
		if oldUser != nil {
			oldValue = oldUser.Username
		}
	}
	newUser, _ := s.userRepo.FindByID(req.UserID)
	if newUser != nil {
		newValue = newUser.Username
	}
	s.createHistory(id, assignedByID, "assigned", "assigned_to", oldValue, newValue)

	// Récupérer le ticket mis à jour
	updatedTicket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket mis à jour")
	}

	ticketDTO := s.ticketToDTO(updatedTicket)
	return &ticketDTO, nil
}

// ChangeStatus change le statut d'un ticket
func (s *ticketService) ChangeStatus(id uint, status string, changedByID uint) (*dto.TicketDTO, error) {
	// Récupérer le ticket
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Valider le statut
	validStatuses := []string{"ouvert", "en_cours", "en_attente", "cloture"}
	valid := false
	for _, vs := range validStatuses {
		if status == vs {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("statut invalide")
	}

	// Enregistrer l'ancien statut pour l'historique
	oldStatus := ticket.Status
	ticket.Status = status

	// Si le ticket est clôturé, enregistrer la date de clôture
	if status == "cloture" && ticket.ClosedAt == nil {
		now := time.Now()
		ticket.ClosedAt = &now
	}

	// Sauvegarder
	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, errors.New("erreur lors du changement de statut")
	}

	// Créer une entrée d'historique
	s.createHistory(id, changedByID, "status_changed", "status", oldStatus, status)

	// Récupérer le ticket mis à jour
	updatedTicket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du ticket mis à jour")
	}

	ticketDTO := s.ticketToDTO(updatedTicket)
	return &ticketDTO, nil
}

// Close ferme un ticket
func (s *ticketService) Close(id uint, closedByID uint) (*dto.TicketDTO, error) {
	return s.ChangeStatus(id, "cloture", closedByID)
}

// Delete supprime un ticket (soft delete)
func (s *ticketService) Delete(id uint) error {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return errors.New("ticket introuvable")
	}

	// Supprimer (soft delete)
	if err := s.ticketRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du ticket")
	}

	return nil
}

// AddComment ajoute un commentaire à un ticket
func (s *ticketService) AddComment(ticketID uint, req dto.CreateTicketCommentRequest, userID uint) (*dto.TicketCommentDTO, error) {
	// Vérifier que le ticket existe
	_, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	// Créer le commentaire
	comment := &models.TicketComment{
		TicketID:   ticketID,
		UserID:     userID,
		Comment:    req.Comment,
		IsInternal: req.IsInternal,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, errors.New("erreur lors de la création du commentaire")
	}

	// Créer une entrée d'historique
	s.createHistory(ticketID, userID, "comment_added", "", "", "Commentaire ajouté")

	// Récupérer le commentaire créé
	createdComment, err := s.commentRepo.FindByID(comment.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du commentaire créé")
	}

	// Convertir en DTO
	commentDTO := s.commentToDTO(createdComment)
	return &commentDTO, nil
}

// GetComments récupère tous les commentaires d'un ticket
func (s *ticketService) GetComments(ticketID uint) ([]dto.TicketCommentDTO, error) {
	comments, err := s.commentRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des commentaires")
	}

	commentDTOs := make([]dto.TicketCommentDTO, len(comments))
	for i, comment := range comments {
		commentDTOs[i] = s.commentToDTO(&comment)
	}

	return commentDTOs, nil
}

// createHistory crée une entrée d'historique pour un ticket
func (s *ticketService) createHistory(ticketID, userID uint, action, fieldName, oldValue, newValue string) {
	history := &models.TicketHistory{
		TicketID:    ticketID,
		UserID:      userID,
		Action:      action,
		FieldName:   fieldName,
		OldValue:    oldValue,
		NewValue:    newValue,
		Description: "",
	}
	s.historyRepo.Create(history)
}

// ticketToDTO convertit un modèle Ticket en DTO TicketDTO
func (s *ticketService) ticketToDTO(ticket *models.Ticket) dto.TicketDTO {

	// Convertir les utilisateurs en DTOs
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

// commentToDTO convertit un modèle TicketComment en DTO TicketCommentDTO
func (s *ticketService) commentToDTO(comment *models.TicketComment) dto.TicketCommentDTO {
	userDTO := s.userToDTO(&comment.User)
	return dto.TicketCommentDTO{
		ID:         comment.ID,
		TicketID:   comment.TicketID,
		User:       userDTO,
		Comment:    comment.Comment,
		IsInternal: comment.IsInternal,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}
}

// userToDTO convertit un modèle User en DTO UserDTO (méthode utilitaire)
func (s *ticketService) userToDTO(user *models.User) dto.UserDTO {
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
