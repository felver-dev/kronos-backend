package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// RequestSourceService interface pour les opérations sur les sources de demande
type RequestSourceService interface {
	GetAll() ([]dto.RequestSourceDTO, error)
	GetByID(id uint) (*dto.RequestSourceDTO, error)
	Create(req dto.CreateRequestSourceRequest, createdByID uint) (*dto.RequestSourceDTO, error)
	Update(id uint, req dto.UpdateRequestSourceRequest, updatedByID uint) (*dto.RequestSourceDTO, error)
	Delete(id uint) error
}

// requestSourceService implémente RequestSourceService
type requestSourceService struct {
	sourceRepo repositories.RequestSourceRepository
}

// NewRequestSourceService crée une nouvelle instance de RequestSourceService
func NewRequestSourceService(sourceRepo repositories.RequestSourceRepository) RequestSourceService {
	return &requestSourceService{
		sourceRepo: sourceRepo,
	}
}

// GetAll récupère toutes les sources de demande
func (s *requestSourceService) GetAll() ([]dto.RequestSourceDTO, error) {
	sources, err := s.sourceRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des sources")
	}

	sourceDTOs := make([]dto.RequestSourceDTO, len(sources))
	for i, source := range sources {
		sourceDTOs[i] = s.sourceToDTO(&source)
	}

	return sourceDTOs, nil
}

// GetByID récupère une source par son ID
func (s *requestSourceService) GetByID(id uint) (*dto.RequestSourceDTO, error) {
	source, err := s.sourceRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("source introuvable")
	}

	sourceDTO := s.sourceToDTO(source)
	return &sourceDTO, nil
}

// Create crée une nouvelle source de demande
func (s *requestSourceService) Create(req dto.CreateRequestSourceRequest, createdByID uint) (*dto.RequestSourceDTO, error) {
	// Vérifier que le code n'existe pas déjà
	existingSource, _ := s.sourceRepo.FindByCode(req.Code)
	if existingSource != nil {
		return nil, errors.New("une source avec ce code existe déjà")
	}

	source := &models.RequestSource{
		Name:      req.Name,
		Code:      req.Code,
		Description: req.Description,
		IsEnabled: req.IsEnabled,
	}

	if err := s.sourceRepo.Create(source); err != nil {
		return nil, errors.New("erreur lors de la création de la source")
	}

	createdSource, err := s.sourceRepo.FindByID(source.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la source créée")
	}

	sourceDTO := s.sourceToDTO(createdSource)
	return &sourceDTO, nil
}

// Update met à jour une source de demande
func (s *requestSourceService) Update(id uint, req dto.UpdateRequestSourceRequest, updatedByID uint) (*dto.RequestSourceDTO, error) {
	source, err := s.sourceRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("source introuvable")
	}

	if req.Name != "" {
		source.Name = req.Name
	}

	if req.Description != "" {
		source.Description = req.Description
	}

	// IsEnabled peut être mis à jour même si false
	if req.IsEnabled != nil {
		source.IsEnabled = *req.IsEnabled
	}

	if err := s.sourceRepo.Update(source); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la source")
	}

	updatedSource, err := s.sourceRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la source mise à jour")
	}

	sourceDTO := s.sourceToDTO(updatedSource)
	return &sourceDTO, nil
}

// Delete supprime une source de demande
func (s *requestSourceService) Delete(id uint) error {
	_, err := s.sourceRepo.FindByID(id)
	if err != nil {
		return errors.New("source introuvable")
	}

	if err := s.sourceRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la source")
	}

	return nil
}

// sourceToDTO convertit un modèle RequestSource en DTO
func (s *requestSourceService) sourceToDTO(source *models.RequestSource) dto.RequestSourceDTO {
	return dto.RequestSourceDTO{
		ID:          source.ID,
		Name:        source.Name,
		Code:        source.Code,
		Description: source.Description,
		IsEnabled:   source.IsEnabled,
		CreatedAt:   source.CreatedAt,
		UpdatedAt:   source.UpdatedAt,
	}
}

