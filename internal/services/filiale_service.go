package services

import (
	"errors"

	"github.com/mcicare/itsm-backend/internal/dto"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/repositories"
)

// FilialeService interface pour les opérations sur les filiales
type FilialeService interface {
	Create(req dto.CreateFilialeRequest) (*dto.FilialeDTO, error)
	GetByID(id uint) (*dto.FilialeDTO, error)
	GetByCode(code string) (*dto.FilialeDTO, error)
	GetAll() ([]dto.FilialeDTO, error)
	GetActive() ([]dto.FilialeDTO, error)
	GetSoftwareProvider() (*dto.FilialeDTO, error)
	Update(id uint, req dto.UpdateFilialeRequest) (*dto.FilialeDTO, error)
	Delete(id uint) error
}

// filialeService implémente FilialeService
type filialeService struct {
	filialeRepo repositories.FilialeRepository
}

// NewFilialeService crée une nouvelle instance de FilialeService
func NewFilialeService(filialeRepo repositories.FilialeRepository) FilialeService {
	return &filialeService{
		filialeRepo: filialeRepo,
	}
}

// Create crée une nouvelle filiale
func (s *filialeService) Create(req dto.CreateFilialeRequest) (*dto.FilialeDTO, error) {
	// Vérifier si le code existe déjà
	existing, _ := s.filialeRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, errors.New("une filiale avec ce code existe déjà")
	}

	filiale := &models.Filiale{
		Code:        req.Code,
		Name:        req.Name,
		Country:     req.Country,
		City:        req.City,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		IsActive:    true,
		IsSoftwareProvider: req.IsSoftwareProvider,
	}

	if err := s.filialeRepo.Create(filiale); err != nil {
		return nil, errors.New("erreur lors de la création de la filiale")
	}

	// Récupérer la filiale créée
	createdFiliale, err := s.filialeRepo.FindByID(filiale.ID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la filiale créée")
	}

	return s.filialeToDTO(createdFiliale), nil
}

// GetByID récupère une filiale par son ID
func (s *filialeService) GetByID(id uint) (*dto.FilialeDTO, error) {
	filiale, err := s.filialeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("filiale introuvable")
	}

	return s.filialeToDTO(filiale), nil
}

// GetByCode récupère une filiale par son code
func (s *filialeService) GetByCode(code string) (*dto.FilialeDTO, error) {
	filiale, err := s.filialeRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("filiale introuvable")
	}

	return s.filialeToDTO(filiale), nil
}

// GetAll récupère toutes les filiales
func (s *filialeService) GetAll() ([]dto.FilialeDTO, error) {
	filiales, err := s.filialeRepo.FindAll()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des filiales")
	}

	var filialeDTOs []dto.FilialeDTO
	for _, filiale := range filiales {
		filialeDTOs = append(filialeDTOs, *s.filialeToDTO(&filiale))
	}

	return filialeDTOs, nil
}

// GetActive récupère toutes les filiales actives
func (s *filialeService) GetActive() ([]dto.FilialeDTO, error) {
	filiales, err := s.filialeRepo.FindActive()
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des filiales actives")
	}

	var filialeDTOs []dto.FilialeDTO
	for _, filiale := range filiales {
		filialeDTOs = append(filialeDTOs, *s.filialeToDTO(&filiale))
	}

	return filialeDTOs, nil
}

// GetSoftwareProvider récupère la filiale fournisseur de logiciels / IT (is_software_provider=true)
func (s *filialeService) GetSoftwareProvider() (*dto.FilialeDTO, error) {
	filiale, err := s.filialeRepo.FindSoftwareProvider()
	if err != nil {
		return nil, errors.New("filiale fournisseur de logiciels introuvable")
	}

	return s.filialeToDTO(filiale), nil
}

// Update met à jour une filiale
func (s *filialeService) Update(id uint, req dto.UpdateFilialeRequest) (*dto.FilialeDTO, error) {
	filiale, err := s.filialeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("filiale introuvable")
	}

	// Mettre à jour les champs fournis
	if req.Name != "" {
		filiale.Name = req.Name
	}
	if req.Country != "" {
		filiale.Country = req.Country
	}
	if req.City != "" {
		filiale.City = req.City
	}
	if req.Address != nil {
		filiale.Address = req.Address
	}
	if req.Phone != "" {
		filiale.Phone = req.Phone
	}
	if req.Email != "" {
		filiale.Email = req.Email
	}
	if req.IsActive != nil {
		filiale.IsActive = *req.IsActive
	}
	if req.IsSoftwareProvider != nil {
		filiale.IsSoftwareProvider = *req.IsSoftwareProvider
	}

	if err := s.filialeRepo.Update(filiale); err != nil {
		return nil, errors.New("erreur lors de la mise à jour de la filiale")
	}

	// Récupérer la filiale mise à jour
	updatedFiliale, err := s.filialeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération de la filiale mise à jour")
	}

	return s.filialeToDTO(updatedFiliale), nil
}

// Delete supprime une filiale
func (s *filialeService) Delete(id uint) error {
	_, err := s.filialeRepo.FindByID(id)
	if err != nil {
		return errors.New("filiale introuvable")
	}

	if err := s.filialeRepo.Delete(id); err != nil {
		return errors.New("erreur lors de la suppression de la filiale")
	}

	return nil
}

// filialeToDTO convertit un modèle Filiale en DTO
func (s *filialeService) filialeToDTO(filiale *models.Filiale) *dto.FilialeDTO {
	return &dto.FilialeDTO{
		ID:          filiale.ID,
		Code:        filiale.Code,
		Name:        filiale.Name,
		Country:     filiale.Country,
		City:        filiale.City,
		Address:     filiale.Address,
		Phone:       filiale.Phone,
		Email:       filiale.Email,
		IsActive:    filiale.IsActive,
		IsSoftwareProvider: filiale.IsSoftwareProvider,
		CreatedAt:   filiale.CreatedAt,
		UpdatedAt:   filiale.UpdatedAt,
	}
}
