package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// FilialeRepository interface pour les opérations sur les filiales
type FilialeRepository interface {
	Create(filiale *models.Filiale) error
	FindByID(id uint) (*models.Filiale, error)
	FindByCode(code string) (*models.Filiale, error)
	FindAll() ([]models.Filiale, error)
	FindActive() ([]models.Filiale, error)
	FindSoftwareProvider() (*models.Filiale, error)
	Update(filiale *models.Filiale) error
	Delete(id uint) error
}

// filialeRepository implémente FilialeRepository
type filialeRepository struct{}

// NewFilialeRepository crée une nouvelle instance de FilialeRepository
func NewFilialeRepository() FilialeRepository {
	return &filialeRepository{}
}

// Create crée une nouvelle filiale
func (r *filialeRepository) Create(filiale *models.Filiale) error {
	return database.DB.Create(filiale).Error
}

// FindByID trouve une filiale par son ID
func (r *filialeRepository) FindByID(id uint) (*models.Filiale, error) {
	var filiale models.Filiale
	err := database.DB.First(&filiale, id).Error
	if err != nil {
		return nil, err
	}
	return &filiale, nil
}

// FindByCode trouve une filiale par son code
func (r *filialeRepository) FindByCode(code string) (*models.Filiale, error) {
	var filiale models.Filiale
	err := database.DB.Where("code = ?", code).First(&filiale).Error
	if err != nil {
		return nil, err
	}
	return &filiale, nil
}

// FindAll récupère toutes les filiales
func (r *filialeRepository) FindAll() ([]models.Filiale, error) {
	var filiales []models.Filiale
	err := database.DB.Order("name ASC").Find(&filiales).Error
	return filiales, err
}

// FindActive récupère toutes les filiales actives
func (r *filialeRepository) FindActive() ([]models.Filiale, error) {
	var filiales []models.Filiale
	err := database.DB.Where("is_active = ?", true).Order("name ASC").Find(&filiales).Error
	return filiales, err
}

// FindSoftwareProvider trouve la filiale marquée comme fournisseur de logiciels / IT
func (r *filialeRepository) FindSoftwareProvider() (*models.Filiale, error) {
	var filiale models.Filiale
	err := database.DB.Where("is_mci_care_ci = ?", true).First(&filiale).Error
	if err != nil {
		return nil, err
	}
	return &filiale, nil
}

// Update met à jour une filiale
func (r *filialeRepository) Update(filiale *models.Filiale) error {
	return database.DB.Save(filiale).Error
}

// Delete supprime une filiale (soft delete)
func (r *filialeRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Filiale{}, id).Error
}
