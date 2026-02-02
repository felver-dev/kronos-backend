package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// SoftwareRepository interface pour les opérations sur les logiciels
type SoftwareRepository interface {
	Create(software *models.Software) error
	FindByID(id uint) (*models.Software, error)
	FindByCode(code string) (*models.Software, error)                       // Premier logiciel avec ce code (toutes versions)
	FindByCodeAndVersion(code, version string) (*models.Software, error)    // Un logiciel précis (code + version)
	FindAll() ([]models.Software, error)
	FindActive() ([]models.Software, error)
	Update(software *models.Software) error
	Delete(id uint) error
}

// softwareRepository implémente SoftwareRepository
type softwareRepository struct{}

// NewSoftwareRepository crée une nouvelle instance de SoftwareRepository
func NewSoftwareRepository() SoftwareRepository {
	return &softwareRepository{}
}

// Create crée un nouveau logiciel
func (r *softwareRepository) Create(software *models.Software) error {
	return database.DB.Create(software).Error
}

// FindByID trouve un logiciel par son ID
func (r *softwareRepository) FindByID(id uint) (*models.Software, error) {
	var software models.Software
	err := database.DB.First(&software, id).Error
	if err != nil {
		return nil, err
	}
	return &software, nil
}

// FindByCode trouve le premier logiciel avec ce code (plusieurs versions possibles)
func (r *softwareRepository) FindByCode(code string) (*models.Software, error) {
	var software models.Software
	err := database.DB.Where("code = ?", code).First(&software).Error
	if err != nil {
		return nil, err
	}
	return &software, nil
}

// FindByCodeAndVersion trouve un logiciel par code et version (version vide et NULL assimilés)
func (r *softwareRepository) FindByCodeAndVersion(code, version string) (*models.Software, error) {
	var software models.Software
	err := database.DB.Where("code = ? AND COALESCE(version, '') = ?", code, version).First(&software).Error
	if err != nil {
		return nil, err
	}
	return &software, nil
}

// FindAll récupère tous les logiciels
func (r *softwareRepository) FindAll() ([]models.Software, error) {
	var software []models.Software
	err := database.DB.Order("name ASC").Find(&software).Error
	return software, err
}

// FindActive récupère tous les logiciels actifs
func (r *softwareRepository) FindActive() ([]models.Software, error) {
	var software []models.Software
	err := database.DB.Where("is_active = ?", true).Order("name ASC").Find(&software).Error
	return software, err
}

// Update met à jour un logiciel
func (r *softwareRepository) Update(software *models.Software) error {
	return database.DB.Save(software).Error
}

// Delete supprime un logiciel (soft delete)
func (r *softwareRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Software{}, id).Error
}
