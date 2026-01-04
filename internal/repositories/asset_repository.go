package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
)

// AssetRepository interface pour les opérations sur les actifs IT
type AssetRepository interface {
	Create(asset *models.Asset) error
	FindByID(id uint) (*models.Asset, error)
	FindAll() ([]models.Asset, error)
	FindByCategory(categoryID uint) ([]models.Asset, error)
	FindByStatus(status string) ([]models.Asset, error)
	FindByAssignedTo(userID uint) ([]models.Asset, error)
	FindBySerialNumber(serialNumber string) (*models.Asset, error)
	Search(query string, category string, limit int) ([]models.Asset, error)
	Update(asset *models.Asset) error
	Delete(id uint) error
}

// AssetCategoryRepository interface pour les opérations sur les catégories d'actifs
type AssetCategoryRepository interface {
	Create(category *models.AssetCategory) error
	FindByID(id uint) (*models.AssetCategory, error)
	FindAll() ([]models.AssetCategory, error)
	FindByParentID(parentID uint) ([]models.AssetCategory, error)
	Update(category *models.AssetCategory) error
	Delete(id uint) error
}

// assetRepository implémente AssetRepository
type assetRepository struct{}

// assetCategoryRepository implémente AssetCategoryRepository
type assetCategoryRepository struct{}

// NewAssetRepository crée une nouvelle instance de AssetRepository
func NewAssetRepository() AssetRepository {
	return &assetRepository{}
}

// NewAssetCategoryRepository crée une nouvelle instance de AssetCategoryRepository
func NewAssetCategoryRepository() AssetCategoryRepository {
	return &assetCategoryRepository{}
}

// Create crée un nouvel actif
func (r *assetRepository) Create(asset *models.Asset) error {
	return database.DB.Create(asset).Error
}

// FindByID trouve un actif par son ID
func (r *assetRepository) FindByID(id uint) (*models.Asset, error) {
	var asset models.Asset
	err := database.DB.Preload("Category").Preload("AssignedTo").Preload("AssignedTo.Role").First(&asset, id).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// FindAll récupère tous les actifs
func (r *assetRepository) FindAll() ([]models.Asset, error) {
	var assets []models.Asset
	err := database.DB.Preload("Category").Preload("AssignedTo").Find(&assets).Error
	return assets, err
}

// FindByCategory récupère les actifs par catégorie
func (r *assetRepository) FindByCategory(categoryID uint) ([]models.Asset, error) {
	var assets []models.Asset
	err := database.DB.Preload("Category").Preload("AssignedTo").Where("category_id = ?", categoryID).Find(&assets).Error
	return assets, err
}

// FindByStatus récupère les actifs par statut
func (r *assetRepository) FindByStatus(status string) ([]models.Asset, error) {
	var assets []models.Asset
	err := database.DB.Preload("Category").Preload("AssignedTo").Where("status = ?", status).Find(&assets).Error
	return assets, err
}

// FindByAssignedTo récupère les actifs assignés à un utilisateur
func (r *assetRepository) FindByAssignedTo(userID uint) ([]models.Asset, error) {
	var assets []models.Asset
	err := database.DB.Preload("Category").Preload("AssignedTo").Where("assigned_to_id = ?", userID).Find(&assets).Error
	return assets, err
}

// FindBySerialNumber trouve un actif par son numéro de série
func (r *assetRepository) FindBySerialNumber(serialNumber string) (*models.Asset, error) {
	var asset models.Asset
	err := database.DB.Preload("Category").Preload("AssignedTo").Where("serial_number = ?", serialNumber).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// Update met à jour un actif
func (r *assetRepository) Update(asset *models.Asset) error {
	return database.DB.Save(asset).Error
}

// Delete supprime un actif (soft delete)
func (r *assetRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Asset{}, id).Error
}

// Search recherche des actifs par nom, description ou numéro de série
func (r *assetRepository) Search(query string, category string, limit int) ([]models.Asset, error) {
	var assets []models.Asset
	searchPattern := "%" + query + "%"
	
	db := database.DB.Preload("Category").Preload("AssignedTo").Preload("AssignedTo.Role").
		Where("name LIKE ? OR description LIKE ? OR serial_number LIKE ?", searchPattern, searchPattern, searchPattern)
	
	if category != "" {
		db = db.Joins("JOIN asset_categories ON assets.category_id = asset_categories.id").
			Where("asset_categories.name = ?", category)
	}
	
	if limit > 0 {
		db = db.Limit(limit)
	}
	
	err := db.Order("created_at DESC").Find(&assets).Error
	return assets, err
}

// Create crée une nouvelle catégorie d'actif
func (r *assetCategoryRepository) Create(category *models.AssetCategory) error {
	return database.DB.Create(category).Error
}

// FindByID trouve une catégorie par son ID
func (r *assetCategoryRepository) FindByID(id uint) (*models.AssetCategory, error) {
	var category models.AssetCategory
	err := database.DB.Preload("Parent").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindAll récupère toutes les catégories
func (r *assetCategoryRepository) FindAll() ([]models.AssetCategory, error) {
	var categories []models.AssetCategory
	err := database.DB.Preload("Parent").Find(&categories).Error
	return categories, err
}

// FindByParentID récupère les catégories enfants d'une catégorie parente
func (r *assetCategoryRepository) FindByParentID(parentID uint) ([]models.AssetCategory, error) {
	var categories []models.AssetCategory
	err := database.DB.Where("parent_id = ?", parentID).Find(&categories).Error
	return categories, err
}

// Update met à jour une catégorie
func (r *assetCategoryRepository) Update(category *models.AssetCategory) error {
	return database.DB.Save(category).Error
}

// Delete supprime une catégorie
func (r *assetCategoryRepository) Delete(id uint) error {
	return database.DB.Delete(&models.AssetCategory{}, id).Error
}
