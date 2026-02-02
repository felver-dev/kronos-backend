package services

import (
	"errors"
	"strings"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// SoftwareService interface pour les opérations sur les logiciels
type SoftwareService interface {
	Create(req dto.CreateSoftwareRequest) (*dto.SoftwareDTO, error)
	GetByID(id uint) (*dto.SoftwareDTO, error)
	GetByCode(code string) (*dto.SoftwareDTO, error)
	GetAll() ([]dto.SoftwareDTO, error)
	GetActive() ([]dto.SoftwareDTO, error)
	Update(id uint, req dto.UpdateSoftwareRequest) (*dto.SoftwareDTO, error)
	Delete(id uint) error
}

// softwareService implémente SoftwareService
type softwareService struct {
	softwareRepo repositories.SoftwareRepository
}

// NewSoftwareService crée une nouvelle instance de SoftwareService
func NewSoftwareService(softwareRepo repositories.SoftwareRepository) SoftwareService {
	return &softwareService{
		softwareRepo: softwareRepo,
	}
}

// Create crée un nouveau logiciel (même code autorisé si version différente)
func (s *softwareService) Create(req dto.CreateSoftwareRequest) (*dto.SoftwareDTO, error) {
	version := strings.TrimSpace(req.Version)
	// Unicité sur (code, version) : vérifier si cette combinaison existe déjà
	existing, _ := s.softwareRepo.FindByCodeAndVersion(req.Code, version)
	if existing != nil {
		return nil, errors.New("un logiciel avec ce code et cette version existe déjà")
	}

	software := &models.Software{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Version:     version,
		IsActive:    true,
	}

	if err := s.softwareRepo.Create(software); err != nil {
		return nil, errors.New("erreur lors de la création du logiciel")
	}

	// Récupérer le logiciel créé
	createdSoftware, err := s.softwareRepo.FindByID(software.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du logiciel créé")
	}

	return s.softwareToDTO(createdSoftware), nil
}

// GetByID récupère un logiciel par son ID
func (s *softwareService) GetByID(id uint) (*dto.SoftwareDTO, error) {
	software, err := s.softwareRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("logiciel introuvable")
	}

	return s.softwareToDTO(software), nil
}

// GetByCode récupère un logiciel par son code
func (s *softwareService) GetByCode(code string) (*dto.SoftwareDTO, error) {
	software, err := s.softwareRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("logiciel introuvable")
	}

	return s.softwareToDTO(software), nil
}

// GetAll récupère tous les logiciels
func (s *softwareService) GetAll() ([]dto.SoftwareDTO, error) {
	software, err := s.softwareRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logiciels")
	}

	var softwareDTOs []dto.SoftwareDTO
	for _, sw := range software {
		softwareDTOs = append(softwareDTOs, *s.softwareToDTO(&sw))
	}

	return softwareDTOs, nil
}

// GetActive récupère tous les logiciels actifs
func (s *softwareService) GetActive() ([]dto.SoftwareDTO, error) {
	software, err := s.softwareRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des logiciels actifs")
	}

	var softwareDTOs []dto.SoftwareDTO
	for _, sw := range software {
		softwareDTOs = append(softwareDTOs, *s.softwareToDTO(&sw))
	}

	return softwareDTOs, nil
}

// Update met à jour un logiciel
func (s *softwareService) Update(id uint, req dto.UpdateSoftwareRequest) (*dto.SoftwareDTO, error) {
	software, err := s.softwareRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("logiciel introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Name != "" {
		software.Name = req.Name
	}
	if req.Description != nil {
		software.Description = req.Description
	}
	if req.Version != "" {
		software.Version = req.Version
	}
	if req.IsActive != nil {
		software.IsActive = *req.IsActive
	}

	if err := s.softwareRepo.Update(software); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du logiciel")
	}

	// Récupérer le logiciel mis à jour
	updatedSoftware, err := s.softwareRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération du logiciel mis à jour")
	}

	return s.softwareToDTO(updatedSoftware), nil
}

// Delete supprime un logiciel
func (s *softwareService) Delete(id uint) error {
	_, err := s.softwareRepo.FindByID(id)
	if err != nil {
		return errors.New("logiciel introuvable")
	}

	if err := s.softwareRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression du logiciel")
	}

	return nil
}

// softwareToDTO convertit un modèle Software en DTO
func (s *softwareService) softwareToDTO(software *models.Software) *dto.SoftwareDTO {
	return &dto.SoftwareDTO{
		ID:          software.ID,
		Code:        software.Code,
		Name:        software.Name,
		Description: software.Description,
		Version:     software.Version,
		IsActive:    software.IsActive,
		CreatedAt:   software.CreatedAt,
		UpdatedAt:   software.UpdatedAt,
	}
}
