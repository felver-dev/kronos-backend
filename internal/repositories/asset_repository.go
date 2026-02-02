package repositories

import (
	"github.com/mcicare/itsm-backend/database"
	"github.com/mcicare/itsm-backend/internal/models"
	"github.com/mcicare/itsm-backend/internal/scope"
)

// AssetRepository interface pour les opérations sur les actifs IT
type AssetRepository interface {
	Create(asset *models.Asset) error
	FindByID(id uint) (*models.Asset, error)
	FindAll(scope interface{}) ([]models.Asset, error) // scope peut être *scope.QueryScope ou nil
	FindByCategory(scope interface{}, categoryID uint) ([]models.Asset, error)
	CountByCategory(categoryID uint) (int64, error)
	FindByStatus(scope interface{}, status string) ([]models.Asset, error)
	FindByAssignedTo(userID uint) ([]models.Asset, error)
	FindBySerialNumber(serialNumber string) (*models.Asset, error)
	Search(scope interface{}, query string, category string, limit int) ([]models.Asset, error) // scope peut être *scope.QueryScope ou nil
	Update(asset *models.Asset) error
	Delete(id uint) error
}

// AssetCategoryRepository interface pour les opérations sur les catégories d'actifs
type AssetCategoryRepository interface {
	Create(category *models.AssetCategory) error
	FindByID(id uint) (*models.AssetCategory, error)
	FindAll() ([]models.AssetCategory, error)
	FindPaginated(page, limit int) ([]models.AssetCategory, int64, error)
	FindByParentID(parentID uint) ([]models.AssetCategory, error)
	CountByParentID(parentID uint) (int64, error)
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
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *assetRepository) FindAll(scopeParam interface{}) ([]models.Asset, error) {
	var assets []models.Asset
	
	// Construire la requête de base
	query := database.DB.Model(&models.Asset{}).Preload("Category").Preload("AssignedTo")
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAssetScope(query, queryScope)
		}
	}
	
	err := query.Find(&assets).Error
	return assets, err
}

// FindByCategory récupère les actifs par catégorie
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *assetRepository) FindByCategory(scopeParam interface{}, categoryID uint) ([]models.Asset, error) {
	var assets []models.Asset
	
	// Construire la requête de base
	query := database.DB.Model(&models.Asset{}).Where("category_id = ?", categoryID)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAssetScope(query, queryScope)
		}
	}
	
	err := query.Find(&assets).Error
	return assets, err
}

// CountByCategory compte le nombre d'actifs par catégorie
func (r *assetRepository) CountByCategory(categoryID uint) (int64, error) {
	var count int64
	// Utiliser une requête explicite pour compter les actifs avec cette catégorie
	// Note: GORM gère automatiquement les soft deletes avec DeletedAt
	err := database.DB.Model(&models.Asset{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindByStatus récupère les actifs par statut
// Le scope est utilisé pour filtrer automatiquement selon les permissions de l'utilisateur
func (r *assetRepository) FindByStatus(scopeParam interface{}, status string) ([]models.Asset, error) {
	var assets []models.Asset
	
	// Construire la requête de base
	query := database.DB.Model(&models.Asset{}).
		Preload("Category").Preload("AssignedTo").
		Where("status = ?", status)
	
	// Appliquer le scope si fourni
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			query = scope.ApplyAssetScope(query, queryScope)
		}
	}
	
	err := query.Find(&assets).Error
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
func (r *assetRepository) Search(scopeParam interface{}, query string, category string, limit int) ([]models.Asset, error) {
	var assets []models.Asset
	searchPattern := "%" + query + "%"
	
	// Construire la requête de base
	db := database.DB.Model(&models.Asset{}).
		Preload("Category").Preload("AssignedTo").Preload("AssignedTo.Role").
		Where("assets.name LIKE ? OR assets.description LIKE ? OR assets.serial_number LIKE ?", searchPattern, searchPattern, searchPattern)
	
	// Appliquer le scope si fourni (doit être fait avant les autres filtres)
	if scopeParam != nil {
		if queryScope, ok := scopeParam.(*scope.QueryScope); ok {
			db = scope.ApplyAssetScope(db, queryScope)
		}
	}
	
	if category != "" {
		db = db.Joins("JOIN asset_categories ON assets.category_id = asset_categories.id").
			Where("asset_categories.name = ?", category)
	}
	
	if limit > 0 {
		db = db.Limit(limit)
	}
	
	err := db.Order("assets.created_at DESC").Find(&assets).Error
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

// FindPaginated récupère les catégories avec pagination
func (r *assetCategoryRepository) FindPaginated(page, limit int) ([]models.AssetCategory, int64, error) {
	var categories []models.AssetCategory
	var total int64

	// Compter le total
	err := database.DB.Model(&models.AssetCategory{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculer l'offset
	offset := (page - 1) * limit

	// Récupérer les catégories avec pagination
	err = database.DB.Preload("Parent").
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// FindByParentID récupère les catégories enfants d'une catégorie parente
func (r *assetCategoryRepository) FindByParentID(parentID uint) ([]models.AssetCategory, error) {
	var categories []models.AssetCategory
	err := database.DB.Where("parent_id = ?", parentID).Find(&categories).Error
	return categories, err
}

// CountByParentID compte le nombre de catégories enfants d'une catégorie parente
func (r *assetCategoryRepository) CountByParentID(parentID uint) (int64, error) {
	var count int64
	// Utiliser une requête explicite pour compter les catégories avec ce parent_id
	err := database.DB.Model(&models.AssetCategory{}).
		Where("parent_id = ?", parentID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Update met à jour une catégorie
func (r *assetCategoryRepository) Update(category *models.AssetCategory) error {
	return database.DB.Save(category).Error
}

// Delete supprime une catégorie
func (r *assetCategoryRepository) Delete(id uint) error {
	return database.DB.Delete(&models.AssetCategory{}, id).Error
}
