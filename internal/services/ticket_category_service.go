package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// TicketCategoryService interface pour les opérations sur les catégories de tickets
type TicketCategoryService interface {
	Create(req dto.CreateTicketCategoryRequest) (*dto.TicketCategoryDTO, error)
	GetByID(id uint) (*dto.TicketCategoryDTO, error)
	GetBySlug(slug string) (*dto.TicketCategoryDTO, error)
	GetAll() ([]dto.TicketCategoryDTO, error)
	GetActive() ([]dto.TicketCategoryDTO, error)
	Update(id uint, req dto.UpdateTicketCategoryRequest) (*dto.TicketCategoryDTO, error)
	Delete(id uint) error
}

// ticketCategoryService implémente TicketCategoryService
type ticketCategoryService struct {
	categoryRepo repositories.TicketCategoryRepository
}

// NewTicketCategoryService crée une nouvelle instance de TicketCategoryService
func NewTicketCategoryService(
	categoryRepo repositories.TicketCategoryRepository,
) TicketCategoryService {
	return &ticketCategoryService{
		categoryRepo: categoryRepo,
	}
}

// Create crée une nouvelle catégorie
func (s *ticketCategoryService) Create(req dto.CreateTicketCategoryRequest) (*dto.TicketCategoryDTO, error) {
	// Vérifier que le slug n'existe pas déjà
	existing, _ := s.categoryRepo.FindBySlug(req.Slug)
	if existing != nil {
		return nil, errors.New("une catégorie avec ce slug existe déjà")
	}

	category := &models.TicketCategory{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		Icon:         req.Icon,
		Color:        req.Color,
		IsActive:     req.IsActive,
		DisplayOrder: req.DisplayOrder,
	}

	// Valeurs par défaut
	if !req.IsActive {
		category.IsActive = true
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, errors.New("erreur lors de la création de la catégorie")
	}

	createdCategory, err := s.categoryRepo.FindByID(category.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la catégorie créée")
	}

	categoryDTO := s.categoryToDTO(createdCategory)
	return &categoryDTO, nil
}

// GetByID récupère une catégorie par son ID
func (s *ticketCategoryService) GetByID(id uint) (*dto.TicketCategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	categoryDTO := s.categoryToDTO(category)
	return &categoryDTO, nil
}

// GetBySlug récupère une catégorie par son slug
func (s *ticketCategoryService) GetBySlug(slug string) (*dto.TicketCategoryDTO, error) {
	category, err := s.categoryRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	categoryDTO := s.categoryToDTO(category)
	return &categoryDTO, nil
}

// GetAll récupère toutes les catégories
func (s *ticketCategoryService) GetAll() ([]dto.TicketCategoryDTO, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des catégories")
	}

	var categoryDTOs []dto.TicketCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, s.categoryToDTO(&category))
	}

	return categoryDTOs, nil
}

// GetActive récupère toutes les catégories actives
func (s *ticketCategoryService) GetActive() ([]dto.TicketCategoryDTO, error) {
	categories, err := s.categoryRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des catégories")
	}

	var categoryDTOs []dto.TicketCategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, s.categoryToDTO(&category))
	}

	return categoryDTOs, nil
}

// Update met à jour une catégorie
func (s *ticketCategoryService) Update(id uint, req dto.UpdateTicketCategoryRequest) (*dto.TicketCategoryDTO, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("catégorie introuvable")
	}

	// Vérifier que le slug n'existe pas déjà (si modifié)
	if req.Slug != "" && req.Slug != category.Slug {
		existing, _ := s.categoryRepo.FindBySlug(req.Slug)
		if existing != nil {
			return nil, errors.New("une catégorie avec ce slug existe déjà")
		}
		category.Slug = req.Slug
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.DisplayOrder != nil {
		category.DisplayOrder = *req.DisplayOrder
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la catégorie")
	}

	updatedCategory, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la catégorie mise à jour")
	}

	categoryDTO := s.categoryToDTO(updatedCategory)
	return &categoryDTO, nil
}

// Delete supprime une catégorie
func (s *ticketCategoryService) Delete(id uint) error {
	// Vérifier que la catégorie existe
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("catégorie introuvable")
	}

	// TODO: Vérifier qu'aucun ticket n'utilise cette catégorie avant de supprimer

	if err := s.categoryRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la catégorie")
	}

	return nil
}

// categoryToDTO convertit un modèle TicketCategory en DTO TicketCategoryDTO
func (s *ticketCategoryService) categoryToDTO(category *models.TicketCategory) dto.TicketCategoryDTO {
	return dto.TicketCategoryDTO{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		Icon:         category.Icon,
		Color:        category.Color,
		IsActive:     category.IsActive,
		DisplayOrder: category.DisplayOrder,
	}
}
