package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// AssetSoftwareRepository interface pour les opérations sur les logiciels installés
type AssetSoftwareRepository interface {
	Create(software *models.AssetSoftware) error
	FindByID(id uint) (*models.AssetSoftware, error)
	FindAll() ([]models.AssetSoftware, error)
	FindByAssetID(assetID uint) ([]models.AssetSoftware, error)
	FindBySoftwareName(softwareName string) ([]models.AssetSoftware, error)
	FindBySoftwareNameAndVersion(softwareName, version string) ([]models.AssetSoftware, error)
	Update(software *models.AssetSoftware) error
	Delete(id uint) error
	GetStatistics() (map[string]interface{}, error)
}

// assetSoftwareRepository implémente AssetSoftwareRepository
type assetSoftwareRepository struct{}

// NewAssetSoftwareRepository crée une nouvelle instance de AssetSoftwareRepository
func NewAssetSoftwareRepository() AssetSoftwareRepository {
	return &assetSoftwareRepository{}
}

// Create crée un nouveau logiciel installé
func (r *assetSoftwareRepository) Create(software *models.AssetSoftware) error {
	return database.DB.Create(software).Error
}

// FindByID trouve un logiciel installé par son ID
func (r *assetSoftwareRepository) FindByID(id uint) (*models.AssetSoftware, error) {
	var software models.AssetSoftware
	err := database.DB.
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.AssignedTo.Role").
		First(&software, id).Error
	if err != nil {
		return nil, err
	}
	return &software, nil
}

// FindAll trouve tous les logiciels installés
func (r *assetSoftwareRepository) FindAll() ([]models.AssetSoftware, error) {
	var software []models.AssetSoftware
	err := database.DB.
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.AssignedTo.Role").
		Order("software_name ASC, version ASC").
		Find(&software).Error
	return software, err
}

// FindByAssetID trouve tous les logiciels installés sur un actif
func (r *assetSoftwareRepository) FindByAssetID(assetID uint) ([]models.AssetSoftware, error) {
	var software []models.AssetSoftware
	err := database.DB.Where("asset_id = ?", assetID).Order("software_name ASC").Find(&software).Error
	return software, err
}

// FindBySoftwareName trouve tous les actifs ayant un logiciel spécifique installé
func (r *assetSoftwareRepository) FindBySoftwareName(softwareName string) ([]models.AssetSoftware, error) {
	var software []models.AssetSoftware
	err := database.DB.
		Where("software_name = ?", softwareName).
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.AssignedTo.Role").
		Find(&software).Error
	return software, err
}

// FindBySoftwareNameAndVersion trouve tous les actifs ayant un logiciel avec une version spécifique
func (r *assetSoftwareRepository) FindBySoftwareNameAndVersion(softwareName, version string) ([]models.AssetSoftware, error) {
	var software []models.AssetSoftware
	err := database.DB.
		Where("software_name = ? AND version = ?", softwareName, version).
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.AssignedTo.Role").
		Find(&software).Error
	return software, err
}

// Update met à jour un logiciel installé
func (r *assetSoftwareRepository) Update(software *models.AssetSoftware) error {
	return database.DB.Save(software).Error
}

// Delete supprime un logiciel installé
func (r *assetSoftwareRepository) Delete(id uint) error {
	return database.DB.Delete(&models.AssetSoftware{}, id).Error
}

// GetStatistics récupère des statistiques sur les logiciels installés
func (r *assetSoftwareRepository) GetStatistics() (map[string]interface{}, error) {
	var results []struct {
		SoftwareName string
		Version      string
		Count        int64
		CategoryName string
	}

	err := database.DB.Table("asset_software").
		Select("asset_software.software_name, asset_software.version, COUNT(*) as count, asset_categories.name as category_name").
		Joins("JOIN assets ON asset_software.asset_id = assets.id").
		Joins("JOIN asset_categories ON assets.category_id = asset_categories.id").
		Group("asset_software.software_name, asset_software.version, asset_categories.name").
		Order("count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})
	stats["by_software"] = results

	// Statistiques par logiciel (toutes versions confondues)
	var softwareStats []struct {
		SoftwareName string
		Count        int64
	}
	err = database.DB.Table("asset_software").
		Select("software_name, COUNT(*) as count").
		Group("software_name").
		Order("count DESC").
		Scan(&softwareStats).Error

	if err == nil {
		stats["by_software_name"] = softwareStats
	}

	return stats, nil
}
