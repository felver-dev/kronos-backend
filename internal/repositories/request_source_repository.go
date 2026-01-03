package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// RequestSourceRepository interface pour les opérations sur les sources de demande
type RequestSourceRepository interface {
	Create(source *models.RequestSource) error
	FindByID(id uint) (*models.RequestSource, error)
	FindByCode(code string) (*models.RequestSource, error)
	FindAll() ([]models.RequestSource, error)
	FindEnabled() ([]models.RequestSource, error)
	Update(source *models.RequestSource) error
	Delete(id uint) error
}

// requestSourceRepository implémente RequestSourceRepository
type requestSourceRepository struct{}

// NewRequestSourceRepository crée une nouvelle instance de RequestSourceRepository
func NewRequestSourceRepository() RequestSourceRepository {
	return &requestSourceRepository{}
}

// Create crée une nouvelle source de demande
func (r *requestSourceRepository) Create(source *models.RequestSource) error {
	return database.DB.Create(source).Error
}

// FindByID trouve une source par son ID
func (r *requestSourceRepository) FindByID(id uint) (*models.RequestSource, error) {
	var source models.RequestSource
	err := database.DB.First(&source, id).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

// FindByCode trouve une source par son code
func (r *requestSourceRepository) FindByCode(code string) (*models.RequestSource, error) {
	var source models.RequestSource
	err := database.DB.Where("code = ?", code).First(&source).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

// FindAll récupère toutes les sources
func (r *requestSourceRepository) FindAll() ([]models.RequestSource, error) {
	var sources []models.RequestSource
	err := database.DB.Order("name ASC").Find(&sources).Error
	return sources, err
}

// FindEnabled récupère toutes les sources activées
func (r *requestSourceRepository) FindEnabled() ([]models.RequestSource, error) {
	var sources []models.RequestSource
	err := database.DB.Where("is_enabled = ?", true).Order("name ASC").Find(&sources).Error
	return sources, err
}

// Update met à jour une source
func (r *requestSourceRepository) Update(source *models.RequestSource) error {
	return database.DB.Save(source).Error
}

// Delete supprime une source
func (r *requestSourceRepository) Delete(id uint) error {
	return database.DB.Delete(&models.RequestSource{}, id).Error
}
