package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TicketSolutionService interface pour les opérations sur les solutions de tickets
type TicketSolutionService interface {
	Create(ticketID uint, req dto.CreateTicketSolutionRequest, createdByID uint) (*dto.TicketSolutionDTO, error)
	GetByID(id uint) (*dto.TicketSolutionDTO, error)
	GetByTicketID(ticketID uint) ([]dto.TicketSolutionDTO, error)
	Update(id uint, req dto.UpdateTicketSolutionRequest, updatedByID uint) (*dto.TicketSolutionDTO, error)
	Delete(id uint) error
	PublishToKB(solutionID uint, req dto.PublishSolutionToKBRequest, publishedByID uint) (*dto.KnowledgeArticleDTO, error)
}

// ticketSolutionService implémente TicketSolutionService
type ticketSolutionService struct {
	solutionRepo   repositories.TicketSolutionRepository
	ticketRepo     repositories.TicketRepository
	userRepo       repositories.UserRepository
	roleRepo       repositories.RoleRepository
	kbArticleRepo  repositories.KnowledgeArticleRepository
	kbCategoryRepo repositories.KnowledgeCategoryRepository
}

// NewTicketSolutionService crée une nouvelle instance de TicketSolutionService
func NewTicketSolutionService(
	solutionRepo repositories.TicketSolutionRepository,
	ticketRepo repositories.TicketRepository,
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	kbArticleRepo repositories.KnowledgeArticleRepository,
	kbCategoryRepo repositories.KnowledgeCategoryRepository,
) TicketSolutionService {
	return &ticketSolutionService{
		solutionRepo:   solutionRepo,
		ticketRepo:     ticketRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		kbArticleRepo:  kbArticleRepo,
		kbCategoryRepo: kbCategoryRepo,
	}
}

// Create crée une nouvelle solution pour un ticket
func (s *ticketSolutionService) Create(ticketID uint, req dto.CreateTicketSolutionRequest, createdByID uint) (*dto.TicketSolutionDTO, error) {
	// Vérifier que le ticket existe et est résolu ou clôturé (résolveurs/assignés peuvent documenter dès que le ticket est résolu)
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	if ticket.Status != "resolu" && ticket.Status != "cloture" {
		return nil, errors.New("seuls les tickets résolus ou clôturés peuvent avoir des solutions documentées")
	}

	// Vérifier que l'utilisateur est assigné ou admin
	user, err := s.userRepo.FindByID(createdByID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	// Vérifier si l'utilisateur est assigné au ticket ou est admin
	isAssigned := false
	if ticket.AssignedToID != nil && *ticket.AssignedToID == createdByID {
		isAssigned = true
	}
	// Vérifier dans les assignees
	for _, assignee := range ticket.Assignees {
		if assignee.UserID == createdByID {
			isAssigned = true
			break
		}
	}

	// Vérifier si l'utilisateur a la permission de modifier les tickets (documentation de solution autorisée)
	perms, _ := s.roleRepo.GetPermissionsByRoleID(user.RoleID)
	canManageTickets := false
	for _, p := range perms {
		if p == "tickets.update" {
			canManageTickets = true
			break
		}
	}

	if !isAssigned && !canManageTickets {
		return nil, errors.New("seuls les assignés au ticket ou les administrateurs peuvent documenter une solution")
	}

	// Créer la solution
	solution := &models.TicketSolution{
		TicketID:    ticketID,
		Solution:    req.Solution,
		CreatedByID: createdByID,
	}

	if err := s.solutionRepo.Create(solution); err != nil {
		return nil, errors.New("erreur lors de la création de la solution")
	}

	// Récupérer la solution créée avec ses relations
	createdSolution, err := s.solutionRepo.FindByID(solution.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la solution créée")
	}

	solutionDTO := s.solutionToDTO(createdSolution)
	return &solutionDTO, nil
}

// GetByID récupère une solution par son ID
func (s *ticketSolutionService) GetByID(id uint) (*dto.TicketSolutionDTO, error) {
	solution, err := s.solutionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("solution introuvable")
	}

	solutionDTO := s.solutionToDTO(solution)
	return &solutionDTO, nil
}

// GetByTicketID récupère toutes les solutions d'un ticket
func (s *ticketSolutionService) GetByTicketID(ticketID uint) ([]dto.TicketSolutionDTO, error) {
	solutions, err := s.solutionRepo.FindByTicketID(ticketID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des solutions")
	}

	var solutionDTOs []dto.TicketSolutionDTO
	for _, solution := range solutions {
		solutionDTOs = append(solutionDTOs, s.solutionToDTO(&solution))
	}

	return solutionDTOs, nil
}

// Update met à jour une solution
func (s *ticketSolutionService) Update(id uint, req dto.UpdateTicketSolutionRequest, updatedByID uint) (*dto.TicketSolutionDTO, error) {
	solution, err := s.solutionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("solution introuvable")
	}

	// Vérifier que l'utilisateur est le créateur, assigné au ticket ou admin
	user, err := s.userRepo.FindByID(updatedByID)
	if err != nil {
		return nil, errors.New("utilisateur introuvable")
	}

	isCreator := solution.CreatedByID == updatedByID
	perms, _ := s.roleRepo.GetPermissionsByRoleID(user.RoleID)
	canManageTickets := false
	for _, p := range perms {
		if p == "tickets.update" {
			canManageTickets = true
			break
		}
	}

	// Vérifier si l'utilisateur est assigné au ticket
	ticket, err := s.ticketRepo.FindByID(solution.TicketID)
	if err != nil {
		return nil, errors.New("ticket introuvable")
	}

	isAssigned := false
	if ticket.AssignedToID != nil && *ticket.AssignedToID == updatedByID {
		isAssigned = true
	}
	for _, assignee := range ticket.Assignees {
		if assignee.UserID == updatedByID {
			isAssigned = true
			break
		}
	}

	if !isCreator && !isAssigned && !canManageTickets {
		return nil, errors.New("vous n'avez pas la permission de modifier cette solution")
	}

	// Mettre à jour la solution
	solution.Solution = req.Solution

	if err := s.solutionRepo.Update(solution); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la solution")
	}

	// Récupérer la solution mise à jour
	updatedSolution, err := s.solutionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la solution mise à jour")
	}

	solutionDTO := s.solutionToDTO(updatedSolution)
	return &solutionDTO, nil
}

// Delete supprime une solution
func (s *ticketSolutionService) Delete(id uint) error {
	_, err := s.solutionRepo.FindByID(id)
	if err != nil {
		return errors.New("solution introuvable")
	}

	if err := s.solutionRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la solution")
	}

	return nil
}

// PublishToKB publie une solution dans la base de connaissances
func (s *ticketSolutionService) PublishToKB(solutionID uint, req dto.PublishSolutionToKBRequest, publishedByID uint) (*dto.KnowledgeArticleDTO, error) {
	// Récupérer la solution
	solution, err := s.solutionRepo.FindByID(solutionID)
	if err != nil {
		return nil, errors.New("solution introuvable")
	}

	// Vérifier que la catégorie KB existe
	_, err = s.kbCategoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("catégorie de base de connaissances introuvable")
	}

	// Créer l'article de base de connaissances
	article := &models.KnowledgeArticle{
		Title:       req.Title,
		Content:     solution.Solution,
		CategoryID:  req.CategoryID,
		AuthorID:    publishedByID,
		IsPublished: true, // Publication directe
		ViewCount:   0,
	}

	if err := s.kbArticleRepo.Create(article); err != nil {
		return nil, errors.New("erreur lors de la création de l'article de base de connaissances")
	}

	// Récupérer l'article créé avec ses relations
	createdArticle, err := s.kbArticleRepo.FindByID(article.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de l'article créé")
	}

	// Convertir en DTO
	articleDTO := s.articleToDTO(createdArticle)
	return &articleDTO, nil
}

// solutionToDTO convertit un modèle TicketSolution en DTO
func (s *ticketSolutionService) solutionToDTO(solution *models.TicketSolution) dto.TicketSolutionDTO {
	solutionDTO := dto.TicketSolutionDTO{
		ID:        solution.ID,
		TicketID:  solution.TicketID,
		Solution:  solution.Solution,
		CreatedAt: solution.CreatedAt,
		UpdatedAt: solution.UpdatedAt,
	}

	// Convertir CreatedBy
	if solution.CreatedBy.ID != 0 {
		userDTO := s.userToDTO(&solution.CreatedBy)
		solutionDTO.CreatedBy = userDTO
	}

	return solutionDTO
}

// userToDTO convertit un modèle User en DTO
func (s *ticketSolutionService) userToDTO(user *models.User) dto.UserDTO {
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

	if user.RoleID != 0 && user.Role.ID != 0 {
		userDTO.Role = user.Role.Name
	}

	if user.LastLogin != nil {
		userDTO.LastLogin = user.LastLogin
	}

	return userDTO
}

// articleToDTO convertit un modèle KnowledgeArticle en DTO
func (s *ticketSolutionService) articleToDTO(article *models.KnowledgeArticle) dto.KnowledgeArticleDTO {
	articleDTO := dto.KnowledgeArticleDTO{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		CategoryID:  article.CategoryID,
		AuthorID:    article.AuthorID,
		IsPublished: article.IsPublished,
		ViewCount:   article.ViewCount,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	// Convertir Category
	if article.Category.ID != 0 {
		categoryDTO := &dto.KnowledgeCategoryDTO{
			ID:          article.Category.ID,
			Name:        article.Category.Name,
			Description: article.Category.Description,
		}
		articleDTO.Category = categoryDTO
	}

	// Convertir Author
	if article.Author.ID != 0 {
		authorDTO := s.userToDTO(&article.Author)
		articleDTO.Author = &authorDTO
	}

	return articleDTO
}
