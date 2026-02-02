package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// FilialeSoftwareRepository interface pour les opérations sur les déploiements de logiciels
type FilialeSoftwareRepository interface {
	Create(deployment *models.FilialeSoftware) error
	FindByID(id uint) (*models.FilialeSoftware, error)
	FindAll() ([]models.FilialeSoftware, error)
	FindByFilialeID(filialeID uint) ([]models.FilialeSoftware, error)
	FindBySoftwareID(softwareID uint) ([]models.FilialeSoftware, error)
	FindByFilialeAndSoftware(filialeID, softwareID uint) (*models.FilialeSoftware, error)
	FindActive() ([]models.FilialeSoftware, error)
	FindActiveByFiliale(filialeID uint) ([]models.FilialeSoftware, error)
	FindActiveBySoftware(softwareID uint) ([]models.FilialeSoftware, error)
	Update(deployment *models.FilialeSoftware) error
	Delete(id uint) error
}

// filialeSoftwareRepository implémente FilialeSoftwareRepository
type filialeSoftwareRepository struct{}

// NewFilialeSoftwareRepository crée une nouvelle instance de FilialeSoftwareRepository
func NewFilialeSoftwareRepository() FilialeSoftwareRepository {
	return &filialeSoftwareRepository{}
}

// Create crée un nouveau déploiement
func (r *filialeSoftwareRepository) Create(deployment *models.FilialeSoftware) error {
	return database.DB.Create(deployment).Error
}

// FindByID trouve un déploiement par son ID
func (r *filialeSoftwareRepository) FindByID(id uint) (*models.FilialeSoftware, error) {
	var deployment models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		First(&deployment, id).Error
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

// FindAll récupère tous les déploiements
func (r *filialeSoftwareRepository) FindAll() ([]models.FilialeSoftware, error) {
	var deployments []models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		Order("created_at DESC").
		Find(&deployments).Error
	return deployments, err
}

// FindByFilialeID récupère tous les déploiements d'une filiale
func (r *filialeSoftwareRepository) FindByFilialeID(filialeID uint) ([]models.FilialeSoftware, error) {
	var deployments []models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		Where("filiale_id = ?", filialeID).
		Order("created_at DESC").
		Find(&deployments).Error
	return deployments, err
}

// FindBySoftwareID récupère tous les déploiements d'un logiciel
func (r *filialeSoftwareRepository) FindBySoftwareID(softwareID uint) ([]models.FilialeSoftware, error) {
	var deployments []models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		Where("software_id = ?", softwareID).
		Order("created_at DESC").
		Find(&deployments).Error
	return deployments, err
}

// FindByFilialeAndSoftware trouve un déploiement spécifique pour une filiale et un logiciel
func (r *filialeSoftwareRepository) FindByFilialeAndSoftware(filialeID, softwareID uint) (*models.FilialeSoftware, error) {
	var deployment models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		Where("filiale_id = ? AND software_id = ?", filialeID, softwareID).
		First(&deployment).Error
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

// FindActive récupère tous les déploiements actifs
func (r *filialeSoftwareRepository) FindActive() ([]models.FilialeSoftware, error) {
	var deployments []models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		Where("is_active = ?", true).
		Order("created_at DESC").
		Find(&deployments).Error
	return deployments, err
}

// FindActiveByFiliale récupère les déploiements actifs d'une filiale
func (r *filialeSoftwareRepository) FindActiveByFiliale(filialeID uint) ([]models.FilialeSoftware, error) {
	var deployments []models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		Where("filiale_id = ? AND is_active = ?", filialeID, true).
		Order("created_at DESC").
		Find(&deployments).Error
	return deployments, err
}

// FindActiveBySoftware récupère les déploiements actifs d'un logiciel
func (r *filialeSoftwareRepository) FindActiveBySoftware(softwareID uint) ([]models.FilialeSoftware, error) {
	var deployments []models.FilialeSoftware
	err := database.DB.Preload("Filiale").Preload("Software").
		Where("software_id = ? AND is_active = ?", softwareID, true).
		Order("created_at DESC").
		Find(&deployments).Error
	return deployments, err
}

// Update met à jour un déploiement
func (r *filialeSoftwareRepository) Update(deployment *models.FilialeSoftware) error {
	return database.DB.Save(deployment).Error
}

// Delete supprime un déploiement (soft delete)
func (r *filialeSoftwareRepository) Delete(id uint) error {
	return database.DB.Delete(&models.FilialeSoftware{}, id).Error
}
